package normal

//
// Workflow
//

type Workflow struct {
	ServiceTemplate *ServiceTemplate `json:"-" yaml:"-"`
	Name            string           `json:"-" yaml:"-"`

	Metadata      map[string]string       `json:"metadata" yaml:"metadata"`
	Description   string                  `json:"description" yaml:"description"`
	Preconditions []*WorkflowPrecondition `json:"preconditions" yaml:"preconditions"`
	Steps         WorkflowSteps           `json:"steps" yaml:"steps"`
	Inputs        Values                  `json:"inputs" yaml:"inputs"`
}

func (self *ServiceTemplate) NewWorkflow(name string) *Workflow {
	workflow := &Workflow{
		ServiceTemplate: self,
		Name:            name,
		Metadata:        make(map[string]string),
		Steps:           make(WorkflowSteps),
		Preconditions:   make([]*WorkflowPrecondition, 0),
		Inputs:          make(Values),
	}
	self.Workflows[name] = workflow
	return workflow
}

//
// Workflows
//

type Workflows map[string]*Workflow

//
// WorkflowPrecondition
//

type WorkflowPrecondition struct {
	TargetNodeTemplate *NodeTemplate `json:"-" yaml:"-"`
	TargetGroup        *Group        `json:"-" yaml:"-"`
	// TODO: Conditions
}

//
// WorkflowStep
//

// TODO: JSON/YAML marshalling

type WorkflowStep struct {
	Workflow           *Workflow     `json:"-" yaml:"-"`
	Name               string        `json:"-" yaml:"-"`
	TargetNodeTemplate *NodeTemplate `json:"-" yaml:"-"`
	TargetGroup        *Group        `json:"-" yaml:"-"`
	// TODO: Filters
	OnSuccessSteps []*WorkflowStep     `json:"-" yaml:"-"`
	OnFailureSteps []*WorkflowStep     `json:"-" yaml:"-"`
	Activities     []*WorkflowActivity `json:"-" yaml:"-"`
	Host           string              `json:"-" yaml:"-"`
}

func (self *Workflow) NewStep(name string) *WorkflowStep {
	step := &WorkflowStep{
		Workflow:       self,
		Name:           name,
		OnSuccessSteps: make([]*WorkflowStep, 0),
		OnFailureSteps: make([]*WorkflowStep, 0),
		Activities:     make([]*WorkflowActivity, 0),
	}
	self.Steps[name] = step
	return step
}

//
// WorkflowSteps
//

type WorkflowSteps map[string]*WorkflowStep

//
// WorkflowActivity
//

type WorkflowActivity struct {
	Step             *WorkflowStep  `json:"-" yaml:"-"`
	DelegateWorkflow *Workflow      `json:"-" yaml:"-"`
	InlineWorkflow   *Workflow      `json:"-" yaml:"-"`
	SetNodeState     string         `json:"setNodeState" yaml:"setNodeState"`
	CallOperation    *CallOperation `json:"callOperation" yaml:"callOperation"`
}

func (self *WorkflowStep) NewActivity() *WorkflowActivity {
	activity := &WorkflowActivity{Step: self}
	self.Activities = append(self.Activities, activity)
	return activity
}

//
// CallOperation
//

type CallOperation struct {
	Activity  *WorkflowActivity `json:"-" yaml:"-"`
	Operation *Operation        `json:"operation" yaml:"operation"`
	Inputs    Values            `json:"inputs" yaml:"inputs"`
}

func (self *WorkflowActivity) NewCallOperation() *CallOperation {
	callOperation := &CallOperation{
		Activity: self,
		Inputs:   make(Values),
	}
	self.CallOperation = callOperation
	return callOperation
}
