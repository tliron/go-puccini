package cloudify_v1_3

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Group
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-groups/]
//

type Group struct {
	*Entity `name:"group"`
	Name    string `namespace:""`

	MemberNodeTemplateNames *[]string     `read:"members" mandatory:""`
	Policies                GroupPolicies `read:"policies,GroupPolicy"`

	MemberNodeTemplates NodeTemplates `lookup:"members,MemberNodeTemplateNames" traverse:"ignore" json:"-" yaml:"-"`
}

func NewGroup(context *parsing.Context) *Group {
	return &Group{
		Entity:   NewEntity(context),
		Name:     context.Name,
		Policies: make(GroupPolicies),
	}
}

// ([parsing.Reader] signature)
func ReadGroup(context *parsing.Context) parsing.EntityPtr {
	self := NewGroup(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

var groupTypeName = "cloudify.Group"
var groupTypes = normal.NewEntityTypes(groupTypeName)

func (self *Group) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.Group {
	logNormalize.Debugf("group: %s", self.Name)

	normalGroup := normalServiceTemplate.NewGroup(self.Name)
	normalGroup.Types = groupTypes

	for _, nodeTemplate := range self.MemberNodeTemplates {
		if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[nodeTemplate.Name]; ok {
			normalGroup.Members = append(normalGroup.Members, normalNodeTemplate)
		}
	}

	// TODO: normalize policies
	// TODO: normalize triggers in policies

	return normalGroup
}

//
// Groups
//

type Groups []*Group

func (self Groups) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, group := range self {
		normalServiceTemplate.Groups[group.Name] = group.Normalize(normalServiceTemplate)
	}
}
