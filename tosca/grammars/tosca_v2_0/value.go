package tosca_v2_0

import (
	"fmt"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/yamlkeys"
)

//
// Value
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.11, 3.6.13
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.10, 3.6.12
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.9, 3.5.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.9, 3.5.11
//

type Value struct {
	*Entity `name:"value"`
	Name    string

	ConstraintClauses ConstraintClauses
	Description       *string

	DataType    *DataType                `traverse:"ignore" json:"-" yaml:"-"`
	Information *normal.ValueInformation `traverse:"ignore" json:"-" yaml:"-"`
	Converter   *tosca.FunctionCall      `traverse:"ignore" json:"-" yaml:"-"`
}

func NewValue(context *tosca.Context) *Value {
	return &Value{
		Entity:      NewEntity(context),
		Name:        context.Name,
		Information: normal.NewValueInformation(),
	}
}

// tosca.Reader signature
func ReadValue(context *tosca.Context) tosca.EntityPtr {
	ToFunctionCall(context)
	return NewValue(context)
}

// tosca.Reader signature
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.12.2.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.11.2.2
func ReadAttributeValue(context *tosca.Context) tosca.EntityPtr {
	self := NewValue(context)

	// Unpack long notation
	if context.Is(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if len(map_) == 2 {
			if description, ok := map_["description"]; ok {
				if value, ok := map_["value"]; ok {
					self.Description = context.FieldChild("description", description).ReadString()
					context.Data = value
				}
			}
		}
	}

	ToFunctionCall(context)

	return self
}

// tosca.Mappable interface
func (self *Value) GetKey() string {
	return self.Name
}

// fmt.Stringer interface
func (self *Value) String() string {
	return yamlkeys.KeyString(self.Context.Data)
}

func (self *Value) RenderDataType(dataTypeName string) {
	if e, ok := self.Context.Namespace.Lookup(dataTypeName); ok {
		if dataType, ok := e.(*DataType); ok {
			lock := dataType.GetEntityLock()
			lock.RLock()
			defer lock.RUnlock()

			self.RenderAttribute(dataType, nil, false, false)
		} else {
			self.Context.ReportUnknownDataType(dataTypeName)
		}
	} else {
		self.Context.ReportUnknownDataType(dataTypeName)
	}
}

