package hot

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Output
//
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#outputs-section]
//

type Output struct {
	*Entity `name:"output"`
	Name    string `namespace:""`

	Description *string    `read:"description"`
	Value       *Value     `read:"value,Value" mandatory:""`
	Condition   *Condition `read:"condition,Condition"`
}

func NewOutput(context *parsing.Context) *Output {
	return &Output{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadOutput(context *parsing.Context) parsing.EntityPtr {
	self := NewOutput(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Mappable] interface)
func (self *Output) GetKey() string {
	return self.Name
}

func (self *Output) Normalize(context *parsing.Context) normal.Value {
	var value *Value
	if self.Value != nil {
		value = self.Value
	} else {
		// Parameters should always appear, even if they have no default value
		value = NewValue(context.MapChild(self.Name, nil))
	}
	return value.Normalize()
}

//
// Outputs
//

type Outputs map[string]*Output

func (self Outputs) Normalize(normalConstrainables normal.Values, context *parsing.Context) {
	for key, output := range self {
		normalConstrainables[key] = output.Value.Normalize()
	}
}
