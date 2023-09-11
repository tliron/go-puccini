package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// InterfaceImplementation
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.16, 3.6.18
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.14
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13.2.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.13.2.3
//

type InterfaceImplementation struct {
	*Entity `name:"interface implementation"`

	Primary       *string   `read:"primary"`
	Dependencies  *[]string `read:"dependencies"`
	Timeout       *int64    `read:"timeout"`        // introduced in TOSCA 1.2
	OperationHost *string   `read:"operation_host"` // introduced in TOSCA 1.2
}

func NewInterfaceImplementation(context *parsing.Context) *InterfaceImplementation {
	return &InterfaceImplementation{Entity: NewEntity(context)}
}

// ([parsing.Reader] signature)
func ReadInterfaceImplementation(context *parsing.Context) parsing.EntityPtr {
	self := NewInterfaceImplementation(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.Primary = context.FieldChild("primary", context.Data).ReadString()
	}

	return self
}

func (self *InterfaceImplementation) Render(definition *InterfaceImplementation) {
	if definition != nil {
		if (self.Primary == nil) && (definition.Primary != nil) {
			self.Primary = definition.Primary
		}
		if (self.Dependencies == nil) && (definition.Dependencies != nil) {
			self.Dependencies = definition.Dependencies
		}
		if (self.Timeout == nil) && (definition.Timeout != nil) {
			self.Timeout = definition.Timeout
		}
		if (self.OperationHost == nil) && (definition.OperationHost != nil) {
			self.OperationHost = definition.OperationHost
		}
	}

	if self.OperationHost != nil {
		path := self.Context.Path[:2].String()
		supported := false
		operationHost := *self.OperationHost
		switch operationHost {
		case "ORCHESTRATOR":
			supported = true
		case "SELF", "HOST":
			if path == "topology_template.node_templates" {
				supported = true
			}
		case "SOURCE", "TARGET":
			if path == "topology_template.relationship_templates" {
				supported = true
			}
		}

		if !supported {
			self.Context.FieldChild("operation_host", operationHost).ReportKeynameUnsupportedValue()
		}
	}
}

func (self *InterfaceImplementation) NormalizeOperation(normalOperation *normal.Operation) {
	if self.Primary != nil {
		normalOperation.Implementation = *self.Primary
	}

	if self.Dependencies != nil {
		normalOperation.Dependencies = *self.Dependencies
	}

	if self.Timeout != nil {
		normalOperation.Timeout = *self.Timeout
	}

	if self.OperationHost != nil {
		normalOperation.Host = *self.OperationHost
	}
}

func (self *InterfaceImplementation) NormalizeNotification(normalNotification *normal.Notification) {
	if self.Primary != nil {
		normalNotification.Implementation = *self.Primary
	}

	if self.Dependencies != nil {
		normalNotification.Dependencies = *self.Dependencies
	}

	if self.Timeout != nil {
		normalNotification.Timeout = *self.Timeout
	}

	if self.OperationHost != nil {
		normalNotification.Host = *self.OperationHost
	}
}
