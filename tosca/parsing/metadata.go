package parsing

import (
	"strings"
)

//
// HasMetadata
//

const (
	MetadataType                  = "puccini.type"
	MetadataConverter             = "puccini.converter"
	MetadataComparer              = "puccini.comparer"
	MetadataQuirks                = "puccini.quirks"
	MetadataDataTypePrefix        = "puccini.data-type:"
	MetadataScriptletPrefix       = "puccini.scriptlet:"
	MetadataScriptletImportPrefix = "puccini.scriptlet.import:"

	MetadataCanonicalName   = "tosca.canonical-name"
	MetadataNormative       = "tosca.normative"
	MetadataFunctionPrefix  = "tosca.function."
	MetadataContraintPrefix = "tosca.constraint."
)

type HasMetadata interface {
	GetDescription() (string, bool)
	GetMetadata() (map[string]string, bool) // should return a copy
	SetMetadata(name string, value string) bool
}

// From [HasMetadata] interface
func GetDescription(entityPtr EntityPtr) (string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetDescription()
	} else {
		return "", false
	}
}

// From [HasMetadata] interface
func GetMetadata(entityPtr EntityPtr) (map[string]string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetMetadata()
	} else {
		return nil, false
	}
}

// From [HasMetadata] interface
func SetMetadata(entityPtr EntityPtr, name string, value string) bool {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		hasMetadata.SetMetadata(name, value)
		return true
	} else {
		return false
	}
}

func GetDataTypeMetadata(metadata map[string]string) map[string]string {
	dataTypeMetadata := make(map[string]string)
	if metadata != nil {
		for key, value := range metadata {
			if strings.HasPrefix(key, MetadataDataTypePrefix) {
				dataTypeMetadata[key[len(MetadataDataTypePrefix):]] = value
			}
		}
	}
	return dataTypeMetadata
}
