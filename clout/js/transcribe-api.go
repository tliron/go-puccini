package js

import (
	"io"
	"os"
	"path/filepath"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/commonjs-goja/api"
	"github.com/tliron/go-transcribe"
	"github.com/tliron/kutil/terminal"
)

// ([commonjs.CreateExtensionFunc] signature)
func (self *Environment) CreateTranscribeExtension(jsContext *commonjs.Context) any {
	return self.NewTranscribeAPI()
}

//
// TranscribeAPI
//

type TranscribeAPI struct {
	*api.Transcribe

	Stdout        io.Writer
	Stderr        io.Writer
	Stdin         io.Writer
	StdoutStylist *terminal.Stylist
	FilePath      string
	Format        string
	Strict        bool
	Pretty        bool
	Base64        bool

	context *Environment
}

func (self *Environment) NewTranscribeAPI() *TranscribeAPI {
	format := self.Format
	if format == "" {
		format = "yaml"
	}

	return &TranscribeAPI{
		Transcribe:    api.NewTranscribe(self.Stdout, self.Stderr),
		Stdin:         self.Stdin,
		StdoutStylist: self.StdoutStylist,
		FilePath:      self.FilePath,
		Format:        format,
		Strict:        self.Strict,
		Pretty:        self.Pretty,
		Base64:        self.Base64,
		context:       self,
	}
}

func (self *TranscribeAPI) Output(data any, path string, dontOverwrite bool) error {
	output := self.context.FilePath

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
