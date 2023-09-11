package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Artifact
//

// ([parsing.Reader] signature)
func ReadArtifact(context *parsing.Context) parsing.EntityPtr {
	context.SetReadTag("ArtifactVersion", "")
	context.SetReadTag("ChecksumAlgorithm", "")
	context.SetReadTag("Checksum", "")

	return tosca_v2_0.ReadArtifact(context)
}
