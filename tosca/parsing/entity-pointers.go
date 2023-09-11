package parsing

import (
	"strings"
)

//
// EntityPtr
//

type EntityPtr = any

//
// EntityPtrs
//

type EntityPtrs []EntityPtr

// ([sort.Interface])
func (self EntityPtrs) Len() int {
	return len(self)
}

// ([sort.Interface])
func (self EntityPtrs) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// ([sort.Interface])
func (self EntityPtrs) Less(i, j int) bool {
	iName := GetContext(self[i]).Path.String()
	jName := GetContext(self[j]).Path.String()
	return strings.Compare(iName, jName) < 0
}

//
// EntityPtrSet
//

type EntityPtrSet map[EntityPtr]struct{}

func (self EntityPtrSet) Add(entityPtr EntityPtr) {
	self[entityPtr] = struct{}{}
}

func (self EntityPtrSet) Contains(entityPtr EntityPtr) bool {
	_, ok := self[entityPtr]
	return ok
}
