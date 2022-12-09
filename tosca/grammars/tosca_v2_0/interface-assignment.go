package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// InterfaceAssignment
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.20
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.16
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.14
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.14
//

type InterfaceAssignment struct {
	*Entity `name:"interface" json:"-" yaml:"-"`
	Name    string

	Inputs        Values                  `read:"inputs,Value"`
	Operations    OperationAssignments    `read:"operations,OperationAssignment"`       // keyword since TOSCA 1.3
	Notifications NotificationAssignments `read:"notifications,NotificationAssignment"` // introduced in TOSCA 1.3
}

func NewInterfaceAssignment(context *tosca.Context) *InterfaceAssignment {
	return &InterfaceAssignment{
		Entity:        NewEntity(context),
		Name:          context.Name,
		Inputs:        make(Values),
		Operations:    make(OperationAssignments),
		Notifications: make(NotificationAssignments),
	}
}

// tosca.Reader signature
func ReadInterfaceAssignment(context *tosca.Context) tosca.EntityPtr {
	self := NewInterfaceAssignment(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *InterfaceAssignment) GetKey() string {
	return self.Name
}

func (self *InterfaceAssignment) GetDefinitionForNodeTemplate(nodeTemplate *NodeTemplate) (*InterfaceDefinition, bool) {
	if nodeTemplate.NodeType == nil {
		return nil, false
	}
	definition, ok := nodeTemplate.NodeType.InterfaceDefinitions[self.Name]
	return definition, ok
}

func (self *InterfaceAssignment) GetDefinitionForGroup(group *Group) (*InterfaceDefinition, bool) {
	if group.GroupType == nil {
		return nil, false
	}
	definition, ok := group.GroupType.InterfaceDefinitions[self.Name]
	return definition, ok
}

func (self *InterfaceAssignment) GetDefinitionForRelationship(relationship *RelationshipAssignment, relationshipDefinition *RelationshipDefinition) (*InterfaceDefinition, bool) {
	relationshipType := relationship.GetType(relationshipDefinition)
	if relationshipType == nil {
		return nil, false
	}
	definition, ok := relationshipType.InterfaceDefinitions[self.Name]
	return definition, ok
}

func (self *InterfaceAssignment) RenderForNodeTemplate(nodeTemplate *NodeTemplate, definition *InterfaceDefinition) {
	self.Inputs.RenderInputs(definition.InputDefinitions, self.Context.FieldChild("inputs", nil))
	self.Operations.Render(definition.OperationDefinitions, self.Context.FieldChild("operations", nil))
	self.Operations.RenderForNodeTemplate(nodeTemplate, definition.OperationDefinitions, self.Context.FieldChild("operations", nil))
	self.Notifications.RenderForNodeTemplate(nodeTemplate, definition.NotificationDefinitions, self.Context.FieldChild("notifications", nil))
}

func (self *InterfaceAssignment) RenderForRelationship(relationship *RelationshipAssignment, definition *InterfaceDefinition) {
	self.Inputs.RenderInputs(definition.InputDefinitions, self.Context.FieldChild("inputs", nil))
	self.Operations.Render(definition.OperationDefinitions, self.Context.FieldChild("operations", nil))
	self.Notifications.RenderForRelationship(relationship, definition.NotificationDefinitions, self.Context.FieldChild("notifications", nil))
}

func (self *InterfaceAssignment) RenderForGroup(definition *InterfaceDefinition) {
	self.Inputs.RenderInputs(definition.InputDefinitions, self.Context.FieldChild("inputs", nil))
	self.Operations.Render(definition.OperationDefinitions, self.Context.FieldChild("operations", nil))
	self.Notifications.RenderForGroup(definition.NotificationDefinitions, self.Context.FieldChild("notifications", nil))
}

func (self *InterfaceAssignment) Normalize(normalInterface *normal.Interface, definition *InterfaceDefinition) {
	logNormalize.Debugf("interface: %s", self.Name)

	if (definition.InterfaceType != nil) && (definition.InterfaceType.Description != nil) {
		normalInterface.Description = *definition.InterfaceType.Description
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, definition.InterfaceType); ok {
		normalInterface.Types = types
	}

	self.Inputs.Normalize(normalInterface.Inputs)
	self.Operations.Normalize(normalInterface)
	self.Notifications.Normalize(normalInterface)
}

//
// InterfaceAssignments
//

type InterfaceAssignments map[string]*InterfaceAssignment

func (self InterfaceAssignments) CopyUnassigned(assignments InterfaceAssignments) {
	for key, assignment := range assignments {
		if selfAssignment, ok := self[key]; ok {
			selfAssignment.Inputs.CopyUnassigned(assignment.Inputs)
			selfAssignment.Operations.CopyUnassigned(assignment.Operations)
			selfAssignment.Notifications.CopyUnassigned(assignment.Notifications)
		} else {
			self[key] = assignment
		}
	}
}

func (self InterfaceAssignments) RenderForNodeTemplate(nodeTemplate *NodeTemplate, definitions InterfaceDefinitions, context *tosca.Context) {
	self.render(definitions, context)
	for name, assignment := range self {
		if definition, ok := definitions[name]; ok {
			assignment.RenderForNodeTemplate(nodeTemplate, definition)
		}
	}
}

func (self InterfaceAssignments) RenderForRelationship(relationship *RelationshipAssignment, definitions InterfaceDefinitions, context *tosca.Context) {
	self.render(definitions, context)
	for name, assignment := range self {
		if definition, ok := definitions[name]; ok {
			assignment.RenderForRelationship(relationship, definition)
		}
	}
}

func (self InterfaceAssignments) RenderForGroup(definitions InterfaceDefinitions, context *tosca.Context) {
	self.render(definitions, context)
	for name, assignment := range self {
		if definition, ok := definitions[name]; ok {
			assignment.RenderForGroup(definition)
		}
	}
}

func (self InterfaceAssignments) render(definitions InterfaceDefinitions, context *tosca.Context) {
	for key := range definitions {
		assignment, ok := self[key]
		if !ok {
			assignment = NewInterfaceAssignment(context.MapChild(key, nil))
			self[key] = assignment
		}
	}

	for key, assignment := range self {
		if _, ok := definitions[key]; !ok {
			assignment.Context.ReportUndeclared("interface")
			delete(self, key)
		}
	}
}

func (self InterfaceAssignments) NormalizeForNodeTemplate(nodeTemplate *NodeTemplate, normalNodeTemplate *normal.NodeTemplate) {
	for key, interface_ := range self {
		if definition, ok := interface_.GetDefinitionForNodeTemplate(nodeTemplate); ok {
			interface_.Normalize(normalNodeTemplate.NewInterface(key), definition)
		}
	}
}

func (self InterfaceAssignments) NormalizeForGroup(group *Group, normalGroup *normal.Group) {
	for key, interface_ := range self {
		if definition, ok := interface_.GetDefinitionForGroup(group); ok {
			interface_.Normalize(normalGroup.NewInterface(key), definition)
		}
	}
}

func (self InterfaceAssignments) NormalizeForRelationship(relationship *RelationshipAssignment, relationshipDefinition *RelationshipDefinition, normalRelationship *normal.Relationship) {
	for key, interface_ := range self {
		if definition, ok := interface_.GetDefinitionForRelationship(relationship, relationshipDefinition); ok {
			interface_.Normalize(normalRelationship.NewInterface(key), definition)
		}
	}
}
