package normal

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// ValueMeta
//

type ValueMeta struct {
	// Element, Key+Value, and Fields cannot be set together

	Type             string            `json:"type,omitempty" yaml:"type,omitempty"`
	Element          *ValueMeta        `json:"element,omitempty" yaml:"element,omitempty"`
	Key              *ValueMeta        `json:"key,omitempty" yaml:"key,omitempty"`
	Value            *ValueMeta        `json:"value,omitempty" yaml:"value,omitempty"`
	Fields           ValueMetaMap      `json:"fields,omitempty" yaml:"fields,omitempty"`
	Validators       FunctionCalls     `json:"validators,omitempty" yaml:"validators,omitempty"`
	Converter        *FunctionCall     `json:"converter,omitempty" yaml:"converter,omitempty"`
	Description      string            `json:"description,omitempty" yaml:"description,omitempty"`
	TypeDescription  string            `json:"typeDescription,omitempty" yaml:"typeDescription,omitempty"`
	LocalDescription string            `json:"localDescription,omitempty" yaml:"localDescription,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	TypeMetadata     map[string]string `json:"typeMetadata,omitempty" yaml:"typeMetadata,omitempty"`
	LocalMetadata    map[string]string `json:"localMetadata,omitempty" yaml:"localMetadata,omitempty"`
}

func NewValueMeta() *ValueMeta {
	return &ValueMeta{
		Fields:        make(ValueMetaMap),
		Metadata:      make(map[string]string),
		TypeMetadata:  make(map[string]string),
		LocalMetadata: make(map[string]string),
	}
}

func (self *ValueMeta) AddValidator(validator *parsing.FunctionCall) {
	self.Validators = append(self.Validators, NewFunctionCall(validator))
}

func (self *ValueMeta) SetConverter(converter *parsing.FunctionCall) {
	self.Converter = NewFunctionCall(converter)
}

func CopyValueMeta(valueMeta *ValueMeta) *ValueMeta {
	if valueMeta != nil {
		self := ValueMeta{
			Type:             valueMeta.Type,
			Element:          CopyValueMeta(valueMeta.Element),
			Key:              CopyValueMeta(valueMeta.Key),
			Value:            CopyValueMeta(valueMeta.Value),
			Fields:           CopyValueMetaMap(valueMeta.Fields),
			Converter:        valueMeta.Converter,
			Description:      valueMeta.Description,
			TypeDescription:  valueMeta.TypeDescription,
			LocalDescription: valueMeta.LocalDescription,
			Metadata:         make(map[string]string),
			TypeMetadata:     make(map[string]string),
			LocalMetadata:    make(map[string]string),
		}

		self.Validators = append(valueMeta.Validators[:0:0], valueMeta.Validators...)

		for key, value := range valueMeta.Metadata {
			self.Metadata[key] = value
		}

		for key, value := range valueMeta.TypeMetadata {
			self.TypeMetadata[key] = value
		}

		for key, value := range valueMeta.LocalMetadata {
			self.LocalMetadata[key] = value
		}

		if !self.Empty() {
			return &self
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (self *ValueMeta) Empty() bool {
	return (self.Type == "") &&
		((self.Element == nil) || self.Element.Empty()) &&
		((self.Key == nil) || self.Key.Empty()) &&
		((self.Value == nil) || self.Value.Empty()) &&
		((self.Fields == nil) || self.Fields.Empty()) &&
		((self.Validators == nil) || (len(self.Validators) == 0)) &&
		(self.Converter == nil) &&
		(self.Description == "") &&
		(self.TypeDescription == "") &&
		(self.LocalDescription == "") &&
		((self.Metadata == nil) || (len(self.Metadata) == 0)) &&
		((self.TypeMetadata == nil) || (len(self.TypeMetadata) == 0)) &&
		((self.LocalMetadata == nil) || (len(self.LocalMetadata) == 0))
}

func (self *ValueMeta) Prune() {
	if self.Element != nil {
		if self.Element.Prune(); self.Element.Empty() {
			self.Element = nil
		}
	}
	if self.Key != nil {
		if self.Key.Prune(); self.Key.Empty() {
			self.Key = nil
		}
	}
	if self.Value != nil {
		if self.Value.Prune(); self.Value.Empty() {
			self.Value = nil
		}
	}
	if self.Fields != nil {
		if self.Fields.Prune(); self.Fields.Empty() {
			self.Fields = nil
		}
	}
	if len(self.Validators) == 0 {
		self.Validators = nil
	}
	if (self.Metadata != nil) && (len(self.Metadata) == 0) {
		self.Metadata = nil
	}
	if (self.TypeMetadata != nil) && (len(self.TypeMetadata) == 0) {
		self.TypeMetadata = nil
	}
	if (self.LocalMetadata != nil) && (len(self.LocalMetadata) == 0) {
		self.LocalMetadata = nil
	}
}

//
// ValueMetaMap
//

type ValueMetaMap map[string]*ValueMeta

func CopyValueMetaMap(valueMetaMap ValueMetaMap) ValueMetaMap {
	if (valueMetaMap == nil) || (len(valueMetaMap) == 0) {
		return nil
	}
	self := make(ValueMetaMap)
	for key, valueMeta := range valueMetaMap {
		self[key] = CopyValueMeta(valueMeta)
	}
	if !self.Empty() {
		return self
	} else {
		return nil
	}
}

func (self ValueMetaMap) Empty() bool {
	for _, valueMeta := range self {
		if !valueMeta.Empty() {
			return false
		}
	}
	return true
}

func (self ValueMetaMap) Prune() {
	for key, valueMeta := range self {
		if valueMeta.Prune(); valueMeta.Empty() {
			delete(self, key)
		}
	}
}
