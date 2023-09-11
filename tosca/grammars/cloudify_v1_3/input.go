package cloudify_v1_3

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Input
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-inputs/]
//

type Input struct {
	*ParameterDefinition `name:"input"`
	Name                 string `namespace:""`

	Value *Value
}

func NewInput(context *parsing.Context) *Input {
	return &Input{
		ParameterDefinition: NewParameterDefinition(context),
		Name:                context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadInput(context *parsing.Context) parsing.EntityPtr {
	self := NewInput(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *Input) Normalize(context *parsing.Context) normal.Value {
	value := self.Value
	if value == nil {
		if self.Default != nil {
			value = self.Default
		} else {
			// Inputs should always appear, even if they have no default value
			// (Note that in Cloudify DSL they are *always* required, so it would be abnormal)
			value = NewValue(context.MapChild(self.Name, nil))
		}
	}
	return value.Normalize()
}

// ([parsing.Renderable] interface)
func (self *Input) Render() {
	self.renderOnce.Do(self.render)
}

func (self *Input) render() {
	logRender.Debugf("input: %s", self.Name)

	if (self.Value == nil) && (self.Default != nil) {
		self.Value = self.Default
	}

	if self.Value == nil {
		self.Context.ReportValueRequired("input")
	}
}

//
// Inputs
//

type Inputs map[string]*Input

func (self Inputs) Normalize(normalConstrainables normal.Values, context *parsing.Context) {
	for key, input := range self {
		normalConstrainables[key] = input.Normalize(context)
	}
}
