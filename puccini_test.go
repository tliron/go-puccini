package main

import (
	contextpkg "context"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/tliron/exturl"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parser"

	_ "github.com/tliron/commonlog/simple"
)

func TestParse(t *testing.T) {
	context := NewContext(t)
	defer context.urlContext.Release()

	context.compileAll()
}

func BenchmarkParse(b *testing.B) {
	context := NewContext(b)
	defer context.urlContext.Release()

	for i := 0; i < b.N; i++ {
		context.compileAll()
	}
}

//
// Context
//

type Context struct {
	tb         testing.TB
	root       string
	urlContext *exturl.Context
	parser     *parser.Parser
}

func NewContext(tb testing.TB) *Context {
	var root string
	var ok bool
	if root, ok = os.LookupEnv("PUCCINI_TEST_ROOT"); !ok {
		var err error
		if root, err = os.Getwd(); err != nil {
			tb.Errorf("%s", err.Error())
		}
	}

	return &Context{
		tb:         tb,
		root:       root,
		urlContext: exturl.NewContext(),
		parser:     parser.NewParser(),
	}
}

func (self *Context) compileAll() {
	self.compile("tosca/artifacts.yaml", nil)
	self.compile("tosca/attributes.yaml", nil)
	self.compile("tosca/copy.yaml", nil)
	self.compile("tosca/data-types.yaml", nil)
	self.compile("tosca/descriptions.yaml", nil)
	self.compile("tosca/dsl-definitions.yaml", nil)
	self.compile("tosca/functions.yaml", nil)
	self.compile("tosca/inputs-and-outputs.yaml", map[string]any{"ram": "1gib"})
	self.compile("tosca/interfaces.yaml", nil)
	self.compile("tosca/metadata.yaml", nil)
	self.compile("tosca/namespaces.yaml", nil)
	self.compile("tosca/policies-and-groups.yaml", nil)
	self.compile("tosca/requirements-and-capabilities.yaml", nil)
	self.compile("tosca/simple-for-nfv.yaml", nil)
	self.compile("tosca/source-and-target.yaml", nil)
	self.compile("tosca/substitution-mapping-client.yaml", nil)
	self.compile("tosca/substitution-mapping.yaml", nil)
	self.compile("tosca/unicode.yaml", nil)
	self.compile("tosca/workflows.yaml", nil)
	self.compile("tosca/legacy/tosca_1_0.yaml", nil)
	self.compile("tosca/legacy/tosca_1_1.yaml", nil)
	self.compile("tosca/legacy/tosca_1_2.yaml", nil)
	self.compile("tosca/future/tosca_2_0.yaml", nil)
	self.compile("javascript/artifacts.yaml", nil)
	self.compile("javascript/constraints.yaml", nil)
	self.compile("javascript/converters.yaml", nil)
	self.compile("javascript/define.yaml", nil)
	self.compile("javascript/exec.yaml", nil)
	self.compile("javascript/functions.yaml", nil)
	self.compile("openstack/hello-world.yaml", nil)
	self.compile("bpmn/open-loop.yaml", nil)
	self.compile("cloudify/advanced-blueprint-example.yaml", map[string]any{
		"host_ip":                "1.2.3.4",
		"agent_user":             "my_user",
		"agent_private_key_path": "my_key",
	})
	self.compile("cloudify/example.yaml", nil)
	self.compile("hot/hello-world.yaml", map[string]any{
		"username": "test",
	})
}

func (self *Context) compile(url string, inputs map[string]any) {
	if t, ok := self.tb.(*testing.T); ok {
		t.Run(url, func(t_ *testing.T) {
			// Running the tests in parallel is not just for speed;
			// it actually helps us to find concurrency bugs
			t_.Parallel()
			self.compile_(t_, url, inputs)
		})
	} else {
		self.compile_(self.tb, url, inputs)
	}
}

func (self *Context) compile_(t testing.TB, url string, inputs map[string]any) {
	var normalServiceTemplate *normal.ServiceTemplate
	var clout *cloutpkg.Clout
	var err error

	url_ := self.urlContext.NewFileURL(path.Join(filepath.ToSlash(self.root), "examples", url))

	parserContext := self.parser.NewContext()
	parserContext.URL = url_
	parserContext.Inputs = inputs
	if normalServiceTemplate, err = parserContext.Parse(contextpkg.TODO()); err != nil {
		t.Errorf("%s\n%s", err.Error(), parserContext.GetProblems().ToString(true))
		return
	}

	problems := parserContext.GetProblems()
	if clout, err = normalServiceTemplate.Compile(); err != nil {
		t.Errorf("%s\n%s", err.Error(), problems.ToString(true))
		return
	}

	execContext := js.ExecContext{
		Clout:      clout,
		Problems:   problems,
		URLContext: self.urlContext,
		History:    true,
		Format:     "yaml",
		Pretty:     true,
	}

	execContext.Resolve()
	if !problems.Empty() {
		t.Errorf("%s", problems.ToString(true))
		return
	}

	execContext.Coerce()
	if !problems.Empty() {
		t.Errorf("%s", problems.ToString(true))
		return
	}
}
