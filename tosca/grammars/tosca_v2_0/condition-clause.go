package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

//
// ConditionClause
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.25
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.21
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.19
//

type ConditionClause struct {
	*Entity `name:"condition clause"`

	// Either an assertion definition (attribute with constraints)
	AttributeName     *string
	ConstraintClauses ConstraintClauses

	// Or one or more child condition clauses
	Operator         *string
	ConditionClauses []*ConditionClause
}

func NewConditionClause(context *parsing.Context) *ConditionClause {
	return &ConditionClause{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadConditionClause(context *parsing.Context) parsing.EntityPtr {
	self := NewConditionClause(context)

	if context.ValidateType(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if len(map_) != 1 {
			context.ReportValueMalformed("condition clause", "map length not 1")
			return self
		}

		for key, value := range map_ {
			name := yamlkeys.KeyString(key)

			if name == "assert" {
				// deprecated in TOSCA 1.3
				name = "and"
			}

			switch name {
			case "and":
				self.Operator = &name
				context.Clone(value).ReadListItems(ReadConditionClause, func(item ard.Value) {
					self.ConditionClauses = append(self.ConditionClauses, item.(*ConditionClause))
				})

			case "or":
				self.Operator = &name
				context.Clone(value).ReadListItems(ReadConditionClause, func(item ard.Value) {
					self.ConditionClauses = append(self.ConditionClauses, item.(*ConditionClause))
				})

			case "not": // introduced in TOSCA 1.3
				self.Operator = &name
				context.Clone(value).ReadListItems(ReadConditionClause, func(item ard.Value) {
					self.ConditionClauses = append(self.ConditionClauses, item.(*ConditionClause))
				})
				if len(self.ConditionClauses) != 1 {
					context.ReportValueMalformed("condition clause", "\"not\" does not have one and only one clause")
				}

			default:
				// Assertion definition
				self.AttributeName = &name
				context.Clone(value).ReadListItems(ReadConstraintClause, func(item ard.Value) {
					self.ConstraintClauses = append(self.ConstraintClauses, item.(*ConstraintClause))
				})
			}

			// We have only one key
			break
		}
	}

	return self
}

// ([parsing.Reader] signature)
func ReadConditionClauseAnd(context *parsing.Context) parsing.EntityPtr {
	self := NewConditionClause(context)

	if context.ValidateType(ard.TypeList) {
		name := "and"
		self.Operator = &name
		context.ReadListItems(ReadConditionClause, func(item ard.Value) {
			self.ConditionClauses = append(self.ConditionClauses, item.(*ConditionClause))
		})
	}

	return self
}
