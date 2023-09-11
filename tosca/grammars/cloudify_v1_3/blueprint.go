package cloudify_v1_3

import (
	"fmt"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

//
// Blueprint
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/]
//

type Blueprint struct {
	*File `name:"blueprint"`

	Description *string `read:"description"` // not in spec, but in code
	Groups      Groups  `read:"groups,Group"`
}

func NewBlueprint(context *parsing.Context) *Blueprint {
	return &Blueprint{File: NewFile(context)}
}

// ([parsing.Reader] signature)
func ReadBlueprint(context *parsing.Context) parsing.EntityPtr {
	self := NewBlueprint(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	return self
}

// parsing.HasInputs interface
func (self *Blueprint) SetInputs(inputs map[string]ard.Value) {
	context := self.Context.FieldChild("inputs", nil)
	for name, data := range inputs {
		childContext := context.MapChild(name, data)
		if input, ok := self.Inputs[name]; ok {
			input.Value = ReadValue(childContext).(*Value)
		} else {
			childContext.ReportUndeclared("input")
		}
	}
}

// normal.Normalizable interface
func (self *Blueprint) NormalizeServiceTemplate() *normal.ServiceTemplate {
	logNormalize.Debug("blueprint")

	normalServiceTemplate := normal.NewServiceTemplate()

	if self.Metadata != nil {
		for key, value := range self.Metadata {
			// TODO: does Cloudify DSL really allow for any kind of value?
			normalServiceTemplate.Metadata[yamlkeys.KeyString(key)] = fmt.Sprintf("%s", value)
		}
	}

	if self.Description != nil {
		normalServiceTemplate.Description = *self.Description
	}

	normalServiceTemplate.ScriptletNamespace = self.Context.ScriptletNamespace

	self.Inputs.Normalize(normalServiceTemplate.Inputs, self.Context.FieldChild("inputs", nil))
	self.Outputs.Normalize(normalServiceTemplate.Outputs)
	self.NodeTemplates.Normalize(normalServiceTemplate)
	self.Groups.Normalize(normalServiceTemplate)
	self.Workflows.Normalize(normalServiceTemplate)
	self.Policies.Normalize(normalServiceTemplate)

	// TODO: normalize plugins
	// TODO: normalize upload resources
	// TODO: normalize capabilities

	return normalServiceTemplate
}
