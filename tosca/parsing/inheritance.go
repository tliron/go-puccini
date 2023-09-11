package parsing

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/tliron/kutil/reflection"
	"github.com/tliron/kutil/terminal"
)

//
// Inherits
//

type Inherits interface {
	Inherit()
}

// From [Inherits] interface
func Inherit(entityPtr EntityPtr) bool {
	if inherits, ok := entityPtr.(Inherits); ok {
		inherits.Inherit()
		return true
	} else {
		return false
	}
}

//
// Hierarchical
//

type Hierarchical interface {
	GetParent() EntityPtr
}

// From [Hierarchical] interface
func GetParent(entityPtr EntityPtr) (EntityPtr, bool) {
	if hierarchical, ok := entityPtr.(Hierarchical); ok {
		parentPtr := hierarchical.GetParent()
		if reflect.ValueOf(parentPtr).IsNil() {
			parentPtr = nil
		}
		return parentPtr, true
	} else {
		return nil, false
	}
}

//
// Hierarchy
//

type Hierarchy struct {
	entityPtr EntityPtr
	parent    *Hierarchy
	children  []*Hierarchy
}

// Keeps track of failed types
type HierarchyContext = EntityPtrSet

func NewHierarchy() *Hierarchy {
	return new(Hierarchy)
}

func NewHierarchyFor(entityPtr EntityPtr, work reflection.EntityWork, hierarchyContext HierarchyContext) *Hierarchy {
	self := NewHierarchy()

	work.TraverseEntities(entityPtr, func(entityPtr EntityPtr) bool {
		if parentPtr, ok := GetParent(entityPtr); ok {
			self.add(entityPtr, parentPtr, hierarchyContext, nil)
		}
		return true
	})

	return self
}

func (self *Hierarchy) Empty() bool {
	return len(self.children) == 0
}

func (self *Hierarchy) Root() *Hierarchy {
	for self != nil {
		if self.parent == nil {
			return self
		}
		self = self.parent
	}

	panic("bad hierarchy")
}

func (self *Hierarchy) GetContext() *Context {
	return GetContext(self.entityPtr)
}

func (self *Hierarchy) Range(f func(EntityPtr, EntityPtr) bool) {
	hierarchy_ := self
	for (hierarchy_ != nil) && (hierarchy_.entityPtr != nil) {
		parent := hierarchy_.parent

		if parent != nil {
			if !f(hierarchy_.entityPtr, parent.entityPtr) {
				return
			}
		} else if !f(hierarchy_.entityPtr, nil) {
			return
		}

		hierarchy_ = parent
	}
}

func (self *Hierarchy) Find(entityPtr EntityPtr) (*Hierarchy, bool) {
	if entityPtr == nil {
		return nil, false
	}

	if self.entityPtr == entityPtr {
		return self, true
	}

	for _, child := range self.children {
		if found, ok := child.Find(entityPtr); ok {
			return found, true
		}
	}

	return nil, false
}

/*
func (self *Hierarchy) Lookup(name string, type_ reflect.Type) (*Hierarchy, bool) {
	return self.lookup(name, type_)
}

func (self *Hierarchy) lookup(name string, type_ reflect.Type) (*Hierarchy, bool) {
	ourType := reflect.TypeOf(self.entityPtr) == type_

	if ourType {
		if self.getContext().Name == name {
			return self, true
		}
	}

	// Recurse only if in our type or in the root node
	if ourType || (self.parent == nil) {
		for _, child := range self.children {
			if found, ok := child.lookup(name, type_); ok {
				return found, true
			}
		}
	}

	return nil, false
}
*/

func (self *Hierarchy) IsCompatible(baseEntityPtr EntityPtr, entityPtr EntityPtr) bool {
	// Trivial case
	if baseEntityPtr == entityPtr {
		return true
	}

	if baseNode, ok := self.Find(baseEntityPtr); ok {
		_, ok = baseNode.Find(entityPtr)
		return ok
	}

	return false
}

