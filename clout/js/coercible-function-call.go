package js

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-transcribe"
)

//
// FunctionCall
//

type FunctionCall struct {
	Name       string        `json:"name" yaml:"name"`
	Arguments  []Coercible   `json:"arguments" yaml:"arguments"`
	Path       string        `json:"path,omitempty" yaml:"path,omitempty"`
	URL        string        `json:"url,omitempty" yaml:"url,omitempty"`
	Row        int           `json:"row" yaml:"row"`
	Column     int           `json:"column" yaml:"column"`
	Validators Validators    `json:"validators,omitempty" yaml:"validators,omitempty"`
	Converter  *FunctionCall `json:"converter,omitempty" yaml:"converter,omitempty"`

	Notation         ard.StringMap     `json:"-" yaml:"-"`
	ExecutionContext *ExecutionContext `json:"-" yaml:"-"`
}

func (self *ExecutionContext) NewFunctionCall(map_ ard.StringMap, notation ard.StringMap, meta ard.StringMap) (*FunctionCall, error) {
	functionCall := FunctionCall{
		Notation:         notation,
		ExecutionContext: self,
	}

	if data, ok := map_["name"]; ok {
		if functionCall.Name, ok = data.(string); !ok {
			return nil, fmt.Errorf("malformed function call, \"name\" not a string: %T", data)
		}
	} else {
		return nil, fmt.Errorf("malformed function call, no \"name\": %+v", map_)
	}

	if data, ok := map_["arguments"]; ok {
		if originalArguments, ok := data.(ard.List); ok {
			functionCall.Arguments = make([]Coercible, len(originalArguments))
			for index, argument := range originalArguments {
				var err error
				if functionCall.Arguments[index], err = self.NewCoercible(argument, nil); err != nil {
					return nil, err
				}
			}
		} else {
			if data != nil {
				return nil, fmt.Errorf("malformed function call, \"arguments\" not a list: %T", data)
			}
		}
	} else {
		return nil, fmt.Errorf("malformed function call, no \"arguments\": %+v", map_)
	}

	if data, ok := map_["path"]; ok {
		if functionCall.Path, ok = data.(string); !ok {
			return nil, fmt.Errorf("malformed function call, \"path\" not a string: %T", data)
		}
	}

	if data, ok := map_["url"]; ok {
		if functionCall.URL, ok = data.(string); !ok {
			return nil, fmt.Errorf("malformed function call, \"url\" not a string: %T", data)
		}
	}

	if data, ok := map_["row"]; ok {
		if functionCall.Row, ok = asInt(data); !ok {
			return nil, fmt.Errorf("malformed function call, \"row\" not an integer: %T", data)
		}
	}

	if data, ok := map_["column"]; ok {
		if functionCall.Column, ok = asInt(data); !ok {
			return nil, fmt.Errorf("malformed function call, \"column\" not an integer: %T", data)
		}
	}

	var err error
	if functionCall.Validators, err = self.NewValidatorsFromMeta(meta); err != nil {
		return nil, err
	}
	if functionCall.Converter, err = self.NewConverter(meta); err != nil {
		return nil, err
	}

	return &functionCall, nil
}

func (self *ExecutionContext) NewConverter(meta ard.StringMap) (*FunctionCall, error) {
	if converter, ok := meta["converter"]; ok {
		if value, err := self.NewCoercible(converter, nil); err == nil {
			if converter_, ok := value.(*FunctionCall); ok {
				return converter_, nil
			} else {
				return nil, fmt.Errorf("malformed converter, not a function call: %+v", converter)
			}
		} else {
			return nil, err
		}
	} else {
		return nil, nil
	}
}

func (self *FunctionCall) Signature(arguments []ard.Value) string {
	s := make([]string, len(arguments))
	for index, argument := range arguments {
		s[index] = encodeArgument(argument)
	}
	return fmt.Sprintf("%s(%s)", self.Name, strings.Join(s, ","))
}

// ([Coercible] interface)
func (self *FunctionCall) Coerce() (ard.Value, error) {
	arguments, err := self.CoerceArguments()
	if err != nil {
		return nil, err
	}

	logEvaluate.Debugf("%s %s", self.Path, self.Signature(arguments))

	data, err := self.ExecutionContext.Call(self.Name, "evaluate", arguments...)
	if err != nil {
		return nil, self.WrapError(arguments, err)
	}

	// TODO: Coerce result?

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
func (self *FunctionCall) AddValidators(validators Validators) {
	self.Validators = append(self.Validators, validators...)
}

// ([Coercible] interface)
func (self *FunctionCall) Unwrap() ard.Value {
	return self.Notation
}

func (self *FunctionCall) CoerceArguments() ([]ard.Value, error) {
	arguments := make([]ard.Value, len(self.Arguments))
	for index, argument := range self.Arguments {
		var err error
		if arguments[index], err = argument.Coerce(); err != nil {
			return nil, err
		}
	}
	return arguments, nil
}

func (self *FunctionCall) Validate(value ard.Value, errorWhenInvalid bool) (bool, error) {
	arguments, err := self.CoerceArguments()
	if err != nil {
		return false, err
	}

	// Prepend value to be first argument
	arguments = append([]ard.Value{value}, arguments...)

	logValidate.Debugf("%s %s", self.Path, self.Signature(arguments))

	r, err := self.ExecutionContext.Call(self.Name, "validate", arguments...)
	if err != nil {
		return false, self.WrapError(arguments, err)
	}

	switch valid := r.(type) {
	case bool:
		if valid {
			return true, nil
		}
	case string:
		return false, self.NewError(arguments, valid, nil)
	default:
		return false, self.WrapError(arguments, errors.New("\"validate\" must return a bool or a string"))
	}

	if errorWhenInvalid {
		return false, self.NewError(arguments, "returned false", nil)
	} else {
		return false, nil
	}
}

func (self *FunctionCall) Convert(value ard.Value) (ard.Value, error) {
	arguments := []ard.Value{value}

	logConvert.Debugf("%s %s", self.Path, self.Signature(arguments))

	if r, err := self.ExecutionContext.Call(self.Name, "convert", arguments...); err == nil {
		return r, nil
	} else {
		return false, self.WrapError(arguments, err)
	}
}

// Utils

func encodeArgument(argument ard.Value) string {
	switch argument_ := argument.(type) {
	case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
		return fmt.Sprintf("%d", argument_)
	case float64, float32:
		return fmt.Sprintf("%g", argument_)
	case bool:
		return fmt.Sprintf("%t", argument_)
	case string:
		argument_ = strings.ReplaceAll(argument_, "\n", "¶")
		return fmt.Sprintf("%q", argument_)
	default:
		argument__, _ := transcribe.NewTranscriber().StringifyJSON(argument)
		return argument__
	}
}
