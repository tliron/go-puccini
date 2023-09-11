package cloudify_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// UploadResources
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-upload-resources/]
//

type UploadResources struct {
	*Entity `name:"upload resources"`

	PluginResources *[]string    `read:"plugin_resources"`
	DSLResources    DSLResources `read:"dsl_resources,[]DSLResource"`
	Parameters      Values       `read:"parameters,Value"`
}

func NewUploadResources(context *parsing.Context) *UploadResources {
	return &UploadResources{
		Entity:     NewEntity(context),
		Parameters: make(Values),
	}
}

// ([parsing.Reader] signature)
func ReadUploadResources(context *parsing.Context) parsing.EntityPtr {
	self := NewUploadResources(context)

	context.ValidateUnsupportedFields(context.ReadFields(self))

	parametersContext := context.FieldChild("parameters", nil)
	for key, value := range self.Parameters {
		childContext := parametersContext.MapChild(key, value.Context.Data)
		switch key {
		case "fetch_timeout":
			childContext.ValidateType(ard.TypeInteger)
		default:
			childContext.ReportKeynameUnsupported()
		}
	}

	return self
}
