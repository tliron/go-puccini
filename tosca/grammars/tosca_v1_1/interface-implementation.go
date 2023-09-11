package tosca_v1_1

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// InterfaceImplementation
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13.2.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.13.2.3
//

// ([parsing.Reader] signature)
func ReadInterfaceImplementation(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Timeout", "")
	context.SetReadTag("OperationHost", "")

	return tosca_v2_0.ReadInterfaceImplementation(context)
}
