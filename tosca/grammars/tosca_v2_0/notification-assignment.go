package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// NotificationAssignment
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.19
//

type NotificationAssignment struct {
	*Entity `name:"notification"`
	Name    string

	Description    *string                  `read:"description"`
	Implementation *InterfaceImplementation `read:"implementation,InterfaceImplementation"`
	Outputs        OutputMappings           `read:"outputs,OutputMapping"`
}

func NewNotificationAssignment(context *parsing.Context) *NotificationAssignment {
	return &NotificationAssignment{
		Entity:  NewEntity(context),
		Name:    context.Name,
		Outputs: make(OutputMappings),
	}
}

// ([parsing.Reader] signature)
func ReadNotificationAssignment(context *parsing.Context) parsing.EntityPtr {
	self := NewNotificationAssignment(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.Implementation = ReadInterfaceImplementation(context.FieldChild("implementation", context.Data)).(*InterfaceImplementation)
	}

	return self
}

// ([parsing.Mappable] interface)
func (self *NotificationAssignment) GetKey() string {
	return self.Name
}

func (self *NotificationAssignment) Normalize(normalInterface *normal.Interface) *normal.Notification {
	logNormalize.Debugf("notification: %s", self.Name)

	normalNotification := normalInterface.NewNotification(self.Name)

	if self.Description != nil {
		normalNotification.Description = *self.Description
	}

	if self.Implementation != nil {
		self.Implementation.NormalizeNotification(normalNotification)
	}

	self.Outputs.Normalize(normalNotification.Outputs)

	return normalNotification
}

//
// NotificationAssignments
//

type NotificationAssignments map[string]*NotificationAssignment

func (self NotificationAssignments) CopyUnassigned(assignments NotificationAssignments) {
	for key, assignment := range assignments {
		if selfAssignment, ok := self[key]; ok {
			selfAssignment.Outputs.CopyUnassigned(assignment.Outputs)
		} else {
			self[key] = assignment
		}
	}
}

func (self NotificationAssignments) RenderForNodeType(nodeType *NodeType, definitions NotificationDefinitions, context *parsing.Context) {
	self.render(definitions, context)
	for _, assignment := range self {
		assignment.Outputs.RenderForNodeType(nodeType)
	}
}

func (self NotificationAssignments) RenderForRelationshipType(relationshipType *RelationshipType, definitions NotificationDefinitions, sourceNodeTemplate *NodeTemplate, context *parsing.Context) {
	self.render(definitions, context)
	for _, assignment := range self {
		assignment.Outputs.RenderForRelationshipType(relationshipType, sourceNodeTemplate)
	}
}

func (self NotificationAssignments) RenderForGroup(definitions NotificationDefinitions, context *parsing.Context) {
	self.render(definitions, context)
	for _, assignment := range self {
		assignment.Outputs.RenderForGroup()
	}
}

func (self NotificationAssignments) render(definitions NotificationDefinitions, context *parsing.Context) {
	for key, definition := range definitions {
		assignment, ok := self[key]

		if !ok {
			assignment = NewNotificationAssignment(context.FieldChild(key, nil))
			self[key] = assignment
		}

		if assignment.Description == nil {
			assignment.Description = definition.Description
		}

		if (assignment.Implementation == nil) && (definition.Implementation != nil) {
			// If the definition has an implementation then we must have one, too
			assignment.Implementation = NewInterfaceImplementation(assignment.Context.FieldChild("implementation", nil))
		}

		if assignment.Implementation != nil {
			assignment.Implementation.Render(definition.Implementation)
		}

		assignment.Outputs.Inherit(definition.Outputs)
	}

	for key, assignment := range self {
		if _, ok := definitions[key]; !ok {
			assignment.Context.ReportUndeclared("notification")
			delete(self, key)
		}
	}
}

func (self NotificationAssignments) Normalize(normalInterface *normal.Interface) {
	for key, notification := range self {
		normalInterface.Notifications[key] = notification.Normalize(normalInterface)
	}
}
