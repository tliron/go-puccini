package hot

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

var DeletionPolicies = map[string]string{
	"Delete":   "Delete",
	"Retain":   "Retain",
	"Snapshot": "Snapshot",
	"delete":   "Delete",
	"retain":   "Retain",
	"snapshot": "Snapshot",
}

func GetDeletionPolicy(policy string) (string, bool) {
	if p, ok := DeletionPolicies[policy]; ok {
		return p, true
	}
	return "", false
}

//
// Resource
//
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#resources-section]
//

type Resource struct {
	*Entity `name:"resource"`
	Name    string `namespace:""`

	Type           *string    `read:"type" mandatory:""`
	Properties     Values     `read:"properties,Value"`
	Metadata       *Data      `read:"metadata,Data"`
	DependsOn      *[]string  `read:"depends_on"`
	UpdatePolicy   *Data      `read:"update_policy,Data"`
	DeletionPolicy *string    `read:"deletion_policy"`
	ExternalID     *string    `read:"external_id"`
	Condition      *Condition `read:"condition,Condition"`

	ToscaType          *string   `json:"-" yaml:"-"`
	DependsOnResources Resources `lookup:"depends_on,DependsOn" traverse:"ignore" json:"-" yaml:"-"`
}

func NewResource(context *parsing.Context) *Resource {
	return &Resource{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Properties: make(Values),
	}
}

// ([parsing.Reader] signature)
func ReadResource(context *parsing.Context) parsing.EntityPtr {
	self := NewResource(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))

	if childContext, ok := context.GetFieldChild("depends_on"); ok {
		self.DependsOn = childContext.ReadStringOrStringList()
	}

	if self.Type != nil {
		if type_, ok := ResourceTypes[*self.Type]; ok {
			self.ToscaType = &type_
		} else {
			context.FieldChild("type", *self.Type).ReportKeynameUnsupportedValue()
		}
	}

	if self.DeletionPolicy != nil {
		if policy, ok := GetDeletionPolicy(*self.DeletionPolicy); ok {
			self.DeletionPolicy = &policy
		} else {
			context.FieldChild("deletion_policy", *self.DeletionPolicy).ReportKeynameUnsupportedValue()
		}
	}

	return self
}

var capabilityTypeName = "Resource"
var capabilityTypes = normal.NewEntityTypes(capabilityTypeName)
var relationshipTypes = normal.NewEntityTypes("DependsOn")

func (self *Resource) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.NodeTemplate {
	logNormalize.Debugf("resource: %s", self.Name)

	normalNodeTemplate := normalServiceTemplate.NewNodeTemplate(self.Name)

	if self.ToscaType != nil {
		normalNodeTemplate.Types = normal.NewEntityTypes(*self.ToscaType)
	}

	self.Properties.Normalize(normalNodeTemplate.Properties)

	capabilityContext := self.Context.FieldChild("capabilities", nil).MapChild("resource", nil)
	normalNodeTemplate.NewCapability("resource", normal.NewLocationForContext(capabilityContext)).Types = capabilityTypes

	return normalNodeTemplate
}

func (self *Resource) NormalizeDependencies(normalServiceTemplate *normal.ServiceTemplate) {
	logNormalize.Debugf("resource dependencies: %s", self.Name)

	normalNodeTemplate := normalServiceTemplate.NodeTemplates[self.Name]
	requirementsContext := self.Context.FieldChild("requirements", nil)

	for index, resource := range self.DependsOnResources {
		normalRequirement := normalNodeTemplate.NewRequirement("depends_on", normal.NewLocationForContext(requirementsContext.ListChild(index, nil)))
		normalRequirement.NodeTemplate = normalServiceTemplate.NodeTemplates[resource.Name]
		normalRequirement.CapabilityTypeName = &capabilityTypeName

		normalRelationship := normalRequirement.NewRelationship()
		normalRelationship.Types = relationshipTypes
	}
}

//
// Resources
//

type Resources []*Resource

func (self Resources) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, resource := range self {
		normalServiceTemplate.NodeTemplates[resource.Name] = resource.Normalize(normalServiceTemplate)
	}

	// Dependencies must be normalized after resources
	// (because they may reference other resources)
	for _, resource := range self {
		resource.NormalizeDependencies(normalServiceTemplate)
	}
}
