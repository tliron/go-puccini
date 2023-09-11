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

var arguments map[string]string

func init() {
	scriptletCommand.AddCommand(execCommand)
	execCommand.Flags().StringVarP(&output, "output", "o", "", "output to file or directory (default is stdout)")
	execCommand.Flags().StringToStringVarP(&arguments, "argument", "a", nil, "specify a scriptlet argument (format is key=value)")
}

var execCommand = &cobra.Command{
	Use:   "exec [NAME or JavaScript PATH or URL] [[Clout PATH or URL]]",
	Short: "Execute JavaScript scriptlet on Clout",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		scriptletName := args[0]

		var url string
		if len(args) == 2 {
			url = args[1]
		}

		urlContext := exturl.NewContext()
		defer urlContext.Release()
		context := contextpkg.TODO()

		clout := LoadClout(context, url, urlContext)

		// Try loading JavaScript from Clout
		scriptlet, err := js.GetScriptlet(scriptletName, clout)

		if err != nil {
			// Try loading JavaScript from path or URL
			scriptletUrl, err := urlContext.NewValidAnyOrFileURL(context, scriptletName, Bases(urlContext))
			util.FailOnError(err)

			scriptlet, err = exturl.ReadString(context, scriptletUrl)
			util.FailOnError(err)

			err = js.SetScriptlet(scriptletName, js.CleanupScriptlet(scriptlet), clout)
			util.FailOnError(err)
		}

		err = Exec(scriptletName, scriptlet, clout, urlContext)
		util.FailOnError(err)
	},
}

func Exec(scriptletName string, scriptlet string, clout *cloutpkg.Clout, urlContext *exturl.Context) error {
	environment := js.NewEnvironment(scriptletName, log, arguments, terminal.Quiet, format, strict, pretty, false, output, urlContext)
	_, err := environment.Require(clout, scriptletName, nil)
	return err
}
