package tosca_v1_1

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// CapabilityDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.1
//

// ([parsing.Reader] signature)
func ReadCapabilityDefinition(context *parsing.Context) parsing.EntityPtr {
	// Handle TOSCA 1.1 specific fields before calling v2.0 reader
	if context.Is(ard.TypeMap) {
		if m, ok := context.Data.(ard.Map); ok {
			// TOSCA 1.1 uses "valid_source_types" which doesn't exist in TOSCA 2.0
			if _, ok := m["valid_source_types"]; ok {
				delete(m, "valid_source_types")
			}

			// TOSCA 1.1 uses "occurrences" which doesn't exist in TOSCA 2.0 CapabilityDefinition
			if _, ok := m["occurrences"]; ok {
				delete(m, "occurrences")
			}
		}
	}

	return tosca_v2_0.ReadCapabilityDefinition(context)
}
