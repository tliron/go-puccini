package tosca_v1_0

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// NodeTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.3
//

// ([parsing.Reader] signature)
func ReadNodeTemplate(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Metadata", "")

	return tosca_v1_3.ReadNodeTemplate(context)
}
