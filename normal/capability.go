package normal

import (
	"encoding/json"
	"math"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/go-ard"
)

//
// Capability
//

type Capability struct {
	NodeTemplate *NodeTemplate
	Name         string

	Description          string
	Types                EntityTypes
	Properties           Values
	Attributes           Values
	MinRelationshipCount uint64
	MaxRelationshipCount uint64
	Location             *Location
}

func (self *NodeTemplate) NewCapability(name string, location *Location) *Capability {
	capability := &Capability{
		NodeTemplate:         self,
		Name:                 name,
		Types:                make(EntityTypes),
		Properties:           make(Values),
		Attributes:           make(Values),
		MaxRelationshipCount: math.MaxUint64,
		Location:             location,
	}
	self.Capabilities[name] = capability
	return capability
}

type MarshalableCapability struct {
	Description          string      `json:"description" yaml:"description"`
	Types                EntityTypes `json:"types" yaml:"types"`
	Properties           Values      `json:"properties" yaml:"properties"`
	Attributes           Values      `json:"attributes" yaml:"attributes"`
	MinRelationshipCount int64       `json:"minRelationshipCount" yaml:"minRelationshipCount"`
	MaxRelationshipCount int64       `json:"maxRelationshipCount" yaml:"maxRelationshipCount"`
	Location             *Location   `json:"location" yaml:"location"`
}

func (self *Capability) Marshalable() any {
	var minRelationshipCount int64
	var maxRelationshipCount int64
	minRelationshipCount = int64(self.MinRelationshipCount)
	if self.MaxRelationshipCount == math.MaxUint64 {
		maxRelationshipCount = -1
	} else {
		maxRelationshipCount = int64(self.MaxRelationshipCount)
	}

	return &MarshalableCapability{
		Description:          self.Description,
		Types:                self.Types,
		Properties:           self.Properties,
		Attributes:           self.Attributes,
		MinRelationshipCount: minRelationshipCount,
		MaxRelationshipCount: maxRelationshipCount,
		Location:             self.Location,
	}
}

// ([json.Marshaler] interface)
func (self *Capability) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// ([yaml.Marshaler] interface)
func (self *Capability) MarshalYAML() (any, error) {
	return self.Marshalable(), nil
}

// ([cbor.Marshaler] interface)
func (self *Capability) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(self.Marshalable())
}

// ([msgpack.Marshaler] interface)
func (self *Capability) MarshalMsgpack() ([]byte, error) {
	return ard.MarshalMessagePack(self.Marshalable())
}

// ([ard.ToARD] interface)
func (self *Capability) ToARD(reflector *ard.Reflector) (any, error) {
	return reflector.Unpack(self.Marshalable())
}

//
// Capabilities
//

type Capabilities map[string]*Capability