// Avoid rendering more than once (can happen if we were copied from PropertyDefinition.Default)
func (self *Value) RenderAttribute(dataType *DataType, definition *AttributeDefinition, bare bool, allowNil bool) {
	self.DataType = dataType

	/*if definition != nil {
		definition.Render()
	}*/

	if !bare {
		if self.Description != nil {
			self.Information.Description = *self.Description
		}

		if definition != nil {
			self.Information.Definition = definition.GetTypeInformation()
		}

		if dataType != nil {
			self.Information.Type = dataType.GetTypeInformation()
		}
	}

	dataType.Complete(self.Context)
	if !bare {
		self.ConstraintClauses.Render(dataType)
		dataType.ConstraintClauses.Render(dataType)
		self.ConstraintClauses = dataType.ConstraintClauses.Append(self.ConstraintClauses)
	}

	if _, ok := self.Context.Data.(*tosca.FunctionCall); ok {
		return
	}

	if allowNil && (self.Context.Data == nil) {
		return
	}

	// Internal types
	if internalTypeName, typeValidator, reader, ok := dataType.GetInternal(); ok {
		if typeValidator != nil {
			if self.Context.Data == nil {
				// Nil data only happens when an attribute is added despite not having a
				// "default" value; we will give it a valid zero value instead
				if self.Context.Data, ok = ScalarUnitTypeZeroes[internalTypeName]; !ok {
					if self.Context.Data, ok = ard.TypeZeroes[internalTypeName]; !ok {
						panic(fmt.Sprintf("unsupported internal type name: %s", internalTypeName))
					}
				}
			}

			if (internalTypeName == ard.TypeString) && self.Context.HasQuirk(tosca.QuirkDataTypesStringPermissive) {
				self.Context.Data = ard.ValueToString(self.Context.Data)
			}

			// Primitive types
			if typeValidator(self.Context.Data) {
				// Render list and map elements according to entry schema
				// (The entry schema may also have additional constraints)
				switch internalTypeName {
				case ard.TypeList, ard.TypeMap:
					if definition == nil {
						return
					} else if definition.EntrySchema == nil {
						return
					} else {
						lock := definition.EntrySchema.GetEntityLock()
						lock.RLock()
						defer lock.RUnlock()
						if definition.EntrySchema.DataType == nil {
							// This problem is reported in AttributeDefinition.Render
							return
						}
					}

					if internalTypeName == ard.TypeList {
						// Information
						entryDataType := definition.EntrySchema.DataType
						if entryDataType != dataType {
							lock := entryDataType.GetEntityLock()
							lock.RLock()
							defer lock.RUnlock()
						}
						entryConstraints := definition.EntrySchema.GetConstraints()
						self.Information.Entry = entryDataType.GetTypeInformation()
						if definition.EntrySchema.Description != nil {
							self.Information.Entry.SchemaDescription = *definition.EntrySchema.Description
						}

						slice := self.Context.Data.(ard.List)
						valueList := NewValueList(definition, len(slice), entryConstraints)

						for index, data := range slice {
							value := ReadAndRenderBareAttribute(self.Context.ListChild(index, data), entryDataType)
							valueList.Set(index, value)
						}

						self.Context.Data = valueList
					} else { // ard.TypeMap
						if definition.KeySchema == nil {
							// This problem is reported in AttributeDefinition.Complete
							return
						}

						lock := definition.KeySchema.GetEntityLock()
						lock.RLock()
						defer lock.RUnlock()

						// Information

						keyDataType := definition.KeySchema.DataType
						if keyDataType != dataType {
							lock1 := keyDataType.GetEntityLock()
							lock1.RLock()
							defer lock1.RUnlock()
						}
						keyConstraints := definition.KeySchema.GetConstraints()
						self.Information.Key = keyDataType.GetTypeInformation()
						if definition.KeySchema.Description != nil {
							self.Information.Key.SchemaDescription = *definition.KeySchema.Description
						}

						valueDataType := definition.EntrySchema.DataType
						if (valueDataType != dataType) && (valueDataType != keyDataType) {
							lock2 := valueDataType.GetEntityLock()
							lock2.RLock()
							defer lock2.RUnlock()
						}
						valueConstraints := definition.EntrySchema.GetConstraints()
						self.Information.Value = valueDataType.GetTypeInformation()
						if definition.EntrySchema.Description != nil {
							self.Information.Value.SchemaDescription = *definition.EntrySchema.Description
						}

						valueMap := NewValueMap(definition, keyConstraints, valueConstraints)

						for key, data := range self.Context.Data.(ard.Map) {
							// Complex keys are stringified for the purpose of the contexts

							// Validate key schema
							keyContext := self.Context.MapChild(key, yamlkeys.KeyData(key))
							key = ReadAndRenderBareAttribute(keyContext, keyDataType)

							context := self.Context.MapChild(key, data)
							value := ReadAndRenderBareAttribute(context, valueDataType)
							value.ConstraintClauses = ConstraintClauses{}
							valueMap.Put(key, value)
						}

						self.Context.Data = valueMap
					}
				}
			} else {
				self.Context.ReportValueWrongType(internalTypeName)
			}
		} else {
			// Special types
			self.Context.Data = reader(self.Context)
		}
	} else if self.Context.ValidateType(ard.TypeMap) {
		// Complex data types

		map_ := self.Context.Data.(ard.Map)

		// All properties must be defined in type
		for key := range map_ {
			name := yamlkeys.KeyString(key)
			if _, ok := dataType.PropertyDefinitions[name]; !ok {
				self.Context.MapChild(name, nil).ReportUndeclared("property")
				delete(map_, key)
			}
		}

		// Render properties
		for key, definition := range dataType.PropertyDefinitions {
			definition.Render()
			if data, ok := map_[key]; ok {
				var value *Value
				if value, ok = data.(*Value); !ok {
					// Convert to value
					value = ReadValue(self.Context.MapChild(key, data)).(*Value)
					map_[key] = value
				}
				if definition.DataType != nil {
					var lock util.RWLocker
					if definition.DataType != dataType {
						lock = definition.DataType.GetEntityLock()
						lock.RLock()
					}
					value.RenderProperty(definition.DataType, definition)
					if lock != nil {
						lock.RUnlock()
					}
				}

				// Grab information
				if (value.Information != nil) && !value.Information.Empty() {
					self.Information.Fields[key] = value.Information
				}
			} else if definition.IsRequired() {
				self.Context.MapChild(key, data).ReportPropertyRequired("property")
			}
		}
	}

	if (definition != nil) && (definition.Metadata != nil) {
		if converter, ok := definition.Metadata["puccini.converter"]; ok {
			self.Converter = self.Context.NewFunctionCall(converter, nil)
		}
	}
	if self.Converter == nil {
		if converter, ok := dataType.GetMetadataValue("puccini.converter"); ok {
			self.Converter = self.Context.NewFunctionCall(converter, nil)
		}
	}

	if comparer, ok := dataType.GetMetadataValue("puccini.comparer"); ok {
		if hasComparer, ok := self.Context.Data.(HasComparer); ok {
			hasComparer.SetComparer(comparer)
		} else {
			panic(fmt.Sprintf("type has \"puccini.comparer\" metadata but does not support HasComparer interface: %T", self.Context.Data))
		}
	}
}

