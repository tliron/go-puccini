package tosca_v1_0

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// ServiceTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.8
//

// ([parsing.Reader] signature)
func ReadServiceTemplate(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("WorkflowDefinitions", "")

	return tosca_v2_0.ReadServiceTemplate(context)
}
