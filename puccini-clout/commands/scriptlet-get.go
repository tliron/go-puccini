package commands

import (
	contextpkg "context"

	"github.com/spf13/cobra"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/clout/js"
)

func init() {
	scriptletCommand.AddCommand(getCommand)
	getCommand.Flags().StringVarP(&output, "output", "o", "", "output to file (default is stdout)")
}

var getCommand = &cobra.Command{
	Use:   "get [NAME] [[Clout PATH or URL]]",
	Short: "Get JavaScript scriptlet from Clout",
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

		clout := LoadClout(contextpkg.TODO(), url, urlContext)

		scriptlet, err := js.GetScriptlet(scriptletName, clout)
		util.FailOnError(err)

		if !terminal.Quiet {
			err = Transcriber().Write(scriptlet)
			util.FailOnError(err)
		}
	},
}
