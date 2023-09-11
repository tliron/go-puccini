package clout

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tliron/go-ard"
)

const Version = "1.0"

//
// Clout
//

type Clout struct {
	Version    string        `json:"version" yaml:"version" ard:"version"`
	Metadata   ard.StringMap `json:"metadata" yaml:"metadata" ard:"metadata"`
	Properties ard.StringMap `json:"properties" yaml:"properties" ard:"properties"`
	Vertexes   Vertexes      `json:"vertexes" yaml:"vertexes" ard:"vertexes"`
}

func NewClout() *Clout {
	return &Clout{
		Version:    Version,
		Metadata:   make(ard.StringMap),
		Properties: make(ard.StringMap),
		Vertexes:   make(Vertexes),
	}
}

type MarshalableCloutStringMaps Clout

func (self *Clout) MarshalableStringMaps() any {
	return &MarshalableCloutStringMaps{
		Version:    self.Version,
		Metadata:   ard.CopyMapsToStringMaps(self.Metadata).(ard.StringMap),
		Properties: ard.CopyMapsToStringMaps(self.Properties).(ard.StringMap),
		Vertexes:   self.Vertexes,
	}
}

// ([json.Marshaler] interface)
func (self *Clout) MarshalJSON() ([]byte, error) {
	// JavaScript requires keys to be strings, so we would lose complex keys
	return json.Marshal(self.MarshalableStringMaps())
}

func (self *Clout) Resolve() error {
	if self.Version == "" {
		return errors.New("no Clout \"Version\"")
	}
	if self.Version != Version {
		return fmt.Errorf("unsupported Clout version: %q", self.Version)
	}

	return self.ResolveEdges()
}

func (self *Clout) ResolveEdges() error {
	for id, vertex := range self.Vertexes {
		vertex.Clout = self
		vertex.ID = id

		for _, edge := range vertex.EdgesOut {
			if edge.Target == nil {
				var ok bool
				if edge.Target, ok = self.Vertexes[edge.TargetID]; !ok {
					return fmt.Errorf("could not resolve Clout, bad TargetID: %q", edge.TargetID)
				}

				edge.Source = vertex
				edge.Target.EdgesIn = append(edge.Target.EdgesIn, edge)
			}
		}
	}
	return nil
}

// Creates a copy of the Clout in which ARD is used for all data.
func (self *Clout) Copy() (*Clout, error) {
	return self.copy(true)
}

// Creates a copy of the Clout in which non-ARD data is left as is.
func (self *Clout) CopyAsIs() *Clout {
	clout, _ := self.copy(false)
	return clout
}

func (self *Clout) copy(toArd bool) (*Clout, error) {
	clout := Clout{
		Version: Version,
	}

	var err error
	if clout.Vertexes, err = self.Vertexes.copy(toArd); err == nil {
		if toArd {
			if metadata, err := ard.ValidCopyMapsToStringMaps(self.Metadata, nil); err == nil {
				if properties, err := ard.ValidCopyMapsToStringMaps(self.Properties, nil); err == nil {
					clout.Metadata = metadata.(ard.StringMap)
					clout.Properties = properties.(ard.StringMap)
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			clout.Metadata = ard.CopyMapsToStringMaps(self.Metadata).(ard.StringMap)
			clout.Properties = ard.CopyMapsToStringMaps(self.Properties).(ard.StringMap)
		}
	} else {
		return nil, err
	}

	if err := clout.ResolveEdges(); err == nil {
		return &clout, nil
	} else {
		return nil, err
	}
}
