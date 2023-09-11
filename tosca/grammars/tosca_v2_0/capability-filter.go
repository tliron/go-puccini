package tosca_v2_0

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// CapabilityFilter
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.5.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.5.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.4.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.4.2
//

type CapabilityFilter struct {
	*Entity `name:"capability filter"`
	Name    string // name or type name

	PropertyFilters PropertyFilters `read:"properties,{}PropertyFilter"`
}

func NewCapabilityFilter(context *parsing.Context) *CapabilityFilter {
	return &CapabilityFilter{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadCapabilityFilter(context *parsing.Context) parsing.EntityPtr {
	self := NewCapabilityFilter(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self CapabilityFilter) Normalize(normalRequirement *normal.Requirement) normal.FunctionCallMap {
	if len(self.PropertyFilters) == 0 {
		return nil
	}

	// TODO: separate maps for by-name vs. by-type-name

	var normalFunctionCallMap normal.FunctionCallMap
	var ok bool
	if normalFunctionCallMap, ok = normalRequirement.CapabilityPropertyValidation[self.Name]; !ok {
		normalFunctionCallMap = make(normal.FunctionCallMap)
		normalRequirement.CapabilityPropertyValidation[self.Name] = normalFunctionCallMap
	}

	self.PropertyFilters.Normalize(normalFunctionCallMap)

	return normalFunctionCallMap
}

//
// CapabilityFilters
//

type CapabilityFilters []*CapabilityFilter

func (self CapabilityFilters) Normalize(normalRequirement *normal.Requirement) {
	for _, capabilityFilter := range self {
		capabilityFilter.Normalize(normalRequirement)
	}
}
