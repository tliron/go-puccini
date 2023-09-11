package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// PolicyType
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-policy-types/]
//

type PolicyType struct {
	*Entity `name:"policy type"`
	Name    string `namespace:""`

	Source              *string             `read:"source" mandatory:""`
	PropertyDefinitions PropertyDefinitions `read:"properties,PropertyDefinition"`
}

func NewPolicyType(context *parsing.Context) *PolicyType {
	return &PolicyType{
		Entity:              NewEntity(context),
		Name:                context.Name,
		PropertyDefinitions: make(PropertyDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadPolicyType(context *parsing.Context) parsing.EntityPtr {
	self := NewPolicyType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

var policyTypeRoot *PolicyType

// ([parsing.Hierarchical] interface)
func (self *PolicyType) GetParent() parsing.EntityPtr {
	return policyTypeRoot
}

//
// PolicyTypes
//

type PolicyTypes []*PolicyType
