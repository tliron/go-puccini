package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
)

//
// Unit
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.9
//

// tosca.Reader signature
func ReadUnit(context *tosca.Context) tosca.EntityPtr {
	context.SetReadTag("Profile", "namespace")

	self := tosca_v2_0.NewUnit(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	if context.HasQuirk(tosca.QuirkImportsTopologyTemplateIgnore) {
		context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions", "topology_template"))
	} else {
		context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	}
	return self
}
