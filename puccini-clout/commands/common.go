package commands

import (
	contextpkg "context"

	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/go-transcribe"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/clout"
	cloutpkg "github.com/tliron/puccini/clout"
)

const toolName = "puccini-clout"

var log = commonlog.GetLogger(toolName)

var output string

func Transcriber() *transcribe.Transcriber {
	return &transcribe.Transcriber{
		File:        output,
		Format:      format,
		ForTerminal: pretty,
		Strict:      strict,
		Base64:      base64,
		InPlace:     true,
	}
}

func Bases(urlContext *exturl.Context) []exturl.URL {
	workingDir, err := urlContext.NewWorkingDirFileURL()
	util.FailOnError(err)
	return []exturl.URL{workingDir}
}

func LoadClout(context contextpkg.Context, url string, urlContext *exturl.Context) *clout.Clout {
	var url_ exturl.URL
	var err error
	if url != "" {
		url_, err = urlContext.NewValidAnyOrFileURL(context, url, Bases(urlContext))
		util.FailOnError(err)
	} else {
		url_, err = urlContext.ReadToInternalURLFromStdin(context, format)
		util.FailOnError(err)
	}

	if format == "" {
		format = inputFormat
		if format == "" {
			format = url_.Format()
		}
	}

	clout, err := cloutpkg.Load(context, url_, inputFormat)
	util.FailOnError(err)
	return clout
}
