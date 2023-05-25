package parser

import (
	contextpkg "context"
	"errors"

	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

func (self *Context) Parse(context contextpkg.Context, url exturl.URL, origins []exturl.URL, stylist *terminal.Stylist, quirks tosca.Quirks, inputs map[string]ard.Value) (*ServiceContext, *normal.ServiceTemplate, *problems.Problems, error) {
	serviceContext := self.NewServiceContext(stylist, quirks)

	// Phase 1: Read
	ok := serviceContext.ReadRoot(context, url, origins, "")
	serviceContext.MergeProblems()
	problems := serviceContext.GetProblems()

	if !problems.Empty() {
		return serviceContext, nil, problems, errors.New("read problems")
	}

	if !ok {
		return serviceContext, nil, nil, errors.New("read error")
	}

	// Phase 2: Namespaces
	serviceContext.AddNamespaces()
	serviceContext.LookupNames()

	// Phase 3: Hierarchies
	serviceContext.AddHierarchies()

	// Phase 4: Inheritance
	serviceContext.Inherit(nil)

	serviceContext.SetInputs(inputs)

	// Phase 5: Rendering
	serviceContext.Render()
	serviceContext.MergeProblems()
	if !problems.Empty() {
		return serviceContext, nil, problems, errors.New("parsing problems")
	}

	// Normalize
	serviceTemplate, ok := serviceContext.Normalize()
	if !ok || !problems.Empty() {
		return serviceContext, nil, problems, errors.New("normalization")
	}

	return serviceContext, serviceTemplate, problems, nil
}
