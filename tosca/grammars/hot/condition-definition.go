package hot

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

//
// ConditionDefinition
//
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#conditions-section]
//

type ConditionDefinition struct {
	*Entity `name:"condition definition"`
	Name    string `namespace:""`

	Operator              string
	Arguments             ard.List
	NativeArgumentIndexes []int
}

func NewConditionDefinition(context *parsing.Context) *ConditionDefinition {
	return &ConditionDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadConditionDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewConditionDefinition(context)

	if context.ValidateType(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if len(map_) != 1 {
			context.ReportValueMalformed("condition definition", "map length not 1")
			return self
		}

		for key, value := range map_ {
			operator := yamlkeys.KeyString(key)

			script, ok := context.ScriptletNamespace.Lookup(operator)
			if !ok {
				context.Clone(operator).ReportValueMalformed("condition definition", "unsupported operator")
				return self
			}

			self.Operator = operator

			if list, ok := value.(ard.List); ok {
				self.Arguments = list
			} else {
				self.Arguments = ard.List{value}
			}

			self.NativeArgumentIndexes = script.NativeArgumentIndexes

			// We have only one key
			break
		}
	}

	return self
}

//
// ConditionDefinitions
//

type ConditionDefinitions []*ConditionDefinition
