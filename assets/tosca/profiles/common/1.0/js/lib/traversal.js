
const tosca = require('tosca.lib.utils');

exports.toCoercibles = function(clout_) {
	if (!clout_)
		clout_ = clout;
	exports.traverseValues(clout_, function(data) {
		return clout_.newCoercible(data.value, data.site, data.source, data.target);
	});
};

exports.unwrapCoercibles = function(clout_) {
	if (!clout_)
		clout_ = clout;
	exports.traverseValues(clout_, function(data) {
		return clout_.unwrap(data.value);
	});
};

exports.coerce = function(clout_) {
	if (!clout_)
		clout_ = clout;
	exports.toCoercibles(clout_);
	exports.traverseValues(clout_, function(data) {
		return clout_.coerce(data.value);
	});
};

exports.getValueInformation = function(clout_) {
	if (!clout_)
		clout_ = clout;
	var information = {};
	exports.traverseValues(clout_, function(data) {
		if (data.value.$information)
			information[data.path.join('.')] = data.value.$information;
		return data.value;
	});
	return information;
};

exports.traverseValues = function(clout_, traverser) {
	if (!clout_)
		clout_ = clout;

	if (tosca.isTosca(clout_)) {
		exports.traverseObjectValues(traverser, ['inputs'], clout_.properties.tosca.inputs);
		exports.traverseObjectValues(traverser, ['outputs'], clout_.properties.tosca.outputs);
	}

	for (var vertexId in clout_.vertexes) {
		var vertex = clout_.vertexes[vertexId];
		if (!tosca.isTosca(vertex))
			continue;

		if (tosca.isNodeTemplate(vertex)) {
			var nodeTemplate = vertex.properties;
			var path = ['nodeTemplates', nodeTemplate.name];

			exports.traverseObjectValues(traverser, copyAndPush(path, 'properties'), nodeTemplate.properties, vertex);
			exports.traverseObjectValues(traverser, copyAndPush(path, 'attributes'), nodeTemplate.attributes, vertex);
			exports.traverseInterfaceValues(traverser, copyAndPush(path, 'interfaces'), nodeTemplate.interfaces, vertex)

			for (var capabilityName in nodeTemplate.capabilities) {
				var capability = nodeTemplate.capabilities[capabilityName];
				var capabilityPath = copyAndPush(path, 'capabilities', capabilityName);
				exports.traverseObjectValues(traverser, copyAndPush(capabilityPath, 'properties'), capability.properties, vertex);
				exports.traverseObjectValues(traverser, copyAndPush(capabilityPath, 'attributes'), capability.attributes, vertex);
			}

			for (var artifactName in nodeTemplate.artifacts) {
				var artifact = nodeTemplate.artifacts[artifactName];
				var artifactPath = copyAndPush(path, 'artifacts', artifactName);
				exports.traverseObjectValues(traverser, copyAndPush(artifactPath, 'properties'), artifact.properties, vertex);
				if (artifact.credential !== null)
					try {
						artifact.credential = traverser({
							path: copyAndPush(artifactPath, 'credential'),
							value: artifact.credential,
							site: vertex
						});
					} catch (x) {
						if ((typeof problems !== 'undefined') && x.value && x.value.error)
							// Unwrap Go error
							problems.reportError(x.value);
						else
							throw x;
					}
			}

			for (var e = 0, l = vertex.edgesOut.length; e < l; e++) {
				var edge = vertex.edgesOut[e];
				if (!tosca.isTosca(edge, 'Relationship'))
					continue;

				var relationship = edge.properties;
				var relationshipPath = copyAndPush(path, 'relationships', relationship.name);
				exports.traverseObjectValues(traverser, copyAndPush(relationshipPath, 'properties'), relationship.properties, edge, vertex, edge.target);
				exports.traverseObjectValues(traverser,copyAndPush(relationshipPath, 'attributes'), relationship.attributes, edge, vertex, edge.target);
				exports.traverseInterfaceValues(traverser, copyAndPush(relationshipPath, 'interfaces'), relationship.interfaces, edge, vertex, edge.target);
			}
		} else if (tosca.isTosca(vertex, 'Group')) {
			var group = vertex.properties;
			var path = ['groups', group.name];

			exports.traverseObjectValues(traverser, copyAndPush(path, 'properties'), group.properties, vertex);
			exports.traverseInterfaceValues(traverser, copyAndPush(path, 'attributes'), group.interfaces, vertex)
		} else if (tosca.isTosca(vertex, 'Policy')) {
			var policy = vertex.properties;
			var path = ['policies', policy.name];

			exports.traverseObjectValues(traverser, copyAndPush(path, 'properties'), policy.properties, vertex);
		}
	}
};

exports.traverseInterfaceValues = function(traverser, path, interfaces, site, source, target) {
	for (var interfaceName in interfaces) {
		var interface_ = interfaces[interfaceName];
		var interfacePath = copyAndPush(path, interfaceName)
		exports.traverseObjectValues(traverser, copyAndPush(interfacePath, 'inputs'), interface_.inputs, site, source, target);
		for (var operationName in interface_.operations)
			exports.traverseObjectValues(traverser, copyAndPush(interfacePath, 'operations', operationName), interface_.operations[operationName].inputs, site, source, target);
	}
};

exports.traverseObjectValues = function(traverser, path, object, site, source, target) {
	for (var key in object)
		try {
			object[key] = traverser({
				path: copyAndPush(path, key),
				value: object[key],
				site: site,
				source: source,
				target: target
			});
		} catch (x) {
			if ((typeof problems !== 'undefined') && x.value && x.value.error)
				// Unwrap Go error
				problems.reportError(x.value);
			else
				throw x;
		}
};

function copyAndPush(array) {
	var array_ = [];
	for (var i = 0, l = array.length; i < l; i++)
		array_.push(array[i]);
	for (var i = 1, l = arguments.length; i < l; i++)
		array_.push(arguments[i]);
	return array_;
}
