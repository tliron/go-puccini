package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// SubstitutionMappings
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

// ([parsing.Reader] signature)
func ReadSubstitutionMappings(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("AttributeMappings", "")
	context.SetReadTag("SubstitutionFilter", "")

	return tosca_v2_0.ReadSubstitutionMappings(context)
}
