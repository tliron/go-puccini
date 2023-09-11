package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Plugin
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-plugins/]
//

type Plugin struct {
	*Entity `name:"plugin"`
	Name    string `namespace:""`

	Executor            *string `read:"executor" mandatory:""`
	Source              *string `read:"source"`
	InstallArguments    *string `read:"install_arguments"`
	Install             *bool   `read:"install"`
	PackageName         *string `read:"package_name"`
	PackageVersion      *string `read:"package_version"`
	SupportedPlatform   *string `read:"supported_platform"`
	Distribution        *string `read:"distribution"`
	DistributionVersion *string `read:"distribution_version"`
	DistributionRelease *string `read:"distribution_release"`
}

func NewPlugin(context *parsing.Context) *Plugin {
	return &Plugin{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadPlugin(context *parsing.Context) parsing.EntityPtr {
	self := NewPlugin(context)

	context.ValidateUnsupportedFields(context.ReadFields(self))

	if self.Install == nil {
		// Default to true
		install := true
		self.Install = &install
	}

	if self.Executor != nil {
		executor := *self.Executor
		switch executor {
		case "central_deployment_agent", "host_agent":
		default:
			context.FieldChild("executor", executor).ReportKeynameUnsupportedValue()
		}
	}

	if *self.Install && (self.Source == nil) && (self.PackageName == nil) {
		context.FieldChild("source", nil).ReportKeynameMissing()
		context.FieldChild("package_name", nil).ReportKeynameMissing()
	}

	return self
}

//
// Plugins
//

type Plugins []*Plugin
