package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

// ([parsing.Reader] signature)
func ReadAttributeValue(context *parsing.Context) parsing.EntityPtr {
	self := tosca_v2_0.NewValue(context)

	// Unpack long notation (only for attributes)
	// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.13.2.2
	// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.12.2.2
	// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.11.2.2
	// [TOSCA-Simple-Profile-YAML-v1.0] @3.5.11.2.2
	if context.Is(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if len(map_) == 2 {
			if description, ok := map_["description"]; ok {
				if value, ok := map_["value"]; ok {
					self.Description = context.FieldChild("description", description).ReadString()
					context.Data = value
				}
			}
		}
	}

	tosca_v2_0.ParseFunctionCall(context)

	return self
}
