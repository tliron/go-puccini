package parser

import (
	contextpkg "context"
	"sync"

	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca/grammars"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Context
//

type Context struct {
	Parser *Parser

	URL     exturl.URL
	Bases   []exturl.URL
	Quirks  parsing.Quirks
	Inputs  map[string]ard.Value
	Stylist *terminal.Stylist

	Root  *File
	Files Files

	readWork  sync.WaitGroup
	filesLock util.RWLocker
}

func (self *Parser) NewContext() *Context {
	return &Context{
		Parser:    self,
		filesLock: util.NewDefaultRWLocker(),
	}
}

func (self *Context) GetProblems() *problems.Problems {
	return self.Root.GetContext().Problems
}

func (self *Context) MergeProblems() {
	self.filesLock.RLock()
	defer self.filesLock.RUnlock()

	// Note: This could happen many times, but because problems are de-duped, everything is OK :)
	for _, file := range self.Files {
		self.GetProblems().Merge(file.GetContext().Problems)
	}
}

func (self *Context) AddFile(file *File) {
	self.filesLock.Lock()
	defer self.filesLock.Unlock()

	self.Files = append(self.Files, file)
}

func (self *Context) AddImportFile(context contextpkg.Context, entityPtr parsing.EntityPtr, container *File, nameTransformer parsing.NameTransformer) *File {
	file := NewFile(entityPtr, container, nameTransformer)

	if container != nil {
		containerContext := container.GetContext()
		if !containerContext.HasQuirk(parsing.QuirkImportsVersionPermissive) {
			fileContext := file.GetContext()
			if !grammars.CompatibleGrammars(containerContext, fileContext) {
				containerContext.ReportImportIncompatible(fileContext.URL)
				return file
			}
		}
	}

	self.AddFile(file)

	self.goReadImports(context, file)

	return file
}

// Print

func (self *Context) PrintImports(indent int) {
	terminal.PrintIndent(indent)
	terminal.Printf("%s\n", terminal.StdoutStylist.Value(self.Root.GetContext().URL.String()))
	self.Root.PrintImports(indent, terminal.TreePrefix{})
}
