package tosca_v2_0

import (
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// WorkflowActivityCallOperation
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.23.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.19.2.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.17.2.3
//

type WorkflowActivityCallOperation struct {
	*Entity `name:"workflow activity call operation"`
	Name    string

	InterfaceAndOperation *string `read:"operation"`
	Inputs                Values  `read:"inputs,Value"` // introduced in TOSCA 1.3

	Interface *InterfaceAssignment `json:"-" yaml:"-"`
	Operation *OperationAssignment `json:"-" yaml:"-"`
}

func NewWorkflowActivityCallOperation(context *parsing.Context) *WorkflowActivityCallOperation {
	return &WorkflowActivityCallOperation{
		Entity: NewEntity(context),
		Name:   context.Name,
		Inputs: make(Values),
	}
}

// ([parsing.Reader] signature)
func ReadWorkflowActivityCallOperation(context *parsing.Context) parsing.EntityPtr {
	self := NewWorkflowActivityCallOperation(context)

	if context.Is(ard.TypeMap) {
		// Long notation (introduced in TOSCA 1.3)
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.InterfaceAndOperation = context.FieldChild("operation", context.Data).ReadString()
	}

	return self
}

func (self *WorkflowActivityCallOperation) Render(stepDefinition *WorkflowStepDefinition) {
	if self.InterfaceAndOperation == nil {
		return
	}

	// Parse operation spec
	s := strings.SplitN(*self.InterfaceAndOperation, ".", 2)
	if len(s) != 2 {
		self.Context.FieldChild("operation", *self.InterfaceAndOperation).ReportValueWrongFormat("interface.operation")
		return
	}

	var ok bool

	// Lookup interface by name
	var interfaceDefinition *InterfaceDefinition
	if stepDefinition.TargetNodeTemplate != nil {
		if self.Interface, ok = stepDefinition.TargetNodeTemplate.Interfaces[s[0]]; ok {
			interfaceDefinition, _ = self.Interface.GetDefinitionForNodeTemplate(stepDefinition.TargetNodeTemplate)
		} else {
			self.Context.FieldChild("operation", s[0]).ReportReferenceNotFound("interface", stepDefinition.TargetNodeTemplate)
			return
		}
	} else if stepDefinition.TargetGroup != nil {
		if self.Interface, ok = stepDefinition.TargetGroup.Interfaces[s[0]]; ok {
			interfaceDefinition, _ = self.Interface.GetDefinitionForGroup(stepDefinition.TargetGroup)
		} else {
			self.Context.FieldChild("operation", s[0]).ReportReferenceNotFound("interface", stepDefinition.TargetGroup)
			return
		}
	} else {
		// There was a lookup problem (neither node template nor group)
		return
	}

	// Lookup operation by name
	var operationDefinition *OperationDefinition
	if self.Operation, ok = self.Interface.Operations[s[1]]; ok {
		if interfaceDefinition != nil {
			operationDefinition, _ = interfaceDefinition.OperationDefinitions[self.Operation.Name]
		}
	} else {
		self.Context.FieldChild("operation", s[1]).ReportReferenceNotFound("operation", self.Interface)
	}

	if operationDefinition != nil {
		self.Inputs.RenderInputs(operationDefinition.InputDefinitions, self.Context.FieldChild("inputs", nil))
	}
}

func (self *WorkflowActivityCallOperation) Normalize(normalWorkflowActivity *normal.WorkflowActivity) {
	logNormalize.Debug("workflow activity call operation")

	if (self.Interface == nil) || (self.Operation == nil) {
		return
	}

	normalCallOperation := normalWorkflowActivity.NewCallOperation()

	var normalInterface *normal.Interface
	if normalWorkflowActivity.Step.TargetNodeTemplate != nil {
		normalInterface = normalWorkflowActivity.Step.TargetNodeTemplate.Interfaces[self.Interface.Name]
	} else if normalWorkflowActivity.Step.TargetGroup != nil {
		normalInterface = normalWorkflowActivity.Step.TargetGroup.Interfaces[self.Interface.Name]
	} else {
		return
	}

	normalCallOperation.Operation, _ = normalInterface.Operations[self.Operation.Name]
	self.Inputs.Normalize(normalCallOperation.Inputs)
}
