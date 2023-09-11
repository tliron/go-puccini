package tosca_v2_0

import (
	"reflect"

	"github.com/tliron/puccini/tosca/parsing"
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

func NewRequirementMapping(context *parsing.Context) *RequirementMapping {
	return &RequirementMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadRequirementMapping(context *parsing.Context) parsing.EntityPtr {
	self := NewRequirementMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.RequirementName = &(*strings)[1]
	}

	return self
}

// ([parsing.Mappable] interface)
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

// ([parsing.Renderable] interface)
func (self *RequirementMapping) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *RequirementMapping) render() {
	logRender.Debug("requirement mapping")

	if (self.NodeTemplateName == nil) || (self.RequirementName == nil) {
		return
	}

	nodeTemplateName := *self.NodeTemplateName
	var nodeTemplateType *NodeTemplate
	if nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, reflect.TypeOf(nodeTemplateType)); ok {
		self.NodeTemplate = nodeTemplate.(*NodeTemplate)

		if self.Context.HasQuirk(parsing.QuirkSubstitutionMappingsRequirementsPermissive) {
			return
		}

		self.NodeTemplate.Render()

		name := *self.RequirementName
		found := false
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
			self.Context.ListChild(1, name).ReportReferenceNotFound("requirement", self.NodeTemplate)
		}
	} else {
		self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
		return
	}
}

//
// RequirementMappings
//

type RequirementMappings map[string]*RequirementMapping
