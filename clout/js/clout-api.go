package js

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dop251/goja"
	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/js"
	cloutpkg "github.com/tliron/puccini/clout"
)

//
// CloutAPI
//

type CloutAPI struct {
	*cloutpkg.Clout

	cloutContext *CloutContext
}

func (self *Context) NewCloutAPI(clout *cloutpkg.Clout, jsContext *js.Context) *CloutAPI {
	return &CloutAPI{
		clout,
		self.NewCloutContext(clout, jsContext),
	}
}

func (self *CloutAPI) Load(data interface{}) (*CloutAPI, error) {
	var clout *cloutpkg.Clout
	var err error

	switch data_ := data.(type) {
	case string:
		if clout, err = cloutpkg.Load(data_, ""); err != nil {
			return nil, err
		}

	case ard.Map:
		if clout, err = cloutpkg.Parse(data_); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("not a URL or clout data: %T", data)
	}

	return self.cloutContext.Context.NewCloutAPI(clout, self.cloutContext.JSContext), nil
}

func (self *CloutAPI) NewKey() string {
	return cloutpkg.NewKey()
}

func (self *CloutAPI) Call(scriptletName string, functionName string, arguments []interface{}) (interface{}, error) {
	return self.cloutContext.CallFunction(scriptletName, functionName, arguments, FunctionCallContext{})
}

func (self *CloutAPI) Define(scriptletName string, scriptlet string) error {
	return SetScriptlet(scriptletName, CleanupScriptlet(scriptlet), self.Clout)
}

func (self *CloutAPI) NewCoercible(value goja.Value, site interface{}, source interface{}, target interface{}) (Coercible, error) {
	if goja.IsUndefined(value) {
		return nil, errors.New("undefined")
	}

	if coercible, err := self.cloutContext.NewCoercible(value.Export(), FunctionCallContext{site, source, target}); err == nil {
		return coercible, nil
	} else {
		return nil, err
	}
}

func (self *CloutAPI) NewConstraints(value goja.Value, site interface{}, source interface{}, target interface{}) (Constraints, error) {
	if goja.IsUndefined(value) {
		return nil, errors.New("undefined")
	}

	exported := value.Export()
	if list_, ok := exported.(ard.List); ok {
		if constraints, err := self.cloutContext.NewConstraints(list_, FunctionCallContext{site, source, target}); err == nil {
			return constraints, nil
		} else {
			return nil, err
		}
	}

	return nil, errors.New("not an array")
}

func (self *CloutAPI) Coerce(value interface{}) (interface{}, error) {
	if coercible, ok := value.(Coercible); ok {
		return coercible.Coerce()
	}

	return value, nil
}

func (self *CloutAPI) Unwrap(value interface{}) interface{} {
	if coercible, ok := value.(Coercible); ok {
		return coercible.Unwrap()
	}

	return value
}

// json.Marshaler interface
func (self *CloutAPI) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Clout)
}

// yaml.Marshaler interface
func (self *CloutAPI) MarshalYAML() (interface{}, error) {
	return self.Clout, nil
}

// cbor.Marshaler interface
func (self *CloutAPI) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(self.Clout)
}
