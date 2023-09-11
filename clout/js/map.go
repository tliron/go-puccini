package js

import (
	"fmt"

	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/util"
)

func (self *ExecutionContext) NewMapValue(list ard.List, notation ard.StringMap, meta ard.StringMap) (*Value, error) {
	var keyMeta ard.StringMap
	var valueMeta ard.StringMap
	var fieldsMeta ard.StringMap
	if meta != nil {
		if data, ok := meta["key"]; ok {
			if map_, ok := data.(ard.StringMap); ok {
				keyMeta = map_
			} else {
				return nil, fmt.Errorf("malformed meta \"key\", not a map: %T", data)
			}
		}

		if data, ok := meta["value"]; ok {
			if map_, ok := data.(ard.StringMap); ok {
				valueMeta = map_
			} else {
				return nil, fmt.Errorf("malformed meta \"value\", not a map: %T", data)
			}
		}

		if data, ok := meta["fields"]; ok {
			if map_, ok := data.(ard.StringMap); ok {
				fieldsMeta = map_
			} else {
				return nil, fmt.Errorf("malformed meta \"fields\", not a map: %T", data)
			}
		}
	}

	if map_, err := self.NewMap(list, keyMeta, valueMeta, fieldsMeta); err == nil {
		return self.NewValue(map_, notation, meta)
	} else {
		return nil, err
	}
}

//
// Map
//

type Map []MapEntry

func (self *ExecutionContext) NewMap(list ard.List, keyMeta ard.StringMap, valueMeta ard.StringMap, fieldsMeta ard.StringMap) (Map, error) {
	map_ := make(Map, len(list))

	if fieldsMeta != nil {
		if (keyMeta != nil) || (valueMeta != nil) {
			return nil, fmt.Errorf("has \"fields\" meta in addition to \"key\" or \"value\" meta")
		}
	}

	for index, data := range list {
		var err error
		if map_[index], err = self.NewMapEntry(data, keyMeta, valueMeta, fieldsMeta); err != nil {
			return nil, err
		}
	}

	return map_, nil
}

func (self Map) Coerce() (ard.Value, error) {
	map_ := make(ard.StringMap)

	for _, entry := range self {
		if key, value, err := entry.Coerce(); err == nil {
			if _, ok := map_[key]; ok {
				return nil, fmt.Errorf("duplicate map key during coercion: %s", key)
			}

			map_[key] = value
		} else {
			return nil, err
		}
	}

	return map_, nil
}

//
// MapEntry
//

type MapEntry struct {
	Key   Coercible `json:"key" yaml:"key"`
	Value Coercible `json:"value" yaml:"value"`
}

func (self *ExecutionContext) NewMapEntry(data any, keyMeta ard.StringMap, valueMeta ard.StringMap, fieldsMeta ard.StringMap) (MapEntry, error) {
	var entry MapEntry

	if map_, ok := data.(ard.StringMap); ok {
		if key, ok := map_["$key"]; ok {
			var err error
			if entry.Key, err = self.NewCoercible(key, keyMeta); err == nil {
				if fieldsMeta != nil {
					// Find field meta
					if key, err := entry.Key.Coerce(); err == nil {
						if key_, ok := key.(string); ok {
							if valueMeta_, ok := fieldsMeta[key_]; ok {
								if valueMeta, ok = valueMeta_.(ard.StringMap); !ok {
									return entry, fmt.Errorf("malformed meta \"fields\", not a map: %T", valueMeta_)
								}
							} else {
								return entry, fmt.Errorf("malformed meta \"fields\", field not found: %q", key)
							}
						} else {
							return entry, fmt.Errorf("malformed field name, not a string: %T", key)
						}
					} else {
						return entry, err
					}
				}

				if entry.Value, err = self.NewCoercible(map_, valueMeta); err == nil {
					return entry, nil
				} else {
					return entry, err
				}
			} else {
				return entry, err
			}
		} else {
			return entry, fmt.Errorf("malformed map entry, no \"$key\": %+v", map_)
		}
	} else {
		return entry, fmt.Errorf("malformed map entry, not a map: %T", data)
	}
}

func (self MapEntry) Coerce() (string, any, error) {
	if key, err := self.Key.Coerce(); err == nil {
		if value, err := self.Value.Coerce(); err == nil {
			return util.ToString(key), value, nil
		} else {
			return "", nil, err
		}
	} else {
		return "", nil, err
	}
}
