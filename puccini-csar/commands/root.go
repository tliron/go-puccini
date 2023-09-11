package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

var logTo string
var verbose int
var cpuProfilePath string

func init() {
	rootCommand.PersistentFlags().BoolVarP(&terminal.Quiet, "quiet", "q", false, "suppress output")
	rootCommand.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	rootCommand.PersistentFlags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	rootCommand.PersistentFlags().BoolVarP(&commonlog.Trace, "trace", "", false, "add stack trace to log messages")
	rootCommand.PersistentFlags().StringVarP(&cpuProfilePath, "cpu-profile", "", "", "CPU profile file path")
}

var rootCommand = &cobra.Command{
	Use:   toolName,
	Short: "CSAR packaging tool",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		util.InitializeCPUProfiling(cpuProfilePath)
		util.InitializeColorization(colorize)
		commonlog.Initialize(verbose, logTo)
	},
}

func Execute() {
	util.FailOnError(rootCommand.Execute())
}
