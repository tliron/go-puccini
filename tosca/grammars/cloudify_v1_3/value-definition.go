package cloudify_v1_3

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// ValueDefinition
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-capabilities/]
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-outputs/]
//

type ValueDefinition struct {
	*Entity `name:"value definition"`
	Name    string `namespace:""`

	Description *string `read:"description"`
	Value       *Value  `read:"value,Value"`
}

func NewValueDefinition(context *parsing.Context) *ValueDefinition {
	return &ValueDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadValueDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewValueDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Mappable] interface)
func (self *ValueDefinition) GetKey() string {
	return self.Name
}

//
// ValueDefinitions
//

type ValueDefinitions map[string]*ValueDefinition

func (self ValueDefinitions) Normalize(c normal.Values) {
	for key, valueDefinition := range self {
		if valueDefinition.Value != nil {
			c[key] = valueDefinition.Value.Normalize()
		}
	}
}
