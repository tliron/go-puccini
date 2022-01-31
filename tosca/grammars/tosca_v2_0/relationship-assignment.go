package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// RelationshipAssignment
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.2.2.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.2.2.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.2.2.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.2.2.3
//

type RelationshipAssignment struct {
	*Entity `name:"relationship"`

	RelationshipTemplateNameOrTypeName *string              `read:"type"`
	Properties                         Values               `read:"properties,Value"`
	Attributes                         Values               `read:"attributes,AttributeValue"` // missing in spec
	Interfaces                         InterfaceAssignments `read:"interfaces,InterfaceAssignment"`

	RelationshipTemplate *RelationshipTemplate `lookup:"type,RelationshipTemplateNameOrTypeName" json:"-" yaml:"-"`
	RelationshipType     *RelationshipType     `lookup:"type,RelationshipTemplateNameOrTypeName" json:"-" yaml:"-"`
}

func NewRelationshipAssignment(context *tosca.Context) *RelationshipAssignment {
	return &RelationshipAssignment{
		Entity:     NewEntity(context),
		Properties: make(Values),
		Attributes: make(Values),
		Interfaces: make(InterfaceAssignments),
	}
}

// tosca.Reader signature
func ReadRelationshipAssignment(context *tosca.Context) tosca.EntityPtr {
	self := NewRelationshipAssignment(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.RelationshipTemplateNameOrTypeName = context.FieldChild("type", context.Data).ReadString()
	}

	return self
}

func (self *RelationshipAssignment) GetType(definition *RelationshipDefinition) *RelationshipType {
	if self.RelationshipTemplate != nil {
		return self.RelationshipTemplate.RelationshipType
	} else if self.RelationshipType != nil {
		return self.RelationshipType
	} else if definition != nil {
		return definition.RelationshipType
	} else {
		return nil
	}
}

func (self *RelationshipAssignment) Render(definition *RelationshipDefinition) {
	relationshipType := self.GetType(definition)
	if relationshipType == nil {
		self.Context.ReportUndeclared("relationship")
		return
	}

	if definition != nil {
		// We will consider the "interfaces" at the definition to take priority over those at the type
		self.Interfaces.RenderForRelationship(self, definition.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))

		// Validate type compatibility
		if (definition.RelationshipType != nil) && !self.Context.Hierarchy.IsCompatible(definition.RelationshipType, relationshipType) {
			self.Context.ReportIncompatibleType(relationshipType, definition.RelationshipType)
			return
		}
	}

	if self.RelationshipTemplate != nil {
		lock1 := self.RelationshipTemplate.GetEntityLock()
		lock1.Lock()
		defer lock1.Unlock()

		self.RelationshipTemplate.Render()
		self.Properties.CopyUnassigned(self.RelationshipTemplate.Properties)
		self.Attributes.CopyUnassigned(self.RelationshipTemplate.Attributes)
		self.Interfaces.CopyUnassigned(self.RelationshipTemplate.Interfaces)

		lock2 := self.RelationshipTemplate.RelationshipType.GetEntityLock()
		lock2.RLock()
		defer lock2.RUnlock()

		self.Interfaces.RenderForRelationship(self, self.RelationshipTemplate.RelationshipType.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))
	} else {
		lock := self.RelationshipType.GetEntityLock()
		lock.RLock()
		defer lock.RUnlock()

		self.Properties.RenderProperties(self.RelationshipType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
		self.Attributes.RenderAttributes(self.RelationshipType.AttributeDefinitions, self.Context.FieldChild("attributes", nil))
		self.Interfaces.RenderForRelationship(self, self.RelationshipType.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))
	}
}

func (self *RelationshipAssignment) Normalize(definition *RelationshipDefinition, normalRelationship *normal.Relationship) {
	if self.RelationshipTemplate != nil {
		lock := self.RelationshipTemplate.GetEntityLock()
		lock.RLock()
		self.RelationshipTemplate.Normalize(normalRelationship)
		lock.RUnlock()
	}

	relationshipType := self.GetType(definition)
	if types, ok := normal.GetTypes(self.Context.Hierarchy, relationshipType); ok {
		normalRelationship.Types = types
	}

	self.Properties.Normalize(normalRelationship.Properties)
	self.Attributes.Normalize(normalRelationship.Attributes)
	self.Interfaces.NormalizeForRelationship(self, definition, normalRelationship)
}
