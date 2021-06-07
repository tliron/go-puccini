package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tliron/kutil/reflection"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca"
)

func (self *Context) GetInheritTasks() Tasks {
	inheritContext := NewInheritContext()
	self.Traverse(logInheritance, func(entityPtr tosca.EntityPtr) bool {
		inheritContext.GetInheritTask(entityPtr)
		return true
	})
	return inheritContext.Tasks
}

//
// InheritContext
//

type InheritContext struct {
	Tasks            Tasks
	TasksForEntities TasksForEntities
	InheritFields    InheritFields
}

func NewInheritContext() *InheritContext {
	return &InheritContext{make(Tasks), make(TasksForEntities), make(InheritFields)}
}

func (self *InheritContext) GetInheritTask(entityPtr tosca.EntityPtr) *Task {
	task, ok := self.TasksForEntities[entityPtr]
	if !ok {
		path := tosca.GetContext(entityPtr).Path.String()
		if path == "" {
			path = "<root>"
		}

		task = NewTask(path)
		self.Tasks.Add(task)
		self.TasksForEntities[entityPtr] = task

		for dependencyEntityPtr := range self.GetDependencies(entityPtr) {
			task.AddDependency(self.GetInheritTask(dependencyEntityPtr))
		}

		task.Executor = self.NewExecutor(entityPtr)

		logInheritance.Debugf("new task: %s (%d)", task.Name, len(task.Dependencies))
	} else {
		logInheritance.Debugf("task cache hit: %s (%d)", task.Name, len(task.Dependencies))
	}
	return task
}

func (self *InheritContext) NewExecutor(entityPtr tosca.EntityPtr) Executor {
	return func(task *Task) {
		defer task.Done()

		logInheritance.Debugf("task: %s", task.Name)

		lock := util.GetLock(entityPtr)
		lock.Lock()
		defer lock.Unlock()

		for _, inheritField := range self.InheritFields.Get(entityPtr) {
			inheritField.Inherit()
		}

		// Custom inheritance after all fields have been inherited
		if inherits, ok := entityPtr.(tosca.Inherits); ok {
			inherits.Inherit()
		}
	}
}

// TODO: rare race condition due to concurrent access of reflect type
func (self *InheritContext) GetDependencies(entityPtr tosca.EntityPtr) map[tosca.EntityPtr]bool {
	dependencies := make(map[tosca.EntityPtr]bool)

	// From "inherit" tags
	for _, inheritField := range self.InheritFields.Get(entityPtr) {
		dependencies[inheritField.FromEntityPtr] = true
	}

	// From field values
	entity := reflect.ValueOf(entityPtr).Elem()
	for _, structField := range reflection.GetStructFields(entity.Type()) {
		// Does this case ever happen?
		// Would conflict with anonymous pointer fields (Go "inheritance")
		//		if reflection.IsPtrToStruct(structField.Type) {
		//			// Compatible with *interface{}
		//			field := entity.FieldByName(structField.Name)
		//			if !field.IsNil() {
		//				e := field.Interface()
		//				// We sometimes have pointers to non-entities, so make sure
		//				if _, ok := e.(tosca.Contextual); ok {
		//					dependencies[e] = true
		//				}
		//			}
		//		}

		if reflection.IsMapOfStringToPtrToStruct(structField.Type) {
			// Compatible with map[string]*interface{}
			field := entity.FieldByName(structField.Name)
			for _, mapKey := range field.MapKeys() {
				element := field.MapIndex(mapKey)
				dependencies[element.Interface()] = true
			}
		} else if reflection.IsSliceOfPtrToStruct(structField.Type) {
			// Compatible with []*interface{}
			field := entity.FieldByName(structField.Name)
			length := field.Len()
			for i := 0; i < length; i++ {
				element := field.Index(i)
				dependencies[element.Interface()] = true
			}
		}
	}

	return dependencies
}

//
// InheritField
//

type InheritField struct {
	Entity        reflect.Value
	FromEntityPtr tosca.EntityPtr
	Key           string
	Field         reflect.Value
	FromField     reflect.Value
}

func (self *InheritField) Inherit() {
	// TODO do we really need all of these? some of them aren't used in TOSCA
	fieldEntityPtr := self.Field.Interface()
	if reflection.IsPtrToString(fieldEntityPtr) {
		self.InheritEntity()
	} else if reflection.IsPtrToInt64(fieldEntityPtr) {
		self.InheritEntity()
	} else if reflection.IsPtrToBool(fieldEntityPtr) {
		self.InheritEntity()
	} else if reflection.IsPtrToSliceOfString(fieldEntityPtr) {
		self.InheritStringsFromSlice()
	} else if reflection.IsPtrToMapOfStringToString(fieldEntityPtr) {
		self.InheritStringsFromMap()
	} else {
		fieldType := self.Field.Type()
		if reflection.IsPtrToStruct(fieldType) {
			self.InheritEntity()
		} else if reflection.IsSliceOfPtrToStruct(fieldType) {
			self.InheritStructsFromSlice()
		} else if reflection.IsMapOfStringToPtrToStruct(fieldType) {
			self.InheritStructsFromMap()
		} else {
			panic(fmt.Sprintf("\"inherit\" tag's field type %q is not supported in struct: %T", fieldType, self.Entity.Interface()))
		}
	}
}

