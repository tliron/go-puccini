package normal

import (
	"github.com/tliron/go-kutil/reflection"
	"github.com/tliron/go-puccini/tosca/parsing"
)

//
// ServiceTemplate
//

type ServiceTemplate struct {
	Description        string                      `json:"description" yaml:"description"`
	NodeTemplates      NodeTemplates               `json:"nodeTemplates" yaml:"nodeTemplates"`
	Groups             Groups                      `json:"groups" yaml:"groups"`
	Policies           Policies                    `json:"policies" yaml:"policies"`
	Inputs             Values                      `json:"inputs" yaml:"inputs"`
	Outputs            Values                      `json:"outputs" yaml:"outputs"`
	Workflows          Workflows                   `json:"workflows" yaml:"workflows"`
	Substitution       *Substitution               `json:"substitution" yaml:"substitution"`
	Metadata           map[string]string           `json:"metadata" yaml:"metadata"`
	ScriptletNamespace *parsing.ScriptletNamespace `json:"scriptletNamespace" yaml:"scriptletNamespace"`
}

func NewServiceTemplate() *ServiceTemplate {
	return &ServiceTemplate{
		NodeTemplates:      make(NodeTemplates),
		Groups:             make(Groups),
		Policies:           make(Policies),
		Inputs:             make(Values),
		Outputs:            make(Values),
		Workflows:          make(Workflows),
		Metadata:           make(map[string]string),
		ScriptletNamespace: parsing.NewScriptletNamespace(),
	}
}

//
// Normalizable
//

type Normalizable interface {
	NormalizeServiceTemplate() *ServiceTemplate
}

// From Normalizable interface
func NormalizeServiceTemplate(entityPtr parsing.EntityPtr) (*ServiceTemplate, bool) {
	var serviceTemplate *ServiceTemplate

	reflection.TraverseEntities(entityPtr, false, func(entityPtr parsing.EntityPtr) bool {
		if normalizable, ok := entityPtr.(Normalizable); ok {
			serviceTemplate = normalizable.NormalizeServiceTemplate()

			// Only one entity should implement the interface
			return false
		} else {
			return true
		}
	})

	if serviceTemplate != nil {
		return serviceTemplate, true
	} else {
		return nil, false
	}
}
