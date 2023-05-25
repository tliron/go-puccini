package js

import (
	contextpkg "context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/js"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/transcribe"
	"github.com/tliron/kutil/util"
)

//
// PucciniAPI
//

type PucciniAPI struct {
	js.UtilAPI
	js.TranscribeAPI
	js.FileAPI

	Arguments map[string]string
	Log       commonlog.Logger
	Stdout    io.Writer
	Stderr    io.Writer
	Stdin     io.Writer
	Stylist   *terminal.Stylist
	Output    string
	Format    string
	Strict    bool
	Pretty    bool

	context *Context
}

func (self *Context) NewPucciniAPI() *PucciniAPI {
	format := self.Format
	if format == "" {
		format = "yaml"
	}
	return &PucciniAPI{
		FileAPI:   js.NewFileAPI(self.URLContext),
		Arguments: self.Arguments,
		Log:       self.Log,
		Stdout:    self.Stdout,
		Stderr:    self.Stderr,
		Stdin:     self.Stdin,
		Stylist:   self.Stylist,
		Output:    self.Output,
		Format:    format,
		Strict:    self.Strict,
		Pretty:    self.Pretty,
		context:   self,
	}
}

func (self *PucciniAPI) NowString() string {
	return self.Now().Format(time.RFC3339Nano)
}

func (self *PucciniAPI) Write(data any, path string, dontOverwrite bool) {
	output := self.context.Output
	if path != "" {
		// Our path is relative to output path
		// (output path is here considered to be a directory)
		output = filepath.Join(output, path)
		var err error
		output, err = filepath.Abs(output)
		self.failOnError(err)
	}

	if output == "" {
		if self.context.Quiet {
			return
		}
	} else {
		_, err := os.Stat(output)
		var message string
		var skip bool
		stylist := self.Stylist
		if stylist == nil {
			stylist = terminal.NewStylist(false)
		}
		if (err == nil) || os.IsExist(err) {
			if dontOverwrite {
				message = stylist.Error("skipping:   ")
				skip = true
			} else {
				message = stylist.Value("overwriting:")
			}
		} else {
			message = stylist.Heading("writing:    ")
		}
		if !self.context.Quiet {
			terminal.Printf("%s %s\n", message, output)
		}
		if skip {
			return
		}
	}

	self.failOnError(transcribe.WriteOrPrint(data, self.Format, self.Stdout, self.Strict, self.Pretty, output))
}

func (self *PucciniAPI) LoadString(url string) (string, error) {
	context := contextpkg.TODO()
	if url_, err := self.context.URLContext.NewValidURL(context, url, nil); err == nil {
		return exturl.ReadString(context, url_)
	} else {
		return "", err
	}
}

func (self *PucciniAPI) Fail(message string) {
	stylist := self.Stylist
	if stylist == nil {
		stylist = terminal.NewStylist(false)
	}
	if !self.context.Quiet {
		terminal.Eprintln(stylist.Error(message))
	}
	util.Exit(1)
}

func (self *PucciniAPI) Failf(format string, args ...any) {
	self.Fail(fmt.Sprintf(format, args...))
}

func (self *PucciniAPI) failOnError(err error) {
	if err != nil {
		self.Fail(err.Error())
	}
}
