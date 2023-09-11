package normal

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/go-ard"
)

//
// Pointer
//

type Pointer struct {
	NodeTemplate *NodeTemplate
	Target       string
}

func NewPointer(target string) *Pointer {
	return &Pointer{
		Target: target,
	}
}

func (self *NodeTemplate) NewPointer(target string) *Pointer {
	return &Pointer{
		NodeTemplate: self,
		Target:       target,
	}
}

type MarshalablePointer struct {
	NodeTemplateName string `json:"nodeTemplateName,omitempty" yaml:"nodeTemplateName,omitempty"`
	Target           string `json:"target,omitempty" yaml:"target,omitempty"`
}

func (self *Pointer) Marshalable() any {
	if self.NodeTemplate != nil {
		return &MarshalablePointer{
			NodeTemplateName: self.NodeTemplate.Name,
			Target:           self.Target,
		}
	} else {
		return &MarshalablePointer{
			Target: self.Target,
		}
	}
}

// ([json.Marshaler] interface)
func (self *Pointer) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// ([yaml.Marshaler] interface)
func (self *Pointer) MarshalYAML() (any, error) {
	return self.Marshalable(), nil
}

// ([cbor.Marshaler] interface)
func (self *Pointer) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(self.Marshalable())
}

// ([msgpack.Marshaler] interface)
func (self *Pointer) MarshalMsgpack() ([]byte, error) {
	return ard.MarshalMessagePack(self.Marshalable())
}

// ([ard.ToARD] interface)
func (self *Pointer) ToARD(reflector *ard.Reflector) (any, error) {
	return reflector.Unpack(self.Marshalable())
}

//
// Pointers
//

type Pointers map[string]*Pointer
