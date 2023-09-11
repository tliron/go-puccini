package tosca_v1_0

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// RelationshipTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.4
//

// ([parsing.Reader] signature)
func ReadRelationshipTemplate(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("Metadata", "")

	return tosca_v2_0.ReadRelationshipTemplate(context)
}
