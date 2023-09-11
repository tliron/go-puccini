package main

import (
	"C"
	contextpkg "context"
	"errors"
	"strings"

	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/go-transcribe"
	"github.com/tliron/kutil/problems"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/normal"
	parserpkg "github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

var parser = parserpkg.NewParser()

//export Compile
func Compile(url *C.char, inputs *C.char, quirks *C.char, resolve C.char, coerce C.char) *C.char {
	context := contextpkg.TODO()

	inputs_ := make(map[string]ard.Value)

	if data, err := yamlkeys.DecodeAll(strings.NewReader(C.GoString(inputs))); err == nil {
		for _, data_ := range data {
			if map_, ok := data_.(ard.Map); ok {
				for key, value := range map_ {
					inputs_[yamlkeys.KeyString(key)] = value
				}
			} else {
				return result(nil, nil, errors.New("malformed inputs"))
			}
		}
	} else {
		return result(nil, nil, err)
	}

	var quirks_ parsing.Quirks

	if data, err := yamlkeys.DecodeAll(strings.NewReader(C.GoString(quirks))); err == nil {
		for _, data_ := range data {
			if list, ok := data_.(ard.List); ok {
				for _, value := range list {
					if value_, ok := value.(string); ok {
						quirks_ = append(quirks_, parsing.Quirk(value_))
					} else {
						return result(nil, nil, errors.New("malformed quirk"))
					}
				}
			} else {
				return result(nil, nil, errors.New("malformed quirks"))
			}
		}
	} else {
		return result(nil, nil, err)
	}

	urlContext := exturl.NewContext()
	defer urlContext.Release()

	var url_ exturl.URL
	var err error
	if url_, err = urlContext.NewValidAnyOrFileURL(context, C.GoString(url), nil); err != nil {
		return result(nil, nil, err)
	}

	parserContext := parser.NewContext()
	parserContext.URL = url_
	parserContext.Quirks = quirks_
	parserContext.Inputs = inputs_
	var normalServiceTemplate *normal.ServiceTemplate
	if normalServiceTemplate, err = parserContext.Parse(context); err != nil {
		return result(nil, parserContext.GetProblems(), err)
	}

	problems := parserContext.GetProblems()

	var clout *cloutpkg.Clout
	if clout, err = normalServiceTemplate.Compile(); err != nil {
		return result(clout, problems, err)
	}

	execContext := js.ExecContext{
		Clout:      clout,
		Problems:   problems,
		URLContext: urlContext,
		History:    true,
		Format:     "yaml",
		Strict:     true,
	}

	if resolve != 0 {
		execContext.Resolve()
		if !problems.Empty() {
			return result(clout, problems, nil)
		}
	}

	if coerce != 0 {
		execContext.Coerce()
		if !problems.Empty() {
			return result(clout, problems, nil)
		}
	}

	return result(clout, problems, nil)
}

func result(clout *cloutpkg.Clout, problems *problems.Problems, err error) *C.char {
	result := make(ard.StringMap)
	if clout != nil {
		result["clout"] = clout
	}
	if (problems != nil) && !problems.Empty() {
		result["problems"] = problems.Problems
	}
	if err != nil {
		result["error"] = err.Error()
	}

	result_, _ := (&transcribe.Transcriber{Indent: "  ", Strict: true}).StringifyYAML(result) // TODO: err
	return C.CString(result_)
}
