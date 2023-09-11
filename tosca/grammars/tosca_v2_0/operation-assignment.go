package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// OperationAssignment
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.17
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.15
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.13
//

type OperationAssignment struct {
	*Entity `name:"operation"`
	Name    string

	Description    *string                  `read:"description"`
	Implementation *InterfaceImplementation `read:"implementation,InterfaceImplementation"`
	Inputs         Values                   `read:"inputs,Value"`
	Outputs        OutputMappings           `read:"outputs,OutputMapping"` // introduced in TOSCA 1.3
}

func NewOperationAssignment(context *parsing.Context) *OperationAssignment {
	return &OperationAssignment{
		Entity:  NewEntity(context),
		Name:    context.Name,
		Inputs:  make(Values),
		Outputs: make(OutputMappings),
	}
}

// ([parsing.Reader] signature)
func ReadOperationAssignment(context *parsing.Context) parsing.EntityPtr {
	self := NewOperationAssignment(context)

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
func (self *OperationAssignment) GetKey() string {
	return self.Name
}

func (self *OperationAssignment) Normalize(normalInterface *normal.Interface) *normal.Operation {
	logNormalize.Debugf("operation: %s", self.Name)

	normalOperation := normalInterface.NewOperation(self.Name)

	if self.Description != nil {
		normalOperation.Description = *self.Description
	}

	if self.Implementation != nil {
		self.Implementation.NormalizeOperation(normalOperation)
	}

	self.Inputs.Normalize(normalOperation.Inputs)
	self.Outputs.Normalize(normalOperation.Outputs)

	return normalOperation
}

//
// OperationAssignments
//

type OperationAssignments map[string]*OperationAssignment

func (self OperationAssignments) CopyUnassigned(assignments OperationAssignments) {
	for key, assignment := range assignments {
		if selfAssignment, ok := self[key]; ok {
			selfAssignment.Inputs.CopyUnassigned(assignment.Inputs)
			selfAssignment.Outputs.CopyUnassigned(assignment.Outputs)
			if selfAssignment.Description == nil {
				selfAssignment.Description = assignment.Description
			}
			if selfAssignment.Implementation == nil {
				selfAssignment.Implementation = assignment.Implementation
			}
		} else {
			self[key] = assignment
		}
	}
}

func (self OperationAssignments) RenderForNodeType(nodeType *NodeType, definitions OperationDefinitions, context *parsing.Context) {
	self.render(definitions, context)
	for _, assignment := range self {
		assignment.Outputs.RenderForNodeType(nodeType)
	}
}

func (self OperationAssignments) RenderForRelationshipType(relationshipType *RelationshipType, definitions OperationDefinitions, sourceNodeTemplate *NodeTemplate, context *parsing.Context) {
	self.render(definitions, context)
	for _, assignment := range self {
		assignment.Outputs.RenderForRelationshipType(relationshipType, sourceNodeTemplate)
	}
}

func (self OperationAssignments) RenderForGroup(definitions OperationDefinitions, context *parsing.Context) {
	self.render(definitions, context)
	for _, assignment := range self {
		assignment.Outputs.RenderForGroup()
	}
}

func (self OperationAssignments) render(definitions OperationDefinitions, context *parsing.Context) {
	for key, definition := range definitions {
		assignment, ok := self[key]

		if !ok {
			assignment = NewOperationAssignment(context.FieldChild(key, nil))
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

		assignment.Inputs.RenderInputs(definition.InputDefinitions, assignment.Context.FieldChild("inputs", nil))
		assignment.Outputs.Inherit(definition.Outputs)
	}

	for key, assignment := range self {
		_, ok := definitions[key]
		if !ok {
			assignment.Context.ReportUndeclared("operation")
			delete(self, key)
		}
	}
}

func (self OperationAssignments) Normalize(normalInterface *normal.Interface) {
	for key, operation := range self {
		normalInterface.Operations[key] = operation.Normalize(normalInterface)
	}
}
