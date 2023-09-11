package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// GroupType
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.11
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.10
//

type GroupType struct {
	*Type `name:"group type"`

	PropertyDefinitions    PropertyDefinitions    `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	CapabilityDefinitions  CapabilityDefinitions  `read:"capabilities,CapabilityDefinition" inherit:"capabilities,Parent"`
	RequirementDefinitions RequirementDefinitions `read:"requirements,{}RequirementDefinition" inherit:"requirements,Parent"` // sequenced list, but we read it into map
	InterfaceDefinitions   InterfaceDefinitions   `inherit:"interfaces,Parent"`                                               // removed in TOSCA 1.3
	MemberNodeTypeNames    *[]string              `read:"members" inherit:"members,Parent"`

	Parent          *GroupType `lookup:"derived_from,ParentName" inherit:"members,Parent" traverse:"ignore" json:"-" yaml:"-"`
	MemberNodeTypes NodeTypes  `lookup:"members,MemberNodeTypeNames" inherit:"members,Parent" traverse:"ignore" json:"-" yaml:"-"`
}

func NewGroupType(context *parsing.Context) *GroupType {
	return &GroupType{
		Type:                   NewType(context),
		PropertyDefinitions:    make(PropertyDefinitions),
		CapabilityDefinitions:  make(CapabilityDefinitions),
		RequirementDefinitions: make(RequirementDefinitions),
		InterfaceDefinitions:   make(InterfaceDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadGroupType(context *parsing.Context) parsing.EntityPtr {
	self := NewGroupType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Hierarchical] interface)
func (self *GroupType) GetParent() parsing.EntityPtr {
	return self.Parent
}

// ([parsing.Inherits] interface)
func (self *GroupType) Inherit() {
	logInherit.Debugf("group type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
	self.CapabilityDefinitions.Inherit(self.Parent.CapabilityDefinitions)
	self.RequirementDefinitions.Inherit(self.Parent.RequirementDefinitions)
	self.InterfaceDefinitions.Inherit(self.Parent.InterfaceDefinitions)

	// (Note we are checking for MemberNodeTypeNames and not MemberNodeTypes, because the latter will never be nil)
	if self.MemberNodeTypeNames == nil {
		self.MemberNodeTypeNames = self.Parent.MemberNodeTypeNames
		self.MemberNodeTypes = self.Parent.MemberNodeTypes
	}
	// We cannot handle the "else" case here, because the node type hierarchy may not have been created yet,
	// So we will do that check in the rendering phase, below
}

// ([parsing.Renderable] interface)
func (self *GroupType) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *GroupType) render() {
	logRender.Debugf("group type: %s", self.Name)

	// (Note we are checking for MemberNodeTypeNames and not MemberNodeTypes, because the latter will never be nil)
	if self.Parent == nil {
		return
	}

	if self.Parent.MemberNodeTypeNames == nil {
		return
	}

	context := self.Context.FieldChild("members", nil)
	self.Parent.MemberNodeTypes.ValidateSubset(self.MemberNodeTypes, context)
}

//
// GroupTypes
//

type GroupTypes []*GroupType

func (self GroupTypes) IsCompatible(groupType *GroupType) bool {
	for _, baseGroupType := range self {
		if baseGroupType.Context.Hierarchy.IsCompatible(baseGroupType, groupType) {
			return true
		}
	}
	return false
}

func (self GroupTypes) ValidateSubset(subset GroupTypes, context *parsing.Context) bool {
	isSubset := true
	for _, subsetGroupType := range subset {
		if !self.IsCompatible(subsetGroupType) {
			context.ReportIncompatibleTypeInSet(subsetGroupType)
			isSubset = false
		}
	}
	return isSubset
}
