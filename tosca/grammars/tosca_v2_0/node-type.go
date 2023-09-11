package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// NodeType
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.9
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.9
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.9
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.8
//

type NodeType struct {
	*Type `name:"node type"`

	PropertyDefinitions    PropertyDefinitions    `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	AttributeDefinitions   AttributeDefinitions   `read:"attributes,AttributeDefinition" inherit:"attributes,Parent"`
	CapabilityDefinitions  CapabilityDefinitions  `read:"capabilities,CapabilityDefinition" inherit:"capabilities,Parent"`
	RequirementDefinitions RequirementDefinitions `read:"requirements,{}RequirementDefinition" inherit:"requirements,Parent"` // sequenced list, but we read it into map
	InterfaceDefinitions   InterfaceDefinitions   `read:"interfaces,InterfaceDefinition" inherit:"interfaces,Parent"`
	ArtifactDefinitions    ArtifactDefinitions    `read:"artifacts,ArtifactDefinition" inherit:"artifacts,Parent"`

	Parent *NodeType `lookup:"derived_from,ParentName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewNodeType(context *parsing.Context) *NodeType {
	return &NodeType{
		Type:                   NewType(context),
		PropertyDefinitions:    make(PropertyDefinitions),
		AttributeDefinitions:   make(AttributeDefinitions),
		CapabilityDefinitions:  make(CapabilityDefinitions),
		RequirementDefinitions: make(RequirementDefinitions),
		InterfaceDefinitions:   make(InterfaceDefinitions),
		ArtifactDefinitions:    make(ArtifactDefinitions),
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

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
	self.AttributeDefinitions.Inherit(self.Parent.AttributeDefinitions)
	self.CapabilityDefinitions.Inherit(self.Parent.CapabilityDefinitions)
	self.RequirementDefinitions.Inherit(self.Parent.RequirementDefinitions)
	self.InterfaceDefinitions.Inherit(self.Parent.InterfaceDefinitions)
	self.ArtifactDefinitions.Inherit(self.Parent.ArtifactDefinitions)
}

//
// NodeTypes
//

type NodeTypes []*NodeType

func (self NodeTypes) IsCompatible(nodeType *NodeType) bool {
	for _, baseNodeType := range self {
		if baseNodeType.Context.Hierarchy.IsCompatible(baseNodeType, nodeType) {
			return true
		}
	}
	return false
}

func (self NodeTypes) ValidateSubset(subset NodeTypes, context *parsing.Context) bool {
	isSubset := true
	for _, subsetNodeType := range subset {
		if !self.IsCompatible(subsetNodeType) {
			context.ReportIncompatibleTypeInSet(subsetNodeType)
			isSubset = false
		}
	}
	return isSubset
}
