package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// ArtifactType
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.4
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.4
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.4
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.3
//

type ArtifactType struct {
	*Type `name:"artifact type"`

	PropertyDefinitions PropertyDefinitions `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	MIMEType            *string             `read:"mime_type" inherit:"mime_type,Parent"`
	FileExtension       *[]string           `read:"file_ext" inherit:"file_ext,Parent"`

	Parent *ArtifactType `lookup:"derived_from,ParentName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewArtifactType(context *parsing.Context) *ArtifactType {
	return &ArtifactType{
		Type:                NewType(context),
		PropertyDefinitions: make(PropertyDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadArtifactType(context *parsing.Context) parsing.EntityPtr {
	self := NewArtifactType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Hierarchical] interface)
func (self *ArtifactType) GetParent() parsing.EntityPtr {
	return self.Parent
}

// ([parsing.Inherits] interface)
func (self *ArtifactType) Inherit() {
	logInherit.Debugf("artifact type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
}

//
// ArtifactTypes
//

type ArtifactTypes []*ArtifactType
