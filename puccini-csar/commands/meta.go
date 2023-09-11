package commands

import (
	contextpkg "context"

	"github.com/spf13/cobra"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca/csar"
)

var output string
var format string
var colorize string
var strict bool
var pretty bool
var base64 bool

func init() {
	rootCommand.AddCommand(metaCommand)

	metaCommand.Flags().StringVarP(&archiveFormat, "archive-format", "a", "", "force archive format (\"tar.gz\", \"tar\", or \"zip\"); leave empty to determine automatically from extension")
	metaCommand.Flags().StringVarP(&output, "output", "o", "", "output metadata to file (leave empty for stdout)")
	metaCommand.Flags().StringVarP(&format, "format", "f", "", "force output format (\"yaml\", \"json\", \"xjson\", \"xml\", \"cbor\", \"messagepack\", or \"go\")")
	metaCommand.Flags().StringVarP(&colorize, "colorize", "z", "true", "colorize output (boolean or \"force\")")
	metaCommand.Flags().BoolVarP(&strict, "strict", "y", false, "strict output (for \"yaml\" format only)")
	metaCommand.Flags().BoolVarP(&pretty, "pretty", "p", true, "prettify output")
	metaCommand.Flags().BoolVarP(&base64, "base64", "", false, "output base64 (for \"cbor\", \"messagepack\" formats)")
}

var metaCommand = &cobra.Command{
	Use:   "meta [CSAR PATH or URL]",
	Short: "Show CSAR metadata",
	Long:  `Parses, validates, and extracts CSAR metadata.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url = args[0]

		Meta(url)
	},
}

func Meta(url string) {
	var err error

	urlContext := exturl.NewContext()
	util.OnExitError(urlContext.Release)
	context := contextpkg.TODO()

	var csarUrl exturl.URL
	csarUrl, err = urlContext.NewValidAnyOrFileURL(context, url, Bases(urlContext))
	util.FailOnError(err)

	var meta *csar.Meta
	meta, err = csar.ReadMetaFromURL(context, csarUrl, archiveFormat)
	util.FailOnError(err)

	if !terminal.Quiet || (output != "") {
		err = Transcriber().Write(meta)
		util.FailOnError(err)
	}
}
