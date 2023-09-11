package tosca_v1_2

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// RequirementAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.2
//

// ([parsing.Reader] signature)
func ReadRequirementAssignment(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Count", "")
	context.SetReadTag("Directives", "")
	context.SetReadTag("Optional", "")
	//context.SetReadTag("Allocation", "")

	self := tosca_v2_0.NewRequirementAssignment(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(append(context.ReadFields(self)))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.TargetNodeTemplateNameOrTypeName = context.FieldChild("node", context.Data).ReadString()
	}

	return self
}
