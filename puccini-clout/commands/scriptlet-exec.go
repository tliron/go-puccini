package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	urlpkg "github.com/tliron/kutil/url"
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

		clout, err := cloutpkg.Load(url, inputFormat)
		util.FailOnError(err)

		// Try loading JavaScript from Clout
		scriptlet, err := js.GetScriptlet(scriptletName, clout)

		urlContext := urlpkg.NewContext()
		defer urlContext.Release()

		if err != nil {
			// Try loading JavaScript from path or URL
			scriptletUrl, err := urlpkg.NewValidURL(scriptletName, nil, urlContext)
			util.FailOnError(err)

			scriptlet, err = urlpkg.ReadString(scriptletUrl)
			util.FailOnError(err)

			err = js.SetScriptlet(scriptletName, js.CleanupScriptlet(scriptlet), clout)
			util.FailOnError(err)
		}

		err = Exec(scriptletName, scriptlet, clout, urlContext)
		util.FailOnError(err)
	},
}

func Exec(scriptletName string, scriptlet string, clout *cloutpkg.Clout, urlContext *urlpkg.Context) error {
	jsContext := js.NewContext(scriptletName, log, arguments, terminal.Quiet, format, strict, timestamps, pretty, output, urlContext)
	err := jsContext.Require(clout, scriptletName, nil)
	return js.UnwrapException(err)
}
