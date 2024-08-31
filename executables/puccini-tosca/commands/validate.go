package commands

import (
	contextpkg "context"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
)

func init() {
	rootCommand.AddCommand(validateCommand)
	validateCommand.Flags().StringSliceVarP(&importPaths, "path", "b", nil, "specify an import path or base URL")
	validateCommand.Flags().StringVarP(&template, "template", "t", "", "select service template in CSAR (leave empty for root, or use \"all\", path, or integer index)")
	validateCommand.Flags().StringToStringVarP(&inputs, "input", "i", nil, "specify input (format is name=value)")
	validateCommand.Flags().StringVarP(&inputsUrl, "inputs", "n", "", "load inputs from a PATH or URL to YAML content")
	validateCommand.Flags().StringVarP(&problemsFormat, "problems-format", "m", "", "problems format (\"yaml\", \"json\", \"xjson\", \"xml\", \"cbor\", \"messagepack\", or \"go\")")
	validateCommand.Flags().StringSliceVarP(&quirks, "quirk", "x", nil, "parser quirk")
	validateCommand.Flags().StringToStringVarP(&urlMappings, "map-url", "u", nil, "map a URL (format is from=to)")

	validateCommand.Flags().BoolVarP(&resolve, "resolve", "r", true, "resolves the topology (attempts to satisfy all requirements with capabilities)")
	validateCommand.Flags().BoolVarP(&coerce, "coerce", "c", false, "coerces all values (calls functions and applies constraints)")
}

var validateCommand = &cobra.Command{
	Use:   "validate [[TOSCA PATH or URL]]",
	Short: "Validate TOSCA",
	Long:  `Validates TOSCA service templates. Equivalent to "compile" without the Clout output.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) == 1 {
			url = args[0]
		}

		dumpPhases = nil
		Compile(contextpkg.TODO(), url)
		if !terminal.Quiet {
			terminal.Eprintln("valid")
		}
	},
}
