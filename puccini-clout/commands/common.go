package commands

import (
	contextpkg "context"

	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/clout"
	cloutpkg "github.com/tliron/puccini/clout"
)

const toolName = "puccini-clout"

var log = commonlog.GetLogger(toolName)

var output string

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

	clout, err := cloutpkg.Load(context, url_)
	util.FailOnError(err)
	return clout
}
