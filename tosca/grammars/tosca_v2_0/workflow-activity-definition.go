package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

//
// WorkflowActivityDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.23
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.19
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.17
//

type WorkflowActivityDefinition struct {
	*Entity `name:"workflow activity definition"`

	DelegateWorkflowDefinitionName *string
	InlineWorkflowDefinitionName   *string
	SetNodeState                   *string
	CallOperation                  *WorkflowActivityCallOperation

	DelegateWorkflowDefinition *WorkflowDefinition `lookup:"delegate,DelegateWorkflowDefinitionName" traverse:"ignore" json:"-" yaml:"-"`
	InlineWorkflowDefinition   *WorkflowDefinition `lookup:"inline,InlineWorkflowDefinitionName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewWorkflowActivityDefinition(context *parsing.Context) *WorkflowActivityDefinition {
	return &WorkflowActivityDefinition{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadWorkflowActivityDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewWorkflowActivityDefinition(context)

	if context.ValidateType(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if len(map_) != 1 {
			context.ReportValueMalformed("workflow activity definition", "map length not 1")
			return self
		}

		for key, value := range map_ {
			operator := yamlkeys.KeyString(key)
			childContext := context.FieldChild(operator, value)

			switch operator {
			case "delegate":
				self.DelegateWorkflowDefinitionName = childContext.ReadString()
			case "inline":
				self.InlineWorkflowDefinitionName = childContext.ReadString()
			case "set_state":
				self.SetNodeState = childContext.ReadString()
			case "call_operation":
				if reader, ok := Grammar.Readers["WorkflowActivityCallOperation"]; ok {
					self.CallOperation = reader(childContext).(*WorkflowActivityCallOperation)
				} else {
					childContext.ReportValueMalformed("workflow activity definition", "unsupported operator")
				}
			default:
				childContext.ReportValueMalformed("workflow activity definition", "unsupported operator")
				return self
			}

			// We have only one key
			break
		}
	}

	return self
}

func (self *WorkflowActivityDefinition) Render(stepDefinition *WorkflowStepDefinition) {
	if self.CallOperation != nil {
		self.CallOperation.Render(stepDefinition)
	}
}

func (self *WorkflowActivityDefinition) Normalize(normalWorkflowStep *normal.WorkflowStep) *normal.WorkflowActivity {
	logNormalize.Debug("workflow activity")

	normalWorkflowActivity := normalWorkflowStep.NewActivity()
	if self.DelegateWorkflowDefinition != nil {
		normalWorkflowActivity.DelegateWorkflow = normalWorkflowStep.Workflow.ServiceTemplate.Workflows[self.DelegateWorkflowDefinition.Name]
	} else if self.InlineWorkflowDefinition != nil {
		normalWorkflowActivity.InlineWorkflow = normalWorkflowStep.Workflow.ServiceTemplate.Workflows[self.InlineWorkflowDefinition.Name]
	} else if self.SetNodeState != nil {
		normalWorkflowActivity.SetNodeState = *self.SetNodeState
	} else if self.CallOperation != nil {
		self.CallOperation.Normalize(normalWorkflowActivity)
	}

	return normalWorkflowActivity
}

//
// WorkflowActivityDefinitions
//

type WorkflowActivityDefinitions []*WorkflowActivityDefinition
