package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Schema
//

type Schema struct {
	*Entity `name:"schema"`

	Metadata          Metadata          `read:"metadata,Metadata"` // introduced in TOSCA 1.3
	Description       *string           `read:"description"`
	DataTypeName      *string           `read:"type" mandatory:""`
	ConstraintClauses ConstraintClauses `read:"constraints,[]ConstraintClause" traverse:"ignore"`
	KeySchema         *Schema           `read:"key_schema,Schema"`   // introduced in TOSCA 1.3
	EntrySchema       *Schema           `read:"entry_schema,Schema"` // mandatory if list or map

	DataType *DataType `lookup:"type,DataTypeName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewSchema(context *parsing.Context) *Schema {
	return &Schema{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadSchema(context *parsing.Context) parsing.EntityPtr {
	self := NewSchema(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.DataTypeName = context.FieldChild("type", context.Data).ReadString()
	}

	return self
}

// ([DataDefinition] interface)
func (self *Schema) ToValueMeta() *normal.ValueMeta {
	return nil
}

// ([DataDefinition] interface)
func (self *Schema) GetDescription() string {
	if self.Description != nil {
		return *self.Description
	} else {
		return ""
	}
}

// ([DataDefinition] interface)
func (self *Schema) GetTypeMetadata() Metadata {
	return self.Metadata
}

// ([DataDefinition] interface)
func (self *Schema) GetConstraintClauses() ConstraintClauses {
	return self.ConstraintClauses
}

// ([DataDefinition] interface)
func (self *Schema) GetKeySchema() *Schema {
	return self.KeySchema
}

// ([DataDefinition] interface)
func (self *Schema) GetEntrySchema() *Schema {
	return self.EntrySchema
}

// ([parsing.Renderable] interface)
func (self *Schema) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *Schema) render() {
	logRender.Debug("schema")

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

// TODO: do we need this?
func (self *Schema) LookupDataType() bool {
	if self.DataTypeName != nil {
		dataTypeName := *self.DataTypeName
		var ok bool
		if self.DataType, ok = LookupDataType(self.Context, dataTypeName); ok {
			return true
		} else {
			self.Context.ReportMissingEntrySchema(dataTypeName)
		}
	}

	return false
}

func (self *Schema) GetConstraints() ConstraintClauses {
	if self.DataType != nil {
		constraints := self.DataType.ConstraintClauses.Append(self.ConstraintClauses)
		for _, constraint := range constraints {
			constraint.DataType = self.DataType
		}
		return constraints
	} else {
		return self.ConstraintClauses.Append(nil)
	}
}
