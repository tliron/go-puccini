package tosca_v1_0

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Policy
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.6
//

// ([parsing.Reader] signature)
func ReadPolicy(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Metadata", "")
	context.SetReadTag("TriggerDefinitions", "")

	return tosca_v2_0.ReadPolicy(context)
}
