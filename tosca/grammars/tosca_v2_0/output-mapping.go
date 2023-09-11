package tosca_v2_0

import (
	"reflect"

	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// OutputMapping
//
// Attaches to notifications and operations
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.15
//

type OutputMapping struct {
	*Entity `name:"output mapping"`
	Name    string

	// EntityName can be a node template name or:
	//
	// * SELF: When in a node template or a relationship
	//         (Groups don't have attributes, so SELF can't be used there)
	// * SOURCE: When in a relationship
	// * TARGET: When in a relationship
	//           (Can only be evaluated after the resolution phase in the Clout)
	//
	// Note that the actual entity is a node template *except* when using SELF in a relationship

	// AttributePath *may* start with a capability name in a node template

	EntityName    *string
	AttributePath []string

	SourceNodeTemplate *NodeTemplate `traverse:"ignore" json:"-" yaml:"-"`
}

func NewOutputMapping(context *parsing.Context) *OutputMapping {
	return &OutputMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadOutputMapping(context *parsing.Context) parsing.EntityPtr {
	self := NewOutputMapping(context)

	if strings := context.ReadStringListMinLength(2); strings != nil {
		self.EntityName = &(*strings)[0]
		self.AttributePath = (*strings)[1:]
	}

	return self
}

// ([parsing.Mappable] interface)
func (self *OutputMapping) GetKey() string {
	return self.Name
}

func (self *OutputMapping) RenderForNodeType(nodeType *NodeType) {
	logRender.Debugf("output mapping: %s", self.Name)

	if (self.EntityName == nil) || (len(self.AttributePath) == 0) {
		return
	}

	entityName := *self.EntityName
	switch entityName {
	case "SELF":
		self.renderForNodeType(nodeType)
	case "SOURCE", "TARGET":
		self.Context.ListChild(0, entityName).ReportValueInvalid("modelable entity name", "cannot be used in node templates")
	default:
		self.renderForNodeTypeByTemplateName(entityName)
	}
}

func (self *OutputMapping) RenderForRelationshipType(relationshipType *RelationshipType, sourceNodeTemplate *NodeTemplate) {
	logRender.Debugf("output mapping: %s", self.Name)

	if (self.EntityName == nil) || (len(self.AttributePath) == 0) {
		return
	}

	switch entityName := *self.EntityName; entityName {
	case "SELF":
		self.renderForRelationshipType(relationshipType)
	case "SOURCE":
		if sourceNodeTemplate != nil {
			self.SourceNodeTemplate = sourceNodeTemplate
			self.renderForNodeType(sourceNodeTemplate.NodeType)
		}
	case "TARGET":
		// The target node type is not known until it is resolved, so there is nothing we can do here
		// (We only know the base node type)
	default:
		self.renderForNodeTypeByTemplateName(entityName)
	}
}

func (self *OutputMapping) RenderForGroup() {
	logRender.Debugf("output mapping: %s", self.Name)

	if (self.EntityName == nil) || (len(self.AttributePath) == 0) {
		return
	}

	switch entityName := *self.EntityName; entityName {
	case "SELF", "SOURCE", "TARGET":
		self.Context.ListChild(0, entityName).ReportValueInvalid("modelable entity name", "cannot be used in groups")
	default:
		self.renderForNodeTypeByTemplateName(entityName)
	}
}

func (self *OutputMapping) renderForNodeTypeByTemplateName(nodeTemplateName string) {
	var nodeTemplateType *NodeTemplate
	if nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, reflect.TypeOf(nodeTemplateType)); ok {
		self.renderForNodeType(nodeTemplate.(*NodeTemplate).NodeType)
	} else {
		self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
		return
	}
}

func (self *OutputMapping) renderForNodeType(nodeType *NodeType) {
	if nodeType == nil {
		return
	}

	// Attribute definitions should already have been rendered

	if len(self.AttributePath) > 1 {
		// Is the first element the name of a capability?
		capabilityName := self.AttributePath[0]
		if capabilityDefinition, ok := nodeType.CapabilityDefinitions[capabilityName]; ok {
			attributeName := self.AttributePath[1]
			if _, ok := capabilityDefinition.AttributeDefinitions[attributeName]; !ok {
				self.Context.ListChild(2, attributeName).ReportReferenceNotFound("attribute", capabilityDefinition)
			}
			return
		}
	}

	attributeName := self.AttributePath[0]
	if _, ok := nodeType.AttributeDefinitions[attributeName]; !ok {
		self.Context.ListChild(1, attributeName).ReportReferenceNotFound("attribute", nodeType)
	}
}

func (self *OutputMapping) renderForRelationshipType(relationshipType *RelationshipType) {
	if relationshipType == nil {
		return
	}

	// Attribute definitions should already have been rendered
	attributeName := self.AttributePath[0]
	if _, ok := relationshipType.AttributeDefinitions[attributeName]; !ok {
		self.Context.ListChild(1, attributeName).ReportReferenceNotFound("attribute", relationshipType)
	}
}

func (self *OutputMapping) Normalize(normalOutputs normal.Values) {
	if (self.EntityName == nil) || (len(self.AttributePath) == 0) {
		return
	}

	list := normal.NewList(len(self.AttributePath) + 1)
	meta := normal.NewValueMeta()
	meta.Type = "list"
	meta.Element = normal.NewValueMeta()
	meta.Element.Type = "string"
	list.SetMeta(meta)

	switch entityName := *self.EntityName; entityName {
	case "SOURCE":
		if self.SourceNodeTemplate != nil {
			list.Set(0, normal.NewPrimitive(self.SourceNodeTemplate.Name))
		} else {
			list.Set(0, normal.NewPrimitive(entityName))
		}
	case "TARGET":
		// Can only be retrieved via a function call
		row, column := self.Context.GetLocation()
		list.Set(0, normal.NewFunctionCall(parsing.NewFunctionCall("tosca.function.$get_target_name", nil, self.Context.URL.String(), row, column, self.Context.Path.String())))
	default:
		list.Set(0, normal.NewPrimitive(entityName))
	}

	for index, pathElement := range self.AttributePath {
		list.Set(index+1, normal.NewPrimitive(pathElement))
	}

	normalOutputs[self.Name] = list
}

//
// OutputMappings
//

type OutputMappings map[string]*OutputMapping

func (self OutputMappings) CopyUnassigned(outputMappings OutputMappings) {
	for key, outputMapping := range outputMappings {
		if _, ok := self[key]; !ok {
			self[key] = outputMapping
		}
	}
}

func (self OutputMappings) Inherit(parent OutputMappings) {
	for name, outputMapping := range parent {
		if _, ok := self[name]; !ok {
			self[name] = outputMapping
		}
	}
}

func (self OutputMappings) RenderForNodeType(nodeType *NodeType) {
	for _, outputMapping := range self {
		outputMapping.RenderForNodeType(nodeType)
	}
}

func (self OutputMappings) RenderForRelationshipType(relationshipType *RelationshipType, sourceNodeTemplate *NodeTemplate) {
	for _, outputMapping := range self {
		outputMapping.RenderForRelationshipType(relationshipType, sourceNodeTemplate)
	}
}

func (self OutputMappings) RenderForGroup() {
	for _, outputMapping := range self {
		outputMapping.RenderForGroup()
	}
}

func (self OutputMappings) Normalize(normalOutputs normal.Values) {
	for _, outputMapping := range self {
		outputMapping.Normalize(normalOutputs)
	}
}
