puccini-tosca
=============

### Format

The default format for output is YAML, but you can select JSON, XML, CBOR, or MessagePack instead with
`--format/-f`. Note that Clout in JSON may lose some type information (e.g. JSON doesn't distinguish
between an integer and a float). For this reason we also support a "compatible JSON" format ("xjson")
that adds that type information. You would need specialized code to be able to consume this format.
XML output uses a bespoke structure for maps and lists, which also must be specially consumed.
(The `puccini-clout` tool supports all these formats as input.)

For YAML you can add the additional `--strict/-y` flag to output a stricter YAML, which adds
scalar type tags (such as `!!str`, `!!int`, `!!timestamp`) and outputs all strings in double quotes
with no `|` or `>` notations. This is useful if you are consuming the YAML output with a
non-compliant or buggy parser.

Another YAML-specific flag is `--timestamps/-w`. By default Puccini will not allow the YAML
`!!timestamp` type in its output, instead emitting a canonical ISO-8601 (RFC-3339) string.
Set this flag to true to emit `!!timestamp`.

### TOSCA Quirks

**pucini-tosca** supports "quirks", via the `--quirk/-x` flag, which are variations on the default
grammar rules. The reason this is required is unfortunate: the low quality of the TOSCA spec,
riddled as it is with gaps, inconsistencies, and errors, means that there's too much room for
varying interpretations of the spec as well as missing functionality. Puccini aims to adhere as
closely as possible to the spec, literally and in spirit, but also must be pragmatic. Quirks allow
Puccini to smooth incompatibilities with other tools and work around a few TOSCA pain points.
Example of use:

    puccini-tosca compile weird.yaml -x data_types.string.permissive

The list of supported quirks is maintained [here](../tosca/parsing/QUIRKS.md).


`compile`
---------

See the [tutorial](../TUTORIAL.md) for more detail.


`validate`
----------

Equivalent to `compile` but without Clout output. Will print "valid" to stderr if valid, otherwise
will print the problems. Use `--quiet` to suppress all output.

Note that both `validate` and `compile` set the exit code: 0 for valid, 1 if there are problems, which
is useful in scripts:

    if puccini-tosca validate service.yaml --quiet; then
      deploy servicel
    fi


`parse`
-------

If you need more diagnostics for TOSCA parsing use the `parse` command. It works similarly to
`compile` but does not emit Clout. Instead, it provides you various flages for examining the
internal workings of Puccini's TOSCA parser.

By default Puccini will attempt all [6 parser phases](../tosca/parser/). This is in order to give
users as complete a problem report as possible. However, if you're getting too many problems it
may be useful to specify `--stop/-s` with a phase number (1-6) at which you wish the to stop.

`--dump/-d` is used to dump the internal data of phases. You may specify multiple phases to dump
using ",", e.g. `-d 2,3,4`. Per phase you will see:

* Phase 1: Read. The hierarchy of imported units starting at the service template URL.
* Phase 2: Namespaces. All names per type. Each imported unit has its own namespace.
* Phase 3: Hierarchies. Tree of all types by fully qualified name. Each imported unit has its own
  type hierarchy.
* Phase 4: Inheritance. A tree of all inheritance tasks and their dependencies by path.  
* Phase 5: Rendering. Dumps the rendered entities.
  More useful, perhaps, would be the `--filter/-r` flag (see below).
* Phase 6: Normalization. Dumps the normalized structures.

The `--filter/-r` flag can be used to filter for specific parsed entities. Each entity is given a
path that more-or-less follows JSON. For example, a path can be:

    topology_template.node_templates["store"].properties["name"]

The flag will search for all paths that contains your string, e.g. `-r properties`. You can even
include one or more "*" wildcards, e.g. `-r 'node*properties*data'`.
