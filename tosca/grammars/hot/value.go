package hot

import (
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

	Constraints Constraints

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

func (self *Value) Normalize() normal.Value {
	var normalValue normal.Value

	switch data := self.Context.Data.(type) {
	case ard.List:
		normalList := normal.NewList(len(data))
		for index, value := range data {
			normalList.Set(index, NewValue(self.Context.ListChild(index, value)).Normalize())
		}
		normalValue = normalList

	case ard.Map:
		normalMap := normal.NewMap()
		for key, value := range data {
			if _, ok := key.(string); !ok {
				// HOT does not support complex keys
				self.Context.MapChild(key, yamlkeys.KeyData(key)).ReportValueWrongType(ard.TypeString)
			}
			name := yamlkeys.KeyString(key)
			normalMap.Put(name, NewValue(self.Context.MapChild(name, value)).Normalize())
		}
		normalValue = normalMap

	case *parsing.FunctionCall:
		NormalizeFunctionCallArguments(data, self.Context)
		normalValue = normal.NewFunctionCall(data)

	default:
		normalValue = normal.NewPrimitive(data)
	}

	self.Constraints.Normalize(self.Context, self.Meta)

	normalValue.SetMeta(self.Meta)

	return normalValue
}

//
// Values
//

type Values map[string]*Value

func (self Values) Normalize(normalConstrainables normal.Values) {
	for key, value := range self {
		normalConstrainables[key] = value.Normalize()
	}
}
