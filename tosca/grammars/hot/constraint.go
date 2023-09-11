package hot

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/assets/tosca/profiles"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

const constraintPathPrefix = "hot/1.0/js/constraints/"

// Built-in constraint functions
var ConstraintScriptlets = map[string]string{
	parsing.MetadataContraintPrefix + "length":            profiles.GetString(constraintPathPrefix + "length.js"),
	parsing.MetadataContraintPrefix + "range":             profiles.GetString(constraintPathPrefix + "range.js"),
	parsing.MetadataContraintPrefix + "modulo":            profiles.GetString(constraintPathPrefix + "modulo.js"),
	parsing.MetadataContraintPrefix + "allowed_values":    profiles.GetString(constraintPathPrefix + "allowed_values.js"),
	parsing.MetadataContraintPrefix + "allowed_pattern":   profiles.GetString(constraintPathPrefix + "allowed_pattern.js"),
	parsing.MetadataContraintPrefix + "custom_constraint": profiles.GetString(constraintPathPrefix + "custom_constraint.js"),
}

var ConstraintNativeArgumentIndexes = map[string][]int{}

//
// Constraint
//
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#parameter-constraints]
//

type Constraint struct {
	*Entity `name:"constraint"`

	Description *string
	Operator    string
	Arguments   ard.List
}

func NewConstraint(context *parsing.Context) *Constraint {
	return &Constraint{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadConstraint(context *parsing.Context) parsing.EntityPtr {
	self := NewConstraint(context)

	if context.ValidateType(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		var length = len(map_)
		if (length != 1) && (length != 2) {
			context.ReportValueMalformed("constraint", "map size not 1 or 2")
			return self
		}

		for key, value := range map_ {
			operator := yamlkeys.KeyString(key)

			if operator == "description" {
				self.Description = context.FieldChild(operator, value).ReadString()
				continue
			}

			scriptletName := parsing.MetadataContraintPrefix + operator
			if _, ok := context.ScriptletNamespace.Lookup(scriptletName); !ok {
				context.Clone(operator).ReportValueMalformed("constraint", "unsupported operator")
				return self
			}

			self.Operator = operator

			if list, ok := value.(ard.List); ok {
				self.Arguments = list
			} else {
				self.Arguments = ard.List{value}
			}
		}
	}

	return self
}

func (self *Constraint) NewFunctionCall(context *parsing.Context) *parsing.FunctionCall {
	return context.NewFunctionCall(parsing.MetadataContraintPrefix+self.Operator, self.Arguments)
}

//
// Constraints
//

type Constraints []*Constraint

func (self Constraints) Normalize(context *parsing.Context, normalDataType *normal.ValueMeta) {
	for _, constraint := range self {
		functionCall := constraint.NewFunctionCall(context)
		NormalizeFunctionCallArguments(functionCall, context)
		normalDataType.AddValidator(functionCall)
		// TODO: normalize constraint description somewhere
	}
}
