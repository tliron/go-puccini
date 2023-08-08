package tosca_v2_0

import (
	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Repository
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.6
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.6
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.5
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.5
//

type Repository struct {
	*Entity `name:"repository"`
	Name    string `namespace:""`

	Description *string `read:"description"`
	URL         *string `read:"url" mandatory:""`
	Credential  *Value  `read:"credential,Value"` // tosca:Credential

	url                exturl.URL
	urlProblemReported bool
}

func NewRepository(context *parsing.Context) *Repository {
	return &Repository{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// parsing.Reader signature
func ReadRepository(context *parsing.Context) parsing.EntityPtr {
	self := NewRepository(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.URL = context.FieldChild("url", context.Data).ReadString()
	}

	return self
}

// parsing.Renderable interface
func (self *Repository) Render() {
	self.renderOnce.Do(self.render)
}

func (self *Repository) render() {
	logRender.Debugf("repository: %s", self.Name)
	if self.Credential != nil {
		self.Credential.RenderDataType("tosca:Credential")
	}
}

func (self *Repository) GetURL() exturl.URL {
	if (self.url == nil) && (self.URL != nil) {
		origin := self.Context.URL.Origin()
		var err error
		if self.url, err = origin.Context().NewURL(*self.URL); err != nil {
			// Avoid reporting more than once
			if !self.urlProblemReported {
				self.Context.ReportError(err)
				self.urlProblemReported = true
			}
		}
	}

	return nil
}

//
// Repositories
//

type Repositories []*Repository
