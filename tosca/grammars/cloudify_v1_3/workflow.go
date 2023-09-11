package cloudify_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Workflow
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-workflows/]
//

type Workflow struct {
	*Entity `name:"workflow"`
	Name    string `namespace:""`

	Mapping              *string              `read:"mapping" mandatory:""`
	ParameterDefinitions ParameterDefinitions `read:"parameters,ParameterDefinition"`
	IsCascading          *bool                `read:"is_cascading"` // See: https://docs.cloudify.co/5.0.5/working_with/service_composition/component/
}

func NewWorkflow(context *parsing.Context) *Workflow {
	return &Workflow{
		Entity:               NewEntity(context),
		Name:                 context.Name,
		ParameterDefinitions: make(ParameterDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadWorkflow(context *parsing.Context) parsing.EntityPtr {
	self := NewWorkflow(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.Mapping = context.FieldChild("mapping", context.Data).ReadString()
	}

	return self
}

func (self *Workflow) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.Workflow {
	logNormalize.Debugf("workflow: %s", self.Name)

	normalWorkflow := normalServiceTemplate.NewWorkflow(self.Name)

	// TODO: mapping

	// TODO: support property definitions
	//self.ParameterDefinitions.Normalize(w.Inputs)

	return normalWorkflow
}

//
// Workflows
//

type Workflows []*Workflow

func (self Workflows) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, workflow := range self {
		normalServiceTemplate.Workflows[workflow.Name] = workflow.Normalize(normalServiceTemplate)
	}
}
