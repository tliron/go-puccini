package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// OperationDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.17
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.15
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.13
//

type OperationDefinition struct {
	*Entity `name:"operation definition"`
	Name    string

	Description      *string                  `read:"description"`
	Implementation   *InterfaceImplementation `read:"implementation,InterfaceImplementation"`
	InputDefinitions ParameterDefinitions     `read:"inputs,ParameterDefinition"`
	Outputs          OutputMappings           `read:"outputs,OutputMapping"` // introduced in TOSCA 1.3
}

func NewOperationDefinition(context *parsing.Context) *OperationDefinition {
	return &OperationDefinition{
		Entity:           NewEntity(context),
		Name:             context.Name,
		InputDefinitions: make(ParameterDefinitions),
		Outputs:          make(OutputMappings),
	}
}

// ([parsing.Reader] signature)
func ReadOperationDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewOperationDefinition(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.Implementation = ReadInterfaceImplementation(context.FieldChild("implementation", context.Data)).(*InterfaceImplementation)
	}

	return self
}

// ([parsing.Mappable] interface)
func (self *OperationDefinition) GetKey() string {
	return self.Name
}

func (self *OperationDefinition) Inherit(parentDefinition *OperationDefinition) {
	logInherit.Debugf("operation definition: %s", self.Name)

	if (self.Description == nil) && (parentDefinition.Description != nil) {
		self.Description = parentDefinition.Description
	}

	self.InputDefinitions.Inherit(parentDefinition.InputDefinitions)
	self.Outputs.Inherit(parentDefinition.Outputs)
}

func (self *OperationDefinition) Normalize(normalOperation *normal.Operation) {
	if self.Description != nil {
		normalOperation.Description = *self.Description
	}

	if self.Implementation != nil {
		self.Implementation.NormalizeOperation(normalOperation)
	}

	// TODO: input definitions
	//self.InputDefinitions.Normalize(o.Inputs)
}

//
// OperationDefinitions
//

type OperationDefinitions map[string]*OperationDefinition

func (self OperationDefinitions) Inherit(parentDefinitions OperationDefinitions) {
	for name, definition := range parentDefinitions {
		if _, ok := self[name]; !ok {
			self[name] = definition
		}
	}

	for name, definition := range self {
		if parentDefinition, ok := parentDefinitions[name]; ok {
			if definition != parentDefinition {
				definition.Inherit(parentDefinition)
			}
		}
	}
}
