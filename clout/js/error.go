package js

import (
	"fmt"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/kutil/terminal"
)

//
// Error
//

type Error struct {
	FunctionCall *FunctionCall
	Arguments    []any
	Message      string
	Cause        error
}

func (self *FunctionCall) NewError(arguments []any, message string, cause error) *Error {
	return &Error{
		FunctionCall: self,
		Arguments:    arguments,
		Message:      message,
		Cause:        cause,
	}
}

func (self *FunctionCall) NewErrorf(arguments []any, format string, arg ...any) *Error {
	return self.NewError(arguments, fmt.Sprintf(format, arg...), nil)
}

func (self *FunctionCall) WrapError(arguments []any, err error) *Error {
	return self.NewError(arguments, "", commonjs.UnwrapJavaScriptException(err))
}

func (self *Error) Signature() string {
	return self.FunctionCall.Signature(self.Arguments)
}

// (error interface)
func (self *Error) Error() string {
	_, _, message, _, _ := self.Problem(nil)
	return message
}

// ([fmt.Stringer] interface)
func (self *Error) String() string {
	return self.Error()
}

// ([problems.Problematic] interface)
func (self *Error) Problem(stylist *terminal.Stylist) (string, string, string, int, int) {
	if stylist == nil {
		stylist = terminal.NewStylist(false)
	}

	message := fmt.Sprintf("%s: ", stylist.Path(self.FunctionCall.Path))

	if self.Message != "" {
		message += fmt.Sprintf("%s in call to %s", self.Message, stylist.Name(self.Signature()))
	} else {
		message += fmt.Sprintf("call to %s failed", stylist.Name(self.Signature()))
	}

	if self.Cause != nil {
		if jsError, ok := self.Cause.(*Error); ok {
			_, _, message_, _, _ := jsError.Problem(stylist)
			message += fmt.Sprintf(" because %s", stylist.Error(message_))
		} else {
			message += fmt.Sprintf(" because %s", stylist.Error(self.Cause.Error()))
		}
	}

	return self.FunctionCall.URL, "", message, self.FunctionCall.Row, self.FunctionCall.Column
}
