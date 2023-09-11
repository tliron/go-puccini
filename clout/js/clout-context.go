package js

import (
	"github.com/tliron/commonjs-goja"
	cloutpkg "github.com/tliron/puccini/clout"
)

//
// CloutContext
//

type CloutContext struct {
	Context   *Environment
	Clout     *cloutpkg.Clout
	JSContext *commonjs.Context
}

func (self *Environment) NewCloutContext(clout *cloutpkg.Clout, jsContext *commonjs.Context) *CloutContext {
	return &CloutContext{
		Context:   self,
		Clout:     clout,
		JSContext: jsContext,
	}
}
