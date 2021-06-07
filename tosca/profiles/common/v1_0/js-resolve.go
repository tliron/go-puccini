// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/common/1.0/js/resolve.js"] = `

const traversal = require('tosca.lib.traversal');
const tosca = require('tosca.lib.utils');

// Remove existing relationships
var nodeTemplateVertexes = [];
for (var vertexId in clout.vertexes) {
	var vertex = clout.vertexes[vertexId];
	if (tosca.isNodeTemplate(vertex)) {
		nodeTemplateVertexes.push(vertex);
		var remove = [];
		for (var e = 0, l = vertex.edgesOut.length; e < l; e++) {
			var edge = vertex.edgesOut[e];
			if (tosca.isTosca(edge, 'Relationship'))
				remove.push(edge);
		}
		for (var e = 0, l = remove.length; e < l; e++)
			remove[e].remove();
	}
}

// For consistent results, we will sort the node templates by name
nodeTemplateVertexes.sort(function(a, b) {
	return a.properties.name < b.properties.name ? -1 : 1;
});

traversal.toCoercibles();

// Resolve all requirements
for (var v = 0, l = nodeTemplateVertexes.length; v < l; v++) {
	var vertex = nodeTemplateVertexes[v];
	var nodeTemplate = vertex.properties;
	var requirements = nodeTemplate.requirements;
	for (var r = 0, ll = requirements.length; r < ll; r++) {
		var requirement = requirements[r];
		resolve(vertex, nodeTemplate, requirement);
	}
}

// Check that all capabilities have their minimum relationship count
for (var v = 0, l = nodeTemplateVertexes.length; v < l; v++) {
	var vertex = nodeTemplateVertexes[v];
	var nodeTemplate = vertex.properties;
	var capabilities = nodeTemplate.capabilities;
	for (var capabilityName in capabilities) {
		var capability = capabilities[capabilityName];
		var relationshipCount = countRelationships(vertex, capabilityName);
		var minRelationshipCount = capability.minRelationshipCount;
		if (relationshipCount < minRelationshipCount)
			notEnoughRelationships(capability.location, relationshipCount, minRelationshipCount)
	}
}

traversal.unwrapCoercibles();

if (puccini.arguments.history !== 'false')
	tosca.addHistory('resolve');
puccini.write(clout)

function resolve(sourceVertex, sourceNodeTemplate, requirement) {
	var location = requirement.location;
	var name = requirement.name;

	if (isSubstituted(sourceNodeTemplate.name, name)) {
		puccini.log.debugf('%s: skipping because in substitution mappings', location.path)
		return;
	}

	var candidates = gatherCandidateNodeTemplates(sourceVertex, requirement);
	if (candidates.length === 0) {
		unsatisfied(location, name, 'there are no candidate node templates');
		return;
	}

	candidates = gatherCandidateCapabilities(requirement, candidates);
	if (candidates.length === 0) {
		unsatisfied(location, name, 'no candidate node template provides required capability');
		return;
	}

	// Gather priority candidates: those that have not yet fulfilled their minimum relationship count
	var priorityCandidates = [];
	for (var c = 0, l = candidates.length; c < l; c++) {
		var candidate = candidates[c];
		if ((candidate.capability.minRelationshipCount !== 0) && (countRelationships(candidate.vertex, candidate.capabilityName) < candidate.capability.minRelationshipCount))
			priorityCandidates.push(candidate);
	}

	var chosen = null;

	if (priorityCandidates.length !== 0)
		// Of the priority candidates, pick the one with the highest minimum relationship count
		// (needs to be fulfilled soonest)
		for (var c = 0, l = priorityCandidates.length; c < l; c++) {
			var candidate = priorityCandidates[c];
			if ((chosen === null) || (candidate.capability.minRelationshipCount > chosen.capability.minRelationshipCount))
				chosen = candidate;
		}
	else
		// Of the candidates, pick the one with highest maximum relationship count
		// (has the most room)
		for (var c = 0, l = candidates.length; c < l; c++) {
			candidate = candidates[c];
			if ((chosen === null) || isMaxCountGreater(candidate.capability.maxRelationshipCount, chosen.capability.maxRelationshipCount))
				chosen = candidate;
		}

	puccini.log.debugf('%s: satisfied "%s" with capability "%s" in node template "%s"', location.path, name, chosen.capabilityName, chosen.nodeTemplateName);
	addRelationship(sourceVertex, requirement, chosen.vertex, chosen.capabilityName);
}

