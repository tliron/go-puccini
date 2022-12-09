package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
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

func NewOperationAssignment(context *tosca.Context) *OperationAssignment {
	return &OperationAssignment{
		Entity:  NewEntity(context),
		Name:    context.Name,
		Inputs:  make(Values),
		Outputs: make(OutputMappings),
	}
}

// tosca.Reader signature
func ReadOperationAssignment(context *tosca.Context) tosca.EntityPtr {
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

// tosca.Mappable interface
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
	if normalInterface.NodeTemplate != nil {
		self.Outputs.NormalizeForNodeTemplate(normalInterface.NodeTemplate.ServiceTemplate, normalOperation.Outputs)
	} else if normalInterface.Relationship != nil {
		self.Outputs.NormalizeForRelationship(normalInterface.Relationship, normalOperation.Outputs)
	} else if normalInterface.Group != nil {
		self.Outputs.NormalizeForGroup(normalInterface.Group.ServiceTemplate, normalOperation.Outputs)
	}

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

func (self OperationAssignments) Render(definitions OperationDefinitions, context *tosca.Context) {
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

func (self OperationAssignments) RenderForNodeTemplate(nodeTemplate *NodeTemplate, definitions OperationDefinitions, context *tosca.Context) {
	for _, assignment := range self {
		assignment.Outputs.RenderForNodeTemplate(nodeTemplate)
	}
}
