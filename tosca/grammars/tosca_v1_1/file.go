package tosca_v1_1

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// File
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.9
//

// ([parsing.Reader] signature)
func ReadFile(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Profile", "")

	self := tosca_v2_0.NewFile(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	ignore := []string{"dsl_definitions"}
	if context.HasQuirk(parsing.QuirkImportsTopologyTemplateIgnore) {
		ignore = append(ignore, "topology_template")
	}
	if context.HasQuirk(parsing.QuirkAnnotationsIgnore) {
		ignore = append(ignore, "annotation_types")
	}
	context.ValidateUnsupportedFields(append(context.ReadFields(self), ignore...))
	return self
}
