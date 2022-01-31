package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// PropertyFilter
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.4
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.4
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.3
//

type PropertyFilter struct {
	*Entity `name:"property filter"`
	Name    string

	ConstraintClauses ConstraintClauses
}

func NewPropertyFilter(context *tosca.Context) *PropertyFilter {
	return &PropertyFilter{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadPropertyFilter(context *tosca.Context) tosca.EntityPtr {
	self := NewPropertyFilter(context)

	if context.Is(ard.TypeList) {
		context.ReadListItems(ReadConstraintClause, func(item ard.Value) {
			self.ConstraintClauses = append(self.ConstraintClauses, item.(*ConstraintClause))
		})
	} else {
		self.ConstraintClauses = ConstraintClauses{ReadConstraintClause(context).(*ConstraintClause)}
	}

	return self
}

func (self *PropertyFilter) Normalize(normalFunctionCallMap normal.FunctionCallMap) normal.FunctionCalls {
	if len(self.ConstraintClauses) == 0 {
		return nil
	}

	normalFunctionCalls := self.ConstraintClauses.Normalize(self.Context)
	if existing, ok := normalFunctionCallMap[self.Name]; ok {
		normalFunctionCallMap[self.Name] = append(existing, normalFunctionCalls...)
	} else {
		normalFunctionCallMap[self.Name] = normalFunctionCalls
	}
	return normalFunctionCalls
}

//
// PropertyFilters
//

type PropertyFilters []*PropertyFilter

func (self PropertyFilters) Normalize(normalFunctionCallMap normal.FunctionCallMap) {
	for _, propertyFilter := range self {
		lock := propertyFilter.GetEntityLock()
		lock.RLock()
		propertyFilter.Normalize(normalFunctionCallMap)
		lock.RUnlock()
	}
}
