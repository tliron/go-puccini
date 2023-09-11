package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// DSLResource
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-upload-resources/]
//

type DSLResource struct {
	*Entity `name:"DSL resource"`

	SourcePath      *string `read:"source_path" mandatory:""`
	DestinationPath *string `read:"destination_path" mandatory:""`
}

func NewDSLResource(context *parsing.Context) *DSLResource {
	return &DSLResource{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadDSLResource(context *parsing.Context) parsing.EntityPtr {
	self := NewDSLResource(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

//
// DSLResources
//

type DSLResources []*DSLResource
