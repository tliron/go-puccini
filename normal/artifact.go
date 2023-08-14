package normal

//
// Artifact
//

type Artifact struct {
	NodeTemplate *NodeTemplate `json:"-" yaml:"-"`
	Name         string        `json:"-" yaml:"-"`

	Description       string      `json:"description" yaml:"description"`
	Types             EntityTypes `json:"types" yaml:"types"`
	Properties        Values      `json:"properties" yaml:"properties"`
	Filename          string      `json:"filename" yaml:"filename"`
	SourcePath        string      `json:"sourcePath" yaml:"sourcePath"`
	TargetPath        string      `json:"targetPath" yaml:"targetPath"`
	Version           string      `json:"version" yaml:"version"`
	ChecksumAlgorithm string      `json:"checksumAlgorithm" yaml:"checksumAlgorithm"`
	Checksum          string      `json:"checksum" yaml:"checksum"`
	Credential        Value       `json:"credential" yaml:"credential"`
}

func (self *NodeTemplate) NewArtifact(name string) *Artifact {
	artifact := &Artifact{
		NodeTemplate: self,
		Name:         name,
		Types:        make(EntityTypes),
		Properties:   make(Values),
	}
	self.Artifacts[name] = artifact
	return artifact
}

//
// Artifacts
//

type Artifacts map[string]*Artifact
