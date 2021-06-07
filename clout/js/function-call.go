package js

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/format"
)

//
// FunctionCallContext
//

type FunctionCallContext struct {
	Site   interface{}
	Source interface{}
	Target interface{}
}

//
// FunctionCall
//

type FunctionCall struct {
	CloutContext        *CloutContext       `json:"-" yaml:"-"`
	FunctionCallContext FunctionCallContext `json:"-" yaml:"-"`
	Notation            ard.StringMap       `json:"-" yaml:"-"`

	Name        string      `json:"name" yaml:"name"`
	Arguments   []Coercible `json:"arguments" yaml:"arguments"`
	Path        string      `json:"path,omitempty" yaml:"path,omitempty"`
	URL         string      `json:"url,omitempty" yaml:"url,omitempty"`
	Row         int         `json:"row" yaml:"row"`
	Column      int         `json:"column" yaml:"column"`
	Constraints Constraints `json:"constraints,omitempty" yaml:"constraints,omitempty"`
}

func (self *CloutContext) NewFunctionCall(map_ ard.StringMap, notation ard.StringMap, functionCallContext FunctionCallContext) (*FunctionCall, error) {
	functionCall := FunctionCall{
		CloutContext:        self,
		Notation:            notation,
		FunctionCallContext: functionCallContext,
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
				if functionCall.Arguments[index], err = self.NewCoercible(argument, functionCallContext); err != nil {
					return nil, err
				}
			}
		} else {
			return nil, fmt.Errorf("malformed function call, \"arguments\" not a list: %T", data)
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
	if functionCall.Constraints, err = self.NewConstraintsFromNotation(notation, "$constraints", functionCallContext); err != nil {
		return nil, err
	}

	return &functionCall, nil
}

func (self *FunctionCall) Signature(arguments []ard.Value) string {
	s := make([]string, len(arguments))
	for index, argument := range arguments {
		s[index] = encodeArgument(argument)
	}
	return fmt.Sprintf("%s(%s)", self.Name, strings.Join(s, ","))
}

// Coercible interface
func (self *FunctionCall) Coerce() (ard.Value, error) {
	arguments, err := self.CoerceArguments()
	if err != nil {
		return nil, err
	}

	logEvaluate.Debugf("%s %s", self.Path, self.Signature(arguments))

	r, err := self.CloutContext.CallFunction(self.Name, "evaluate", arguments, self.FunctionCallContext)
	if err != nil {
		return nil, self.WrapError(arguments, err)
	}

	// TODO: Coerce result?

	return self.Constraints.Apply(r)
}

// Coercible interface
func (self *FunctionCall) SetConstraints(constraints Constraints) {
	self.Constraints = constraints
}

// Coercible interface
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

	r, err := self.CloutContext.CallFunction(self.Name, "validate", arguments, self.FunctionCallContext)
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
		return false, self.NewError(arguments, "", nil)
	} else {
		return false, nil
	}
}

// Utils

func encodeArgument(argument ard.Value) string {
	var encodedArgument string
	switch argument.(type) {
	case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
		return fmt.Sprintf("%d", argument)
	case float64, float32:
		return fmt.Sprintf("%g", argument)
	case bool:
		return fmt.Sprintf("%t", argument)
	case ard.Map, ard.StringMap, ard.List:
		encodedArgument, _ = format.EncodeYAML(argument, "", false)
		encodedArgument = strings.TrimSuffix(encodedArgument, "\n")
	default:
		encodedArgument = fmt.Sprintf("%s", argument)
	}

	encodedArgument = strings.ReplaceAll(encodedArgument, "\n", "¶")
	encodedArgument = strings.ReplaceAll(encodedArgument, "\"", "\\\"")
	return fmt.Sprintf("%q", encodedArgument)
}

func asInt(value interface{}) (int, bool) {
	switch value_ := value.(type) {
	case int64:
		return int(value_), true
	case int:
		return value_, true
	}
	return 0, false
}
