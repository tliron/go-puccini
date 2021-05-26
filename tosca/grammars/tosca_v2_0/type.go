package tosca_v2_0

import (
	"sync"

	"github.com/tliron/puccini/tosca"
)

//
// Type
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.1
//

type Type struct {
	*Entity `json:"-" yaml:"-"`
	Name    string `namespace:""`

	ParentName  *string  `read:"derived_from"`
	Version     *Value   `read:"version,Value"`
	Metadata    Metadata `read:"metadata,!Metadata"`
	Description *string  `read:"description"`

	metadataLock sync.RWMutex
}

func NewType(context *tosca.Context) *Type {
	return &Type{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.HasMetadata interface
func (self *Type) GetDescription() (string, bool) {
	self.metadataLock.RLock()
	defer self.metadataLock.RUnlock()

	if self.Description != nil {
		return *self.Description, true
	}
	return "", false
}

// tosca.HasMetadata interface
func (self *Type) GetMetadata() (map[string]string, bool) {
	self.metadataLock.RLock()
	defer self.metadataLock.RUnlock()

	metadata := make(map[string]string)
	if self.Metadata != nil {
		for key, value := range self.Metadata {
			metadata[key] = value
		}
	}
	return metadata, true
}

// tosca.HasMetadata interface
func (self *Type) SetMetadata(name string, value string) bool {
	self.metadataLock.Lock()
	defer self.metadataLock.Unlock()

	self.Metadata[name] = value
	return true
}

// parser.Renderable interface
func (self *Type) Render() {
	logRender.Debugf("type: %s", self.Name)
	if self.Version != nil {
		self.Version.RenderDataType("version")
	}
}

func (self *Type) GetMetadataValue(key string) (string, bool) {
	if self.Metadata != nil {
		value, ok := self.Metadata[key]
		return value, ok
	}
	return "", false
}
