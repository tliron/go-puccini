package hot

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Data
//

type Data struct {
	*Entity `name:"data"`

	Data ard.Value
}

func NewData(context *parsing.Context) *Data {
	return &Data{
		Entity: NewEntity(context),
		Data:   context.Data,
	}
}

// ([parsing.Reader] signature)
func ReadData(context *parsing.Context) parsing.EntityPtr {
	return NewData(context)
}
