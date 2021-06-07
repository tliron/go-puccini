package js

import (
	"fmt"

	"github.com/dop251/goja"
)

func CallFunction(runtime *goja.Runtime, function goja.Value, this interface{}, arguments []interface{}) (interface{}, error) {
	function_, ok := goja.AssertFunction(function)
	if !ok {
		return nil, fmt.Errorf("not a function: %v", function)
	}

	values := make([]goja.Value, len(arguments))
	for index, argument := range arguments {
		values[index] = runtime.ToValue(argument)
	}

	r, err := function_(runtime.ToValue(this), values...)
	if err != nil {
		return nil, UnwrapException(err)
	}

	return r.Export(), nil
}
