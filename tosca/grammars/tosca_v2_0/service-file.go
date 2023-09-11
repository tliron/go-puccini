package tosca_v2_0

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// ServiceFile
//
// See File
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.9
//

type ServiceFile struct {
	*File `name:"service file"`

	ServiceTemplate *ServiceTemplate `read:"service_template,ServiceTemplate"`
}

func NewServiceFile(context *parsing.Context) *ServiceFile {
	return &ServiceFile{File: NewFile(context)}
}

// ([parsing.Reader] signature)
func ReadServiceFile(context *parsing.Context) parsing.EntityPtr {
	context.FunctionPrefix = "$"
	self := NewServiceFile(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	ignore := []string{"dsl_definitions"}
	if context.HasQuirk(parsing.QuirkAnnotationsIgnore) {
		ignore = append(ignore, "annotation_types")
	}
	context.ValidateUnsupportedFields(append(context.ReadFields(self), ignore...))
	if self.Profile != nil {
		context.CanonicalNamespace = self.Profile
	}
	return self
}

// normal.Normalizable interface
func (self *ServiceFile) NormalizeServiceTemplate() *normal.ServiceTemplate {
	logNormalize.Debug("service file")

	normalServiceTemplate := normal.NewServiceTemplate()

	if self.Description != nil {
		normalServiceTemplate.Description = *self.Description
	}

	normalServiceTemplate.ScriptletNamespace = self.Context.ScriptletNamespace

	self.File.Normalize(normalServiceTemplate)
	if self.ServiceTemplate != nil {
		self.ServiceTemplate.Normalize(normalServiceTemplate)
	}

	return normalServiceTemplate
}
