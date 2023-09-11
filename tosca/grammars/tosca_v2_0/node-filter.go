package tosca_v2_0

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// NodeFilter
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.5
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.4
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.4
//

type NodeFilter struct {
	*Entity `name:"node filter"`

	PropertyFilters   PropertyFilters   `read:"properties,{}PropertyFilter"`
	CapabilityFilters CapabilityFilters `read:"capabilities,{}CapabilityFilter"`
}

func NewNodeFilter(context *parsing.Context) *NodeFilter {
	return &NodeFilter{
		Entity: NewEntity(context),
	}
}

// ([parsing.Reader] signature)
func ReadNodeFilter(context *parsing.Context) parsing.EntityPtr {
	self := NewNodeFilter(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *NodeFilter) Normalize(normalRequirement *normal.Requirement) {
	self.PropertyFilters.Normalize(normalRequirement.NodeTemplatePropertyValidation)
	self.CapabilityFilters.Normalize(normalRequirement)
}
