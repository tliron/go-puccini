package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Policy
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.6
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.6
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.6
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.6
//

type Policy struct {
	*Entity `name:"policy"`
	Name    string `namespace:""`

	PolicyTypeName                 *string            `read:"type" require:""`
	Metadata                       Metadata           `read:"metadata,Metadata"` // introduced in TOSCA 1.1
	Description                    *string            `read:"description"`
	Properties                     Values             `read:"properties,Value"`
	TargetNodeTemplateOrGroupNames *[]string          `read:"targets"`
	TriggerDefinitions             TriggerDefinitions `read:"triggers,TriggerDefinition" inherit:"triggers,PolicyType"` // introduced in TOSCA 1.1

	PolicyType          *PolicyType   `lookup:"type,PolicyTypeName" json:"-" yaml:"-"`
	TargetNodeTemplates NodeTemplates `lookup:"targets,TargetNodeTemplateOrGroupNames" json:"-" yaml:"-"`
	TargetGroups        Groups        `lookup:"targets,TargetNodeTemplateOrGroupNames" json:"-" yaml:"-"`
}

func NewPolicy(context *tosca.Context) *Policy {
	return &Policy{
		Entity:             NewEntity(context),
		Name:               context.Name,
		Properties:         make(Values),
		TriggerDefinitions: make(TriggerDefinitions),
	}
}

// tosca.Reader signature
func ReadPolicy(context *tosca.Context) tosca.EntityPtr {
	self := NewPolicy(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *Policy) GetKey() string {
	return self.Name
}

// parser.Renderable interface
func (self *Policy) Render() {
	self.renderOnce.Do(self.render)
}

func (self *Policy) render() {
	logRender.Debugf("policy: %s", self.Name)

	if self.PolicyType == nil {
		return
	}

	lock := self.PolicyType.GetEntityLock()
	lock.RLock()
	defer lock.RUnlock()

	self.Properties.RenderProperties(self.PolicyType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))

	// Validate targets

	if len(self.PolicyType.TargetNodeTypes) > 0 {
		for index, nodeTemplate := range self.TargetNodeTemplates {
			lock := nodeTemplate.GetEntityLock()
			lock.RLock()

			compatible := false
			for _, nodeType := range self.PolicyType.TargetNodeTypes {
				if self.Context.Hierarchy.IsCompatible(nodeType, nodeTemplate.NodeType) {
					compatible = true
					break
				}
			}

			if !compatible {
				childContext := self.Context.FieldChild("targets", nil).ListChild(index, nil)
				childContext.ReportIncompatible(nodeTemplate.Name, "policy", "target")
			}

			lock.RUnlock()
		}
	}

	if len(self.PolicyType.TargetGroupTypes) > 0 {
		for index, group := range self.TargetGroups {
			lock := group.GetEntityLock()
			lock.Lock()

			compatible := false
			for _, groupType := range self.PolicyType.TargetGroupTypes {
				if self.Context.Hierarchy.IsCompatible(groupType, group.GroupType) {
					compatible = true
					break
				}
			}

			if !compatible {
				childContext := self.Context.FieldChild("targets", nil).ListChild(index, nil)
				childContext.ReportIncompatible(group.Name, "policy", "target")
			}

			lock.Unlock()
		}
	}
}

func (self *Policy) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.Policy {
	logNormalize.Debugf("policy: %s", self.Name)

	normalPolicy := normalServiceTemplate.NewPolicy(self.Name)

	normalPolicy.Metadata = self.Metadata

	if self.Description != nil {
		normalPolicy.Description = *self.Description
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.PolicyType); ok {
		normalPolicy.Types = types
	}

	self.Properties.Normalize(normalPolicy.Properties)

	for _, nodeTemplate := range self.TargetNodeTemplates {
		lock := nodeTemplate.GetEntityLock()
		lock.RLock()
		if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[nodeTemplate.Name]; ok {
			normalPolicy.NodeTemplateTargets = append(normalPolicy.NodeTemplateTargets, normalNodeTemplate)
		}
		lock.RUnlock()
	}

	for _, group := range self.TargetGroups {
		lock := group.GetEntityLock()
		lock.RLock()
		if normalGroup, ok := normalServiceTemplate.Groups[group.Name]; ok {
			normalPolicy.GroupTargets = append(normalPolicy.GroupTargets, normalGroup)
		}
		lock.RUnlock()
	}

	self.TriggerDefinitions.Normalize(normalPolicy)

	return normalPolicy
}

//
// Policies
//

type Policies []*Policy

func (self Policies) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, policy := range self {
		lock := policy.GetEntityLock()
		lock.RLock()
		normalServiceTemplate.Policies[policy.Name] = policy.Normalize(normalServiceTemplate)
		lock.RUnlock()
	}
}
