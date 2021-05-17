package tosca_v2_0

import (
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
	SubstitutionFilter  NodeFilter          `read:"substitution_filter,NodeFilter"` // introduced in TOSCA 1.3

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
		context.SetReadTag("RequirementMappings", "requirements,{}RequirementMapping")
	}

	self := NewSubstitutionMappings(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *SubstitutionMappings) IsRequirementMapped(nodeTemplate *NodeTemplate, requirementName string) bool {
	for _, mapping := range self.RequirementMappings {
		if mapping.NodeTemplate == nodeTemplate {
			if (mapping.RequirementName != nil) && (*mapping.RequirementName == requirementName) {
				return true
			}
		}
	}
	return false
}

func (self *SubstitutionMappings) Render(inputDefinitions ParameterDefinitions) {
	logRender.Debug("substitution mappings")

	if self.NodeType == nil {
		return
	}

	for name, mapping := range self.CapabilityMappings {
		if definition, ok := self.NodeType.CapabilityDefinitions[name]; ok {
			if mappedDefinition, ok := mapping.GetCapabilityDefinition(); ok {
				if (definition.CapabilityType != nil) && (mappedDefinition.CapabilityType != nil) {
					if !self.Context.Hierarchy.IsCompatible(definition.CapabilityType, mappedDefinition.CapabilityType) {
						self.Context.ReportIncompatibleType(definition.CapabilityType, mappedDefinition.CapabilityType)
					}
				}
			}
		} else {
			mapping.Context.Clone(name).ReportReferenceNotFound("capability", self.NodeType)
		}
	}

	for name, mapping := range self.RequirementMappings {
		if _, ok := self.NodeType.RequirementDefinitions[name]; !ok {
			if ! mapping.Context.HasQuirk(tosca.QuirkSubstitutionMappingsRequirementsAllowDangling) {
				mapping.Context.Clone(name).ReportReferenceNotFound("requirement", self.NodeType)
			}
		}
	}

	self.PropertyMappings.Render(inputDefinitions)
	for name, mapping := range self.PropertyMappings {
		if definition, ok := self.NodeType.PropertyDefinitions[name]; ok {
			if mapping.InputDefinition != nil {
				// Input mapping
				if (definition.DataType != nil) && (mapping.InputDefinition.DataType != nil) {
					if !self.Context.Hierarchy.IsCompatible(definition.DataType, mapping.InputDefinition.DataType) {
						self.Context.ReportIncompatibleType(definition.DataType, mapping.InputDefinition.DataType)
					}
				}
			} else if mapping.Property != nil {
				// Property mapping (deprecated in TOSCA 1.3)
				if (definition.DataType != nil) && (mapping.Property.DataType != nil) {
					if !self.Context.Hierarchy.IsCompatible(definition.DataType, mapping.Property.DataType) {
						self.Context.ReportIncompatibleType(definition.DataType, mapping.Property.DataType)
					}
				}
			}
		} else {
			mapping.Context.Clone(name).ReportReferenceNotFound("property", self.NodeType)
		}
	}

	self.AttributeMappings.EnsureRender()
	for name, mapping := range self.AttributeMappings {
		if definition, ok := self.NodeType.AttributeDefinitions[name]; ok {
			if (definition.DataType != nil) && (mapping.Attribute != nil) && (mapping.Attribute.DataType != nil) {
				if !self.Context.Hierarchy.IsCompatible(definition.DataType, mapping.Attribute.DataType) {
					self.Context.ReportIncompatibleType(definition.DataType, mapping.Attribute.DataType)
				}
			}
		} else {
			mapping.Context.Clone(name).ReportReferenceNotFound("attribute", self.NodeType)
		}
	}

	for name, mapping := range self.InterfaceMappings {
		if definition, ok := self.NodeType.InterfaceDefinitions[name]; ok {
			if mappedDefinition, ok := mapping.GetInterfaceDefinition(); ok {
				if (definition.InterfaceType != nil) && (mappedDefinition.InterfaceType != nil) {
					if !self.Context.Hierarchy.IsCompatible(definition.InterfaceType, mappedDefinition.InterfaceType) {
						self.Context.ReportIncompatibleType(definition.InterfaceType, mappedDefinition.InterfaceType)
					}
				}
			}
		} else {
			mapping.Context.Clone(name).ReportReferenceNotFound("interface", self.NodeType)
		}
	}
}

func (self *SubstitutionMappings) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.Substitution {
	logNormalize.Debug("substitution mappings")

	if self.NodeType == nil {
		return nil
	}

	normalSubstitution := normalServiceTemplate.NewSubstitution()

	normalSubstitution.Type = tosca.GetCanonicalName(self.NodeType)

	if metadata, ok := self.NodeType.GetMetadata(); ok {
		normalSubstitution.TypeMetadata = metadata
	}

	for _, mapping := range self.CapabilityMappings {
		if (mapping.NodeTemplate != nil) && (mapping.CapabilityName != nil) {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.CapabilityMappings[mapping.Name] = normalNodeTemplate.NewMapping("capability", *mapping.CapabilityName)
			}
		}
	}

	for _, mapping := range self.RequirementMappings {
		if (mapping.NodeTemplate != nil) && (mapping.RequirementName != nil) {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.RequirementMappings[mapping.Name] = normalNodeTemplate.NewMapping("requirement", *mapping.RequirementName)
			}
		}
	}

	for _, mapping := range self.PropertyMappings {
		if (mapping.NodeTemplate != nil) && (mapping.PropertyName != nil) {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.PropertyMappings[mapping.Name] = normalNodeTemplate.NewMapping("property", *mapping.PropertyName)
			}
		} else if mapping.InputName != nil {
			normalSubstitution.PropertyMappings[mapping.Name] = normal.NewMapping("input", *mapping.InputName)
		}
	}

	for _, mapping := range self.AttributeMappings {
		if (mapping.NodeTemplate != nil) && (mapping.AttributeName != nil) {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.AttributeMappings[mapping.Name] = normalNodeTemplate.NewMapping("attribute", *mapping.AttributeName)
			}
		}
	}

	for _, mapping := range self.InterfaceMappings {
		if (mapping.NodeTemplate != nil) && (mapping.InterfaceName != nil) {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.InterfaceMappings[mapping.Name] = normalNodeTemplate.NewMapping("interface", *mapping.InterfaceName)
			}
		}
	}

	return normalSubstitution
}
