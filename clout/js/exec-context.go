package js

import (
	"github.com/dop251/goja"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/exturl"
	problemspkg "github.com/tliron/kutil/problems"
	cloutpkg "github.com/tliron/puccini/clout"
)

//
// ExecContext
//

type ExecContext struct {
	Clout      *cloutpkg.Clout
	Problems   *problemspkg.Problems
	URLContext *exturl.Context
	History    bool
	Format     string
	Strict     bool
	Pretty     bool
	Base64     bool
}

func (self *ExecContext) NewEnvironment(scriptletName string, arguments map[string]string) *Environment {
	return NewEnvironment(scriptletName, log, arguments, true, self.Format, self.Strict, self.Pretty, self.Base64, "", self.URLContext)
}

func (self *ExecContext) Exec(scriptletName string, arguments map[string]string) *goja.Object {
	// commonjs.CreateExtensionFunc signature
	createProblemsExtension := func(jsContext *commonjs.Context) any {
		return self.Problems
	}

	context := self.NewEnvironment(scriptletName, arguments)
	if r, err := context.Require(self.Clout, scriptletName, map[string]commonjs.CreateExtensionFunc{"problems": createProblemsExtension}); err == nil {
		return r
	} else {
		self.Problems.ReportError(err)
		return nil
	}
}

func (self *ExecContext) ExecWithHistory(scriptletName string) *goja.Object {
	var arguments map[string]string
	if !self.History {
		arguments = make(map[string]string)
		arguments["history"] = "false"
	}
	return self.Exec(scriptletName, arguments)
}

func (self *ExecContext) Resolve() {
	self.ExecWithHistory("tosca.resolve")
}

func (self *ExecContext) Coerce() {
	self.ExecWithHistory("tosca.coerce")
}

func (self *ExecContext) Outputs() *goja.Object {
	return self.Exec("tosca.outputs", nil)
}
