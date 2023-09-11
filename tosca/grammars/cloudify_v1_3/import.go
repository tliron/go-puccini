package cloudify_v1_3

import (
	contextpkg "context"
	"fmt"
	"strings"

	"github.com/tliron/exturl"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Import
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-imports/]
//

type Import struct {
	*Entity `name:"import" json:"-" yaml:"-"`

	File *string
}

func NewImport(context *parsing.Context) *Import {
	return &Import{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadImport(context *parsing.Context) parsing.EntityPtr {
	self := NewImport(context)
	self.File = context.ReadString()
	return self
}

func (self *Import) NewImportSpec(unit *File) (*parsing.ImportSpec, bool) {
	if self.File == nil {
		return nil, false
	}

	file := *self.File

	if strings.HasPrefix(file, "plugin:") {
		return nil, false
	}

	var nameTransformer parsing.NameTransformer
	if s := strings.SplitN(file, "--", 1); len(s) == 2 {
		if strings.Contains(s[0], "-") {
			self.Context.ReportValueMalformed("namespace", "contains '-'")
		}
		nameTransformer = newImportNameTransformer(s[0])
		file = s[1]
	}

	base := self.Context.URL.Base()
	var bases = []exturl.URL{base}
	url, err := base.Context().NewValidAnyOrFileURL(contextpkg.TODO(), file, bases)
	if err != nil {
		self.Context.ReportError(err)
		return nil, false
	}

	importSpec := &parsing.ImportSpec{
		URL:             url,
		NameTransformer: nameTransformer,
		Implicit:        false,
	}
	return importSpec, true
}

func newImportNameTransformer(prefix string) parsing.NameTransformer {
	return func(name string, entityPtr parsing.EntityPtr) []string {
		var names []string

		// Prefixed name
		names = append(names, fmt.Sprintf("%s--%s", prefix, name))

		return names
	}
}

//
// Imports
//

type Imports []*Import
