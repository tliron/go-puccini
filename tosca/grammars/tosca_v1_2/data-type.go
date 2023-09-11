package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// DataType
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.6
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.6
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.5
//

// ([parsing.Reader] signature)
func ReadDataType(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("KeySchema", "")
	context.SetReadTag("EntrySchema", "")

	return tosca_v2_0.ReadDataType(context)
}
