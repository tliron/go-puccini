package commands

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/go-kutil/util"
	"github.com/tliron/go-transcribe"
)

const toolName = "puccini-csar"

var (
	log = commonlog.GetLogger(toolName)

	archiveFormat string
)

func Transcriber() *transcribe.Transcriber {
	return &transcribe.Transcriber{
		File:        output,
		Format:      format,
		ForTerminal: pretty,
		Strict:      strict,
		Base64:      base64,
	}
}

func Bases(urlContext *exturl.Context) []exturl.URL {
	workingDir, err := urlContext.NewWorkingDirFileURL()
	util.FailOnError(err)
	return []exturl.URL{workingDir}
}
