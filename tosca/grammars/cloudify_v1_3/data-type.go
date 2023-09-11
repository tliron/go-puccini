package cloudify_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// DataType
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-data-types/]
//

type DataType struct {
	*Type `name:"data type"`

	Description         *string             `read:"description"`
	PropertyDefinitions PropertyDefinitions `read:"properties,PropertyDefinition" inherit:"properties,Parent"`

	Parent *DataType `lookup:"derived_from,ParentName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewDataType(context *parsing.Context) *DataType {
	return &DataType{
		Type:                NewType(context),
		PropertyDefinitions: make(PropertyDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadDataType(context *parsing.Context) parsing.EntityPtr {
	self := NewDataType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Hierarchical] interface)
func (self *DataType) GetParent() parsing.EntityPtr {
	return self.Parent
}

// ([parsing.Inherits] interface)
func (self *DataType) Inherit() {
	logInherit.Debugf("data type: %s", self.Name)

	if _, ok := self.GetInternalTypeName(); ok && (len(self.PropertyDefinitions) > 0) {
		// Doesn't make sense to be an internal type (non-complex) and also have properties (complex)
		self.Context.ReportPrimitiveType()
		self.PropertyDefinitions = make(PropertyDefinitions)
		return
	}

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
}

// ([parsing.Renderable] interface)
func (self *DataType) Render() {
	self.renderOnce.Do(self.render)
}

func (self *DataType) render() {
	logRender.Debugf("data type: %s", self.Name)

	if internalTypeName, ok := self.GetInternalTypeName(); ok {
		if _, ok := ard.TypeValidators[internalTypeName]; !ok {
			if _, ok := self.Context.Grammar.Readers[string(internalTypeName)]; !ok {
				self.Context.ReportUnsupportedType()
			}
		}
	}
}

func (self *DataType) GetInternalTypeName() (ard.TypeName, bool) {
	switch self.Name {
	case "boolean":
		return ard.TypeBoolean, true
	case "integer":
		return ard.TypeInteger, true
	case "float":
		return ard.TypeFloat, true
	case "string":
		return ard.TypeString, true
	case "list":
		return ard.TypeList, true
	case "dict":
		return ard.TypeMap, true
	default:
		return ard.NoType, false
	}
}

func (self *DataType) GetInternal() (ard.TypeName, ard.TypeValidator, parsing.Reader, bool) {
	if internalTypeName, ok := self.GetInternalTypeName(); ok {
		if typeValidator, ok := ard.TypeValidators[internalTypeName]; ok {
			return internalTypeName, typeValidator, nil, true
		} else if reader, ok := self.Context.Grammar.Readers[string(internalTypeName)]; ok {
			return internalTypeName, nil, reader, true
		}
	}
	return ard.NoType, nil, nil, false
}

func (self *DataType) NewValueMeta() *normal.ValueMeta {
	normalValueMeta := normal.NewValueMeta()
	normalValueMeta.Type = parsing.GetCanonicalName(self)
	if self.Description != nil {
		normalValueMeta.TypeDescription = *self.Description
	}
	return normalValueMeta
}

//
// DataTypes
//

type DataTypes []*DataType
