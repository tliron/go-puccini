package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// PolicyTriggerType
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-policy-triggers/]
//

type PolicyTriggerType struct {
	*Entity `name:"policy trigger type"`
	Name    string `namespace:""`

	Source     *string             `read:"source" mandatory:""`
	Parameters PropertyDefinitions `read:"parameters,PropertyDefinition"`
}

func NewPolicyTriggerType(context *parsing.Context) *PolicyTriggerType {
	return &PolicyTriggerType{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Parameters: make(PropertyDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadPolicyTriggerType(context *parsing.Context) parsing.EntityPtr {
	self := NewPolicyTriggerType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

var policyTriggerTypeRoot *PolicyTriggerType

// ([parsing.Hierarchical] interface)
func (self *PolicyTriggerType) GetParent() parsing.EntityPtr {
	return policyTriggerTypeRoot
}

//
// PolicyTriggerTypes
//

type PolicyTriggerTypes []*PolicyTriggerType
