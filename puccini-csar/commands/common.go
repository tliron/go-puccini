package commands

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/go-transcribe"
	"github.com/tliron/kutil/util"
)

const toolName = "puccini-csar"

var log = commonlog.GetLogger(toolName)

var archiveFormat string

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
