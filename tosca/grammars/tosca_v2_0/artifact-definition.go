package tosca_v2_0

import (
	contextpkg "context"
	"strings"

	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// ArtifactDefinition
//
// Attaches to NodeType
//
// (See Artifact for a variation that attaches to NodeTemplate)
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.7
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.7
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.6
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.6
//

type ArtifactDefinition struct {
	*Entity `name:"artifact definition"`
	Name    string

	ArtifactTypeName  *string `read:"type"` // mandatory only if cannot be inherited
	Description       *string `read:"description"`
	ArtifactVersion   *string `read:"artifact_version"` // introduced in TOSCA 1.3
	Properties        Values  `read:"properties,Value"` // ERRATUM: ommited in TOSCA 1.0-1.2 (appears in artifact type)
	RepositoryName    *string `read:"repository"`
	File              *string `read:"file"` // mandatory only if cannot be inherited
	DeployPath        *string `read:"deploy_path"`
	ChecksumAlgorithm *string `read:"checksum_algorithm"` // introduced in TOSCA 1.3
	Checksum          *string `read:"checksum"`           // introduced in TOSCA 1.3

	ArtifactType *ArtifactType `lookup:"type,ArtifactTypeName" traverse:"ignore" json:"-" yaml:"-"`
	Repository   *Repository   `lookup:"repository,RepositoryName" traverse:"ignore" json:"-" yaml:"-"`

	url                exturl.URL
	urlProblemReported bool
}

func NewArtifactDefinition(context *parsing.Context) *ArtifactDefinition {
	return &ArtifactDefinition{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Properties: make(Values),
	}
}

// ([parsing.Reader] signature)
func ReadArtifactDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewArtifactDefinition(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.File = context.FieldChild("file", context.Data).ReadString()
		// TODO: infer ArtifactTypeName from content's URI
	}

	return self
}

func (self *ArtifactDefinition) GetURL(context contextpkg.Context) exturl.URL {
	if self.File == nil {
		return nil
	}

	if self.url == nil {
		if self.Repository != nil {
			if url := self.Repository.GetURL(); url != nil {
				self.url = url.Relative(*self.File)
			}
		} else {
			base := self.Context.URL.Base()
			bases := []exturl.URL{base}
			var err error
			if self.url, err = base.Context().NewValidAnyOrFileURL(context, *self.File, bases); err != nil {
				// Avoid reporting more than once
				if !self.urlProblemReported {
					self.Context.ReportError(err)
					self.urlProblemReported = true
				}
			}
		}
	}

	return self.url
}

func (self *ArtifactDefinition) GetExtension() string {
	if self.File == nil {
		return ""
	}
	file := *self.File
	if dot := strings.Index(file, "."); dot != -1 {
		// Note: filepath.Ext will return the last extension only
		return file[dot+1:]
	} else {
		return ""
	}
}

// ([parsing.Mappable] interface)
func (self *ArtifactDefinition) GetKey() string {
	return self.Name
}

func (self *ArtifactDefinition) Inherit(parentDefinition *ArtifactDefinition) {
	logInherit.Debugf("artifact definition: %s", self.Name)

	// Validate type compatibility
	if (self.ArtifactType != nil) && (parentDefinition.ArtifactType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.ArtifactType, self.ArtifactType) {
		self.Context.ReportIncompatibleType(self.ArtifactType, parentDefinition.ArtifactType)
	}

	if (self.ArtifactTypeName == nil) && (parentDefinition.ArtifactTypeName != nil) {
		self.ArtifactTypeName = parentDefinition.ArtifactTypeName
	}
	if (self.Description == nil) && (parentDefinition.Description != nil) {
		self.Description = parentDefinition.Description
	}
	if (self.ArtifactVersion == nil) && (parentDefinition.ArtifactVersion != nil) {
		self.ArtifactVersion = parentDefinition.ArtifactVersion
	}
	if (self.Properties == nil) && (parentDefinition.Properties != nil) {
		self.Properties = parentDefinition.Properties
	}
	if (self.RepositoryName == nil) && (parentDefinition.RepositoryName != nil) {
		self.RepositoryName = parentDefinition.RepositoryName
	}
	if (self.File == nil) && (parentDefinition.File != nil) {
		self.File = parentDefinition.File
	}
	if (self.DeployPath == nil) && (parentDefinition.DeployPath != nil) {
		self.DeployPath = parentDefinition.DeployPath
	}
	if (self.ChecksumAlgorithm == nil) && (parentDefinition.ChecksumAlgorithm != nil) {
		self.ChecksumAlgorithm = parentDefinition.ChecksumAlgorithm
	}
	if (self.Checksum == nil) && (parentDefinition.Checksum != nil) {
		self.Checksum = parentDefinition.Checksum
	}
	if (self.ArtifactType == nil) && (parentDefinition.ArtifactType != nil) {
		self.ArtifactType = parentDefinition.ArtifactType
	}
	if (self.Repository == nil) && (parentDefinition.Repository != nil) {
		self.Repository = parentDefinition.Repository
	}
}

//
// ArtifactDefinitions
//

type ArtifactDefinitions map[string]*ArtifactDefinition

func (self ArtifactDefinitions) Inherit(parentDefinitions ArtifactDefinitions) {
	for name, definition := range parentDefinitions {
		if _, ok := self[name]; !ok {
			self[name] = definition
		}
	}

	for name, definition := range self {
		if parentDefinition, ok := parentDefinitions[name]; ok {
			if definition != parentDefinition {
				definition.Inherit(parentDefinition)
			}
		}
	}
}
