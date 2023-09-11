package hot

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// ParameterGroup
//
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#parameter-groups-section]
//

type ParameterGroup struct {
	*Entity `name:"parameter group"`

	Label       *string   `read:"label"`
	Description *string   `read:"description"`
	Parameters  []*string `read:"parameters"`
}

func NewParameterGroup(context *parsing.Context) *ParameterGroup {
	return &ParameterGroup{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadParameterGroup(context *parsing.Context) parsing.EntityPtr {
	self := NewParameterGroup(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

//
// ParameterGroups
//

type ParameterGroups []*ParameterGroup
