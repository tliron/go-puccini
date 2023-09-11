package tosca_v2_0

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// NodeTemplate
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.3
//

type NodeTemplate struct {
	*Entity `name:"node template"`
	Name    string `namespace:""`

	Directives                   *[]string              `read:"directives"`
	CopyNodeTemplateName         *string                `read:"copy"`
	NodeTypeName                 *string                `read:"type" mandatory:""`
	Metadata                     Metadata               `read:"metadata,Metadata"` // introduced in TOSCA 1.1
	Description                  *string                `read:"description"`
	Properties                   Values                 `read:"properties,Value"`
	Attributes                   Values                 `read:"attributes,AttributeValue"`
	Capabilities                 CapabilityAssignments  `read:"capabilities,CapabilityAssignment"`
	Requirements                 RequirementAssignments `read:"requirements,{}RequirementAssignment"`
	RequirementTargetsNodeFilter *NodeFilter            `read:"node_filter,NodeFilter"`
	Interfaces                   InterfaceAssignments   `read:"interfaces,InterfaceAssignment"`
	Artifacts                    Artifacts              `read:"artifacts,Artifact"`

	CopyNodeTemplate *NodeTemplate `lookup:"copy,CopyNodeTemplateName" traverse:"ignore" json:"-" yaml:"-"`
	NodeType         *NodeType     `lookup:"type,NodeTypeName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewNodeTemplate(context *parsing.Context) *NodeTemplate {
	return &NodeTemplate{
		Entity:       NewEntity(context),
		Name:         context.Name,
		Properties:   make(Values),
		Attributes:   make(Values),
		Capabilities: make(CapabilityAssignments),
		Interfaces:   make(InterfaceAssignments),
		Artifacts:    make(Artifacts),
	}
}

// ([parsing.Reader] signature)
func ReadNodeTemplate(context *parsing.Context) parsing.EntityPtr {
	self := NewNodeTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	switch self.Name {
	case "SELF", "SOURCE", "TARGET":
		context.Clone(self.Name).ReportValueInvalid("node template name", "reserved")
	}
	return self
}

// ([parsing.PreReadable] interface)
func (self *NodeTemplate) PreRead() {
	CopyTemplate(self.Context)
}

// ([parsing.Renderable] interface)
func (self *NodeTemplate) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *NodeTemplate) render() {
	logRender.Debugf("node template: %s", self.Name)

	if self.NodeType == nil {
		return
	}

	self.Properties.RenderProperties(self.NodeType.PropertyDefinitions, self.Context.FieldChild("properties", nil))
	self.Attributes.RenderAttributes(self.NodeType.AttributeDefinitions, self.Context.FieldChild("attributes", nil))
	self.Capabilities.Render(self.NodeType.CapabilityDefinitions, self.Context.FieldChild("capabilities", nil))
	self.Requirements.Render(self, self.Context.FieldChild("requirements", nil))
	self.Interfaces.RenderForNodeType(self.NodeType, self.Context.FieldChild("interfaces", nil))
	self.Artifacts.Render(self.NodeType.ArtifactDefinitions, self.Context.FieldChild("artifacts", nil))
}

func (self *NodeTemplate) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.NodeTemplate {
	logNormalize.Debugf("node template: %s", self.Name)

	normalNodeTemplate := normalServiceTemplate.NewNodeTemplate(self.Name)

	normalNodeTemplate.Metadata = self.Metadata

	if self.Description != nil {
		normalNodeTemplate.Description = *self.Description
	}

	if types, ok := normal.GetEntityTypes(self.Context.Hierarchy, self.NodeType); ok {
		normalNodeTemplate.Types = types
	}

	if self.Directives != nil {
		normalNodeTemplate.Directives = *self.Directives
	}

	self.Properties.Normalize(normalNodeTemplate.Properties)
	self.Attributes.Normalize(normalNodeTemplate.Attributes)
	self.Capabilities.Normalize(self, normalNodeTemplate)
	self.Interfaces.NormalizeForNodeTemplate(self, normalNodeTemplate)
	self.Artifacts.Normalize(normalNodeTemplate)

	return normalNodeTemplate
}

//
// NodeTemplates
//

type NodeTemplates []*NodeTemplate

func (self NodeTemplates) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, nodeTemplate := range self {
		normalServiceTemplate.NodeTemplates[nodeTemplate.Name] = nodeTemplate.Normalize(normalServiceTemplate)
	}

	// Requirements must be normalized after node templates
	// (because they may reference other node templates)
	for _, nodeTemplate := range self {
		if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[nodeTemplate.Name]; ok {
			nodeTemplate.Requirements.Normalize(nodeTemplate, normalNodeTemplate)
		}
	}
}
