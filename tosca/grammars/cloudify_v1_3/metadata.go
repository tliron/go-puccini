package cloudify_v1_3

import (
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Metadata
//
// Note: not in spec
//

type Metadata map[string]ard.Value

// ([parsing.Reader] signature)
func ReadMetadata(context *parsing.Context) parsing.EntityPtr {
	var self map[string]ard.Value

	if context.ValidateType(ard.TypeMap) {
		metadata := context.ReadStringMap()
		if metadata != nil {
			self = *metadata
		}
	}

	if self != nil {
		for key, value := range self {
			if strings.HasPrefix(key, parsing.MetadataScriptletImportPrefix) {
				if v, ok := value.(string); ok {
					context.ImportScriptlet(key[len(parsing.MetadataScriptletImportPrefix):], v)
					delete(self, key)
				} else {
					context.MapChild(key, value).ReportValueWrongType(ard.TypeString)
				}
			} else if strings.HasPrefix(key, parsing.MetadataScriptletPrefix) {
				if v, ok := value.(string); ok {
					context.EmbedScriptlet(key[len(parsing.MetadataScriptletPrefix):], v)
					delete(self, key)
				} else {
					context.MapChild(key, value).ReportValueWrongType(ard.TypeString)
				}
			}
		}
	}

	return self
}
