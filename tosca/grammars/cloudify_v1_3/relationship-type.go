package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// RelationshipType
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-relationships/]
//

type RelationshipType struct {
	*Type `name:"relationship type"`

	SourceInterfaceDefinitions InterfaceDefinitions `read:"source_interfaces,InterfaceDefinition" inherit:"source_interfaces,Parent"`
	TargetInterfaceDefinitions InterfaceDefinitions `read:"target_interfaces,InterfaceDefinition" inherit:"target_interfaces,Parent"`
	PropertyDefinitions        PropertyDefinitions  `read:"properties,PropertyDefinition" inherit:"properties,Parent"`

	Parent *RelationshipType `lookup:"derived_from,ParentName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewRelationshipType(context *parsing.Context) *RelationshipType {
	return &RelationshipType{
		Type:                       NewType(context),
		SourceInterfaceDefinitions: make(InterfaceDefinitions),
		TargetInterfaceDefinitions: make(InterfaceDefinitions),
		PropertyDefinitions:        make(PropertyDefinitions),
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

	self.SourceInterfaceDefinitions.Inherit(self.Parent.SourceInterfaceDefinitions)
	self.TargetInterfaceDefinitions.Inherit(self.Parent.TargetInterfaceDefinitions)
}

//
// RelationshipTypes
//

type RelationshipTypes []*RelationshipType
