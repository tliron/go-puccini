package tosca_v1_1

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// SubstitutionMappings
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

// ([parsing.Reader] signature)
func ReadSubstitutionMappings(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("PropertyMappings", "")
	context.SetReadTag("InterfaceMappings", "")

	return tosca_v1_2.ReadSubstitutionMappings(context)
}
