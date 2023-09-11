package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// InterfaceType
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.5
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.4
//

// ([parsing.Reader] signature)
func ReadInterfaceType(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("OperationDefinitions", "?,OperationDefinition")
	context.SetReadTag("NotificationDefinitions", "")

	self := tosca_v2_0.NewInterfaceType(context)
	context.ReadFields(self)
	return self
}
