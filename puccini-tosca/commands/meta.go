package commands

import (
	"io"

	"github.com/spf13/cobra"
	formatpkg "github.com/tliron/kutil/format"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca/csar"

	urlpkg "github.com/tliron/kutil/url"
)

func init() {
	rootCommand.AddCommand(metaCommand)
}

var metaCommand = &cobra.Command{
	Use:   "meta [CSAR PATH or URL]",
	Short: "Show CSAR metadata",
	Long:  `Parses and validates a CSAR an extracts the metadata.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url = args[0]

		Meta(url)
	},
}

func Meta(url string) {
	var err error

	urlContext := urlpkg.NewContext()
	util.OnExit(func() {
		if err := urlContext.Release(); err != nil {
			log.Errorf("%s", err.Error())
		}
	})

	var csarUrl urlpkg.URL
	csarUrl, err = urlpkg.NewValidURL(url, nil, urlContext)
	util.FailOnError(err)

	var zipUrl *urlpkg.ZipURL
	zipUrl, err = urlpkg.NewValidZipURL("TOSCA-Metadata/TOSCA.meta", csarUrl)
	util.FailOnError(err)

	var zipReader io.ReadCloser
	zipReader, err = zipUrl.Open()
	util.FailOnError(err)
	defer zipReader.Close()

	var meta *csar.Meta
	meta, err = csar.ReadMeta(zipReader)
	util.FailOnError(err)

	if !terminal.Quiet || (output != "") {
		err = formatpkg.WriteOrPrint(meta, format, terminal.Stdout, strict, pretty, output)
		util.FailOnError(err)
	}
}
