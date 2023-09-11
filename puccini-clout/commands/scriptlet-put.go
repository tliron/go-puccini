package commands

import (
	contextpkg "context"

	"github.com/spf13/cobra"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/clout/js"
)

func init() {
	scriptletCommand.AddCommand(putCommand)
	putCommand.Flags().StringVarP(&output, "output", "o", "", "output Clout to file (default is stdout)")
}

var putCommand = &cobra.Command{
	Use:   "put [NAME] [JavaScript PATH or URL] [[Clout PATH or URL]]",
	Short: "Put JavaScript scriptlet in Clout",
	Long:  ``,
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		scriptletName := args[0]
		scriptletUrl := args[1]

		var url string
		if len(args) == 3 {
			url = args[2]
		}

		urlContext := exturl.NewContext()
		defer urlContext.Release()
		context := contextpkg.TODO()

		clout := LoadClout(context, url, urlContext)

		scriptletUrl_, err := urlContext.NewValidAnyOrFileURL(context, scriptletUrl, Bases(urlContext))
		util.FailOnError(err)

		scriptlet, err := exturl.ReadString(context, scriptletUrl_)
		util.FailOnError(err)

		err = js.SetScriptlet(scriptletName, js.CleanupScriptlet(scriptlet), clout)
		util.FailOnError(err)

		err = Transcriber().Write(clout)
		util.FailOnError(err)
	},
}