function gatherCandidateNodeTemplates(sourceVertex, requirement) {
	var path = requirement.location.path;
	var nodeTemplateName = requirement.nodeTemplateName;
	var nodeTypeName = requirement.nodeTypeName;
	var nodeTemplatePropertyConstraints = requirement.nodeTemplatePropertyConstraints;
	var capabilityPropertyConstraintsMap = requirement.capabilityPropertyConstraints;

	var candidates = [];
	for (var v = 0, l = nodeTemplateVertexes.length; v < l; v++) {
		var vertex = nodeTemplateVertexes[v];
		var candidateNodeTemplate = vertex.properties;
		var candidateNodeTemplateName = candidateNodeTemplate.name;

		if ((nodeTemplateName !== '') && (nodeTemplateName !== candidateNodeTemplateName)) {
			puccini.log.debugf('%s: node template "%s" is not named "%s"', path, candidateNodeTemplateName, nodeTemplateName);
			continue;
		}

		if ((nodeTypeName !== '') && !(nodeTypeName in candidateNodeTemplate.types)) {
			puccini.log.debugf('%s: node template "%s" is not of type "%s"', path, candidateNodeTemplateName, nodeTypeName);
			continue;
		}

		// Node filter
		if ((nodeTemplatePropertyConstraints.length !== 0) && !arePropertiesValid(path, sourceVertex, 'node template', candidateNodeTemplateName, candidateNodeTemplate, nodeTemplatePropertyConstraints)) {
			puccini.log.debugf('%s: properties of node template "%s" do not match constraints', path, candidateNodeTemplateName);
			continue;
		}

		var candidateCapabilities = candidateNodeTemplate.capabilities;

		// Capability filter
		if (capabilityPropertyConstraintsMap.length !== 0) {
			var valid = true;
			for (var candidateCapabilityName in candidateCapabilities) {
				var candidateCapability = candidateCapabilities[candidateCapabilityName];
				var capabilityPropertyConstraints = capabilityPropertyConstraintsMap[candidateCapabilityName];
				if ((capabilityPropertyConstraints !== undefined) && (capabilityPropertyConstraints.length !== 0) && !arePropertiesValid(path, sourceVertex, 'capability', candidateCapabilityName, candidateCapability, capabilityPropertyConstraints)) {
					puccini.log.debugf('%s: properties of capability "%s" in node template "%s" do not match constraints', path, candidateCapabilityName, candidateNodeTemplateName);
					valid = false;
					break;
				}
			}
			if (!valid)
				continue;
		}

		candidates.push({
			vertex: vertex,
			nodeTemplateName: candidateNodeTemplateName,
			capabilities: candidateCapabilities
		});
	}

	return candidates;
}

function gatherCandidateCapabilities(requirement, candidateNodeTemplates) {
	var path = requirement.location.path;
	var capabilityName = requirement.capabilityName;
	var capabilityTypeName = requirement.capabilityTypeName;

	var candidates = [];
	for (var c = 0, l = candidateNodeTemplates.length; c < l; c++) {
		var candidate = candidateNodeTemplates[c];
		var candidateVertex = candidate.vertex;
		var candidateNodeTemplateName = candidate.nodeTemplateName;

		var candidateCapabilities = [];
		for (var candidateCapabilityName in candidate.capabilities) {
			candidateCapabilities.push({
				name: candidateCapabilityName,
				capability: candidate.capabilities[candidateCapabilityName]
			});
		}

		// For consistent results, we will sort the candidate capabilities by name
		candidateCapabilities.sort(function(a, b) {
			return a.name < b.name ? -1 : 1;
		});

		for (var cc = 0, ll = candidateCapabilities.length; cc < ll; cc++) {
			var candidateCapabilityName = candidateCapabilities[cc].name;

			if ((capabilityName !== '') && (capabilityName !== candidateCapabilityName)) {
				puccini.log.debugf('%s: capability "%s" in node template "%s" is not named "%s"', path, candidateCapabilityName, candidateNodeTemplateName, capabilityName);
				continue;
			}

			var candidateCapability = candidateCapabilities[cc].capability;

			if ((capabilityTypeName !== '') && !(capabilityTypeName in candidateCapability.types)) {
				puccini.log.debugf('%s: capability "%s" in node template "%s" is not of type "%s"', path, candidateCapabilityName, candidateNodeTemplateName, capabilityTypeName);
				continue;
			}

			var maxRelationshipCount = candidateCapability.maxRelationshipCount;
			if ((maxRelationshipCount !== -1) && (countRelationships(candidateVertex, candidateCapabilityName) === maxRelationshipCount)) {
				puccini.log.debugf('%s: capability "%s" in node template "%s" already has %d relationships, the maximum allowed', path, candidateCapabilityName, candidateNodeTemplateName, maxRelationshipCount);
				continue;
			}

			candidates.push({
				vertex: candidateVertex,
				nodeTemplateName: candidateNodeTemplateName,
				capability: candidateCapability,
				capabilityName: candidateCapabilityName
			});
		}
	}

	return candidates;
}

