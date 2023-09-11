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

	context.compile("tosca/artifacts.yaml", nil)
	context.compile("tosca/attributes.yaml", nil)
	context.compile("tosca/copy.yaml", nil)
	context.compile("tosca/data-types.yaml", nil)
	context.compile("tosca/descriptions.yaml", nil)
	context.compile("tosca/dsl-definitions.yaml", nil)
	context.compile("tosca/functions.yaml", nil)
	context.compile("tosca/inputs-and-outputs.yaml", map[string]any{"ram": "1gib"})
	context.compile("tosca/interfaces.yaml", nil)
	context.compile("tosca/metadata.yaml", nil)
	context.compile("tosca/namespaces.yaml", nil)
	context.compile("tosca/policies-and-groups.yaml", nil)
	context.compile("tosca/requirements-and-capabilities.yaml", nil)
	context.compile("tosca/simple-for-nfv.yaml", nil)
	context.compile("tosca/source-and-target.yaml", nil)
	context.compile("tosca/substitution-mapping-client.yaml", nil)
	context.compile("tosca/substitution-mapping.yaml", nil)
	context.compile("tosca/unicode.yaml", nil)
	context.compile("tosca/workflows.yaml", nil)
	context.compile("tosca/legacy/tosca_1_0.yaml", nil)
	context.compile("tosca/legacy/tosca_1_1.yaml", nil)
	context.compile("tosca/legacy/tosca_1_2.yaml", nil)
	context.compile("tosca/future/tosca_2_0.yaml", nil)
	context.compile("javascript/artifacts.yaml", nil)
	context.compile("javascript/constraints.yaml", nil)
	context.compile("javascript/converters.yaml", nil)
	context.compile("javascript/define.yaml", nil)
	context.compile("javascript/exec.yaml", nil)
	context.compile("javascript/functions.yaml", nil)
	context.compile("openstack/hello-world.yaml", nil)
	context.compile("bpmn/open-loop.yaml", nil)
	context.compile("cloudify/advanced-blueprint-example.yaml", map[string]any{
		"host_ip":                "1.2.3.4",
		"agent_user":             "my_user",
		"agent_private_key_path": "my_key",
	})
	context.compile("cloudify/example.yaml", nil)
	context.compile("hot/hello-world.yaml", map[string]any{
		"username": "test",
	})
}

//
// Context
//

type Context struct {
	t          *testing.T
	root       string
	urlContext *exturl.Context
	parser     *parser.Parser
}

func NewContext(t *testing.T) *Context {
	var root string
	var ok bool
	if root, ok = os.LookupEnv("PUCCINI_TEST_ROOT"); !ok {
		var err error
		if root, err = os.Getwd(); err != nil {
			t.Errorf("%s", err.Error())
		}
	}

	return &Context{
		t:          t,
		root:       root,
		urlContext: exturl.NewContext(),
		parser:     parser.NewParser(),
	}
}

func (self *Context) compile(url string, inputs map[string]any) {
	self.t.Run(url, func(t *testing.T) {
		// Running the tests in parallel is not just for speed;
		// it actually helps us to find concurrency bugs
		t.Parallel()

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
	})
}
