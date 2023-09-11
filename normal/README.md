Normal
======

These structs together comprise a "normalized" TOSCA-compatible topology, which is a flat,
serializable data structure.

They are meant to be generic and independent of TOSCA grammar, especially any specific version of
the TOSCA grammar.

Though they are usually created by a TOSCA parser as its final result, it's entirely possible to
write code that uses a different approach, for example by translating from a non-TOSCA topology
descriptor into the normal structures.
