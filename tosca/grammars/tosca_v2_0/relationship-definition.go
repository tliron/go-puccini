package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// RelationshipDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.3.2.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.3.2.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.3.2.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.2.2.3
//

type RelationshipDefinition struct {
	*Entity `name:"relationship definition"`

	RelationshipTypeName *string              `read:"type"` // mandatory only if cannot be inherited
	InterfaceDefinitions InterfaceDefinitions `read:"interfaces,InterfaceDefinition" inherit:"interfaces,RelationshipType"`

	RelationshipType *RelationshipType `lookup:"type,RelationshipTypeName" traverse:"ignore" json:"-" yaml:"-"`

	typeMissingProblemReported bool
}

func NewRelationshipDefinition(context *parsing.Context) *RelationshipDefinition {
	return &RelationshipDefinition{
		Entity:               NewEntity(context),
		InterfaceDefinitions: make(InterfaceDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadRelationshipDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewRelationshipDefinition(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.RelationshipTypeName = context.FieldChild("type", context.Data).ReadString()
	}

	return self
}

func (self *RelationshipDefinition) NewDefaultAssignment(context *parsing.Context) *RelationshipAssignment {
	assignment := NewRelationshipAssignment(context)
	assignment.RelationshipTemplateNameOrTypeName = self.RelationshipTypeName
	assignment.RelationshipType = self.RelationshipType
	return assignment
}

func (self *RelationshipDefinition) Inherit(parentDefinition *RelationshipDefinition) {
	logInherit.Debug("relationship definition")

	if (self.RelationshipTypeName == nil) && (parentDefinition.RelationshipTypeName != nil) {
		self.RelationshipTypeName = parentDefinition.RelationshipTypeName
	}
	if (self.RelationshipType == nil) && (parentDefinition.RelationshipType != nil) {
		self.RelationshipType = parentDefinition.RelationshipType
	}

	// Validate type compatibility
	if (self.RelationshipType != nil) && (parentDefinition.RelationshipType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.RelationshipType, self.RelationshipType) {
		self.Context.ReportIncompatibleType(self.RelationshipType, parentDefinition.RelationshipType)
		return
	}

	self.InterfaceDefinitions.Inherit(parentDefinition.InterfaceDefinitions)
}

// ([parsing.Renderable] interface)
func (self *RelationshipDefinition) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *RelationshipDefinition) render() {
	logRender.Debug("relationship definition")

	if self.RelationshipTypeName == nil {
		// Avoid reporting more than once
		if !self.typeMissingProblemReported {
			self.Context.FieldChild("type", nil).ReportKeynameMissing()
			self.typeMissingProblemReported = true
		}
	}
}
