package tosca_v1_0

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Group
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.5
//

// ([parsing.Reader] signature)
func ReadGroup(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Metadata", "")

	return tosca_v1_2.ReadGroup(context)
}
