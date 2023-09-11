package tosca_v2_0

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// WorkflowStepDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.27
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.23
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.21
//

type WorkflowStepDefinition struct {
	*Entity `name:"workflow step definition"`
	Name    string

	TargetNodeTemplateOrGroupName *string                     `read:"target" mandatory:""`
	TargetNodeRequirementName     *string                     `read:"target_relationship"`
	OperationHost                 *string                     `read:"operation_host"`
	FilterConditionClauses        *ConditionClause            `read:"filter,ConditionClauseAnd"` // spec is wrong, says constraint clause
	ActivityDefinitions           WorkflowActivityDefinitions `read:"activities,[]WorkflowActivityDefinition" mandatory:""`
	OnSuccessStepNames            *[]string                   `read:"on_success"`
	OnFailureStepNames            *[]string                   `read:"on_failure"`

	TargetNodeTemplate *NodeTemplate             `lookup:"target,TargetNodeTemplateOrGroupName" traverse:"ignore" json:"-" yaml:"-"`
	TargetGroup        *Group                    `lookup:"target,TargetNodeTemplateOrGroupName" traverse:"ignore" json:"-" yaml:"-"`
	OnSuccessSteps     []*WorkflowStepDefinition // custom lookup
	OnFailureSteps     []*WorkflowStepDefinition // custom lookup
}

func NewWorkflowStepDefinition(context *parsing.Context) *WorkflowStepDefinition {
	return &WorkflowStepDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadWorkflowStepDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewWorkflowStepDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Mappable] interface)
func (self *WorkflowStepDefinition) GetKey() string {
	return self.Name
}

func (self *WorkflowStepDefinition) Render(definitions WorkflowStepDefinitions) {
	logInherit.Debugf("workflow step definition: %s", self.Name)

	if self.OnSuccessStepNames != nil {
		for index, name := range *self.OnSuccessStepNames {
			if definition, ok := definitions[name]; ok {
				self.OnSuccessSteps = append(self.OnSuccessSteps, definition)
			} else {
				self.Context.ListChild(index, name).ReportUnknown("step")
			}
		}
	}

	if self.OnFailureStepNames != nil {
		for index, name := range *self.OnFailureStepNames {
			if definition, ok := definitions[name]; ok {
				self.OnFailureSteps = append(self.OnFailureSteps, definition)
			} else {
				self.Context.ListChild(index, name).ReportUnknown("step")
			}
		}
	}

	for _, activity := range self.ActivityDefinitions {
		activity.Render(self)
	}

	// TODO: validate OperationHost
}

func (self *WorkflowStepDefinition) Normalize(normalWorkflow *normal.Workflow) *normal.WorkflowStep {
	logNormalize.Debugf("workflow step definition: %s", self.Name)

	normalWorkflowStep := normalWorkflow.NewStep(self.Name)

	if self.TargetNodeTemplate != nil {
		if normalNodeTemplate, ok := normalWorkflow.ServiceTemplate.NodeTemplates[self.TargetNodeTemplate.Name]; ok {
			normalWorkflowStep.TargetNodeTemplate = normalNodeTemplate
		}
	} else if self.TargetGroup != nil {
		if normalGroup, ok := normalWorkflow.ServiceTemplate.Groups[self.TargetGroup.Name]; ok {
			normalWorkflowStep.TargetGroup = normalGroup
		}
	}

	for _, activity := range self.ActivityDefinitions {
		activity.Normalize(normalWorkflowStep)
	}

	return normalWorkflowStep
}

func (self *WorkflowStepDefinition) NormalizeNext(normalWorkflowStep *normal.WorkflowStep, normalWorkflow *normal.Workflow) {
	for _, next := range self.OnSuccessSteps {
		normalWorkflowStep.OnSuccessSteps = append(normalWorkflowStep.OnSuccessSteps, normalWorkflow.Steps[next.Name])
	}

	for _, next := range self.OnFailureSteps {
		normalWorkflowStep.OnFailureSteps = append(normalWorkflowStep.OnFailureSteps, normalWorkflow.Steps[next.Name])
	}

	if self.OperationHost != nil {
		normalWorkflowStep.Host = *self.OperationHost
	}
}

//
// WorkflowStepDefinitions
//

type WorkflowStepDefinitions map[string]*WorkflowStepDefinition

func (self WorkflowStepDefinitions) Render() {
	for _, step := range self {
		step.Render(self)
	}
}

func (self WorkflowStepDefinitions) Normalize(normalWorkflow *normal.Workflow) {
	steps := make(normal.WorkflowSteps)
	for name, step := range self {
		steps[name] = step.Normalize(normalWorkflow)
	}
	for name, step := range self {
		step.NormalizeNext(steps[name], normalWorkflow)
	}
}
