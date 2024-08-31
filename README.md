Puccini
=======

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/puccini.svg)](https://github.com/tliron/puccini/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/tliron/puccini.svg)](https://pkg.go.dev/github.com/tliron/puccini)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/puccini)](https://goreportcard.com/report/github.com/tliron/puccini)

Cloud topology management and deployment tools based on
[TOSCA](https://www.oasis-open.org/committees/tosca/).

Puccini is primarily a TOSCA processor. It parses TOSCA service templates and compiles them into
the [Clout format](clout/).

Let's dive in!

* Head to the [tutorial](TUTORIAL.md). It will likely answer most of your questions.
* Want to try it out without installing it?
  [Here's a live demo of Puccini TOSCA running in a browser](https://web.puccini.cloud/).

Note that Puccini is intentionally *not* an orchestrator. This is a "BYOO" kind of establishment
("Bring Your Own Orchestrator"). Available orchestration integrations:

* [Khutulun](https://khutulun.puccini.cloud/) is an extensible TOSCA-native orchestrator
* [Kubernetes](https://kubernetes.io/): [Turandot](https://turandot.puccini.cloud/) packages Puccini
  as a Kubernetes operator and includes TOSCA types generated directly from the Kubernetes API
* [Ansible](https://www.ansible.com/): Puccini is [packaged for Ansible Galaxy](wrappers/ansible/)
* [OpenStack](https://www.openstack.org/): see [these examples using Ansible](examples/openstack/)
  (Puccini also natively supports OpenStack's Heat Orchestration Template language)
* Process automation: see [these examples for BPMN integration](examples/bpmn/)

Puccini is also available as a [GitHub action](https://github.com/marketplace/actions/setup-puccini-tosca),
allowing you to validate and otherwise use TOSCA in your git workflow.

For a TOSCA development environment, check out the
[TOSCA Visual Studio Code Extension](https://github.com/tliron/puccini-vscode/), which is based
on Puccini (work in progress).


Get It
------

[![Download](assets/media/download.png "Download")](https://github.com/tliron/puccini/releases)

To build Puccini yourself see the [build guide](scripts/).


Overview
--------

### Distribution

Each tool is a self-contained executable file, allowing them to be easily distributed with and
embedded in toolchains, orchestration, and development environments. They are coded in 100%
[Go](https://golang.org/) and are very portable, even available for
[WebAssembly](https://webassembly.org/), which is how the in-browser demo linked above works.

You can also embed Puccini into your program as a library. Puccini is immediately usable from Go,
but can be used in many other programming languages via self-contained shared C libraries. See
included wrappers and examples for [Java](wrappers/java/), [Python](wrappers/python/), and
[Ruby](wrappers/ruby/).

### Dialects

Puccini can parse all versions of TOSCA:

* [TOSCA 2.0](http://docs.oasis-open.org/tosca/TOSCA/v2.0/TOSCA-v2.0.html) (in progress)
* [TOSCA 1.3](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.3/TOSCA-Simple-Profile-YAML-v1.3.html)
* [TOSCA 1.2](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.2/TOSCA-Simple-Profile-YAML-v1.2.html)
* [TOSCA 1.1](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html)
* [TOSCA 1.0](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.0/TOSCA-Simple-Profile-YAML-v1.0.html)

Additionally, Puccini can parse the following TOSCA-like dialects:

* [Cloudify DSL 1.3](https://docs.cloudify.co/6.3.0/developer/blueprints/)
* [OpenStack Heat Orchestration Template language (HOT) 2021-04-16](https://docs.openstack.org/heat/wallaby/template_guide/hot_guide.html)

TOSCA is a complex object-oriented language. We put considerable effort into adhering to every
aspect of the grammar, especially in regards to data type checking and type inheritance contracts,
which are key to delivering the object-oriented promise of extensibility while maintaining reliable
base-type compatibility. Unfortunately, in earlier versions of TOSCA some grammatical features and
even some syntax have been specified in ways that are open to interpretation. Puccini picks one
interpretation by default, but also supports [quirk modes](tosca/parsing/QUIRKS.md) that enable
alternative behaviors.

### Packaging

The TOSCA source can be accessed by URL, either on the local file systems or via HTTP/HTTPS, as
individual files as well as packaged in
[CSAR files](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.3/TOSCA-Simple-Profile-YAML-v1.3.html#_Toc302251718).

Puccini also comes with a simple CSAR creation tool that can be used independently of the
other tools, [`puccini-csar`](executables/puccini-csar/).

### Design Principles

Puccini's TOSCA parser comprises 6 multi-threaded [phases](tosca/parser/) that handle validation,
inheritance, assignment, and normalization of TOSCA's many types and templates, resulting in a
[flat, serializable data structure](normal/) that can be directly consumed by your program.
Validation error messages are precise and useful.


FAQ
===

Please [read it](FAQ.md)!
