package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// InterfaceDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.20
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.16
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.14
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.14
//

type InterfaceDefinition struct {
	*Entity `name:"interface definition" json:"-" yaml:"-"`
	Name    string

	InterfaceTypeName         *string                 `read:"type"` // mandatory only if cannot be inherited
	InputDefinitions          ParameterDefinitions    `read:"inputs,ParameterDefinition" inherit:"inputs,InterfaceType"`
	OperationDefinitions      OperationDefinitions    `read:"operations,OperationDefinition" inherit:"operations,InterfaceType"`          // keyword since TOSCA 1.3
	NotificationDefinitions   NotificationDefinitions `read:"notifications,NotificationDefinition" inherit:"notifications,InterfaceType"` // introduced in TOSCA 1.3
	ExtraOperationDefinitions OperationDefinitions    `json:"-" yaml:"-"`

	InterfaceType *InterfaceType `lookup:"type,InterfaceTypeName" traverse:"ignore" json:"-" yaml:"-"`

	typeMissingProblemReported bool
}

func NewInterfaceDefinition(context *parsing.Context) *InterfaceDefinition {
	return &InterfaceDefinition{
		Entity:                    NewEntity(context),
		Name:                      context.Name,
		InputDefinitions:          make(ParameterDefinitions),
		OperationDefinitions:      make(OperationDefinitions),
		NotificationDefinitions:   make(NotificationDefinitions),
		ExtraOperationDefinitions: make(OperationDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadInterfaceDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewInterfaceDefinition(context)

	if context.HasQuirk(parsing.QuirkInterfacesOperationsPermissive) {
		context.SetReadTag("ExtraOperationDefinitions", "?,OperationDefinition")
		context.ReadFields(self)
		for name, definition := range self.ExtraOperationDefinitions {
			self.OperationDefinitions[name] = definition
		}
	} else {
		context.ValidateUnsupportedFields(context.ReadFields(self))
	}

	return self
}

// ([parsing.Mappable] interface)
func (self *InterfaceDefinition) GetKey() string {
	return self.Name
}

func (self *InterfaceDefinition) Inherit(parentDefinition *InterfaceDefinition) {
	logInherit.Debugf("interface definition: %s", self.Name)

	// Validate type compatibility
	if (self.InterfaceType != nil) && (parentDefinition.InterfaceType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.InterfaceType, self.InterfaceType) {
		self.Context.ReportIncompatibleType(self.InterfaceType, parentDefinition.InterfaceType)
		return
	}

	if (self.InterfaceTypeName == nil) && (parentDefinition.InterfaceTypeName != nil) {
		self.InterfaceTypeName = parentDefinition.InterfaceTypeName
	}
	if (self.InterfaceType == nil) && (parentDefinition.InterfaceType != nil) {
		self.InterfaceType = parentDefinition.InterfaceType
	}

	self.InputDefinitions.Inherit(parentDefinition.InputDefinitions)
	self.OperationDefinitions.Inherit(parentDefinition.OperationDefinitions)
	self.NotificationDefinitions.Inherit(parentDefinition.NotificationDefinitions)
}

// ([parsing.Renderable] interface)
func (self *InterfaceDefinition) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *InterfaceDefinition) render() {
	logRender.Debugf("interface definition: %s", self.Name)

	if self.InterfaceTypeName == nil {
		// Avoid reporting more than once
		if !self.typeMissingProblemReported {
			self.Context.FieldChild("type", nil).ReportKeynameMissing()
			self.typeMissingProblemReported = true
		}
	}
}

//
// InterfaceDefinitions
//

type InterfaceDefinitions map[string]*InterfaceDefinition

func (self InterfaceDefinitions) Inherit(parentDefinitions InterfaceDefinitions) {
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
