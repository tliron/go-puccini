package js

import (
	"fmt"

	"github.com/tliron/kutil/js"
	cloutpkg "github.com/tliron/puccini/clout"
)

//
// CloutContext
//

type CloutContext struct {
	Context   *Context
	Clout     *cloutpkg.Clout
	JSContext *js.Context
}

func (self *Context) NewCloutContext(clout *cloutpkg.Clout, jsContext *js.Context) *CloutContext {
	return &CloutContext{
		Context:   self,
		Clout:     clout,
		JSContext: jsContext,
	}
}

func (self *CloutContext) CallFunction(scriptletName string, functionName string, arguments []interface{}, functionCallContext FunctionCallContext) (interface{}, error) {
	if exports, err := self.JSContext.Environment.RequireID(scriptletName); err == nil {
		if function := exports.Get(functionName); function != nil {
			return CallFunction(self.JSContext.Environment.Runtime, function, functionCallContext, arguments)
		} else {
			return nil, fmt.Errorf("function not exported from module: %s", functionName)
		}
	} else {
		return nil, err
	}
}
