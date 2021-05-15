package tosca_v2_0

import (
	"reflect"

	"github.com/tliron/puccini/tosca"
)

//
// RequirementMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.8
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.8
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

type RequirementMapping struct {
	*Entity `name:"requirement mapping"`
	Name    string

	NodeTemplateName *string
	RequirementName  *string

	NodeTemplate *NodeTemplate          `traverse:"ignore" json:"-" yaml:"-"`
	Requirement  *RequirementAssignment `traverse:"ignore" json:"-" yaml:"-"`
}

func NewRequirementMapping(context *tosca.Context) *RequirementMapping {
	return &RequirementMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadRequirementMapping(context *tosca.Context) tosca.EntityPtr {
	self := NewRequirementMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.RequirementName = &(*strings)[1]
	}

	return self
}

// tosca.Mappable interface
func (self *RequirementMapping) GetKey() string {
	return self.Name
}

func (self *RequirementMapping) GetRequirementDefinition() (*RequirementDefinition, bool) {
	if (self.Requirement != nil) && (self.NodeTemplate != nil) {
		return self.Requirement.GetDefinition(self.NodeTemplate)
	} else {
		return nil, false
	}
}

// parser.Renderable interface
func (self *RequirementMapping) Render() {
	logRender.Debug("requirement mapping")

	if (self.NodeTemplateName == nil) || (self.RequirementName == nil) {
		return
	}

	nodeTemplateName := *self.NodeTemplateName
	var nodeTemplateType *NodeTemplate
	if nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, reflect.TypeOf(nodeTemplateType)); ok {
		self.NodeTemplate = nodeTemplate.(*NodeTemplate)
	} else {
		self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
		return
	}

	name := *self.RequirementName
	found := false
	self.NodeTemplate.Render()
	for _, requirement := range self.NodeTemplate.Requirements {
		if requirement.Name == name {
			if found {
				self.Context.ListChild(1, name).ReportReferenceAmbiguous("requirement", self.NodeTemplate)
				break
			} else {
				self.Requirement = requirement
				found = true
			}
		}
	}

	if !found {
		if ! self.Context.HasQuirk(tosca.QuirkSubstitutionMappingsRequirementsAllowDangling) {
			self.Context.ListChild(1, name).ReportReferenceNotFound("requirement", self.NodeTemplate)
		}
	}
}

//
// RequirementMappings
//

type RequirementMappings map[string]*RequirementMapping
