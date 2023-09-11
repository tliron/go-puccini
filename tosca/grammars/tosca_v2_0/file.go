package tosca_v2_0

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// File
//
// See ServiceTemplate
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.9
//

type File struct {
	*Entity `name:"file"`

	ToscaDefinitionsVersion *string           `read:"tosca_definitions_version" mandatory:""`
	Profile                 *string           `read:"profile"` // introduced in TOSCA 1.2 as "namespace", renamed in TOSCA 2.0
	Metadata                Metadata          `read:"metadata,!Metadata"`
	Description             *string           `read:"description"`
	Repositories            Repositories      `read:"repositories,Repository"`
	Imports                 Imports           `read:"imports,[]Import"`
	ArtifactTypes           ArtifactTypes     `read:"artifact_types,ArtifactType" hierarchy:""`
	CapabilityTypes         CapabilityTypes   `read:"capability_types,CapabilityType" hierarchy:""`
	DataTypes               DataTypes         `read:"data_types,DataType" hierarchy:""`
	GroupTypes              GroupTypes        `read:"group_types,GroupType" hierarchy:""`
	InterfaceTypes          InterfaceTypes    `read:"interface_types,InterfaceType" hierarchy:""`
	NodeTypes               NodeTypes         `read:"node_types,NodeType" hierarchy:""`
	PolicyTypes             PolicyTypes       `read:"policy_types,PolicyType" hierarchy:""`
	RelationshipTypes       RelationshipTypes `read:"relationship_types,RelationshipType" hierarchy:""`
}

func NewFile(context *parsing.Context) *File {
	return &File{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadFile(context *parsing.Context) parsing.EntityPtr {
	context.FunctionPrefix = "$"
	self := NewFile(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	ignore := []string{"dsl_definitions"}
	if context.HasQuirk(parsing.QuirkImportsTopologyTemplateIgnore) {
		ignore = append(ignore, "topology_template")
	}
	if context.HasQuirk(parsing.QuirkAnnotationsIgnore) {
		ignore = append(ignore, "annotation_types")
	}
	context.ValidateUnsupportedFields(append(context.ReadFields(self), ignore...))
	if self.Profile != nil {
		context.CanonicalNamespace = self.Profile
	}
	return self
}

// ([parsing.Importer] interface)
func (self *File) GetImportSpecs() []*parsing.ImportSpec {
	// TODO: importing should also import repositories

	var importSpecs = make([]*parsing.ImportSpec, 0, len(self.Imports))
	for _, import_ := range self.Imports {
		if importSpec, ok := import_.NewImportSpec(self); ok {
			importSpecs = append(importSpecs, importSpec)
		}
	}
	return importSpecs
}

func (self *File) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	logNormalize.Debug("file")

	if self.Metadata != nil {
		for k, v := range self.Metadata {
			normalServiceTemplate.Metadata[k] = v
		}
	}

	if len(self.Context.Quirks) > 0 {
		normalServiceTemplate.Metadata[parsing.MetadataQuirks] = self.Context.Quirks.String()
	}
}
