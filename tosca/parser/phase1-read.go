package parser

import (
	contextpkg "context"
	"fmt"
	"sort"

	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca/csar"
	"github.com/tliron/puccini/tosca/grammars"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

func (self *ServiceContext) ReadRoot(context contextpkg.Context, url exturl.URL, origins []exturl.URL, template string) bool {
	toscaContext := parsing.NewContext(self.Stylist, self.Quirks)
	toscaContext.Origins = origins

	toscaContext.URL = url

	var ok bool

	self.readWork.Add(1)
	self.Root, ok = self.read(context, nil, toscaContext, nil, nil, "$Root", template)
	self.readWork.Wait()

	self.filesLock.Lock()
	sort.Sort(self.Files)
	self.filesLock.Unlock()

	return ok
}

func (self *ServiceContext) read(context contextpkg.Context, promise util.Promise, toscaContext *parsing.Context, container *File, nameTransformer parsing.NameTransformer, readerName string, serviceTemplateName string) (*File, bool) {
	defer self.readWork.Done()
	if promise != nil {
		// For the goroutines waiting for our cached entityPtr
		defer promise.Release()
	}

	logRead.Infof("%s: %s", readerName, toscaContext.URL.Key())

	// TODO: allow override of CSAR format
	if format := toscaContext.URL.Format(); csar.IsValidFormat(format) {
		var err error
		if toscaContext.URL, err = csar.GetServiceTemplateURL(context, toscaContext.URL, format, serviceTemplateName); err != nil {
			toscaContext.ReportError(err)
			file := NewFileNoEntity(toscaContext, container, nameTransformer)
			self.AddFile(file)
			return file, false
		}
	}

	// Read ARD
	var err error
	if toscaContext.Data, toscaContext.Locator, err = ard.ReadURL(context, toscaContext.URL, true); err != nil {
		if decodeError, ok := err.(*yamlkeys.DecodeError); ok {
			err = NewYAMLDecodeError(decodeError)
		}
		toscaContext.ReportError(err)
		file := NewFileNoEntity(toscaContext, container, nameTransformer)
		self.AddFile(file)
		return file, false
	}

	// Detect grammar
	if !grammars.Detect(toscaContext) {
		file := NewFileNoEntity(toscaContext, container, nameTransformer)
		self.AddFile(file)
		return file, false
	}

	// Read entityPtr
	read, ok := toscaContext.Grammar.Readers[readerName]
	if !ok {
		panic(fmt.Sprintf("grammar does not support reader %q", readerName))
	}
	entityPtr := read(toscaContext)
	if entityPtr == nil {
		// Even if there are problems, the reader should return an entityPtr
		panic(fmt.Sprintf("reader %q returned a non-entity: %T", reflection.GetFunctionName(read), entityPtr))
	}

	// Validate required fields
	reflection.TraverseEntities(entityPtr, false, parsing.ValidateRequiredFields)

	self.Context.readCache.Store(toscaContext.URL.Key(), entityPtr)

	return self.AddImportFile(context, entityPtr, container, nameTransformer), true
}

// From tosca.Importer interface
func (self *ServiceContext) goReadImports(context contextpkg.Context, container *File) {
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
		if cached, inCache := self.Context.readCache.LoadOrStore(key, promise); inCache {
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

func (self *ServiceContext) waitForPromise(context contextpkg.Context, promise util.Promise, key string, container *File, nameTransformer parsing.NameTransformer) {
	defer self.readWork.Done()
	if err := promise.Wait(context); err != nil {
		logRead.Debugf("promise interrupted: %s, %s", key, err.Error())
	}

	if cached, inCache := self.Context.readCache.Load(key); inCache {
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
