package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Import
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.8
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.8
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.7
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.7
//

// ([parsing.Reader] signature)
func ReadImport(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("URL", "file")
	context.SetReadTag("Namespace", "namespace_prefix")
	context.SetReadTag("NamespaceURI", "namespace_uri")

	self := tosca_v2_0.NewImport(context)

	if context.Is(ard.TypeMap) {
		if context.HasQuirk(parsing.QuirkImportsSequencedList) {
			map_ := context.Data.(ard.Map)
			if len(map_) == 1 {
				for _, data := range map_ {
					if data_, ok := data.(ard.Map); ok {
						context.Data = data_
					}
					break
				}
			}
		}

		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.URL = context.FieldChild("file", context.Data).ReadString()
	}

	return self
}
