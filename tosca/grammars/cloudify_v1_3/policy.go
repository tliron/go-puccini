package cloudify_v1_3

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Policy
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-policies/]
//

type Policy struct {
	*Entity `name:"policy"`
	Name    string `namespace:""`

	PolicyTypeName   *string   `read:"type" mandatory:""`
	Properties       Values    `read:"properties,Value"`
	TargetGroupNames *[]string `read:"targets" mandatory:""`

	PolicyType   *PolicyType `lookup:"type,PolicyTypeName" traverse:"ignore" json:"-" yaml:"-"`
	TargetGroups Groups      `lookup:"targets,TargetGroupNames" traverse:"ignore" json:"-" yaml:"-"`
}

func NewPolicy(context *parsing.Context) *Policy {
	return &Policy{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Properties: make(Values),
	}
}

// ([parsing.Reader] signature)
func ReadPolicy(context *parsing.Context) parsing.EntityPtr {
	self := NewPolicy(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *Policy) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.Policy {
	logNormalize.Debugf("policy: %s", self.Name)

	normalPolicy := normalServiceTemplate.NewPolicy(self.Name)

	if types, ok := normal.GetEntityTypes(self.Context.Hierarchy, self.PolicyType); ok {
		normalPolicy.Types = types
	}

	self.Properties.Normalize(normalPolicy.Properties, "")

	for _, group := range self.TargetGroups {
		if normalGroup, ok := normalServiceTemplate.Groups[group.Name]; ok {
			normalPolicy.GroupTargets = append(normalPolicy.GroupTargets, normalGroup)
		}
	}

	return normalPolicy
}

//
// Policies
//

type Policies []*Policy

func (self Policies) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, policy := range self {
		normalServiceTemplate.Policies[policy.Name] = policy.Normalize(normalServiceTemplate)
	}
}
