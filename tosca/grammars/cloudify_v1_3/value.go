package cloudify_v1_3

import (
	"fmt"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

//
// Value
//

type Value struct {
	*Entity `name:"value"`
	Name    string

	Meta *normal.ValueMeta `traverse:"ignore" json:"-" yaml:"-"`
}

func NewValue(context *parsing.Context) *Value {
	return &Value{
		Entity: NewEntity(context),
		Name:   context.Name,
		Meta:   normal.NewValueMeta(),
	}
}

// ([parsing.Reader] signature)
func ReadValue(context *parsing.Context) parsing.EntityPtr {
	ParseFunctionCalls(context)
	return NewValue(context)
}

// ([parsing.Mappable] interface)
func (self *Value) GetKey() string {
	return self.Name
}

func (self *Value) RenderProperty(dataType *DataType, definition *PropertyDefinition) {
	if definition == nil {
		self.RenderParameter(dataType, nil, true, false)
	} else {
		self.RenderParameter(dataType, definition.ParameterDefinition, true, false)
	}
}

func (self *Value) RenderParameter(dataType *DataType, definition *ParameterDefinition, validateRequire bool, allowNil bool) {
	if dataType != nil {
		self.Meta = dataType.NewValueMeta()
	}

	if _, ok := self.Context.Data.(*parsing.FunctionCall); ok {
		return
	}

	if allowNil && (self.Context.Data == nil) {
		return
	}

	// TODO: dataType.Complete(self.Context.Data, self.Context)

	// Internal types
	if internalTypeName, typeValidator, reader, ok := dataType.GetInternal(); ok {
		if typeValidator != nil {
			if self.Context.Data == nil {
				// Nil data only happens when a parameter is added despite not having a
				// "default" value; we will give it a valid zero value instead
				if self.Context.Data, ok = ard.TypeZeroes[internalTypeName]; !ok {
					panic(fmt.Sprintf("unsupported internal type name: %s", internalTypeName))
				}
			}

			// Primitive types
			if !typeValidator(self.Context.Data) {
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
			if data, ok := map_[key]; ok {
				var value *Value
				if value, ok = data.(*Value); !ok {
					// Convert to value
					value = NewValue(self.Context.MapChild(key, data))
					map_[key] = value
				}
				if definition.DataType != nil {
					value.RenderProperty(definition.DataType, definition)
				}
			} else if validateRequire {
				if definition.IsRequired() {
					self.Context.MapChild(key, data).ReportValueRequired("property")
				}
			}
		}
	}
}

func (self *Value) Normalize() normal.Value {
	return self.normalize(true)
}

func (self *Value) normalize(withMeta bool) normal.Value {
	var normalValue normal.Value

	switch data := self.Context.Data.(type) {
	case ard.List:
		normalList := normal.NewList(len(data))
		for index, value := range data {
			normalList.Set(index, NewValue(self.Context.ListChild(index, value)).normalize(false))
		}
		normalValue = normalList

	case ard.Map:
		normalMap := normal.NewMap()
		for key, value := range data {
			if _, ok := key.(string); !ok {
				// Cloudify DSL does not support complex keys
				self.Context.MapChild(key, yamlkeys.KeyData(key)).ReportValueWrongType(ard.TypeString)
			}
			name := yamlkeys.KeyString(key)
			normalMap.Put(name, NewValue(self.Context.MapChild(name, value)).normalize(false))
		}
		normalValue = normalMap

	case *parsing.FunctionCall:
		NormalizeFunctionCallArguments(data, self.Context)
		normalValue = normal.NewFunctionCall(data)

	default:
		normalValue = normal.NewPrimitive(data)
	}

	if withMeta {
		normalValue.SetMeta(self.Meta)
	}

	return normalValue
}

//
// Values
//

type Values map[string]*Value

func (self Values) SetIfNil(context *parsing.Context, key string, data ard.Value) {
	if _, ok := self[key]; !ok {
		self[key] = NewValue(context.MapChild(key, data))
	}
}

func (self Values) RenderMissingValue(definition *ParameterDefinition, kind string, required bool, context *parsing.Context) {
	if definition.Default != nil {
		self[definition.Name] = definition.Default
	} else if required {
		context.MapChild(definition.Name, nil).ReportValueRequired(kind)
	} else if kind == "attribute" {
		// Parameters should always appear, even if they have no default value
		self[definition.Name] = NewValue(context.MapChild(definition.Name, nil))
	}
}

func (self Values) RenderProperties(definitions PropertyDefinitions, kind string, context *parsing.Context) {
	for key, definition := range definitions {
		if value, ok := self[key]; !ok {
			self.RenderMissingValue(definition.ParameterDefinition, kind, definition.IsRequired(), context)
			// (If the above assigns the "default" value -- it has already been rendered elsewhere)
		} else if definition.DataType != nil {
			value.RenderProperty(definition.DataType, definition)
		}
	}
}

func (self Values) RenderParameters(definitions ParameterDefinitions, kind string, context *parsing.Context) {
	for key, definition := range definitions {
		if _, ok := self[key]; !ok {
			self.RenderMissingValue(definition, kind, false, context)
		}
	}

	for key, value := range self {
		if definition, ok := definitions[key]; !ok {
			value.Context.ReportUndeclared(kind)
			delete(self, key)
		} else if definition.DataType != nil {
			value.RenderParameter(definition.DataType, definition, true, true)
		}
	}
}

func (self Values) Normalize(normalConstrainables normal.Values, prefix string) {
	for key, value := range self {
		normalConstrainables[prefix+key] = value.Normalize()
	}
}