function addRelationship(sourceVertex, requirement, targetVertex, capabilityName) {
	var edge = sourceVertex.newEdgeTo(targetVertex);
	edge.metadata['puccini'] = {
		version: '1.0',
		kind: 'Relationship'
	};

	var relationship = requirement.relationship;
	if (relationship)
		edge.properties = {
			name: requirement.name,
			description: relationship.description,
			types: relationship.types,
			properties: relationship.properties,
			attributes: relationship.attributes,
			interfaces: relationship.interfaces,
			capability: capabilityName
		};
	else
		// Untyped relationship
		edge.properties = {
			name: requirement.name,
			description: '',
			types: {},
			properties: {},
			attributes: {},
			interfaces: {},
			capability: capabilityName
		};
}

function countRelationships(vertex, capabilityName) {
	var count = 0;
	for (var e = 0, l = vertex.edgesIn.length; e < l; e++) {
		var edge = vertex.edgesIn[e];
		if (tosca.isTosca(edge, 'Relationship') && (edge.properties.capability === capabilityName))
			count++;
	}
	return count;
}

function arePropertiesValid(path, sourceVertex, kind, name, entity, constraintsMap) {
	var valid = true;

	var properties = entity.properties;
	for (var propertyName in constraintsMap) {
		puccini.log.debugf('%s: applying constraints to property "%s" of %s "%s"', path, propertyName, kind, name);

		var property = properties[propertyName];
		if (property === undefined) {
			// return false; GOJA: returning from inside for-loop is broken
			valid = false;
			break;
		}

		var constraints = constraintsMap[propertyName];
		constraints = clout.newConstraints(constraints, sourceVertex, sourceVertex, entity)
		if (!constraints.validate(property)) {
			// return false; GOJA: returning from inside for-loop is broken
			valid = false;
			break;
		}
	}

	return valid;
}

function isSubstituted(nodeTemplateName, requirementName) {
	for (var vertexId in clout.vertexes) {
		var vertex = clout.vertexes[vertexId];
		if (tosca.isTosca(vertex, 'Substitution')) {
			for (var e = 0, l = vertex.edgesOut.length; e < l; e++) {
				var edge = vertex.edgesOut[e];
				if (!tosca.isTosca(edge, 'RequirementMapping'))
					continue;

				if ((edge.target.properties.name === nodeTemplateName) && (edge.properties.target === requirementName))
					return true;
			}

			// There's only ever one substitution
			return false;
		}
	}

	return false;
}

function isMaxCountGreater(a, b) {
	if (a == -1)
		return b !== -1;
	else if (b == -1)
		return false;
	return a > b;
}

function unsatisfied(location, name, message) {
	if (typeof problems === 'undefined')
		throw puccini.sprintf('%s: could not satisfy "%s" because %s', location.path, name, message);
	else
		problems.reportFull(11, 'Resolution', location.path, puccini.sprintf('could not satisfy "%s" because %s', name, message), location.row, location.column);
}

function notEnoughRelationships(location, relationshipCount, minRelationshipCount) {
	if (typeof problems === 'undefined')
		throw puccini.sprintf('%s: not enough relationships: %d < %d', location.path, relationshipCount, minRelationshipCount);
	else
		problems.reportFull(11, 'Resolution', location.path, puccini.sprintf('not enough relationships: %d < %d', relationshipCount, minRelationshipCount), location.row, location.column);
}
`
}
