package js

import (
	contextpkg "context"
	"io"
	"os"
	"sync"

	"github.com/dop251/goja"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/commonjs-goja/api"
	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/terminal"
	cloutpkg "github.com/tliron/puccini/clout"
)

//
// Environment
//

type Environment struct {
	Arguments     map[string]string
	Quiet         bool
	Format        string
	Strict        bool
	Pretty        bool
	Base64        bool
	Output        string
	Log           commonlog.Logger
	Stdout        io.Writer
	Stderr        io.Writer
	Stdin         io.Writer
	StdoutStylist *terminal.Stylist
	URLContext    *exturl.Context

	programCache sync.Map
}

func NewEnvironment(name string, log commonlog.Logger, arguments map[string]string, quiet bool, format string, strict bool, pretty bool, base64 bool, output string, urlContext *exturl.Context) *Environment {
	if arguments == nil {
		arguments = make(map[string]string)
	}

	return &Environment{
		Arguments:     arguments,
		Quiet:         quiet,
		Format:        format,
		Strict:        strict,
		Pretty:        pretty,
		Base64:        base64,
		Output:        output,
		Log:           commonlog.NewScopeLogger(log, name),
		Stdout:        os.Stdout,
		Stderr:        os.Stderr,
		Stdin:         os.Stdin,
		StdoutStylist: terminal.StdoutStylist,
		URLContext:    urlContext,
	}
}

func (self *Environment) Require(clout *cloutpkg.Clout, scriptletName string, apis map[string]any) (*goja.Object, error) {
	environment := self.NewJsEnvironment(clout, apis)
	if r, err := environment.Require(scriptletName); err == nil {
		return r, nil
	} else {
		return r, commonjs.UnwrapJavaScriptException(err)
	}
}

func (self *Environment) NewJsEnvironment(clout *cloutpkg.Clout, apis map[string]any) *commonjs.Environment {
	environment := commonjs.NewEnvironment(self.URLContext)

	environment.CreateResolver = func(url exturl.URL, jsContext *commonjs.Context) commonjs.ResolveFunc {
		// commonjs.ResolveFunc signature
		return func(context contextpkg.Context, id string, raw bool) (exturl.URL, error) {
			if scriptlet, err := GetScriptlet(id, clout); err == nil {
				url := self.URLContext.NewInternalURL(id)
				url.SetContent(scriptlet)
				return url, nil
			} else {
				return nil, err
			}
		}
	}

	environment.Extensions = append(
		api.CreateDefaultExtensions(false, self.Log, self.URLContext),
		commonjs.Extension{
			Name:   "puccini",
			Create: self.CreatePucciniExtension,
		},
		commonjs.Extension{
			Name:   "clout",
			Create: self.CreateCloutExtension(clout),
		},
	)

	environment.Extensions = append(environment.Extensions, commonjs.NewExtensions(apis)...)

	return environment
}
