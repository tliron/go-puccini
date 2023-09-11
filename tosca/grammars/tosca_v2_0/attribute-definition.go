package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// AttributeDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.12
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.10
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.10
//

type AttributeDefinition struct {
	*Entity `name:"attribute definition"`
	Name    string

	Metadata          Metadata          `read:"metadata,Metadata"` // introduced in TOSCA 1.3
	Description       *string           `read:"description"`
	DataTypeName      *string           `read:"type"`                                             // mandatory only if cannot be inherited or discovered
	ConstraintClauses ConstraintClauses `read:"constraints,[]ConstraintClause" traverse:"ignore"` // introduced in TOSCA 2.0
	KeySchema         *Schema           `read:"key_schema,Schema"`                                // introduced in TOSCA 1.3
	EntrySchema       *Schema           `read:"entry_schema,Schema"`                              // mandatory if list or map
	Default           *Value            `read:"default,Value"`
	Status            *string           `read:"status"`

	DataType *DataType `lookup:"type,DataTypeName" traverse:"ignore" json:"-" yaml:"-"`

	looseType bool
}

func NewAttributeDefinition(context *parsing.Context) *AttributeDefinition {
	return &AttributeDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadAttributeDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewAttributeDefinition(context)
	var ignore []string
	if context.HasQuirk(parsing.QuirkAnnotationsIgnore) {
		ignore = append(ignore, "annotations")
	}
	context.ValidateUnsupportedFields(append(context.ReadFields(self), ignore...))
	return self
}

// ([parsing.Mappable] interface)
func (self *AttributeDefinition) GetKey() string {
	return self.Name
}

// ([DataDefinition] interface)
func (self *AttributeDefinition) ToValueMeta() *normal.ValueMeta {
	information := normal.NewValueMeta()
	information.Metadata = parsing.GetDataTypeMetadata(self.Metadata)
	if self.Description != nil {
		information.TypeDescription = *self.Description
	}
	return information
}

// ([DataDefinition] interface)
func (self *AttributeDefinition) GetDescription() string {
	if self.Description != nil {
		return *self.Description
	} else {
		return ""
	}
}

// ([DataDefinition] interface)
func (self *AttributeDefinition) GetTypeMetadata() Metadata {
	return self.Metadata
}

// ([DataDefinition] interface)
func (self *AttributeDefinition) GetConstraintClauses() ConstraintClauses {
	return self.ConstraintClauses
}

// ([DataDefinition] interface)
func (self *AttributeDefinition) GetKeySchema() *Schema {
	return self.KeySchema
}

// ([DataDefinition] interface)
func (self *AttributeDefinition) GetEntrySchema() *Schema {
	return self.EntrySchema
}

func (self *AttributeDefinition) Inherit(parentDefinition *AttributeDefinition) {
	logInherit.Debugf("attribute definition: %s", self.Name)

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
	if (self.KeySchema == nil) && (parentDefinition.KeySchema != nil) {
		self.KeySchema = parentDefinition.KeySchema
	}
	if (self.EntrySchema == nil) && (parentDefinition.EntrySchema != nil) {
		self.EntrySchema = parentDefinition.EntrySchema
	}
	if (self.Default == nil) && (parentDefinition.Default != nil) {
		self.Default = parentDefinition.Default
	}
	if (self.Status == nil) && (parentDefinition.Status != nil) {
		self.Status = parentDefinition.Status
	}
	if (self.DataType == nil) && (parentDefinition.DataType != nil) {
		self.DataType = parentDefinition.DataType
	}
}

// ([parsing.Renderable] interface)
func (self *AttributeDefinition) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *AttributeDefinition) render() {
	logRender.Debugf("attribute definition: %s", self.Name)

	self.doRender()

	if (self.Default != nil) && (self.DataType != nil) {
		// The "default" value must be a valid value of the type
		self.Default.Render(self.DataType, self, false, false)
	}
}

func (self *AttributeDefinition) doRender() {
	if !self.looseType && (self.DataTypeName == nil) {
		self.Context.FieldChild("type", nil).ReportKeynameMissing()
		return
	}

	if self.DataType == nil {
		return
	}

	if internalTypeName, ok := self.DataType.GetInternalTypeName(); ok {
		switch internalTypeName {
		case ard.TypeList, ard.TypeMap:
			if self.EntrySchema == nil {
				self.EntrySchema = self.DataType.EntrySchema
			}

			// Make sure we have an entry schema
			if (self.EntrySchema == nil) || (self.EntrySchema.DataType == nil) {
				self.Context.ReportMissingEntrySchema(self.DataType.Name)
				return
			}

			if internalTypeName == ard.TypeMap {
				if self.KeySchema == nil {
					self.KeySchema = self.DataType.KeySchema
				}

				if self.KeySchema == nil {
					// Default to "string" for key schema
					self.KeySchema = ReadSchema(self.Context.FieldChild("key_schema", "string")).(*Schema)
					if !self.KeySchema.LookupDataType() {
						panic("missing \"string\" type")
					}
				}
			}
		}
	}
}

//
// AttributeDefinitions
//

type AttributeDefinitions map[string]*AttributeDefinition

func (self AttributeDefinitions) Inherit(parentDefinitions AttributeDefinitions) {
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
