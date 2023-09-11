package hot

import (
	"sync"

	"github.com/tliron/puccini/tosca/parsing"
)

//
// Entity
//

type Entity struct {
	Context *parsing.Context `traverse:"ignore" json:"-" yaml:"-"`

	renderOnce sync.Once
}

func NewEntity(context *parsing.Context) *Entity {
	return &Entity{
		Context: context,
	}
}

// ([parsing.Contextual] interface)
func (self *Entity) GetContext() *parsing.Context {
	return self.Context
}
