package cloudify_v1_3

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// ParameterDefinition
//

type ParameterDefinition struct {
	*Entity `name:"parameter definition"`
	Name    string

	Description  *string `read:"description"`
	DataTypeName *string `read:"type"`
	Default      *Value  `read:"default,Value"`

	DataType *DataType `lookup:"type,DataTypeName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewParameterDefinition(context *parsing.Context) *ParameterDefinition {
	return &ParameterDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadParameterDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewParameterDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Mappable] interface)
func (self *ParameterDefinition) GetKey() string {
	return self.Name
}

func (self *ParameterDefinition) Inherit(parentDefinition *ParameterDefinition) {
	logInherit.Debugf("parameter definition: %s", self.Name)

	// Validate type compatibility
	if (self.DataType != nil) && (parentDefinition.DataType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.DataType, self.DataType) {
		self.Context.ReportIncompatibleType(self.DataType, parentDefinition.DataType)
		return
	}

	if (self.Description == nil) && (parentDefinition.Description != nil) {
		self.Description = parentDefinition.Description
	}
	if (self.DataTypeName == nil) && (parentDefinition.DataTypeName != nil) {
		self.DataTypeName = parentDefinition.DataTypeName
	}
	if (self.Default == nil) && (parentDefinition.Default != nil) {
		self.Default = parentDefinition.Default
	}
	if (self.DataType == nil) && (parentDefinition.DataType != nil) {
		self.DataType = parentDefinition.DataType
	}
}

// ([parsing.Renderable] interface)
func (self *ParameterDefinition) Render() {
	self.renderOnce.Do(self.render)
}

func (self *ParameterDefinition) render() {
	logRender.Debugf("parameter definition: %s", self.Name)

	if self.DataType == nil {
		return
	}

	if self.Default != nil {
		// The "default" value must be a valid value of the type
		self.Default.RenderParameter(self.DataType, self, false, false)
	}
}

func (self *ParameterDefinition) GetNormalDataType() *normal.ValueMeta {
	normalDataType := normal.NewValueMeta()
	if self.Description != nil {
		normalDataType.TypeDescription = *self.Description
	}
	return normalDataType
}

//
// ParameterDefinitions
//

type ParameterDefinitions map[string]*ParameterDefinition

func (self ParameterDefinitions) Inherit(parentDefinitions ParameterDefinitions) {
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
