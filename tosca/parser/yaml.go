package parser

import (
	"fmt"

	"github.com/tliron/kutil/terminal"
	"github.com/tliron/yamlkeys"
)

//
// YAMLDecodeError
//

type YAMLDecodeError struct {
	DecodeError *yamlkeys.DecodeError
}

func NewYAMLDecodeError(decodeError *yamlkeys.DecodeError) *YAMLDecodeError {
	return &YAMLDecodeError{decodeError}
}

// (error interface)
func (self *YAMLDecodeError) Error() string {
	return self.DecodeError.Error()
}

// problems.Problematic interface
func (self *YAMLDecodeError) Problem(stylist terminal.Stylist) (string, string, string, int, int) {
	return "", "", fmt.Sprintf("malformed YAML, %s", self.DecodeError.Message), self.DecodeError.Line, self.DecodeError.Column
}
