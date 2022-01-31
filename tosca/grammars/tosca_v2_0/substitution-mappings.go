package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// SubstitutionMappings
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.13, 2.10, 2.11, 2.12
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.12, 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

type SubstitutionMappings struct {
	*Entity `name:"substitution mappings"`

	NodeTypeName        *string             `read:"node_type" require:""`
	CapabilityMappings  CapabilityMappings  `read:"capabilities,CapabilityMapping"`
	RequirementMappings RequirementMappings `read:"requirements,RequirementMapping"`
	PropertyMappings    PropertyMappings    `read:"properties,PropertyMapping"`     // introduced in TOSCA 1.2
	AttributeMappings   AttributeMappings   `read:"attributes,AttributeMapping"`    // introduced in TOSCA 1.3
	InterfaceMappings   InterfaceMappings   `read:"interfaces,InterfaceMapping"`    // introduced in TOSCA 1.2
	SubstitutionFilter  *NodeFilter         `read:"substitution_filter,NodeFilter"` // introduced in TOSCA 1.3

	NodeType *NodeType `lookup:"node_type,NodeTypeName" json:"-" yaml:"-"`
}

func NewSubstitutionMappings(context *tosca.Context) *SubstitutionMappings {
	return &SubstitutionMappings{
		Entity:              NewEntity(context),
		CapabilityMappings:  make(CapabilityMappings),
		RequirementMappings: make(RequirementMappings),
		PropertyMappings:    make(PropertyMappings),
		AttributeMappings:   make(AttributeMappings),
		InterfaceMappings:   make(InterfaceMappings),
	}
}

