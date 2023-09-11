package js

import (
	"github.com/tliron/go-ard"
)

//
// Value
//

type Value struct {
	Data       any           `json:"data" yaml:"data"` // List, Map, or ard.Value
	Validators Validators    `json:"validators,omitempty" yaml:"validators,omitempty"`
	Converter  *FunctionCall `json:"converter,omitempty" yaml:"converter,omitempty"`

	Notation ard.StringMap `json:"-" yaml:"-"`
}

func (self *ExecutionContext) NewValue(data ard.Value, notation ard.StringMap, meta ard.StringMap) (*Value, error) {
	value := Value{
		Data:     data,
		Notation: notation,
	}

	var err error
	if value.Validators, err = self.NewValidatorsFromMeta(meta); err != nil {
		return nil, err
	}
	if value.Converter, err = self.NewConverter(meta); err != nil {
		return nil, err
	}

	return &value, nil
}

// ([Coercible] interface)
func (self *Value) Coerce() (ard.Value, error) {
	data := self.Data

	var err error
	switch data_ := data.(type) {
	case List:
		if data, err = data_.Coerce(); err != nil {
			return nil, err
		}

	case Map:
		if data, err = data_.Coerce(); err != nil {
			return nil, err
		}
	}

	if err := self.Validators.Apply(data); err == nil {
		if self.Converter != nil {
			return self.Converter.Convert(data)
		} else {
			return data, nil
		}
	} else {
		return nil, err
	}
}

// ([Coercible] interface)
func (self *Value) AddValidators(validators Validators) {
	self.Validators = append(self.Validators, validators...)
}

// ([Coercible] interface)
func (self *Value) Unwrap() ard.Value {
	return self.Notation
}
