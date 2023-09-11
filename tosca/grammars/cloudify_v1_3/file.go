package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// File
//
// See Blueprint
//

type File struct {
	*Entity `name:"file"`

	ToscaDefinitionsVersion *string            `read:"tosca_definitions_version" mandatory:""`
	Metadata                Metadata           `read:"metadata,!Metadata"` // not in spec, but in code
	Imports                 Imports            `read:"imports,[]Import"`
	Inputs                  Inputs             `read:"inputs,Input"`
	NodeTemplates           NodeTemplates      `read:"node_templates,NodeTemplate"`
	NodeTypes               NodeTypes          `read:"node_types,NodeType" hierarchy:""`
	Capabilities            []*ValueDefinition `read:"capability,ValueDefinition"`
	Outputs                 ValueDefinitions   `read:"outputs,ValueDefinition"`
	RelationshipTypes       RelationshipTypes  `read:"relationships,RelationshipType" hierarchy:""`
	Plugins                 Plugins            `read:"plugins,Plugin"`
	Workflows               Workflows          `read:"workflows,Workflow"`
	DataTypes               DataTypes          `read:"data_types,DataType" hierarchy:""`
	Policies                Policies           `read:"policies,Policy"`
	PolicyTypes             PolicyTypes        `read:"policy_types,PolicyType" hierarchy:""`
	PolicyTriggerTypes      PolicyTriggerTypes `read:"policy_triggers,PolicyTriggerType" hierarchy:""`
	UploadResources         *UploadResources   `read:"upload_resources,UploadResources"`
}

func NewFile(context *parsing.Context) *File {
	return &File{
		Entity:  NewEntity(context),
		Inputs:  make(Inputs),
		Outputs: make(ValueDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadFile(context *parsing.Context) parsing.EntityPtr {
	self := NewFile(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	return self
}

// ([parsing.Importer] interface)
func (self *File) GetImportSpecs() []*parsing.ImportSpec {
	var importSpecs = make([]*parsing.ImportSpec, 0, len(self.Imports))
	for _, import_ := range self.Imports {
		if importSpec, ok := import_.NewImportSpec(self); ok {
			importSpecs = append(importSpecs, importSpec)
		}
	}
	return importSpecs
}
