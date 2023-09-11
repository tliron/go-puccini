package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// GroupPolicyTriggerTrigger
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-groups/]
//

type GroupPolicyTrigger struct {
	*Entity `name:"group policy trigger"`

	PolicyTriggerTypeName *string `read:"type" mandatory:""`
	Parameters            Values  `read:"parameters,Value"`

	PolicyTriggerType *PolicyTriggerType `lookup:"type,PolicyTriggerTypeName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewGroupPolicyTrigger(context *parsing.Context) *GroupPolicyTrigger {
	return &GroupPolicyTrigger{
		Entity:     NewEntity(context),
		Parameters: make(Values),
	}
}

// ([parsing.Reader] signature)
func ReadGroupPolicyTrigger(context *parsing.Context) parsing.EntityPtr {
	self := NewGroupPolicyTrigger(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

//
// GroupPolicyTriggers
//

type GroupPolicyTriggers map[string]*GroupPolicyTrigger
