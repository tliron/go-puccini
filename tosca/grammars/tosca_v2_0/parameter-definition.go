package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// ParameterDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.14
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.13
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.12
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.12
//

type ParameterDefinition struct {
	*PropertyDefinition `name:"parameter definition"`

	Value *Value `read:"value,Value"`
}

func NewParameterDefinition(context *tosca.Context) *ParameterDefinition {
	return &ParameterDefinition{PropertyDefinition: NewPropertyDefinition(context)}
}

// tosca.Reader signature
func ReadParameterDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewParameterDefinition(context)
	var ignore []string
	if context.HasQuirk(tosca.QuirkAnnotationsIgnore) {
		ignore = append(ignore, "annotations")
	}
	context.ValidateUnsupportedFields(append(context.ReadFields(self), ignore...))
	return self
}

func (self *ParameterDefinition) Render(kind string, mapped []string) {
	logRender.Debugf("parameter definition: %s", self.Name)

	if self.DataTypeName == nil {
		self.Context.FieldChild("type", nil).ReportFieldMissing()
	}

	if self.Value == nil {
		self.Value = self.Default
	}

	if self.Value == nil {
		isMapped := false
		for _, mapped_ := range mapped {
			if self.Name == mapped_ {
				isMapped = true
				break
			}
		}

		if !isMapped && self.IsRequired() {
			self.Context.ReportPropertyRequired(kind)
			return
		}
	} else if self.DataType != nil {
		lock := self.DataType.GetEntityLock()
		lock.RLock()
		defer lock.RUnlock()

		self.Value.RenderProperty(self.DataType, self.PropertyDefinition)
	}
}

func (self *ParameterDefinition) Normalize(context *tosca.Context) normal.Constrainable {
	var value *Value
	if self.Value != nil {
		value = self.Value
	} else {
		// Parameters should always appear, even if they have no default value
		value = NewValue(context.MapChild(self.Name, nil))
	}
	lock := value.GetEntityLock()
	lock.RLock()
	defer lock.RUnlock()
	return value.Normalize()
}

//
// ParameterDefinitions
//

type ParameterDefinitions map[string]*ParameterDefinition

func (self ParameterDefinitions) Render(kind string, mapped []string, context *tosca.Context) {
	for _, definition := range self {
		lock := definition.GetEntityLock()
		lock.Lock()
		definition.Render(kind, mapped)
		lock.Unlock()
	}
}

func (self ParameterDefinitions) Normalize(c normal.Constrainables, context *tosca.Context) {
	for key, definition := range self {
		lock := definition.GetEntityLock()
		lock.RLock()
		c[key] = definition.Normalize(context)
		lock.RUnlock()
	}
}
