package tosca

import (
	"strings"
)

const (
	METADATA_INFORMATION_PREFIX      = "puccini.information:"
	METADATA_SCRIPTLET_PREFIX        = "puccini.scriptlet:"
	METADATA_SCRIPTLET_IMPORT_PREFIX = "puccini.scriptlet.import:"
	METADATA_CANONICAL_NAME          = "tosca.canonical-name"
	METADATA_NORMATIVE               = "tosca.normative"
)

//
// HasMetadata
//

// Must be thread-safe!
type HasMetadata interface {
	GetDescription() (string, bool)
	GetMetadata() (map[string]string, bool) // should return a copy
	SetMetadata(name string, value string) bool
}

// From HasMetadata interface
func GetDescription(entityPtr EntityPtr) (string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetDescription()
	}
	return "", false
}

// From HasMetadata interface
func GetMetadata(entityPtr EntityPtr) (map[string]string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetMetadata()
	}
	return nil, false
}

// From HasMetadata interface
func SetMetadata(entityPtr EntityPtr, name string, value string) bool {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		hasMetadata.SetMetadata(name, value)
		return true
	}
	return false
}

func GetInformationMetadata(metadata map[string]string) map[string]string {
	informationMetadata := make(map[string]string)
	if metadata != nil {
		for key, value := range metadata {
			if strings.HasPrefix(key, METADATA_INFORMATION_PREFIX) {
				informationMetadata[key[len(METADATA_INFORMATION_PREFIX):]] = value
			}
		}
	}
	return informationMetadata
}
