package commands

import (
	contextpkg "context"

	"github.com/spf13/cobra"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
)

var output string
var resolve bool
var coerce bool
var exec string
var arguments map[string]string

func init() {
	rootCommand.AddCommand(compileCommand)
	compileCommand.Flags().StringSliceVarP(&importPaths, "path", "b", nil, "specify an import path or base URL")
	compileCommand.Flags().StringVarP(&template, "template", "t", "", "select service template in CSAR (leave empty for root, or use \"all\", path, or integer index)")
	compileCommand.Flags().StringToStringVarP(&inputs, "input", "i", nil, "specify input (format is name=value)")
	compileCommand.Flags().StringVarP(&inputsUrl, "inputs", "n", "", "load inputs from a PATH or URL to YAML content")
	compileCommand.Flags().StringVarP(&problemsFormat, "problems-format", "m", "", "problems format (\"yaml\", \"json\", \"xjson\", \"xml\", \"cbor\", \"messagepack\", or \"go\")")
	compileCommand.Flags().StringSliceVarP(&quirks, "quirk", "x", nil, "parser quirk")
	compileCommand.Flags().StringToStringVarP(&urlMappings, "map-url", "u", nil, "map a URL (format is from=to)")

	compileCommand.Flags().StringVarP(&output, "output", "o", "", "output Clout to file (leave empty for stdout)")
	compileCommand.Flags().BoolVarP(&resolve, "resolve", "r", true, "resolves the topology (attempts to satisfy all requirements with capabilities)")
	compileCommand.Flags().BoolVarP(&coerce, "coerce", "c", false, "coerces all values (calls functions and applies constraints)")
	compileCommand.Flags().StringVarP(&exec, "exec", "e", "", "execute JavaScript scriptlet")
	compileCommand.Flags().StringToStringVarP(&arguments, "argument", "a", nil, "used with --exec to specify a scriptlet argument (format is key=value)")
}

var compileCommand = &cobra.Command{
	Use:   "compile [[TOSCA PATH or URL]]",
	Short: "Compile TOSCA to Clout",
	Long:  `Parses TOSCA service templates and compiles the normalized output to Clout.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) == 1 {
			url = args[0]
		}

		dumpPhases = nil
		Compile(contextpkg.TODO(), url)
	},
}

func Compile(context contextpkg.Context, url string) {
	// Parse
	serviceContext, serviceTemplate := Parse(context, url)
	problems := serviceContext.GetProblems()
	urlContext := serviceContext.Root.GetContext().URL.Context()

	// Compile
	clout, err := serviceTemplate.Compile()
	util.FailOnError(err)

	execContext := js.ExecContext{
		Clout:      clout,
		Problems:   problems,
		URLContext: urlContext,
		History:    true,
		Format:     format,
		Strict:     strict,
		Pretty:     pretty,
	}

	// Resolve
	if resolve {
		execContext.Resolve()
		FailOnProblems(problems)
	}

	// Coerce
	if coerce {
		execContext.Coerce()
		FailOnProblems(problems)
	}

	if exec != "" {
		err = Exec(context, exec, arguments, clout, urlContext)
		util.FailOnError(err)
	} else if !terminal.Quiet || (output != "") {
		err = Transcriber().Write(clout)
		util.FailOnError(err)
	}
}

func Exec(context contextpkg.Context, scriptletName string, arguments map[string]string, clout *cloutpkg.Clout, urlContext *exturl.Context) error {
	// Try loading JavaScript from Clout
	scriptlet, err := js.GetScriptlet(scriptletName, clout)

	if err != nil {
		// Try loading JavaScript from path or URL
		url, err := urlContext.NewValidAnyOrFileURL(context, scriptletName, Bases(urlContext, false))
		util.FailOnError(err)

		scriptlet, err = exturl.ReadString(context, url)
		util.FailOnError(err)

		err = js.SetScriptlet(exec, js.CleanupScriptlet(scriptlet), clout)
		util.FailOnError(err)
	}

	environment := js.NewEnvironment(scriptletName, log, arguments, terminal.Quiet, format, strict, pretty, false, output, urlContext)
	_, err = environment.Require(clout, scriptletName, nil)
	return err
}
