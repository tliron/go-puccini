package tosca_v2_0

import (
	"math"

	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// CapabilityAssignment
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.1
//

type CapabilityAssignment struct {
	*Entity `name:"capability"`
	Name    string

	Properties  Values       `read:"properties,Value"`
	Attributes  Values       `read:"attributes,AttributeValue"`
	Occurrences *RangeEntity `read:"occurrences,RangeEntity"` // introduced in TOSCA 1.3
}

func NewCapabilityAssignment(context *parsing.Context) *CapabilityAssignment {
	return &CapabilityAssignment{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Properties: make(Values),
		Attributes: make(Values),
	}
}

// ([parsing.Reader] signature)
func ReadCapabilityAssignment(context *parsing.Context) parsing.EntityPtr {
	self := NewCapabilityAssignment(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Mappable] interface)
func (self *CapabilityAssignment) GetKey() string {
	return self.Name
}

func (self *CapabilityAssignment) GetDefinition(nodeTemplate *NodeTemplate) (*CapabilityDefinition, bool) {
	if nodeTemplate.NodeType == nil {
		return nil, false
	}
	definition, ok := nodeTemplate.NodeType.CapabilityDefinitions[self.Name]
	return definition, ok
}

func (self *CapabilityAssignment) Normalize(normalNodeTemplate *normal.NodeTemplate, definition *CapabilityDefinition) *normal.Capability {
	logNormalize.Debugf("capability: %s", self.Name)

	normalCapability := normalNodeTemplate.NewCapability(self.Name, normal.NewLocationForContext(self.Context))

	if definition.Description != nil {
		normalCapability.Description = *definition.Description
	}

	if self.Occurrences != nil {
		normalCapability.MinRelationshipCount = self.Occurrences.Range.Lower
		normalCapability.MaxRelationshipCount = self.Occurrences.Range.Upper
	} else {
		// Default occurrences is [ 0, UNBOUNDED ]
		normalCapability.MinRelationshipCount = 0
		normalCapability.MaxRelationshipCount = math.MaxUint64
	}

	if types, ok := normal.GetEntityTypes(self.Context.Hierarchy, definition.CapabilityType); ok {
		normalCapability.Types = types
	}

	self.Properties.Normalize(normalCapability.Properties)
	self.Attributes.Normalize(normalCapability.Attributes)

	return normalCapability
}

//
// CapabilityAssignments
//

type CapabilityAssignments map[string]*CapabilityAssignment

func (self *CapabilityAssignment) Render(definition *CapabilityDefinition) {
	self.Properties.RenderProperties(definition.PropertyDefinitions, self.Context.FieldChild("properties", nil))
	self.Attributes.RenderAttributes(definition.AttributeDefinitions, self.Context.FieldChild("attributes", nil))
	if self.Occurrences == nil {
		self.Occurrences = definition.Occurrences
	}
}

func (self CapabilityAssignments) Render(definitions CapabilityDefinitions, context *parsing.Context) {
	for key, definition := range definitions {
		assignment, ok := self[key]
		if !ok {
			assignment = NewCapabilityAssignment(context.MapChild(key, nil))
			self[key] = assignment
		}
		assignment.Render(definition)
	}

	for key, assignment := range self {
		if _, ok := definitions[key]; !ok {
			assignment.Context.ReportUndeclared("capability")
			delete(self, key)
		}
	}
}

func (self CapabilityAssignments) Normalize(nodeTemplate *NodeTemplate, normalNodeTemplate *normal.NodeTemplate) {
	for key, capability := range self {
		if definition, ok := capability.GetDefinition(nodeTemplate); ok {
			normalNodeTemplate.Capabilities[key] = capability.Normalize(normalNodeTemplate, definition)
		}
	}
}
