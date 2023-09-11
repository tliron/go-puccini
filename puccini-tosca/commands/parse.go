package commands

import (
	contextpkg "context"
	"sort"

	"github.com/spf13/cobra"
	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/normal"
	parserpkg "github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

var stopAtPhase uint32
var dumpPhases []uint
var filter string

func init() {
	rootCommand.AddCommand(parseCommand)
	parseCommand.Flags().StringSliceVarP(&importPaths, "path", "b", nil, "specify an import path or base URL")
	parseCommand.Flags().StringVarP(&template, "template", "t", "", "select service template in CSAR (leave empty for root, or use path or integer index)")
	parseCommand.Flags().StringToStringVarP(&inputs, "input", "i", nil, "specify an input (format is name=YAML)")
	parseCommand.Flags().StringVarP(&inputsUrl, "inputs", "n", "", "load inputs from a PATH or URL to YAML content")
	parseCommand.Flags().StringVarP(&problemsFormat, "problems-format", "m", "", "problems format (\"yaml\", \"json\", \"xjson\", \"xml\", \"cbor\", \"messagepack\", or \"go\")")
	parseCommand.Flags().StringSliceVarP(&quirks, "quirk", "x", nil, "parser quirk")
	parseCommand.Flags().StringToStringVarP(&urlMappings, "map-url", "u", nil, "map a URL (format is from=to)")

	parseCommand.Flags().Uint32VarP(&stopAtPhase, "stop", "s", 6, "parser phase at which to end")
	parseCommand.Flags().UintSliceVarP(&dumpPhases, "dump", "d", []uint{6}, "dump phase internals")
	parseCommand.Flags().StringVarP(&filter, "filter", "r", "", "filter output by entity path; use '*' for wildcard matching (disables --stop and --dump)")
}

var parseCommand = &cobra.Command{
	Use:   "parse [[TOSCA PATH or URL]]",
	Short: "Parse TOSCA",
	Long:  `Parses and validates a TOSCA service template and reports problems if there are any. Provides access to phase diagnostics and entities.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) == 1 {
			url = args[0]
		}

		if stopAtPhase > 6 {
			util.Failf("--stop cannot be > 6: %d", stopAtPhase)
		}

		for _, dumpPhase := range dumpPhases {
			if dumpPhase > 6 {
				util.Failf("--dump entry cannot be > 6: %d", dumpPhase)
			}
		}

		if filter != "" {
			stopAtPhase = 6
			dumpPhases = nil
		}

		Parse(contextpkg.TODO(), url)
	},
}

var parser = parserpkg.NewParser()

func Parse(context contextpkg.Context, url string) (*parserpkg.Context, *normal.ServiceTemplate) {
	urlContext := exturl.NewContext()
	util.OnExitError(urlContext.Release)

	ParseInputs(context, urlContext)

	// URL mappings
	for fromUrl, toUrl := range urlMappings {
		urlContext.Map(fromUrl, toUrl)
	}

	var url_ exturl.URL
	var err error
	if url == "" {
		log.Info("parsing stdin")
		url_, err = urlContext.ReadToInternalURLFromStdin(context, "yaml")
	} else {
		log.Infof("parsing %q", url)
		url_, err = urlContext.NewValidAnyOrFileURL(context, url, Bases(urlContext, false))
	}
	util.FailOnError(err)

	parserContext := parser.NewContext()
	parserContext.Quirks = parsing.NewQuirks(quirks...)
	parserContext.Stylist = terminal.StdoutStylist
	if problemsFormat != "" {
		parserContext.Stylist = terminal.NewStylist(false)
	}

	if stopAtPhase == 0 {
		return nil, nil
	}

	// Phase 1: Read
	ok := parserContext.ReadRoot(context, url_, Bases(urlContext, true), template)

	parserContext.MergeProblems()
	problems := parserContext.GetProblems()
	FailOnProblems(problems)

	if !ok {
		// Stop here if failed to read
		util.Exit(1)
	}

	if ToPrintPhase(1) {
		if len(dumpPhases) > 1 {
			terminal.Printf("%s\n", terminal.StdoutStylist.Heading("Imports"))
			parserContext.PrintImports(1)
		} else {
			parserContext.PrintImports(0)
		}
	}

	if stopAtPhase == 1 {
		return parserContext, nil
	}

	// Phase 2: Namespaces
	parserContext.AddNamespaces()
	parserContext.LookupNames()
	if ToPrintPhase(2) {
		if len(dumpPhases) > 1 {
			terminal.Printf("%s\n", terminal.StdoutStylist.Heading("Namespaces"))
			parserContext.PrintNamespaces(1)
		} else {
			parserContext.PrintNamespaces(0)
		}
	}

	if stopAtPhase == 2 {
		return parserContext, nil
	}

	// Phase 3: Hieararchies
	parserContext.AddHierarchies()
	if ToPrintPhase(3) {
		if len(dumpPhases) > 1 {
			terminal.Printf("%s\n", terminal.StdoutStylist.Heading("Hierarchies"))
			parserContext.PrintHierarchies(1)
		} else {
			parserContext.PrintHierarchies(0)
		}
	}

	if stopAtPhase == 3 {
		return parserContext, nil
	}

	// Phase 4: Inheritance
	if ToPrintPhase(4) {
		parserContext.Inherit(func(tasks parserpkg.Tasks) {
			if len(dumpPhases) > 1 {
				terminal.Printf("%s\n", terminal.StdoutStylist.Heading("Inheritance Tasks"))
				tasks.Print(1)
			} else {
				tasks.Print(0)
			}
		})
	} else {
		parserContext.Inherit(nil)
	}

	if stopAtPhase == 4 {
		return parserContext, nil
	}

	if parserContext.Root == nil {
		return parserContext, nil
	}

	parserContext.SetInputs(inputValues)

	// Phase 5: Rendering
	entityPtrs := parserContext.Render()
	if ToPrintPhase(5) {
		sort.Sort(entityPtrs)
		if len(dumpPhases) > 1 {
			terminal.Printf("%s\n", terminal.StdoutStylist.Heading("Rendering"))
		}
		for _, entityPtr := range entityPtrs {
			terminal.Printf("%s:\n", terminal.StdoutStylist.Path(parsing.GetContext(entityPtr).Path.String()))
			err = Transcriber().Write(entityPtr)
			util.FailOnError(err)
		}
	}

	if stopAtPhase == 5 {
		return parserContext, nil
	}

	if filter != "" {
		entityPtrs := parserContext.Gather(filter)
		if len(entityPtrs) == 0 {
			util.Failf("No paths found matching filter: %q\n", filter)
		} else if !terminal.Quiet {
			for _, entityPtr := range entityPtrs {
				terminal.Printf("%s\n", terminal.StdoutStylist.Path(parsing.GetContext(entityPtr).Path.String()))
				err = Transcriber().Write(entityPtr)
				util.FailOnError(err)
			}
		}
		return parserContext, nil
	}

	parserContext.MergeProblems()
	FailOnProblems(problems)

	// Phase 6: Normalization
	if serviceTemplate, ok := parserContext.Normalize(); ok {
		FailOnProblems(problems)
		if ToPrintPhase(6) {
			if len(dumpPhases) > 1 {
				terminal.Printf("%s\n", terminal.StdoutStylist.Heading("Normalization"))
			}
			Transcriber().Write(serviceTemplate)
		}
		return parserContext, serviceTemplate
	} else {
		util.Fail("grammar does not support normalization")
		return parserContext, nil
	}
}

func ToPrintPhase(phase uint) bool {
	if !terminal.Quiet {
		for _, phase_ := range dumpPhases {
			if phase_ == phase {
				return true
			}
		}
	}
	return false
}

func ParseInputs(context contextpkg.Context, urlContext *exturl.Context) {
	if inputsUrl != "" {
		log.Infof("load inputs from %q", inputsUrl)

		url, err := urlContext.NewValidAnyOrFileURL(context, inputsUrl, Bases(urlContext, false))
		util.FailOnError(err)
		reader, err := url.Open(context)
		util.FailOnError(err)
		reader = util.NewContextualReadCloser(context, reader)
		defer commonlog.CallAndLogWarning(reader.Close, "ParseInputs", log)
		data, err := yamlkeys.DecodeAll(reader)
		util.FailOnError(err)
		for _, data_ := range data {
			if map_, ok := data_.(ard.Map); ok {
				for key, value := range map_ {
					inputValues[yamlkeys.KeyString(key)] = value
				}
			} else {
				util.Failf("malformed inputs in %q", inputsUrl)
			}
		}
	}

	for name, input := range inputs {
		input_, _, err := ard.DecodeYAML(util.StringToBytes(input), false)
		util.FailOnError(err)
		inputValues[name] = input_
	}
}
