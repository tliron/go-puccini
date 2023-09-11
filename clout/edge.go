package clout

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/go-ard"
	"github.com/vmihailenco/msgpack/v5"
)

//
// Edge
//

type Edge struct {
	Metadata   ard.StringMap
	Properties ard.StringMap
	Source     *Vertex
	TargetID   string
	Target     *Vertex
}

func (self *Vertex) NewEdgeTo(target *Vertex) *Edge {
	edge := &Edge{
		Metadata:   make(ard.StringMap),
		Properties: make(ard.StringMap),
		Source:     self,
		TargetID:   target.ID,
		Target:     target,
	}
	self.EdgesOut = append(self.EdgesOut, edge)
	target.EdgesIn = append(target.EdgesIn, edge)
	return edge
}

func (self *Vertex) NewEdgeToID(targetId string) *Edge {
	edge := &Edge{
		Metadata:   make(ard.StringMap),
		Properties: make(ard.StringMap),
		Source:     self,
		TargetID:   targetId,
	}
	self.EdgesOut = append(self.EdgesOut, edge)
	return edge
}

func (self *Edge) Remove() {
	self.Source.EdgesOut = self.Source.EdgesOut.remove(self)
	self.Target.EdgesIn = self.Target.EdgesIn.remove(self)
}

// Entity interface
func (self *Edge) GetMetadata() ard.StringMap {
	return self.Metadata
}

// Entity interface
func (self *Edge) GetProperties() ard.StringMap {
	return self.Properties
}

type MarshalableEdge struct {
	Metadata   ard.StringMap `json:"metadata" yaml:"metadata" ard:"metadata"`
	Properties ard.StringMap `json:"properties" yaml:"properties" ard:"properties"`
	TargetID   string        `json:"targetID" yaml:"targetID" ard:"targetID"`
}

type MarshalableEdgeStringMaps struct {
	Metadata   ard.StringMap `json:"metadata"`
	Properties ard.StringMap `json:"properties"`
	TargetID   string        `json:"targetID"`
}

func (self *Edge) Marshalable(stringMaps bool) any {
	var targetID string
	if self.Target != nil {
		targetID = self.Target.ID
	}

	if stringMaps {
		return &MarshalableEdgeStringMaps{
			Metadata:   ard.CopyMapsToStringMaps(self.Metadata).(ard.StringMap),
			Properties: ard.CopyMapsToStringMaps(self.Properties).(ard.StringMap),
			TargetID:   targetID,
		}
	} else {
		return &MarshalableEdge{
			Metadata:   self.Metadata,
			Properties: self.Properties,
			TargetID:   targetID,
		}
	}
}

func (self *Edge) Unmarshal(f func(m *MarshalableEdge) error) error {
	var m MarshalableEdge
	if err := f(&m); err != nil {
		return err
	}
	self.Metadata = m.Metadata
	self.Properties = m.Properties
	self.TargetID = m.TargetID
	return nil
}

// ([json.Marshaler] interface)
func (self *Edge) MarshalJSON() ([]byte, error) {
	// JavaScript requires keys to be strings, so we would lose complex keys
	return json.Marshal(self.Marshalable(true))
}

// ([yaml.Marshaler] interface)
func (self *Edge) MarshalYAML() (any, error) {
	return self.Marshalable(false), nil
}

// ([cbor.Marshaler] interface)
func (self *Edge) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(self.Marshalable(false))
}

// ([msgpack.Marshaler] interface)
func (self *Edge) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal(self.Marshalable(false))
}

// ([ard.ToARD] interface)
func (self *Edge) ToARD(reflector *ard.Reflector) (any, error) {
	return reflector.Unpack(self.Marshalable(false))
}

// json.Unmarshaler interface
func (self *Edge) UnmarshalJSON(data []byte) error {
	return self.Unmarshal(func(m *MarshalableEdge) error {
		return json.Unmarshal(data, m)
	})
}

// yaml.Unmarshaler interface
func (self *Edge) UnmarshalYAML(unmarshal func(any) error) error {
	return self.Unmarshal(func(m *MarshalableEdge) error {
		return unmarshal(m)
	})
}

// cbor.Unmarshaler interface
func (self *Edge) UnmarshalCBOR(data []byte) error {
	return self.Unmarshal(func(m *MarshalableEdge) error {
		return cbor.Unmarshal(data, m)
	})
}

// msgpack.Unmarshaler interface
func (self *Edge) UnmarshalMsgpack(data []byte) error {
	return self.Unmarshal(func(m *MarshalableEdge) error {
		return ard.UnmarshalMessagePack(data, m)
	})
}

func (self *Edge) copy(toArd bool) (*Edge, error) {
	edge := Edge{
		TargetID: self.TargetID,
	}

	if toArd {
		if metadata, err := ard.ValidCopyMapsToStringMaps(self.Metadata, nil); err == nil {
			if properties, err := ard.ValidCopyMapsToStringMaps(self.Properties, nil); err == nil {
				edge.Metadata = metadata.(ard.StringMap)
				edge.Properties = properties.(ard.StringMap)
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		edge.Metadata = ard.CopyMapsToStringMaps(self.Metadata).(ard.StringMap)
		edge.Properties = ard.CopyMapsToStringMaps(self.Properties).(ard.StringMap)
	}

	return &edge, nil
}

//
// Edges
//

type Edges []*Edge

// Note: ".length" will not work reliably in JavaScript because once the value
// is bound it will not reflect changes to the struct's field
func (self Edges) Size() int {
	return len(self)
}

func (self Edges) copy(toArd bool) (Edges, error) {
	edges := make(Edges, len(self))
	var err error
	for index, edge := range self {
		if edges[index], err = edge.copy(toArd); err != nil {
			return nil, err
		}
	}
	return edges, nil
}

func (self Edges) remove(edge *Edge) Edges {
	for index, e := range self {
		if e == edge {
			return append(self[:index], self[index+1:]...)
		}
	}
	return self
}
