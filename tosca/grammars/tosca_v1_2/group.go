package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Group
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.5
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.5
//

// ([parsing.Reader] signature)
func ReadGroup(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Interfaces", "interfaces,InterfaceAssignment")

	return tosca_v2_0.ReadGroup(context)
}
