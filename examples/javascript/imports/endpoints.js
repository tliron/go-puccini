// This scriptlet gathers all endpoint capabilities and generates a report

traversal = require('tosca.lib.traversal');
tosca = require('tosca.lib.utils');

// "traversal.coerce" calls all intrinsic functions and validates all constraints
traversal.coerce();

var endpoints = [];

for (var v in clout.vertexes) {
	var vertex = clout.vertexes[v];

	// We'll skip vertexes that are not TOSCA node templates
	if (!tosca.isNodeTemplate(vertex))
		continue;

	var nodeTemplate = vertex.properties;

	for (var c in nodeTemplate.capabilities) {
		var capability = nodeTemplate.capabilities[c];

		// We'll skip capabilities that do not inherit from Endpoint
		if (!('tosca::Endpoint' in capability.types))
			continue;

		// Add endpoint to the report
		endpoints.push({
			name: nodeTemplate.name + '.' + c,
			protocol: capability.properties.protocol,
			port: capability.properties.port,
		});
	}
}

// "puccini.write" will use either YAML (the default), JSON, or XML according to the format selected
// in the command line (use --format to change it)
puccini.write(endpoints);
