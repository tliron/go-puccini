package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// RelationshipType
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.10
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.9
//

type RelationshipType struct {
	*Type `name:"relationship type"`

	PropertyDefinitions            PropertyDefinitions  `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	AttributeDefinitions           AttributeDefinitions `read:"attributes,AttributeDefinition" inherit:"attributes,Parent"`
	InterfaceDefinitions           InterfaceDefinitions `read:"interfaces,InterfaceDefinition" inherit:"interfaces,Parent"`
	ValidTargetCapabilityTypeNames *[]string            `read:"valid_target_types" inherit:"valid_target_types,Parent"`

	Parent                     *RelationshipType `lookup:"derived_from,ParentName" traverse:"ignore" json:"-" yaml:"-"`
	ValidTargetCapabilityTypes CapabilityTypes   `lookup:"valid_target_types,ValidTargetCapabilityTypeNames" inherit:"valid_target_types,Parent" traverse:"ignore" json:"-" yaml:"-"`
}

func NewRelationshipType(context *parsing.Context) *RelationshipType {
	return &RelationshipType{
		Type:                 NewType(context),
		PropertyDefinitions:  make(PropertyDefinitions),
		AttributeDefinitions: make(AttributeDefinitions),
		InterfaceDefinitions: make(InterfaceDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadRelationshipType(context *parsing.Context) parsing.EntityPtr {
	self := NewRelationshipType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Hierarchical] interface)
func (self *RelationshipType) GetParent() parsing.EntityPtr {
	return self.Parent
}

// ([parsing.Inherits] interface)
func (self *RelationshipType) Inherit() {
	logInherit.Debugf("relationship type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
	self.AttributeDefinitions.Inherit(self.Parent.AttributeDefinitions)
	self.InterfaceDefinitions.Inherit(self.Parent.InterfaceDefinitions)
}

//
// RelationshipTypes
//

type RelationshipTypes []*RelationshipType
