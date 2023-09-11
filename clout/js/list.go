package js

import (
	"fmt"

	"github.com/tliron/go-ard"
)

func (self *ExecutionContext) NewListValue(list ard.List, notation ard.StringMap, meta ard.StringMap) (*Value, error) {
	var elementMeta ard.StringMap
	if meta != nil {
		if data, ok := meta["element"]; ok {
			if map_, ok := data.(ard.StringMap); ok {
				elementMeta = map_
			} else {
				return nil, fmt.Errorf("malformed meta \"element\", not a map: %T", data)
			}
		}
	}

	if list_, err := self.NewList(list, elementMeta); err == nil {
		return self.NewValue(list_, notation, meta)
	} else {
		return nil, err
	}
}

//
// List
//

type List []Coercible

func (self *ExecutionContext) NewList(list ard.List, elementMeta ard.StringMap) (List, error) {
	list_ := make(List, len(list))

	for index, data := range list {
		var err error
		if list_[index], err = self.NewCoercible(data, elementMeta); err != nil {
			return nil, err
		}
	}

	return list_, nil
}

func (self List) Coerce() (ard.Value, error) {
	list := make(ard.List, len(self))

	for index, element := range self {
		var err error
		if list[index], err = element.Coerce(); err != nil {
			return nil, err
		}
	}

	return list, nil
}
