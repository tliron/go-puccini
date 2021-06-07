package js

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	formatpkg "github.com/tliron/kutil/format"
	"github.com/tliron/kutil/js"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/kutil/util"
)

//
// PucciniAPI
//

type PucciniAPI struct {
	js.UtilAPI
	js.FormatAPI
	js.FileAPI

	Arguments       map[string]string
	Log             logging.Logger
	Stdout          io.Writer
	Stderr          io.Writer
	Stdin           io.Writer
	Stylist         *terminal.Stylist
	Output          string
	Format          string
	Strict          bool
	AllowTimestamps bool // TODO
	Pretty          bool

	context *Context
}

func (self *Context) NewPucciniAPI() *PucciniAPI {
	format := self.Format
	if format == "" {
		format = "yaml"
	}
	return &PucciniAPI{
		FileAPI:         js.NewFileAPI(self.URLContext),
		Arguments:       self.Arguments,
		Log:             self.Log,
		Stdout:          self.Stdout,
		Stderr:          self.Stderr,
		Stdin:           self.Stdin,
		Stylist:         self.Stylist,
		Output:          self.Output,
		Format:          format,
		Strict:          self.Strict,
		AllowTimestamps: self.AllowTimestamps,
		Pretty:          self.Pretty,
		context:         self,
	}
}

func (self *PucciniAPI) Write(data interface{}, path string, dontOverwrite bool) {
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
			fmt.Fprintf(self.Stdout, "%s %s\n", message, output)
		}
		if skip {
			return
		}
	}

	self.failOnError(formatpkg.WriteOrPrint(data, self.Format, self.Stdout, self.Strict, self.Pretty, output))
}

func (self *PucciniAPI) LoadString(url string) (string, error) {
	if url_, err := urlpkg.NewValidURL(url, nil, self.context.URLContext); err == nil {
		return urlpkg.ReadString(url_)
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
		fmt.Fprintln(self.Stderr, stylist.Error(message))
	}
	util.Exit(1)
}

func (self *PucciniAPI) Failf(format string, args ...interface{}) {
	self.Fail(fmt.Sprintf(format, args...))
}

func (self *PucciniAPI) failOnError(err error) {
	if err != nil {
		self.Fail(err.Error())
	}
}
