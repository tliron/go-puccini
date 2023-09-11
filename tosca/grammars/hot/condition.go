package hot

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

//
// Condition
//

type Condition struct {
	*Entity `name:"condition"`

	Operator              string
	Arguments             ard.List
	NativeArgumentIndexes []int

	DefinitionName *string

	Value *bool
}

func NewCondition(context *parsing.Context) *Condition {
	return &Condition{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadCondition(context *parsing.Context) parsing.EntityPtr {
	self := NewCondition(context)

	if context.Is(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if len(map_) != 1 {
			context.ReportValueMalformed("condition", "map length not 1")
			return self
		}

		for key, value := range map_ {
			operator := yamlkeys.KeyString(key)

			script, ok := context.ScriptletNamespace.Lookup(operator)
			if !ok {
				context.Clone(operator).ReportValueMalformed("condition", "unsupported operator")
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
	} else if context.Is(ard.TypeString) {
		self.DefinitionName = context.ReadString()
	} else if context.ValidateType(ard.TypeMap, ard.TypeString, ard.TypeBoolean) {
		self.Value = context.ReadBoolean()
	}

	return self
}
