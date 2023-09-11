package cloudify_v1_3

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// InterfaceAssignment
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-interfaces/]
//

type InterfaceAssignment struct {
	*Entity `name:"interface assignment"`
	Name    string

	Operations OperationAssignments `read:"?,OperationAssignment"`
}

func NewInterfaceAssignment(context *parsing.Context) *InterfaceAssignment {
	return &InterfaceAssignment{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Operations: make(OperationAssignments),
	}
}

// ([parsing.Reader] signature)
func ReadInterfaceAssignment(context *parsing.Context) parsing.EntityPtr {
	self := NewInterfaceAssignment(context)
	context.ReadFields(self)
	return self
}

// ([parsing.Mappable] interface)
func (self *InterfaceAssignment) GetKey() string {
	return self.Name
}

func (self *InterfaceAssignment) GetDefinitionForNodeTemplate(nodeTemplate *NodeTemplate) (*InterfaceDefinition, bool) {
	if nodeTemplate.NodeType == nil {
		return nil, false
	}
	definition, ok := nodeTemplate.NodeType.InterfaceDefinitions[self.Name]
	return definition, ok
}

func (self *InterfaceAssignment) GetDefinitionForRelationshipSource(relationship *RelationshipAssignment) (*InterfaceDefinition, bool) {
	if relationship.RelationshipType == nil {
		return nil, false
	}
	definition, ok := relationship.RelationshipType.SourceInterfaceDefinitions[self.Name]
	return definition, ok
}

func (self *InterfaceAssignment) GetDefinitionForRelationshipTarget(relationship *RelationshipAssignment) (*InterfaceDefinition, bool) {
	if relationship.RelationshipType == nil {
		return nil, false
	}
	definition, ok := relationship.RelationshipType.TargetInterfaceDefinitions[self.Name]
	return definition, ok
}

func (self *InterfaceAssignment) Render(definition *InterfaceDefinition) {
	self.Operations.Render(definition.OperationDefinitions, self.Context)
}

func (self *InterfaceAssignment) Normalize(normalInterface *normal.Interface, definition *InterfaceDefinition) {
	logNormalize.Debugf("interface: %s", self.Name)
	self.Operations.Normalize(normalInterface)
}

//
// InterfaceAssignments
//

type InterfaceAssignments map[string]*InterfaceAssignment

func (self InterfaceAssignments) Render(definitions InterfaceDefinitions, context *parsing.Context) {
	for key, definition := range definitions {
		assignment, ok := self[key]
		if !ok {
			assignment = NewInterfaceAssignment(context.MapChild(key, nil))
			self[key] = assignment
		}
		assignment.Render(definition)
	}

	for key, assignment := range self {
		_, ok := definitions[key]
		if !ok {
			assignment.Context.ReportUndeclared("interface")
			delete(self, key)
		}
	}
}

func (self InterfaceAssignments) NormalizeForNodeTemplate(nodeTemplate *NodeTemplate, normalNodeTemplate *normal.NodeTemplate) {
	for key, interface_ := range self {
		if definition, ok := interface_.GetDefinitionForNodeTemplate(nodeTemplate); ok {
			interface_.Normalize(normalNodeTemplate.NewInterface(key), definition)
		}
	}
}

func (self InterfaceAssignments) NormalizeForRelationshipSource(relationship *RelationshipAssignment, normalRelationship *normal.Relationship) {
	for key, interface_ := range self {
		if definition, ok := interface_.GetDefinitionForRelationshipSource(relationship); ok {
			normalInterface := normalRelationship.NewInterface(key)
			interface_.Normalize(normalInterface, definition)
			for _, operation := range normalInterface.Operations {
				operation.Host = "SOURCE"
			}
		}
	}
}

func (self InterfaceAssignments) NormalizeForRelationshipTarget(relationship *RelationshipAssignment, normalRelationship *normal.Relationship) {
	for key, interface_ := range self {
		if definition, ok := interface_.GetDefinitionForRelationshipTarget(relationship); ok {
			normalInterface := normalRelationship.NewInterface(key)
			interface_.Normalize(normalInterface, definition)
			for _, operation := range normalInterface.Operations {
				operation.Host = "TARGET"
			}
		}
	}
}