func (self *Value) RenderProperty(dataType *DataType, definition *PropertyDefinition) {
	if definition == nil {
		self.RenderAttribute(dataType, nil, false, false)
	} else {
		self.ConstraintClauses.Render(dataType)
		if definition.DataType != dataType {
			lock := definition.DataType.GetEntityLock()
			lock.RLock()
			defer lock.RUnlock()
		}
		definition.ConstraintClauses.Render(definition.DataType)
		self.ConstraintClauses = definition.ConstraintClauses.Append(self.ConstraintClauses)
		self.RenderAttribute(dataType, definition.AttributeDefinition, false, false)
		//definition.ConstraintClauses.Prepend(&self.ConstraintClauses, dataType)
	}
}

func ReadAndRenderBareAttribute(context *tosca.Context, dataType *DataType) *Value {
	self := ReadValue(context).(*Value)
	self.RenderAttribute(dataType, nil, true, false)
	return self
}

func (self *Value) Normalize() normal.Constrainable {
	return self.normalize(true)
}

func (self *Value) normalize(withInformation bool) normal.Constrainable {
	var normalConstrainable normal.Constrainable

	switch data := self.Context.Data.(type) {
	case ard.Map:
		// This is for complex types (the "map" type is a ValueMap, below)
		normalMap := normal.NewMap()
		for key, value := range data {
			if v, ok := value.(*Value); ok {
				normalMap.Put(key, v.normalize(true))
			} else {
				normalMap.Put(key, normal.NewValue(value))
			}
		}
		normalConstrainable = normalMap

	case *ValueList:
		normalConstrainable = data.Normalize(self.Context)

	case *ValueMap:
		normalConstrainable = data.Normalize(self.Context)

	case *tosca.FunctionCall:
		NormalizeFunctionCallArguments(data, self.Context)
		normalConstrainable = normal.NewFunctionCall(data)

	default:
		value := normal.NewValue(data)
		normalConstrainable = value
	}

	if withInformation {
		normalConstrainable.SetInformation(self.Information)
	}

	self.ConstraintClauses.NormalizeConstrainable(self.Context, normalConstrainable)

	if self.Converter != nil {
		normalConstrainable.SetConverter(self.Converter)
	}

	return normalConstrainable
}

//
// Values
//

type Values map[string]*Value

func (self Values) CopyUnassigned(values Values) {
	for key, value := range values {
		if _, ok := self[key]; !ok {
			self[key] = value
		}
	}
}

func (self Values) RenderMissingValue(definition *AttributeDefinition, kind string, required bool, context *tosca.Context) {
	if definition.Default != nil {
		// Note: it doesn't make sense if required=false for properties in this case,
		// but we will just ignore it (future versions of TOSCA may specifically disallow it)
		self[definition.Name] = definition.Default
	} else if required {
		// Attributes are always required=false
		context.MapChild(definition.Name, nil).ReportPropertyRequired(kind)
	} else if kind == "attribute" {
		// Attributes should always appear, even if they have no default value
		self[definition.Name] = NewValue(context.MapChild(definition.Name, nil))
	}
}

