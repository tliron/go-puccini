package tosca_v2_0

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// RelationshipTemplate
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.4
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.4
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.4
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.4
//

type RelationshipTemplate struct {
	*Entity `name:"relationship template"`
	Name    string `namespace:""`

	CopyRelationshipTemplateName *string              `read:"copy"`
	RelationshipTypeName         *string              `read:"type" mandatory:""`
	Metadata                     Metadata             `read:"metadata,Metadata"` // introduced in TOSCA 1.1
	Description                  *string              `read:"description"`
	Properties                   Values               `read:"properties,Value"`
	Attributes                   Values               `read:"attributes,AttributeValue"`
	Interfaces                   InterfaceAssignments `read:"interfaces,InterfaceAssignment"`

	CopyRelationshipTemplate *RelationshipTemplate `lookup:"copy,CopyRelationshipTemplateName" traverse:"ignore" json:"-" yaml:"-"`
	RelationshipType         *RelationshipType     `lookup:"type,RelationshipTypeName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewRelationshipTemplate(context *parsing.Context) *RelationshipTemplate {
	return &RelationshipTemplate{
		Entity:     NewEntity(context),
		Properties: make(Values),
		Attributes: make(Values),
		Name:       context.Name,
		Interfaces: make(InterfaceAssignments),
	}
}

// ([parsing.Reader] signature)
func ReadRelationshipTemplate(context *parsing.Context) parsing.EntityPtr {
	self := NewRelationshipTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.PreReadable] interface)
func (self *RelationshipTemplate) PreRead() {
	CopyTemplate(self.Context)
}

// ([parsing.Renderable] interface)
func (self *RelationshipTemplate) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *RelationshipTemplate) render() {
	logRender.Debugf("relationship template: %s", self.Name)

	if self.RelationshipType == nil {
		return
	}

	self.Properties.RenderProperties(self.RelationshipType.PropertyDefinitions, self.Context.FieldChild("properties", nil))
	self.Attributes.RenderAttributes(self.RelationshipType.AttributeDefinitions, self.Context.FieldChild("attributes", nil))
}

func (self *RelationshipTemplate) Normalize(normalRelationship *normal.Relationship) {
	normalRelationship.Metadata = self.Metadata

	if self.Description != nil {
		normalRelationship.Description = *self.Description
	}
}

//
// RelationshipTemplates
//

type RelationshipTemplates []*RelationshipTemplate
