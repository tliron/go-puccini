package tosca_v2_0

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// WorkflowDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.7
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.7
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.7
//

type WorkflowDefinition struct {
	*Entity `name:"workflow definition"`
	Name    string `namespace:""`

	Metadata                Metadata                        `read:"metadata,Metadata"`
	Description             *string                         `read:"description"`
	InputDefinitions        PropertyDefinitions             `read:"inputs,PropertyDefinition"`
	PreconditionDefinitions WorkflowPreconditionDefinitions `read:"preconditions,WorkflowPreconditionDefinition"`
	StepDefinitions         WorkflowStepDefinitions         `read:"steps,WorkflowStepDefinition"`
}

func NewWorkflowDefinition(context *parsing.Context) *WorkflowDefinition {
	return &WorkflowDefinition{
		Entity:           NewEntity(context),
		Name:             context.Name,
		InputDefinitions: make(PropertyDefinitions),
		StepDefinitions:  make(WorkflowStepDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadWorkflowDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewWorkflowDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Mappable] interface)
func (self *WorkflowDefinition) GetKey() string {
	return self.Name
}

// ([parsing.Renderable] interface)
func (self *WorkflowDefinition) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *WorkflowDefinition) render() {
	logRender.Debugf("workflow definition: %s", self.Name)

	self.StepDefinitions.Render()
}

func (self *WorkflowDefinition) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.Workflow {
	logNormalize.Debugf("workflow definition: %s", self.Name)

	normalWorkflow := normalServiceTemplate.NewWorkflow(self.Name)

	normalWorkflow.Metadata = self.Metadata

	if self.Description != nil {
		normalWorkflow.Description = *self.Description
	}

	// TODO: PreconditionDefinitions

	// TODO: support property definitions
	//self.InputDefinitions.Normalize(w.Inputs)

	self.StepDefinitions.Normalize(normalWorkflow)

	return normalWorkflow
}

//
// WorkflowDefinitions
//

type WorkflowDefinitions map[string]*WorkflowDefinition

func (self WorkflowDefinitions) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, workflowDefinition := range self {
		normalServiceTemplate.Workflows[workflowDefinition.Name] = workflowDefinition.Normalize(normalServiceTemplate)
	}
}
