package parsing

import "github.com/tliron/exturl"

//
// Importer
//

type Importer interface {
	GetImportSpecs() []*ImportSpec
}

// From [Importer] interface
func GetImportSpecs(entityPtr EntityPtr) []*ImportSpec {
	if importer, ok := entityPtr.(Importer); ok {
		return importer.GetImportSpecs()
	} else {
		return nil
	}
}

//
// ImportSpec
//

type ImportSpec struct {
	URL             exturl.URL
	NameTransformer NameTransformer
	Implicit        bool
}
