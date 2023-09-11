package clout

import (
	contextpkg "context"

	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
)

func Load(context contextpkg.Context, url exturl.URL, forceFormat string) (*Clout, error) {
	if reader, err := url.Open(context); err == nil {
		reader = util.NewContextualReadCloser(context, reader)
		defer commonlog.CallAndLogWarning(reader.Close, "clout.Load", log)

		format := forceFormat
		for format == "" {
			format = url.Format()
		}

		return Read(reader, format)
	} else {
		return nil, err
	}
}
