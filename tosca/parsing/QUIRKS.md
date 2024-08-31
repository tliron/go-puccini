Puccini TOSCA Quirks
====================

These are activated via the `--quirk/-x` switch for
[**puccini-tosca**](../../executables/puccini-tosca/):

* **imports.implicit.disable**: In TOSCA 1.0-1.3 the Simple Profile is implicitly imported by
  default. This quirk will disable implicit imports.

* **imports.version.permissive**: By default Puccini will report an error if a file imports
  another file with an incompatible grammar. This quirk will disable the check.

* **imports.topology_template.ignore**: Allows imported files to contain a `topology_template`
  section, which is ignored.

* **imports.sequencedlist**: Allows the "import" syntax to be a sequenced list, in which the
  name is ignored.

* **data_types.string.permissive**: By default Puccini is strict about "string"-typed values
  and will consider integers, floats, and boolean values to be problems. This quirk will accept
  such values and convert them as sensibly as possible to strings. This includes accepting floats
  and integers for the TOSCA "version" primitive type. Note that string conversions may very well
  *not* be identical to the literal YAML. For example, `1.0000` in YAML (a float) would become
  the string `1` in TOSCA.

* **data_types.timestamp.permissive**: By default Puccini requires all "timestamp" values to be
  specified as strings in the ISO 8601 format. However, some YAML environments may support the
  optional !!timestamp type. This quirk will allow such values. Note that such values will not have
  the "$originalString" key, because the literal YAML is not preserved by the YAML parser.

* **capabilities.occurrences.permissive**: By default Puccini will ensure that capabilities have
  the minimum number of incoming relationships. This quirk will disable that validation.

* **namespace.normative.ignore**: This will ignore any type that is has the
  "tosca.normative: true" metadata.

* **namespace.normative.shortcuts.disable**: In TOSCA 1.0-1.3 all the normative types have long
  names, such as "tosca.nodes.Compute", prefixed names ("tosca:Compute"), and also short names
  ("Compute"). Those short names might be annoying because it means you can't use those names for
  your own types. This quirk disables the short names (the prefixed names remain).

* **substitution_mappings.requirements.list**: According to the examples in the TOSCA 1.0-2.0 specs,
  the `requirements` key under `substitution_mappings` is syntactically a map. However, this syntax
  is inconsistent because it doesn't match the syntax in node templates, which is a sequenced list.
  (In node types, too, it is a sequenced list, although grammatically it works like a map.) This
  quirk allows the expected syntax to be a sequenced list.

* **substitution_mappings.requirements.permissive**: Normally the `requirements` under
  `substitution_mappings` must be mapped to an assigned requirement in a node template. This quirk
  allows unassigned requirements to be mapped.

* **annotations.ignore**: Ignores the "annotation_types" keyword in service templates and the
  "annotations" keyword in parameter definitions.

* **interfaces.operations.permissive**: Allows interface types, definitions, and assignments to
  refer to operations directly in addition to using the "operations" keyname. This allows TOSCA 1.3
  and 2.0 to support the TOSCA 1.2 grammar.

Combination Quirks
------------------

* **etsinfv**: Combines "imports.topology_template.ignore", "data_types.string.permissive",
  "capabilities.occurrences.permissive", "substitution_mappings.requirements.permissive",
  "substitution_mappings.requirements.list"
* **onap**: Combines "annotations.ignore", "imports.sequencedlist", "imports.version.permissive"
