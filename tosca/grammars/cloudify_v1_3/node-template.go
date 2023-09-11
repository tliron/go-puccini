package cloudify_v1_3

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// NodeTemplate
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-node-templates/]
//

type NodeTemplate struct {
	*Entity `name:"node template"`
	Name    string `namespace:""`

	NodeTypeName  *string                  `read:"type" mandatory:""`
	Properties    Values                   `read:"properties,Value"`
	Instances     *NodeTemplateInstances   `read:"instances,NodeTemplateInstances"` // deprecated in Cloudify DSL 1.3
	Interfaces    InterfaceAssignments     `read:"interfaces,InterfaceAssignment"`
	Relationships RelationshipAssignments  `read:"relationships,[]RelationshipAssignment"`
	Capabilities  NodeTemplateCapabilities `read:"capabilities,NodeTemplateCapability"`

	NodeType *NodeType `lookup:"type,NodeTypeName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewNodeTemplate(context *parsing.Context) *NodeTemplate {
	return &NodeTemplate{
		Entity:       NewEntity(context),
		Name:         context.Name,
		Properties:   make(Values),
		Interfaces:   make(InterfaceAssignments),
		Capabilities: make(NodeTemplateCapabilities),
	}
}

// ([parsing.Reader] signature)
func ReadNodeTemplate(context *parsing.Context) parsing.EntityPtr {
	self := NewNodeTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	self.Capabilities.Validate(context, self.Instances)
	return self
}

// ([parsing.Renderable] interface)
func (self *NodeTemplate) Render() {
	self.renderOnce.Do(self.render)
}

func (self *NodeTemplate) render() {
	logRender.Debugf("node template: %s", self.Name)

	if self.NodeType == nil {
		return
	}

	self.Properties.RenderProperties(self.NodeType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
	self.Interfaces.Render(self.NodeType.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))
}

var capabilityTypeName = "cloudify.Node"
var capabilityTypes = normal.NewEntityTypes(capabilityTypeName)

func (self *NodeTemplate) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.NodeTemplate {
	logNormalize.Debugf("node template: %s", self.Name)

	normalNodeTemplate := normalServiceTemplate.NewNodeTemplate(self.Name)

	if types, ok := normal.GetEntityTypes(self.Context.Hierarchy, self.NodeType); ok {
		normalNodeTemplate.Types = types
	}

	self.Properties.Normalize(normalNodeTemplate.Properties, "")
	self.Interfaces.NormalizeForNodeTemplate(self, normalNodeTemplate)

	capabilityContext := self.Context.FieldChild("capabilities", nil).MapChild("node", nil)
	normalCapability := normalNodeTemplate.NewCapability("node", normal.NewLocationForContext(capabilityContext))
	normalCapability.Types = capabilityTypes
	for _, capability := range self.Capabilities {
		capability.Properties.Normalize(normalCapability.Properties, capability.Context.Name+".")
	}

	return normalNodeTemplate
}

func (self *NodeTemplate) NormalizeRelationships(normalServiceTemplate *normal.ServiceTemplate) {
	logNormalize.Debugf("node template relationships: %s", self.Name)

	normalNodeTemplate := normalServiceTemplate.NodeTemplates[self.Name]
	for _, relationship := range self.Relationships {
		relationship.Normalize(self, normalNodeTemplate)
	}
}

//
// NodeTemplates
//

type NodeTemplates []*NodeTemplate

func (self NodeTemplates) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, nodeTemplate := range self {
		normalServiceTemplate.NodeTemplates[nodeTemplate.Name] = nodeTemplate.Normalize(normalServiceTemplate)
	}

	// Relationships must be normalized after node templates
	// (because they may reference other node templates)
	for _, nodeTemplate := range self {
		nodeTemplate.NormalizeRelationships(normalServiceTemplate)
	}
}
