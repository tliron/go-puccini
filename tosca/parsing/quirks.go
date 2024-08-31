package parsing

import (
	"strings"
)

//
// Quirk
//

type Quirk string

const (
	// In TOSCA 1.0-1.3 the Simple Profile is implicitly imported by default. This quirk will disable
	// implicit imports.
	QuirkImportsImplicitDisable Quirk = "imports.implicit.disable"

	// Allows imported files to contain a `topology_template`
	// section, which is ignored.
	QuirkImportsTopologyTemplateIgnore Quirk = "imports.topology_template.ignore"

	// By default Puccini will report an error if a file imports
	// another file with an incompatible grammar. This quirk will disable the check.
	QuirkImportsVersionPermissive Quirk = "imports.version.permissive"

	// Allows the "import" syntax to be a sequenced list, in which the
	// name is ignored.
	QuirkImportsSequencedList Quirk = "imports.sequencedlist"

	// By default Puccini is strict about "string"-typed values
	// and will consider integers, floats, and boolean values to be problems. This quirk will accept
	// such values and convert them as sensibly as possible to strings. This includes accepting floats
	// and integers for the TOSCA "version" primitive type. Note that string conversions may very well
	// *not* be identical to the literal YAML. For example, `1.0000` in YAML (a float) would become
	// the string `1` in TOSCA.
	QuirkDataTypesStringPermissive Quirk = "data_types.string.permissive"

	// By default Puccini requires all "timestamp" values to be
	// specified as strings in the ISO 8601 format. However, some YAML environments may support the
	// optional !!timestamp type. This quirk will allow such values. Note that such values will not have
	// the "$originalString" key, because the literal YAML is not preserved by the YAML parser.
	QuirkDataTypesTimestampPermissive Quirk = "data_types.timestamp.permissive"

	// By default Puccini will ensure that capabilities have
	// the minimum number of incoming relationships. This quirk will disable that validation.
	QuirkCapabilitiesOccurrencesPermissive Quirk = "capabilities.occurrences.permissive"

	// This will ignore any type that is has the `tosca.normative: 'true'` metadata.
	QuirkNamespaceNormativeIgnore Quirk = "namespace.normative.ignore"

	// In TOSCA 1.0-1.3 all the normative types have long
	// names, such as "tosca.nodes.Compute", prefixed names ("tosca:Compute"), and also short names
	// ("Compute"). Those short names might be annoying because it means you can't use those names for
	// your own types. This quirk disables the short names (the prefixed names remain).
	QuirkNamespaceNormativeShortcutsDisable Quirk = "namespace.normative.shortcuts.disable"

	// According to the examples in the TOSCA 1.0-1.3 specs,
	// the `requirements` key under `substitution_mappings` is syntactically a map. However, this syntax
	// is inconsistent because it doesn't match the syntax in node templates, which is a sequenced list.
	// (In node types, too, it is a sequenced list, although grammatically it works like a map.) This
	// quirk allows the expected syntax to be a sequenced list.
	QuirkSubstitutionMappingsRequirementsList Quirk = "substitution_mappings.requirements.list"

	// Normally the `requirements` under
	// `substitution_mappings` must be mapped to an assigned requirement in a node template. This quirk
	// allows unassigned requirements to be mapped.
	QuirkSubstitutionMappingsRequirementsPermissive Quirk = "substitution_mappings.requirements.permissive"

	// Ignores the "annotation_types" keyword in service templates and the
	// "annotations" keyword in parameter definitions.
	QuirkAnnotationsIgnore Quirk = "annotations.ignore"

	// Allows interface types, definitions, and assignments to
	// refer to operations directly in addition to using the "operations" keyname. This allows TOSCA 1.3
	// and 2.0 to support the TOSCA 1.2 grammar.
	QuirkInterfacesOperationsPermissive Quirk = "interfaces.operations.permissive"

	// Combines "imports.topology_template.ignore", "data_types.string.permissive",
	// "capabilities.occurrences.permissive", "substitution_mappings.requirements.permissive",
	// "substitution_mappings.requirements.list"
	QuirkETSINFV Quirk = "etsinfv"

	// Combines "annotations.ignore", "imports.sequencedlist", "imports.version.permissive"
	QuirkONAP Quirk = "onap"
)

var combinationQuirks = map[Quirk][]Quirk{
	QuirkETSINFV: {
		QuirkImportsTopologyTemplateIgnore,
		QuirkDataTypesStringPermissive,
		QuirkCapabilitiesOccurrencesPermissive,
		QuirkSubstitutionMappingsRequirementsPermissive,
		QuirkSubstitutionMappingsRequirementsList,
	},
	QuirkONAP: {
		QuirkAnnotationsIgnore,
		QuirkImportsSequencedList,
		QuirkImportsVersionPermissive,
	},
}

//
// Quirks
//

type Quirks []Quirk

func NewQuirks(quirks ...string) Quirks {
	var self Quirks
	for _, quirk := range quirks {
		quirk_ := Quirk(quirk)
		if quirks_, ok := combinationQuirks[quirk_]; ok {
			self = append(self, quirks_...)
		} else {
			self = append(self, quirk_)
		}
	}
	return self
}

func (self Quirks) Has(quirk Quirk) bool {
	for _, quirk_ := range self {
		if quirk_ == quirk {
			return true
		}
	}
	return false
}

// fmt.Stringify interface
func (self Quirks) String() string {
	if len(self) > 0 {
		var b strings.Builder
		b.WriteString(string(self[0]))
		for _, quirk := range self[1:] {
			b.WriteRune(',')
			b.WriteString(string(quirk))
		}
		return b.String()
	} else {
		return ""
	}
}
