package tosca_v2_0

import (
	"reflect"

	"github.com/tliron/puccini/tosca/parsing"
)

//
// AttributeMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.15
//

type AttributeMapping struct {
	*Entity `name:"attribute mapping"`
	Name    string

	NodeTemplateName *string
	AttributeName    *string

	NodeTemplate *NodeTemplate `traverse:"ignore" json:"-" yaml:"-"`
	Attribute    *Value        `traverse:"ignore" json:"-" yaml:"-"`
}

func NewAttributeMapping(context *parsing.Context) *AttributeMapping {
	return &AttributeMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadAttributeMapping(context *parsing.Context) parsing.EntityPtr {
	self := NewAttributeMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.AttributeName = &(*strings)[1]
	}

	return self
}

// ([parsing.Mappable] interface)
func (self *AttributeMapping) GetKey() string {
	return self.Name
}

func (self *AttributeMapping) EnsureRender() {
	logRender.Debug("attribute mapping")

	if (self.NodeTemplateName == nil) || (self.AttributeName == nil) {
		return
	}

	nodeTemplateName := *self.NodeTemplateName
	var nodeTemplateType *NodeTemplate
	if nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, reflect.TypeOf(nodeTemplateType)); ok {
		self.NodeTemplate = nodeTemplate.(*NodeTemplate)

		self.NodeTemplate.Render()

		name := *self.AttributeName
		var ok bool
		if self.Attribute, ok = self.NodeTemplate.Attributes[name]; !ok {
			self.Context.ListChild(1, name).ReportReferenceNotFound("attribute", self.NodeTemplate)
		}
	} else {
		self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
	}
}

//
// AttributeMappings
//

type AttributeMappings map[string]*AttributeMapping

func (self AttributeMappings) EnsureRender() {
	for _, mapping := range self {
		mapping.EnsureRender()
	}
}
