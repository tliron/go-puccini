package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// NodeType
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-node-types/]
//

type NodeType struct {
	*Type `name:"node type"`

	InterfaceDefinitions InterfaceDefinitions `read:"interfaces,InterfaceDefinition" inherit:"interfaces,Parent"`
	PropertyDefinitions  PropertyDefinitions  `read:"properties,PropertyDefinition" inherit:"properties,Parent"`

	Parent *NodeType `lookup:"derived_from,ParentName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewNodeType(context *parsing.Context) *NodeType {
	return &NodeType{
		Type:                 NewType(context),
		InterfaceDefinitions: make(InterfaceDefinitions),
		PropertyDefinitions:  make(PropertyDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadNodeType(context *parsing.Context) parsing.EntityPtr {
	self := NewNodeType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Hierarchical] interface)
func (self *NodeType) GetParent() parsing.EntityPtr {
	return self.Parent
}

// ([parsing.Inherits] interface)
func (self *NodeType) Inherit() {
	logInherit.Debugf("node type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.InterfaceDefinitions.Inherit(self.Parent.InterfaceDefinitions)
	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
}

//
// NodeTypes
//

type NodeTypes []*NodeType
