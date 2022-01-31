package parser

import (
	"sort"
	"strings"

	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca"
)

//
// NoEntity
//

type NoEntity struct {
	Context *tosca.Context
}

func NewNoEntity(toscaContext *tosca.Context) *NoEntity {
	return &NoEntity{toscaContext}
}

// tosca.Contextual interface
func (self *NoEntity) GetContext() *tosca.Context {
	return self.Context
}

//
// Unit
//

type Unit struct {
	EntityPtr       tosca.EntityPtr
	Container       *Unit
	Imports         Units
	NameTransformer tosca.NameTransformer

	importsLock util.RWLocker
}

func NewUnitNoEntity(toscaContext *tosca.Context, container *Unit, nameTransformer tosca.NameTransformer) *Unit {
	return NewUnit(NewNoEntity(toscaContext), container, nameTransformer)
}

func NewUnit(entityPtr tosca.EntityPtr, container *Unit, nameTransformer tosca.NameTransformer) *Unit {
	self := Unit{
		EntityPtr:       entityPtr,
		Container:       container,
		NameTransformer: nameTransformer,
		importsLock:     util.NewDebugRWLocker(),
	}

	if container != nil {
		container.AddImport(&self)
	}

	return &self
}

func (self *Unit) AddImport(import_ *Unit) {
	self.importsLock.Lock()
	defer self.importsLock.Unlock()

	self.Imports = append(self.Imports, import_)
}

func (self *Unit) GetContext() *tosca.Context {
	return tosca.GetContext(self.EntityPtr)
}

// Print

func (self *Unit) PrintImports(indent int, treePrefix terminal.TreePrefix) {
	self.importsLock.RLock()
	length := len(self.Imports)
	imports := make(Units, length)
	copy(imports, self.Imports)
	self.importsLock.RUnlock()

	last := length - 1

	// Sort
	sort.Sort(imports)

	for i, unit := range imports {
		isLast := i == last
		unit.PrintNode(indent, treePrefix, isLast)
		unit.PrintImports(indent, append(treePrefix, isLast))
	}
}

func (self *Unit) PrintNode(indent int, treePrefix terminal.TreePrefix, last bool) {
	treePrefix.Print(indent, last)
	terminal.Printf("%s\n", terminal.Stylize.Value(self.GetContext().URL.String()))
}

//
// Units
//

type Units []*Unit

// sort.Interface

func (self Units) Len() int {
	return len(self)
}

func (self Units) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self Units) Less(i, j int) bool {
	iName := self[i].GetContext().URL.String()
	jName := self[j].GetContext().URL.String()
	return strings.Compare(iName, jName) < 0
}
