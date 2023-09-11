package tosca_v2_0

import (
	"reflect"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// InterfaceMapping
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.12
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.11
//

type InterfaceMapping struct {
	*Entity `name:"interface mapping"`
	Name    string

	NodeTemplateName *string
	InterfaceName    *string

	NodeTemplate *NodeTemplate        `traverse:"ignore" json:"-" yaml:"-"`
	Interface    *InterfaceAssignment `traverse:"ignore" json:"-" yaml:"-"`
}

func NewInterfaceMapping(context *parsing.Context) *InterfaceMapping {
	return &InterfaceMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadInterfaceMapping(context *parsing.Context) parsing.EntityPtr {
	self := NewInterfaceMapping(context)
	if context.ValidateType(ard.TypeList) {
		strings := context.ReadStringListFixed(2)
		if strings != nil {
			self.NodeTemplateName = &(*strings)[0]
			self.InterfaceName = &(*strings)[1]
		}
	}
	return self
}

// ([parsing.Mappable] interface)
func (self *InterfaceMapping) GetKey() string {
	return self.Name
}

func (self *InterfaceMapping) GetInterfaceDefinition() (*InterfaceDefinition, bool) {
	if (self.Interface != nil) && (self.NodeTemplate != nil) {
		return self.Interface.GetDefinitionForNodeTemplate(self.NodeTemplate)
	} else {
		return nil, false
	}
}

// ([parsing.Renderable] interface)
func (self *InterfaceMapping) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *InterfaceMapping) render() {
	logRender.Debug("interface mapping")

	if (self.NodeTemplateName == nil) || (self.InterfaceName == nil) {
		return
	}

	nodeTemplateName := *self.NodeTemplateName
	var nodeTemplateType *NodeTemplate
	if nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, reflect.TypeOf(nodeTemplateType)); ok {
		self.NodeTemplate = nodeTemplate.(*NodeTemplate)

		self.NodeTemplate.Render()

		name := *self.InterfaceName
		var ok bool
		if self.Interface, ok = self.NodeTemplate.Interfaces[name]; !ok {
			self.Context.ListChild(1, name).ReportReferenceNotFound("interface", self.NodeTemplate)
		}
	} else {
		self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
	}
}

//
// InterfaceMappings
//

type InterfaceMappings map[string]*InterfaceMapping
