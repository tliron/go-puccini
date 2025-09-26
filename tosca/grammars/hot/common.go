package hot

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/go-puccini/tosca/parsing"
)

var log = commonlog.GetLogger("puccini.grammars.hot")
var logRender = commonlog.NewScopeLogger(log, "render")
var logNormalize = commonlog.NewScopeLogger(log, "normalize")

var Grammar = parsing.NewGrammar()

var DefaultScriptletNamespace = parsing.NewScriptletNamespace()

func init() {
	// https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html
	Grammar.RegisterVersion("heat_template_version", "wallaby", "")
	Grammar.RegisterVersion("heat_template_version", "train", "") // not mentioned in spec, but probably supported
	Grammar.RegisterVersion("heat_template_version", "stein", "") // not mentioned in spec, but probably supported
	Grammar.RegisterVersion("heat_template_version", "rocky", "")
	Grammar.RegisterVersion("heat_template_version", "queens", "")
	Grammar.RegisterVersion("heat_template_version", "pike", "")
	Grammar.RegisterVersion("heat_template_version", "newton", "")
	Grammar.RegisterVersion("heat_template_version", "ocata", "")
	Grammar.RegisterVersion("heat_template_version", "2021-04-16", "") // wallaby
	Grammar.RegisterVersion("heat_template_version", "2018-08-31", "") // train, stein, rocky
	Grammar.RegisterVersion("heat_template_version", "2018-03-02", "") // queens
	Grammar.RegisterVersion("heat_template_version", "2017-09-01", "") // pike
	Grammar.RegisterVersion("heat_template_version", "2017-02-24", "") // ocata
	Grammar.RegisterVersion("heat_template_version", "2016-10-14", "") // newton
	Grammar.RegisterVersion("heat_template_version", "2016-04-08", "") // mitaka
	Grammar.RegisterVersion("heat_template_version", "2015-10-15", "") // liberty
	Grammar.RegisterVersion("heat_template_version", "2015-04-30", "") // kilo
	Grammar.RegisterVersion("heat_template_version", "2014-10-16", "") // juno
	Grammar.RegisterVersion("heat_template_version", "2013-05-23", "") // icehouse

	Grammar.RegisterReader("$Root", ReadTemplate)

	Grammar.RegisterReader("Condition", ReadCondition)
	Grammar.RegisterReader("ConditionDefinition", ReadConditionDefinition)
	Grammar.RegisterReader("Constraint", ReadConstraint)
	Grammar.RegisterReader("Data", ReadData)
	Grammar.RegisterReader("Output", ReadOutput)
	Grammar.RegisterReader("Parameter", ReadParameter)
	Grammar.RegisterReader("ParameterGroup", ReadParameterGroup)
	Grammar.RegisterReader("Resource", ReadResource)
	Grammar.RegisterReader("Template", ReadTemplate)
	Grammar.RegisterReader("Value", ReadValue)

	DefaultScriptletNamespace.RegisterScriptlets(FunctionScriptlets, nil)
	DefaultScriptletNamespace.RegisterScriptlets(ConstraintScriptlets, ConstraintNativeArgumentIndexes)
}
