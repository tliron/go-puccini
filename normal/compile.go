package normal

import (
	contextpkg "context"
	"fmt"
	"time"

	"github.com/tliron/go-ard"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/tosca/parsing"
)

var emptyMap = make(ard.StringMap)

func (serviceTemplate *ServiceTemplate) Compile() (*cloutpkg.Clout, error) {
	clout := cloutpkg.NewClout()

	puccini := make(ard.StringMap)
	puccini["version"] = VERSION

	scriptlets := make(ard.StringMap)
	var err error = nil
	serviceTemplate.ScriptletNamespace.Range(func(name string, scriptlet *parsing.Scriptlet) bool {
		if scriptlet_, err_ := scriptlet.Read(contextpkg.TODO()); err_ == nil {
			if !ard.With(scriptlets).ForceGetPath(name, ".").Set(js.CleanupScriptlet(scriptlet_)) {
				err = fmt.Errorf("could not set scriptlet at path %q", name)
				return false
			}
		} else {
			err = err_
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	puccini["scriptlets"] = scriptlets

	clout.Metadata["puccini"] = puccini

	history := ard.List{ard.StringMap{
		"timestamp":   time.Now().Format(time.RFC3339Nano),
		"description": "compile",
	}}
	clout.Metadata["history"] = history

	tosca := make(ard.StringMap)
	tosca["description"] = serviceTemplate.Description
	if serviceTemplate.Metadata != nil {
		tosca["metadata"] = serviceTemplate.Metadata
	} else {
		tosca["metadata"] = emptyMap
	}
	tosca["inputs"] = serviceTemplate.Inputs
	tosca["outputs"] = serviceTemplate.Outputs
	clout.Properties["tosca"] = tosca

	nodeTemplates := make(map[string]*cloutpkg.Vertex)

	// Node templates
	for _, nodeTemplate := range serviceTemplate.NodeTemplates {
		vertex := clout.NewVertex(cloutpkg.NewKey())

		nodeTemplates[nodeTemplate.Name] = vertex

		SetMetadata(vertex, "NodeTemplate")
		vertex.Properties["name"] = nodeTemplate.Name
		if nodeTemplate.Metadata != nil {
			vertex.Properties["metadata"] = nodeTemplate.Metadata
		} else {
			vertex.Properties["metadata"] = emptyMap
		}
		vertex.Properties["description"] = nodeTemplate.Description
		vertex.Properties["types"] = nodeTemplate.Types
		vertex.Properties["directives"] = nodeTemplate.Directives
		vertex.Properties["properties"] = nodeTemplate.Properties
		vertex.Properties["attributes"] = nodeTemplate.Attributes
		vertex.Properties["requirements"] = nodeTemplate.Requirements
		vertex.Properties["capabilities"] = nodeTemplate.Capabilities
		vertex.Properties["interfaces"] = nodeTemplate.Interfaces
		vertex.Properties["artifacts"] = nodeTemplate.Artifacts
	}

	groups := make(map[string]*cloutpkg.Vertex)

	// Groups
	for _, group := range serviceTemplate.Groups {
		vertex := clout.NewVertex(cloutpkg.NewKey())

		groups[group.Name] = vertex

		SetMetadata(vertex, "Group")
		vertex.Properties["name"] = group.Name
		if group.Metadata != nil {
			vertex.Properties["metadata"] = group.Metadata
		} else {
			vertex.Properties["metadata"] = emptyMap
		}
		vertex.Properties["description"] = group.Description
		vertex.Properties["types"] = group.Types
		vertex.Properties["properties"] = group.Properties
		vertex.Properties["interfaces"] = group.Interfaces

		for _, nodeTemplate := range group.Members {
			nodeTemplateVertex := nodeTemplates[nodeTemplate.Name]
			edge := vertex.NewEdgeTo(nodeTemplateVertex)

			SetMetadata(edge, "Member")
		}
	}

	workflows := make(map[string]*cloutpkg.Vertex)

	// Workflows
	for _, workflow := range serviceTemplate.Workflows {
		vertex := clout.NewVertex(cloutpkg.NewKey())

		workflows[workflow.Name] = vertex

		SetMetadata(vertex, "Workflow")
		vertex.Properties["name"] = workflow.Name
		vertex.Properties["description"] = workflow.Description
	}

	// Workflow steps
	for name, workflow := range serviceTemplate.Workflows {
		vertex := workflows[name]

		steps := make(map[string]*cloutpkg.Vertex)

		for _, step := range workflow.Steps {
			stepVertex := clout.NewVertex(cloutpkg.NewKey())

			steps[step.Name] = stepVertex

			SetMetadata(stepVertex, "WorkflowStep")
			stepVertex.Properties["name"] = step.Name

			edge := vertex.NewEdgeTo(stepVertex)
			SetMetadata(edge, "WorkflowStep")

			if step.TargetNodeTemplate != nil {
				nodeTemplateVertex := nodeTemplates[step.TargetNodeTemplate.Name]
				edge = stepVertex.NewEdgeTo(nodeTemplateVertex)
				SetMetadata(edge, "NodeTemplateTarget")
			} else if step.TargetGroup != nil {
				groupVertex := groups[step.TargetGroup.Name]
				edge = stepVertex.NewEdgeTo(groupVertex)
				SetMetadata(edge, "GroupTarget")
			} else {
				// This would happen only if there was a parsing error
				continue
			}

			// Workflow activities
			for sequence, activity := range step.Activities {
				activityVertex := clout.NewVertex(cloutpkg.NewKey())

				edge = stepVertex.NewEdgeTo(activityVertex)
				SetMetadata(edge, "WorkflowActivity")
				edge.Properties["sequence"] = sequence

				SetMetadata(activityVertex, "WorkflowActivity")
				if activity.DelegateWorkflow != nil {
					workflowVertex := workflows[activity.DelegateWorkflow.Name]
					edge = activityVertex.NewEdgeTo(workflowVertex)
					SetMetadata(edge, "DelegateWorkflow")
				} else if activity.InlineWorkflow != nil {
					workflowVertex := workflows[activity.InlineWorkflow.Name]
					edge = activityVertex.NewEdgeTo(workflowVertex)
					SetMetadata(edge, "InlineWorkflow")
				} else if activity.SetNodeState != "" {
					activityVertex.Properties["setNodeState"] = activity.SetNodeState
				} else if activity.CallOperation != nil {
					map_ := make(ard.StringMap)
					if activity.CallOperation.Operation != nil {
						map_["interface"] = activity.CallOperation.Operation.Interface.Name
						map_["operation"] = activity.CallOperation.Operation.Name
						map_["inputs"] = activity.CallOperation.Inputs
					}
					activityVertex.Properties["callOperation"] = map_
				}
			}
		}

		for _, step := range workflow.Steps {
			stepVertex := steps[step.Name]

			for _, next := range step.OnSuccessSteps {
				nextStepVertex := steps[next.Name]
				edge := stepVertex.NewEdgeTo(nextStepVertex)
				SetMetadata(edge, "OnSuccess")
			}

			for _, next := range step.OnFailureSteps {
				nextStepVertex := steps[next.Name]
				edge := stepVertex.NewEdgeTo(nextStepVertex)
				SetMetadata(edge, "OnFailure")
			}
		}
	}

	// Policies
	for _, policy := range serviceTemplate.Policies {
		vertex := clout.NewVertex(cloutpkg.NewKey())

		SetMetadata(vertex, "Policy")
		vertex.Properties["name"] = policy.Name
		if policy.Metadata != nil {
			vertex.Properties["metadata"] = policy.Metadata
		} else {
			vertex.Properties["metadata"] = emptyMap
		}
		vertex.Properties["description"] = policy.Description
		vertex.Properties["types"] = policy.Types
		vertex.Properties["properties"] = policy.Properties

		for _, nodeTemplate := range policy.NodeTemplateTargets {
			nodeTemplateVertex := nodeTemplates[nodeTemplate.Name]
			edge := vertex.NewEdgeTo(nodeTemplateVertex)

			SetMetadata(edge, "NodeTemplateTarget")
		}

		for _, group := range policy.GroupTargets {
			groupVertex := groups[group.Name]
			edge := vertex.NewEdgeTo(groupVertex)

			SetMetadata(edge, "GroupTarget")
		}

		for _, trigger := range policy.Triggers {
			if trigger.Operation != nil {
				toVertex := clout.NewVertex(cloutpkg.NewKey())

				SetMetadata(toVertex, "Operation")
				toVertex.Properties["description"] = trigger.Operation.Description
				toVertex.Properties["implementation"] = trigger.Operation.Implementation
				toVertex.Properties["dependencies"] = trigger.Operation.Dependencies
				toVertex.Properties["inputs"] = trigger.Operation.Inputs

				edge := vertex.NewEdgeTo(toVertex)
				SetMetadata(edge, "PolicyTriggerOperation")
			} else if trigger.Workflow != nil {
				workflowVertex := workflows[trigger.Workflow.Name]

				edge := vertex.NewEdgeTo(workflowVertex)
				SetMetadata(edge, "PolicyTriggerWorkflow")
			}
		}
	}

	// Substitution
	if serviceTemplate.Substitution != nil {
		vertex := clout.NewVertex(cloutpkg.NewKey())
		inputs := make(ard.StringMap)
		properties := make(ard.StringMap)

		SetMetadata(vertex, "Substitution")
		vertex.Properties["type"] = serviceTemplate.Substitution.Type
		if serviceTemplate.Substitution.TypeMetadata != nil {
			vertex.Properties["typeMetadata"] = serviceTemplate.Substitution.TypeMetadata
		} else {
			vertex.Properties["typeMetadata"] = emptyMap
		}
		vertex.Properties["inputs"] = inputs
		vertex.Properties["properties"] = properties

		for name, pointer := range serviceTemplate.Substitution.InputPointers {
			inputs[name] = pointer.Target
		}

		for name, pointer := range serviceTemplate.Substitution.CapabilityPointers {
			NewNodeTemplatePointerEdge("Capability", name, pointer, vertex, nodeTemplates)
		}

		for name, pointer := range serviceTemplate.Substitution.RequirementPointer {
			NewNodeTemplatePointerEdge("Requirement", name, pointer, vertex, nodeTemplates)
		}

		for name, pointer := range serviceTemplate.Substitution.PropertyPointers {
			NewNodeTemplatePointerEdge("Property", name, pointer, vertex, nodeTemplates)
		}

		for name, value := range serviceTemplate.Substitution.PropertyValues {
			properties[name] = value
		}

		for name, pointer := range serviceTemplate.Substitution.AttributePointers {
			NewNodeTemplatePointerEdge("Attribute", name, pointer, vertex, nodeTemplates)
		}

		for name, pointer := range serviceTemplate.Substitution.InterfacePointers {
			NewNodeTemplatePointerEdge("Interface", name, pointer, vertex, nodeTemplates)
		}
	}

	// Make agnostic
	clout, err = clout.Copy()
	if err != nil {
		return clout, err
	}

	// TODO: call JavaScript plugins

	return clout, nil
}

func SetMetadata(entity cloutpkg.Entity, kind string) {
	metadata := make(ard.StringMap)
	metadata["version"] = VERSION
	metadata["kind"] = kind
	entity.GetMetadata()["puccini"] = metadata
}

func NewNodeTemplatePointerEdge(type_ string, name string, pointer *Pointer, substitutionVertex *cloutpkg.Vertex, nodeTemplates map[string]*cloutpkg.Vertex) {
	nodeTemplateVertex := nodeTemplates[pointer.NodeTemplate.Name]
	edge := substitutionVertex.NewEdgeTo(nodeTemplateVertex)

	SetMetadata(edge, type_+"Pointer")
	edge.Properties["name"] = name
	edge.Properties["target"] = pointer.Target
}
