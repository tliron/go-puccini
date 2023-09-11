package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/assets/tosca/profiles"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

const constraintPathPrefix = "implicit/2.0/js/constraints/"

// Built-in constraint functions
var ConstraintClauseScriptlets = map[string]string{
	parsing.MetadataContraintPrefix + "equal":            profiles.GetString(constraintPathPrefix + "equal.js"),
	parsing.MetadataContraintPrefix + "greater_than":     profiles.GetString(constraintPathPrefix + "greater_than.js"),
	parsing.MetadataContraintPrefix + "greater_or_equal": profiles.GetString(constraintPathPrefix + "greater_or_equal.js"),
	parsing.MetadataContraintPrefix + "less_than":        profiles.GetString(constraintPathPrefix + "less_than.js"),
	parsing.MetadataContraintPrefix + "less_or_equal":    profiles.GetString(constraintPathPrefix + "less_or_equal.js"),
	parsing.MetadataContraintPrefix + "in_range":         profiles.GetString(constraintPathPrefix + "in_range.js"),
	parsing.MetadataContraintPrefix + "valid_values":     profiles.GetString(constraintPathPrefix + "valid_values.js"),
	parsing.MetadataContraintPrefix + "length":           profiles.GetString(constraintPathPrefix + "length.js"),
	parsing.MetadataContraintPrefix + "min_length":       profiles.GetString(constraintPathPrefix + "min_length.js"),
	parsing.MetadataContraintPrefix + "max_length":       profiles.GetString(constraintPathPrefix + "max_length.js"),
	parsing.MetadataContraintPrefix + "pattern":          profiles.GetString(constraintPathPrefix + "pattern.js"),
	parsing.MetadataContraintPrefix + "schema":           profiles.GetString(constraintPathPrefix + "schema.js"), // introduced in TOSCA 1.3
}

var ConstraintClauseNativeArgumentIndexes = map[string][]int{
	parsing.MetadataContraintPrefix + "equal":            {0},
	parsing.MetadataContraintPrefix + "greater_than":     {0},
	parsing.MetadataContraintPrefix + "greater_or_equal": {0},
	parsing.MetadataContraintPrefix + "less_than":        {0},
	parsing.MetadataContraintPrefix + "less_or_equal":    {0},
	parsing.MetadataContraintPrefix + "in_range":         {0, 1},
	parsing.MetadataContraintPrefix + "valid_values":     {-1}, // -1 means all
}

//
// ConstraintClause
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2
//

type ConstraintClause struct {
	*Entity `name:"constraint clause"`

	Operator              string
	Arguments             ard.List
	NativeArgumentIndexes []int
	DataType              *DataType      `traverse:"ignore" json:"-" yaml:"-"` // TODO: unncessary, this entity should never be traversed
	Definition            DataDefinition `traverse:"ignore" json:"-" yaml:"-"`
}

func NewConstraintClause(context *parsing.Context) *ConstraintClause {
	return &ConstraintClause{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadConstraintClause(context *parsing.Context) parsing.EntityPtr {
	self := NewConstraintClause(context)

	if context.ValidateType(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if len(map_) != 1 {
			context.ReportValueMalformed("constraint clause", "map size not 1")
			return self
		}

		for key, value := range map_ {
			operator := yamlkeys.KeyString(key)

			scriptletName := parsing.MetadataContraintPrefix + operator
			scriptlet, ok := context.ScriptletNamespace.Lookup(scriptletName)
			if !ok {
				context.Clone(operator).ReportValueMalformed("constraint clause", "unsupported operator")
				return self
			}

			self.Operator = operator

			if list, ok := value.(ard.List); ok {
				self.Arguments = list
			} else {
				self.Arguments = ard.List{value}
			}

			self.NativeArgumentIndexes = scriptlet.NativeArgumentIndexes

			// We have only one key
			break
		}
	}

	return self
}

func (self *ConstraintClause) ToFunctionCall(context *parsing.Context, strict bool) *parsing.FunctionCall {
	// Special case: "in_range" for a "range" accepts a range (two integers) rather than two ranges
	var isRangeInRange bool
	if self.Operator == "in_range" {
		if rangeType, ok := self.Context.Namespace.Lookup("range"); ok {
			if isRangeInRange = self.Context.Hierarchy.IsCompatible(rangeType, self.DataType); isRangeInRange {
				range_ := ReadRange(context.Clone(self.Arguments)).(*Range)
				return context.NewFunctionCall(parsing.MetadataContraintPrefix+self.Operator, []any{
					ReadValue(context.ListChild(0, range_.Lower)).(*Value),
					ReadValue(context.ListChild(1, range_.Upper)).(*Value),
				})
			}
		}
	}

	arguments := make([]any, len(self.Arguments))
	for index, argument := range self.Arguments {
		if self.IsNativeArgument(index) {
			if _, ok := argument.(*Value); !ok {
				if self.DataType != nil {
					value := ReadValue(self.Context.ListChild(index, argument)).(*Value)
					value.Render(self.DataType, self.Definition, true, false) // bare
					argument = value
				} else if strict {
					panic("no data type for native argument")
				}
			}
		}

		arguments[index] = argument
	}

	return context.NewFunctionCall(parsing.MetadataContraintPrefix+self.Operator, arguments)
}

func (self *ConstraintClause) IsNativeArgument(index int) bool {
	for _, i := range self.NativeArgumentIndexes {
		if (i == -1) || (i == index) {
			return true
		}
	}
	return false
}

//
// ConstraintClauses
//

type ConstraintClauses []*ConstraintClause

func (self ConstraintClauses) Append(constraints ConstraintClauses) ConstraintClauses {
	r := append(self[:0:0], self...)
	return append(r, constraints...)
}

func (self ConstraintClauses) Normalize(context *parsing.Context) normal.FunctionCalls {
	var normalFunctionCalls normal.FunctionCalls
	for _, constraintClause := range self {
		functionCall := constraintClause.ToFunctionCall(context, false)
		NormalizeFunctionCallArguments(functionCall, context)
		normalFunctionCalls = append(normalFunctionCalls, normal.NewFunctionCall(functionCall))
	}
	return normalFunctionCalls
}

func (self ConstraintClauses) AddToMeta(context *parsing.Context, normalValueMeta *normal.ValueMeta) {
	for _, constraintClause := range self {
		functionCall := constraintClause.ToFunctionCall(context, true)
		NormalizeFunctionCallArguments(functionCall, context)
		normalValueMeta.AddValidator(functionCall)
	}
}
