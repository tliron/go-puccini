package normal

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// FunctionCall
//

type FunctionCall struct {
	Key       Value      `json:"$key,omitempty" yaml:"$key,omitempty"`
	ValueMeta *ValueMeta `json:"$meta,omitempty" yaml:"$meta,omitempty"`

	FunctionCall *parsing.FunctionCall `json:"$functionCall" yaml:"$functionCall"`
}

func NewFunctionCall(functionCall *parsing.FunctionCall) *FunctionCall {
	return &FunctionCall{FunctionCall: functionCall}
}

// Value interface
func (self *FunctionCall) SetKey(key Value) {
	self.Key = key
}

// Value interface
func (self *FunctionCall) SetMeta(valueMeta *ValueMeta) {
	self.ValueMeta = CopyValueMeta(valueMeta)
}

//
// FunctionCalls
//

type FunctionCalls []*FunctionCall

//
// FunctionCallMap
//

type FunctionCallMap map[string]FunctionCalls

//
// FunctionCallMapMap
//

type FunctionCallMapMap map[string]FunctionCallMap
