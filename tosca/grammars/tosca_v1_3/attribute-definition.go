package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// AttributeDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.12
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.10
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.10
//

// ([parsing.Reader] signature)
func ReadAttributeDefinition(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("ConstraintClauses", "")

	return tosca_v2_0.ReadAttributeDefinition(context).(*tosca_v2_0.AttributeDefinition)
}
