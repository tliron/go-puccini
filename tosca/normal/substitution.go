package normal

//
// Substitution
//

type Substitution struct {
	ServiceTemplate *ServiceTemplate `json:"-" yaml:"-"`

	Type                string            `json:"type" yaml:"type"`
	TypeMetadata        map[string]string `json:"typeMetadata" yaml:"typeMetadata"`
	CapabilityMappings  Mappings          `json:"capabilityMappings" yaml:"capabilityMappings"`
	RequirementMappings Mappings          `json:"requirementMappings" yaml:"requirementMappings"`
	PropertyMappings    Mappings          `json:"propertyMappings" yaml:"propertyMappings"`
	AttributeMappings   Mappings          `json:"attributeMappings" yaml:"attributeMappings"`
	InterfaceMappings   Mappings          `json:"interfaceMappings" yaml:"interfaceMappings"`
}

func (self *ServiceTemplate) NewSubstitution() *Substitution {
	substitutionMappings := &Substitution{
		ServiceTemplate:     self,
		TypeMetadata:        make(map[string]string),
		CapabilityMappings:  make(Mappings),
		RequirementMappings: make(Mappings),
		PropertyMappings:    make(Mappings),
		AttributeMappings:   make(Mappings),
		InterfaceMappings:   make(Mappings),
	}
	self.Substitution = substitutionMappings
	return substitutionMappings
}