// Field is compatible with *interface{}
func (self *InheritField) InheritEntity() {
	if self.Field.IsNil() && !self.FromField.IsNil() {
		self.Field.Set(self.FromField)
	}
}

// Field is *[]string
func (self *InheritField) InheritStringsFromSlice() {
	slicePtr := self.Field.Interface().(*[]string)
	fromSlicePtr := self.FromField.Interface().(*[]string)

	if fromSlicePtr == nil {
		return
	}

	fromSlice := *fromSlicePtr
	length := len(fromSlice)
	if length == 0 {
		return
	}

	var slice []string
	if slicePtr != nil {
		slice = *slicePtr
	} else {
		slice = make([]string, 0, length)
	}

	for _, s := range fromSlice {
		slice = append(slice, s)
	}

	self.Field.Set(reflect.ValueOf(&slice))
}

// Field is compatible with []*interface{}
func (self *InheritField) InheritStructsFromSlice() {
	slice := self.Field

	length := self.FromField.Len()
	for i := 0; i < length; i++ {
		element := self.FromField.Index(i)

		if _, ok := element.Interface().(tosca.Mappable); ok {
			// For mappable elements only, *don't* inherit the same key
			// (We'll merge everything else)
			key := tosca.GetKey(element.Interface())
			if ii, ok := getSliceElementIndexForKey(self.Field, key); ok {
				e := self.Field.Index(ii)
				logInheritance.Debugf("override: %s", tosca.GetContext(e.Interface()).Path)
				continue
			}
		}

		slice = reflect.Append(slice, element)
	}

	self.Field.Set(slice)
}

// Field is *map[string]string
func (self *InheritField) InheritStringsFromMap() {
	mapPtr := self.Field.Interface().(*map[string]string)
	fromMapPtr := self.FromField.Interface().(*map[string]string)

	if fromMapPtr == nil {
		return
	}

	fromMap := *fromMapPtr
	length := len(fromMap)
	if length == 0 {
		return
	}

	var m map[string]string
	if mapPtr != nil {
		m = *mapPtr
	} else {
		m = make(map[string]string)
	}

	for k, v := range fromMap {
		_, ok := m[k]
		if !ok {
			m[k] = v
		}
	}

	self.Field.Set(reflect.ValueOf(&m))
}

// Field is compatible with map[string]*interface{}
func (self *InheritField) InheritStructsFromMap() {
	for _, mapKey := range self.FromField.MapKeys() {
		element := self.FromField.MapIndex(mapKey)
		e := self.Field.MapIndex(mapKey)
		if e.IsValid() {
			// We are overriding this element, so don't inherit it
			logInheritance.Debugf("override: %s", tosca.GetContext(e.Interface()).Path)
		} else {
			self.Field.SetMapIndex(mapKey, element)
		}
	}
}

func getSliceElementIndexForKey(slice reflect.Value, key string) (int, bool) {
	length := slice.Len()
	for i := 0; i < length; i++ {
		element := slice.Index(i)
		if tosca.GetKey(element.Elem()) == key {
			return i, true
		}
	}
	return -1, false
}

//
// InheritFields
//

type InheritFields map[tosca.EntityPtr][]*InheritField

// From "inherit" tags
func NewInheritFields(entityPtr tosca.EntityPtr) []*InheritField {
	var inheritFields []*InheritField

	entity := reflect.ValueOf(entityPtr).Elem()
	for fieldName, tag := range reflection.GetFieldTagsForValue(entity, "inherit") {
		key, referenceFieldName := parseInheritTag(tag)

		referenceField, referredField, ok := reflection.GetReferredField(entity, referenceFieldName, fieldName)
		if !ok {
			continue
		}

		field := entity.FieldByName(fieldName)

		inheritFields = append(inheritFields, &InheritField{entity, referenceField.Interface(), key, field, referredField})
	}

	return inheritFields
}

// Cache these, because we call twice for each entity
func (self InheritFields) Get(entityPtr tosca.EntityPtr) []*InheritField {
	inheritFields, ok := self[entityPtr]
	if !ok {
		inheritFields = NewInheritFields(entityPtr)
		self[entityPtr] = inheritFields
	}
	return inheritFields
}

func parseInheritTag(tag string) (string, string) {
	t := strings.Split(tag, ",")
	if len(t) != 2 {
		panic("must be 2")
	}

	key := t[0]
	referenceFieldName := t[1]

	return key, referenceFieldName
}
