package parser

import (
	contextpkg "context"
	"fmt"
	"sort"

	"github.com/tliron/exturl"
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca/csar"
	"github.com/tliron/puccini/tosca/grammars"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

func (self *Context) ReadRoot(context contextpkg.Context, url exturl.URL, bases []exturl.URL, serviceTemplateName string) bool {
	parsingContext := parsing.NewContext(self.Stylist, self.Quirks)
	parsingContext.Bases = bases

	parsingContext.URL = url

	var ok bool

	self.readWork.Add(1)
	self.Root, ok = self.read(context, nil, parsingContext, nil, nil, "$Root", serviceTemplateName)
	self.readWork.Wait()

	self.filesLock.Lock()
	sort.Sort(self.Files)
	self.filesLock.Unlock()

	return ok
}

func (self *Context) read(context contextpkg.Context, promise util.Promise, parsingContext *parsing.Context, container *File, nameTransformer parsing.NameTransformer, readerName string, serviceTemplateName string) (*File, bool) {
	defer self.readWork.Done()
	if promise != nil {
		// For the goroutines waiting for our cached entityPtr
		defer promise.Release()
	}

	logRead.Infof("%s: %s", readerName, parsingContext.URL.Key())

	// TODO: allow override of CSAR format
	if format := parsingContext.URL.Format(); csar.IsValidFormat(format) {
		var err error
		if parsingContext.URL, err = csar.GetServiceTemplateURL(context, parsingContext.URL, format, serviceTemplateName); err != nil {
			parsingContext.ReportError(err)
			file := NewEmptyFile(parsingContext, container, nameTransformer)
			self.AddFile(file)
			return file, false
		}
	}

	// Read ARD
	var err error
	if parsingContext.Data, parsingContext.Locator, err = parsingContext.Read(context); err != nil {
		if decodeError, ok := err.(*yamlkeys.DecodeError); ok {
			err = NewYAMLDecodeError(decodeError)
		}
		parsingContext.ReportError(err)
		file := NewEmptyFile(parsingContext, container, nameTransformer)
		self.AddFile(file)
		return file, false
	}

	// Detect grammar
	if !grammars.DetectGrammar(parsingContext) {
		file := NewEmptyFile(parsingContext, container, nameTransformer)
		self.AddFile(file)
		return file, false
	}

	// Read entityPtr
	read, ok := parsingContext.Grammar.Readers[readerName]
	if !ok {
		panic(fmt.Sprintf("grammar does not support reader %q", readerName))
	}
	entityPtr := read(parsingContext)
	if entityPtr == nil {
		// Even if there are problems, the reader should return an entityPtr
		panic(fmt.Sprintf("reader %q returned a non-entity: %T", reflection.GetFunctionName(read), entityPtr))
	}

	// Validate required fields
	reflection.TraverseEntities(entityPtr, false, parsing.ValidateRequiredFields)

	self.Parser.readCache.Store(parsingContext.URL.Key(), entityPtr)

	return self.AddImportFile(context, entityPtr, container, nameTransformer), true
}

// ([parsing.Importer] interface)
func (self *Context) goReadImports(context contextpkg.Context, container *File) {
	importSpecs := parsing.GetImportSpecs(container.EntityPtr)

	// Implicit import
	if !container.GetContext().HasQuirk(parsing.QuirkImportsImplicitDisable) {
		if implicitImportSpec, ok := grammars.GetImplicitImportSpec(container.GetContext()); ok {
			importSpecs = append(importSpecs, implicitImportSpec)
		}
	}

	for _, importSpec := range importSpecs {
		key := importSpec.URL.Key()

		// Skip if causes import loop
		skip := false
		for container_ := container; container_ != nil; container_ = container_.Container {
			url := container_.GetContext().URL
			if url.Key() == key {
				if !importSpec.Implicit {
					// Import loops are considered errors
					container.GetContext().ReportImportLoop(url)
				}
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		promise := util.NewPromise()
		if cached, inCache := self.Parser.readCache.LoadOrStore(key, promise); inCache {
			switch cached_ := cached.(type) {
			case util.Promise:
				// Wait for promise
				logRead.Debugf("wait for promise: %s", key)
				self.readWork.Add(1)
				go self.waitForPromise(context, cached_, key, container, importSpec.NameTransformer)

			default: // entityPtr
				// Cache hit
				logRead.Debugf("cache hit: %s", key)
				self.AddImportFile(context, cached, container, importSpec.NameTransformer)
			}
		} else {
			importToscaContext := container.GetContext().NewImportContext(importSpec.URL)

			// Read (concurrently)
			self.readWork.Add(1)
			go self.read(context, promise, importToscaContext, container, importSpec.NameTransformer, "$File", "")
		}
	}
}

func (self *Context) waitForPromise(context contextpkg.Context, promise util.Promise, key string, container *File, nameTransformer parsing.NameTransformer) {
	defer self.readWork.Done()
	if err := promise.Wait(context); err != nil {
		logRead.Debugf("promise interrupted: %s, %s", key, err.Error())
	}

	if cached, inCache := self.Parser.readCache.Load(key); inCache {
		switch cached.(type) {
		case util.Promise:
			logRead.Debugf("promise broken: %s", key)

		default: // entityPtr
			// Cache hit
			logRead.Debugf("promise kept: %s", key)
			self.AddImportFile(context, cached, container, nameTransformer)
		}
	} else {
		logRead.Debugf("promise broken (empty): %s", key)
	}
}
