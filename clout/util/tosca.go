package util

import (
	"github.com/tliron/go-ard"
	cloutpkg "github.com/tliron/puccini/clout"
)

func IsTosca(metadata ard.Value, kind string) bool {
	var metadata_ struct {
		Puccini struct {
			Version string `ard:"version"`
			Kind    string `ard:"kind"`
		} `ard:"puccini"`
	}

	if err := ard.NewReflector().Pack(metadata, &metadata_); err == nil {
		if metadata_.Puccini.Version == "1.0" {
			if kind != "" {
				return metadata_.Puccini.Kind == kind
			} else {
				return true
			}
		}
	}

	return false
}

func IsToscaType(entity ard.Value, type_ string) bool {
	if types, ok := ard.With(entity).Get("types").StringMap(); ok {
		if _, ok := types[type_]; ok {
			return true
		}
	}
	return false
}

func GetToscaNodeTemplates(clout *cloutpkg.Clout, type_ string) cloutpkg.Vertexes {
	vertexes := make(cloutpkg.Vertexes)
	for key, vertex := range clout.Vertexes {
		if IsTosca(vertex.Metadata, "NodeTemplate") {
			ok := true
			if type_ != "" {
				ok = IsToscaType(vertex.Properties, type_)
			}
			if ok {
				vertexes[key] = vertex
			}
		}
	}
	return vertexes
}

func GetToscaCapabilities(vertex *cloutpkg.Vertex, type_ string) ard.StringMap {
	capabilities := make(ard.StringMap)
	if capabilities_, ok := ard.With(vertex.Properties).Get("capabilities").StringMap(); ok {
		for name, capability := range capabilities_ {
			ok := true
			if type_ != "" {
				ok = IsToscaType(capability, type_)
			}
			if ok {
				capabilities[name] = capability
			}
		}
	}
	return capabilities
}

func GetToscaRelationships(vertex *cloutpkg.Vertex, type_ string) cloutpkg.Edges {
	var edges cloutpkg.Edges
	for _, edge := range vertex.EdgesOut {
		if IsTosca(edge.Metadata, "Relationship") {
			ok := true
			if type_ != "" {
				ok = IsToscaType(edge.Properties, type_)
			}
			if ok {
				edges = append(edges, edge)
			}
		}
	}
	return edges
}

func GetToscaOutputs(properties ard.Value) (ard.StringMap, bool) {
	outputs, ok := ard.With(properties).Get("tosca", "outputs").StringMap()
	return outputs, ok
}