// tosca.Reader signature
func ReadSubstitutionMappings(context *tosca.Context) tosca.EntityPtr {
	if context.HasQuirk(tosca.QuirkSubstitutionMappingsRequirementsList) {
		if map_, ok := context.Data.(ard.Map); ok {
			if requirements, ok := map_["requirements"]; ok {
				if _, ok := requirements.(ard.List); ok {
					context.SetReadTag("RequirementMappings", "requirements,{}RequirementMapping")
				}
			}
		}
	}

	self := NewSubstitutionMappings(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *SubstitutionMappings) IsRequirementMapped(nodeTemplate *NodeTemplate, requirementName string) bool {
	for _, mapping := range self.RequirementMappings {
		lock := mapping.GetEntityLock()
		lock.RLock()
		if mapping.NodeTemplate == nodeTemplate {
			if (mapping.RequirementName != nil) && (*mapping.RequirementName == requirementName) {
				return true
			}
		}
		lock.RUnlock()
	}
	return false
}

func (self *SubstitutionMappings) Render(inputDefinitions ParameterDefinitions) {
	logRender.Debug("substitution mappings")

	if self.NodeType == nil {
		return
	}

	lock := self.NodeType.GetEntityLock()
	lock.RLock()
	defer lock.RUnlock()

	for name, mapping := range self.CapabilityMappings {
		lock1 := mapping.GetEntityLock()
		lock1.RLock()
		if definition, ok := self.NodeType.CapabilityDefinitions[name]; ok {
			lock2 := definition.GetEntityLock()
			lock2.RLock()
			if mappedDefinition, ok := mapping.GetCapabilityDefinition(); ok {
				lock3 := mappedDefinition.GetEntityLock()
				lock3.RLock()
				if (definition.CapabilityType != nil) && (mappedDefinition.CapabilityType != nil) {
					if !self.Context.Hierarchy.IsCompatible(definition.CapabilityType, mappedDefinition.CapabilityType) {
						self.Context.ReportIncompatibleType(definition.CapabilityType, mappedDefinition.CapabilityType)
					}
				}
				lock3.RUnlock()
			}
			lock2.RUnlock()
		} else {
			mapping.Context.Clone(name).ReportReferenceNotFound("capability", self.NodeType)
		}
		lock1.RUnlock()
	}

	for name, mapping := range self.RequirementMappings {
		if _, ok := self.NodeType.RequirementDefinitions[name]; !ok {
			mapping.Context.Clone(name).ReportReferenceNotFound("requirement", self.NodeType)
		}
	}

	self.PropertyMappings.Render(inputDefinitions)
	for name, mapping := range self.PropertyMappings {
		lock1 := mapping.GetEntityLock()
		lock1.RLock()
		if definition, ok := self.NodeType.PropertyDefinitions[name]; ok {
			definition.Render()
			lock2 := definition.GetEntityLock()
			lock2.RLock()
			if mapping.InputDefinition != nil {
				lock3 := mapping.InputDefinition.GetEntityLock()
				lock3.RLock()
				// Input mapping
				if (definition.DataType != nil) && (mapping.InputDefinition.DataType != nil) {
					if !self.Context.Hierarchy.IsCompatible(definition.DataType, mapping.InputDefinition.DataType) {
						self.Context.ReportIncompatibleType(definition.DataType, mapping.InputDefinition.DataType)
					}
				}
				lock3.RUnlock()
			} else if mapping.Property != nil {
				// Property mapping (deprecated in TOSCA 1.3)
				lock3 := mapping.Property.GetEntityLock()
				lock3.RLock()
				if definition.DataType != nil {
					if mapping.Property.DataType != nil {
						if !self.Context.Hierarchy.IsCompatible(definition.DataType, mapping.Property.DataType) {
							self.Context.ReportIncompatibleType(definition.DataType, mapping.Property.DataType)
						}
					} else {
						mapping.Property.RenderProperty(definition.DataType, definition)
					}
				}
				lock3.RUnlock()
			}
			lock2.RUnlock()
		} else {
			mapping.Context.Clone(name).ReportReferenceNotFound("property", self.NodeType)
		}
		lock1.RUnlock()
	}

	self.AttributeMappings.EnsureRender()
	for name, mapping := range self.AttributeMappings {
		lock1 := mapping.GetEntityLock()
		lock1.RLock()
		if definition, ok := self.NodeType.AttributeDefinitions[name]; ok {
			lock2 := definition.GetEntityLock()
			lock2.RLock()
			if (definition.DataType != nil) && (mapping.Attribute != nil) {
				lock3 := mapping.Attribute.GetEntityLock()
				lock3.RLock()
				if mapping.Attribute.DataType != nil {
					if !self.Context.Hierarchy.IsCompatible(definition.DataType, mapping.Attribute.DataType) {
						self.Context.ReportIncompatibleType(definition.DataType, mapping.Attribute.DataType)
					}
				}
				lock3.RUnlock()
			}
			lock2.RUnlock()
		} else {
			mapping.Context.Clone(name).ReportReferenceNotFound("attribute", self.NodeType)
		}
		lock1.RUnlock()
	}

	for name, mapping := range self.InterfaceMappings {
		lock1 := mapping.GetEntityLock()
		lock1.RLock()
		if definition, ok := self.NodeType.InterfaceDefinitions[name]; ok {
			lock2 := definition.GetEntityLock()
			lock2.RLock()
			if mappedDefinition, ok := mapping.GetInterfaceDefinition(); ok {
				lock3 := mappedDefinition.GetEntityLock()
				lock3.RLock()
				if (definition.InterfaceType != nil) && (mappedDefinition.InterfaceType != nil) {
					if !self.Context.Hierarchy.IsCompatible(definition.InterfaceType, mappedDefinition.InterfaceType) {
						self.Context.ReportIncompatibleType(definition.InterfaceType, mappedDefinition.InterfaceType)
					}
				}
				lock3.RUnlock()
			}
			lock2.RUnlock()
		} else {
			mapping.Context.Clone(name).ReportReferenceNotFound("interface", self.NodeType)
		}
		lock1.RUnlock()
	}
}

func (self *SubstitutionMappings) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.Substitution {
	logNormalize.Debug("substitution mappings")

	if self.NodeType == nil {
		return nil
	}

	lock := self.NodeType.GetEntityLock()
	lock.RLock()
	defer lock.RUnlock()

	normalSubstitution := normalServiceTemplate.NewSubstitution()

	normalSubstitution.Type = tosca.GetCanonicalName(self.NodeType)

	if metadata, ok := self.NodeType.GetMetadata(); ok {
		normalSubstitution.TypeMetadata = metadata
	}

	for _, mapping := range self.CapabilityMappings {
		lock1 := mapping.GetEntityLock()
		lock1.RLock()
		if (mapping.NodeTemplate != nil) && (mapping.CapabilityName != nil) {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.CapabilityMappings[mapping.Name] = normalNodeTemplate.NewMapping("capability", *mapping.CapabilityName)
			}
		}
		lock1.RUnlock()
	}

	for _, mapping := range self.RequirementMappings {
		lock1 := mapping.GetEntityLock()
		lock1.RLock()
		if (mapping.NodeTemplate != nil) && (mapping.RequirementName != nil) {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.RequirementMappings[mapping.Name] = normalNodeTemplate.NewMapping("requirement", *mapping.RequirementName)
			}
		}
		lock1.RUnlock()
	}

	for _, mapping := range self.PropertyMappings {
		lock1 := mapping.GetEntityLock()
		lock1.RLock()
		if mapping.NodeTemplate != nil {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				if mapping.PropertyName != nil {
					normalSubstitution.PropertyMappings[mapping.Name] = normalNodeTemplate.NewMapping("property", *mapping.PropertyName)
				}
			}
		} else if mapping.Property != nil {
			normalSubstitution.PropertyMappings[mapping.Name] = normal.NewMappingValue("property", mapping.Property.Normalize())
		} else if mapping.InputName != nil {
			normalSubstitution.PropertyMappings[mapping.Name] = normal.NewMapping("input", *mapping.InputName)
		}
		lock1.RUnlock()
	}

	for _, mapping := range self.AttributeMappings {
		lock1 := mapping.GetEntityLock()
		lock1.RLock()
		if (mapping.NodeTemplate != nil) && (mapping.AttributeName != nil) {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.AttributeMappings[mapping.Name] = normalNodeTemplate.NewMapping("attribute", *mapping.AttributeName)
			}
		}
		lock1.RUnlock()
	}

	for _, mapping := range self.InterfaceMappings {
		lock1 := mapping.GetEntityLock()
		lock1.RLock()
		if (mapping.NodeTemplate != nil) && (mapping.InterfaceName != nil) {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.InterfaceMappings[mapping.Name] = normalNodeTemplate.NewMapping("interface", *mapping.InterfaceName)
			}
		}
		lock1.RUnlock()
	}

	return normalSubstitution
}
