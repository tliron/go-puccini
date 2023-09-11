package clout

import (
	contextpkg "context"

	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
)

func Load(context contextpkg.Context, url exturl.URL) (*Clout, error) {
	if reader, err := url.Open(context); err == nil {
		reader = util.NewContextualReadCloser(context, reader)
		defer commonlog.CallAndLogWarning(reader.Close, "clout.Load", log)
		return Read(reader, url.Format())
	} else {
		return nil, err
	}
}
