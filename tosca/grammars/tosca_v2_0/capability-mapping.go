package tosca_v2_0

import (
	"reflect"

	"github.com/tliron/puccini/tosca/parsing"
)

//
// CapabilityMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.9
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

type CapabilityMapping struct {
	*Entity `name:"capability mapping"`
	Name    string

	NodeTemplateName *string
	CapabilityName   *string

	NodeTemplate *NodeTemplate         `traverse:"ignore" json:"-" yaml:"-"`
	Capability   *CapabilityAssignment `traverse:"ignore" json:"-" yaml:"-"`
}

func NewCapabilityMapping(context *parsing.Context) *CapabilityMapping {
	return &CapabilityMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadCapabilityMapping(context *parsing.Context) parsing.EntityPtr {
	self := NewCapabilityMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.CapabilityName = &(*strings)[1]
	}

	return self
}

// ([parsing.Mappable] interface)
func (self *CapabilityMapping) GetKey() string {
	return self.Name
}

func (self *CapabilityMapping) GetCapabilityDefinition() (*CapabilityDefinition, bool) {
	if (self.Capability != nil) && (self.NodeTemplate != nil) {
		return self.Capability.GetDefinition(self.NodeTemplate)
	} else {
		return nil, false
	}
}

// ([parsing.Renderable] interface)
func (self *CapabilityMapping) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *CapabilityMapping) render() {
	logRender.Debug("capability mapping")

	if (self.NodeTemplateName == nil) || (self.CapabilityName == nil) {
		return
	}

	nodeTemplateName := *self.NodeTemplateName
	var nodeTemplateType *NodeTemplate
	if nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, reflect.TypeOf(nodeTemplateType)); ok {
		self.NodeTemplate = nodeTemplate.(*NodeTemplate)

		self.NodeTemplate.Render()

		name := *self.CapabilityName
		var ok bool
		if self.Capability, ok = self.NodeTemplate.Capabilities[name]; !ok {
			self.Context.ListChild(1, name).ReportReferenceNotFound("capability", self.NodeTemplate)
		}
	} else {
		self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
	}
}

//
// CapabilityMappings
//

type CapabilityMappings map[string]*CapabilityMapping
