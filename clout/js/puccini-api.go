package js

import (
	"io"
	"os"
	"path/filepath"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/commonlog"
	"github.com/tliron/go-transcribe"
	"github.com/tliron/kutil/terminal"
)

// ([commonjs.CreateExtensionFunc] signature)
func (self *Environment) CreatePucciniExtension(jsContext *commonjs.Context) any {
	return self.NewPucciniAPI()
}

//
// PucciniAPI
//

type PucciniAPI struct {
	Arguments     map[string]string
	Log           commonlog.Logger
	Stdout        io.Writer
	Stderr        io.Writer
	Stdin         io.Writer
	StdoutStylist *terminal.Stylist
	Output        string
	Format        string
	Strict        bool
	Pretty        bool
	Base64        bool

	context *Environment
}

func (self *Environment) NewPucciniAPI() *PucciniAPI {
	format := self.Format
	if format == "" {
		format = "yaml"
	}
	return &PucciniAPI{
		Arguments:     self.Arguments,
		Log:           self.Log,
		Stdout:        self.Stdout,
		Stderr:        self.Stderr,
		Stdin:         self.Stdin,
		StdoutStylist: self.StdoutStylist,
		Output:        self.Output,
		Format:        format,
		Strict:        self.Strict,
		Pretty:        self.Pretty,
		Base64:        self.Base64,
		context:       self,
	}
}

func (self *PucciniAPI) Write(data any, path string, dontOverwrite bool) error {
	output := self.context.Output

	if path != "" {
		// Our path is relative to output path
		// (output path is here considered to be a directory)
		output = filepath.Join(output, path)
		var err error
		output, err = filepath.Abs(output)
		if err != nil {
			return err
		}
	}

	if output == "" {
		if self.context.Quiet {
			return nil
		}
	} else {
		stylist := self.StdoutStylist
		if stylist == nil {
			stylist = terminal.NewStylist(false)
		}

		var message string
		var skip bool
		_, err := os.Stat(output)
		if (err == nil) || os.IsExist(err) {
			// File exists
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
			return nil
		}
	}

	transcriber := transcribe.Transcriber{
		File:        output,
		Writer:      self.Stdout,
		Format:      self.Format,
		Strict:      self.Strict,
		ForTerminal: self.Pretty,
		Base64:      self.Base64,
	}

	return transcriber.Write(data)
}
