package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// CapabilityType
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.7
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.7
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.7
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.6
//

type CapabilityType struct {
	*Type `name:"capability type"`

	PropertyDefinitions      PropertyDefinitions  `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	AttributeDefinitions     AttributeDefinitions `read:"attributes,AttributeDefinition" inherit:"attributes,Parent"`
	ValidSourceNodeTypeNames *[]string            `read:"valid_source_types" inherit:"valid_source_types,Parent"`

	Parent               *CapabilityType `lookup:"derived_from,ParentName" traverse:"ignore" json:"-" yaml:"-"`
	ValidSourceNodeTypes NodeTypes       `lookup:"valid_source_types,ValidSourceNodeTypeNames" inherit:"valid_source_types,Parent" traverse:"ignore" json:"-" yaml:"-"`
}

func NewCapabilityType(context *parsing.Context) *CapabilityType {
	return &CapabilityType{
		Type:                 NewType(context),
		PropertyDefinitions:  make(PropertyDefinitions),
		AttributeDefinitions: make(AttributeDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadCapabilityType(context *parsing.Context) parsing.EntityPtr {
	self := NewCapabilityType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Hierarchical] interface)
func (self *CapabilityType) GetParent() parsing.EntityPtr {
	return self.Parent
}

// ([parsing.Inherits] interface)
func (self *CapabilityType) Inherit() {
	logInherit.Debugf("capability type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
	self.AttributeDefinitions.Inherit(self.Parent.AttributeDefinitions)
}

//
// CapabilityTypes
//

type CapabilityTypes []*CapabilityType
