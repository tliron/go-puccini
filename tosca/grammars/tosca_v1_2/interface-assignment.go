package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// InterfaceAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.16
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.14
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.14
//

// ([parsing.Reader] signature)
func ReadInterfaceAssignment(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Operations", "?,OperationAssignment")
	context.SetReadTag("Notifications", "")

	self := tosca_v2_0.NewInterfaceAssignment(context)
	context.ReadFields(self)
	return self
}