func (self *Hierarchy) add(entityPtr EntityPtr, parentEntityPtr EntityPtr, hierarchyContext HierarchyContext, descendants EntityPtrs) (*Hierarchy, bool) {
	// Several imports may try to add the same entity to their hierarchies, so let's avoid multiple problem reports
	if hierarchyContext.Contains(entityPtr) {
		return nil, false
	}

	root := self.Root()

	// Have we already added this entity?
	if found, ok := root.Find(entityPtr); ok {
		return found, true
	}

	child := Hierarchy{entityPtr: entityPtr}

	if parentEntityPtr == nil {
		// We are a root node
		root.addChild(&child)
		return &child, true
	}

	context := GetContext(entityPtr)

	// Check for inheritance loop
	for _, descendant := range descendants {
		if descendant == entityPtr {
			context.ReportInheritanceLoop(parentEntityPtr)
			hierarchyContext.Add(entityPtr)
			return nil, false
		}
	}

	grandparentEntityPtr, ok := GetParent(parentEntityPtr)
	if !ok {
		panic(fmt.Sprintf("parent is somehow of the wrong type (it doesn't have a \"parent\" tag): %s", context.Path))
	}

	// Make sure parent node has been added first (recursively)
	parentNode, ok := self.add(parentEntityPtr, grandparentEntityPtr, hierarchyContext, append(descendants, entityPtr))
	if !ok {
		// Check if we already reported a failure (one report is enough)
		_, alreadyFailed := hierarchyContext[entityPtr]
		if !alreadyFailed {
			context.ReportTypeIncomplete(parentEntityPtr)
			hierarchyContext.Add(entityPtr)
		}
		return nil, false
	}

	// Add ourself to parent node
	parentNode.addChild(&child)

	return &child, true
}

func (self *Hierarchy) Merge(hierarchy *Hierarchy, hierarchyContext HierarchyContext) {
	if (hierarchy == nil) || (self == hierarchy) {
		return
	}

	self.merge(hierarchy, hierarchyContext)
}

func (self *Hierarchy) merge(hierarchy *Hierarchy, hierarchyContext HierarchyContext) {
	if hierarchy.entityPtr != nil {
		parentPtr, ok := GetParent(hierarchy.entityPtr)
		if ok {
			self.add(hierarchy.entityPtr, parentPtr, hierarchyContext, nil)
		}
	}

	for _, child := range hierarchy.children {
		self.merge(child, hierarchyContext)
	}
}

func (self *Hierarchy) addChild(hierarchy *Hierarchy) {
	hierarchy.parent = self
	self.children = append(self.children, hierarchy)
}

// TODO: Do we need this?

// Into "hierarchy" tags
func (self *Hierarchy) AddTo(entityPtr EntityPtr) {
	for _, field := range reflection.GetTaggedFields(entityPtr, "hierarchy") {
		type_ := field.Type()
		if reflection.IsSliceOfPointerToStruct(type_) {
			type_ = type_.Elem()
			self.appendTypeToSlice(field, type_)
		} else {
			panic(fmt.Sprintf("\"hierarchy\" tag is incompatible with []*struct{}: %v", type_))
		}
	}
}

func (self *Hierarchy) appendTypeToSlice(field reflect.Value, type_ reflect.Type) {
	if reflect.TypeOf(self.entityPtr) == type_ {
		// Don't add if it's already there
		found := false
		length := field.Len()
		for index := 0; index < length; index++ {
			element := field.Index(index)
			if element.Interface() == self.entityPtr {
				found = true
				break
			}
		}

		if !found {
			field.Set(reflect.Append(field, reflect.ValueOf(self.entityPtr)))
		}
	}

	for _, child := range self.children {
		child.appendTypeToSlice(field, type_)
	}
}

// Print

// Note that the same name could be printed out twice in the hierarchy, even under the same
// parent! That's because we are printing the local name of the type, and types imported from
// other files can have the same name (though you would need a namespace_prefix to avoid a
// namespace error)
func (self *Hierarchy) Print(indent int) {
	self.PrintChildren(indent, terminal.TreePrefix{})
}

func (self *Hierarchy) PrintChildren(indent int, treePrefix terminal.TreePrefix) {
	length := len(self.children)
	last := length - 1

	// Sort
	hierarchy := Hierarchy{children: append(self.children[:0:0], self.children...)}
	sort.Sort(hierarchy)

	for i, child := range hierarchy.children {
		isLast := i == last
		child.PrintChild(indent, treePrefix, isLast)
		child.PrintChildren(indent, append(treePrefix, isLast))
	}
}

func (self *Hierarchy) PrintChild(indent int, treePrefix terminal.TreePrefix, last bool) {
	treePrefix.Print(indent, last)
	if self.entityPtr != nil {
		terminal.Printf("%s\n", terminal.StdoutStylist.TypeName(self.GetContext().Name))
	}
}

// ([sort.Interface])
func (self Hierarchy) Len() int {
	return len(self.children)
}

// ([sort.Interface])
func (self Hierarchy) Swap(i, j int) {
	self.children[i], self.children[j] = self.children[j], self.children[i]
}

// ([sort.Interface])
func (self Hierarchy) Less(i, j int) bool {
	iName := self.children[i].GetContext().Name
	jName := self.children[j].GetContext().Name
	return strings.Compare(iName, jName) < 0
}
