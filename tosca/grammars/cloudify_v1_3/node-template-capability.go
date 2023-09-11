package cloudify_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// NodeTemplateCapability
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-node-templates/]
//

type NodeTemplateCapability struct {
	*Entity `name:"node template capability"`

	Properties Values `read:"properties,Value"`
}

func NewNodeTemplateCapability(context *parsing.Context) *NodeTemplateCapability {
	return &NodeTemplateCapability{
		Entity:     NewEntity(context),
		Properties: make(Values),
	}
}

// ([parsing.Reader] signature)
func ReadNodeTemplateCapability(context *parsing.Context) parsing.EntityPtr {
	self := NewNodeTemplateCapability(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *NodeTemplateCapability) ValidateScalableProperties(instances *NodeTemplateInstances) {
	for key, value := range self.Properties {
		switch key {
		case "default_instances":
			value.Context.ValidateType(ard.TypeInteger)
		case "min_instances":
			value.Context.ValidateType(ard.TypeInteger)
		case "max_instances":
			value.Context.ValidateType(ard.TypeInteger, ard.TypeString)
		default:
			value.Context.ReportKeynameUnsupported()
		}
	}

	var defaultInstances int64 = 1
	if (instances != nil) && (instances.Deploy != nil) {
		defaultInstances = *instances.Deploy
	}

	propertiesContext := self.Context.FieldChild("properties", nil)
	self.Properties.SetIfNil(propertiesContext, "default_instances", defaultInstances)
	self.Properties.SetIfNil(propertiesContext, "min_instances", int64(0))
	self.Properties.SetIfNil(propertiesContext, "max_instances", "UNBOUNDED")
}

//
// NodeTemplateCapabilities
//

type NodeTemplateCapabilities map[string]*NodeTemplateCapability

func (self NodeTemplateCapabilities) Validate(context *parsing.Context, instances *NodeTemplateInstances) {
	for capabilityName, capability := range self {
		switch capabilityName {
		case "scalable":
			capability.ValidateScalableProperties(instances)
		default:
			capability.Context.ReportKeynameUnsupported()
		}
	}

	if _, ok := self["scalable"]; !ok {
		scalable := NewNodeTemplateCapability(context.FieldChild("capabilities", nil).MapChild("scalable", nil))
		self["scalable"] = scalable
		scalable.ValidateScalableProperties(instances)
	}
}