func (self Values) RenderProperties(definitions PropertyDefinitions, kind string, context *tosca.Context) {
	for key, definition := range definitions {
		definition.Render()
		lock1 := definition.GetEntityLock()
		lock1.RLock()
		if value, ok := self[key]; !ok {
			self.RenderMissingValue(definition.AttributeDefinition, kind, definition.IsRequired(), context)
			// (If the above assigns the "default" value -- it has already been rendered elsewhere)
		} else if definition.DataType != nil {
			lock2 := value.GetEntityLock()
			lock2.Lock()
			lock3 := definition.DataType.GetEntityLock()
			lock3.RLock()
			value.RenderProperty(definition.DataType, definition)
			lock3.RUnlock()
			lock2.Unlock()
		}
		lock1.RUnlock()
	}

	for key, value := range self {
		if _, ok := definitions[key]; !ok {
			value.Context.ReportUndeclared(kind)
			delete(self, key)
		}
	}
}

func (self Values) RenderAttributes(definitions AttributeDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		if _, ok := self[key]; !ok {
			lock := definition.GetEntityLock()
			lock.RLock()
			self.RenderMissingValue(definition, "attribute", false, context)
			lock.RUnlock()
		}
	}

	for key, value := range self {
		if definition, ok := definitions[key]; !ok {
			value.Context.ReportUndeclared("attribute")
			delete(self, key)
		} else {
			definition.Render()
			lock1 := definition.GetEntityLock()
			lock1.RLock()
			if definition.DataType != nil {
				lock2 := value.GetEntityLock()
				lock2.Lock()
				lock3 := definition.DataType.GetEntityLock()
				lock3.RLock()
				value.RenderAttribute(definition.DataType, definition, false, true)
				lock3.RUnlock()
				lock2.Unlock()
			}
			lock1.RUnlock()
		}
	}
}

func (self Values) Normalize(normalConstrainables normal.Constrainables) {
	for key, value := range self {
		lock := value.GetEntityLock()
		lock.RLock()
		normalConstrainables[key] = value.Normalize()
		lock.RUnlock()
	}
}

//
// ValueList
//

type ValueList struct {
	EntryConstraints ConstraintClauses
	Slice            []interface{}
}

func NewValueList(definition *AttributeDefinition, length int, entryConstraints ConstraintClauses) *ValueList {
	return &ValueList{
		EntryConstraints: entryConstraints,
		Slice:            make([]interface{}, length),
	}
}

func (self *ValueList) Set(index int, value interface{}) {
	self.Slice[index] = value
}

func (self *ValueList) Normalize(context *tosca.Context) *normal.List {
	normalList := normal.NewList(len(self.Slice))

	self.EntryConstraints.NormalizeListEntries(context, normalList)

	for index, value := range self.Slice {
		if v, ok := value.(*Value); ok {
			lock := v.GetEntityLock()
			lock.RLock()
			normalList.Set(index, v.normalize(false))
			lock.RUnlock()
		} else {
			normalList.Set(index, normal.NewValue(value))
		}
	}

	return normalList
}

//
// ValueMap
//

type ValueMap struct {
	KeyConstraints   ConstraintClauses
	ValueConstraints ConstraintClauses
	Map              ard.Map
}

func NewValueMap(definition *AttributeDefinition, keyConstraints ConstraintClauses, valueConstraints ConstraintClauses) *ValueMap {
	return &ValueMap{
		KeyConstraints:   keyConstraints,
		ValueConstraints: valueConstraints,
		Map:              make(ard.Map),
	}
}

func (self *ValueMap) Put(key interface{}, value interface{}) {
	self.Map[key] = value
}

func (self *ValueMap) Normalize(context *tosca.Context) *normal.Map {
	normalMap := normal.NewMap()

	self.KeyConstraints.NormalizeMapKeys(context, normalMap)
	self.ValueConstraints.NormalizeMapValues(context, normalMap)

	for key, value := range self.Map {
		if k, ok := key.(*Value); ok {
			key = k.normalize(false)
		}
		if v, ok := value.(*Value); ok {
			lock := v.GetEntityLock()
			lock.RLock()
			normalMap.Put(key, v.normalize(false))
			lock.RUnlock()
		} else {
			normalMap.Put(key, normal.NewValue(value))
		}
	}

	return normalMap
}
