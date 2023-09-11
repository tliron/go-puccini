package tosca_v2_0

import (
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Metadata
//

type Metadata map[string]string

// ([parsing.Reader] signature)
func ReadMetadata(context *parsing.Context) parsing.EntityPtr {
	var self Metadata

	if context.Is(ard.TypeMap) {
		metadata := context.ReadStringStringMap()
		if metadata != nil {
			self = *metadata
		}
	}

	if self != nil {
		for key, value := range self {
			if strings.HasPrefix(key, parsing.MetadataScriptletImportPrefix) {
				context.ImportScriptlet(key[len(parsing.MetadataScriptletImportPrefix):], value)
				delete(self, key)
			} else if strings.HasPrefix(key, parsing.MetadataScriptletPrefix) {
				context.EmbedScriptlet(key[len(parsing.MetadataScriptletPrefix):], value)
				delete(self, key)
			}
		}
	}

	return self
}
