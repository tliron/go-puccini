
const traversal = require('tosca.lib.traversal');
const tosca = require('tosca.lib.utils');

const enforceCapabilityOccurrences = !traversal.hasQuirk(clout, 'capabilities.occurrences.permissive');

// Remove existing relationships
let nodeTemplateVertexes = [];
for (let vertexId in clout.vertexes) {
	let vertex = clout.vertexes[vertexId];
	if (tosca.isNodeTemplate(vertex)) {
		nodeTemplateVertexes.push(vertex);
		let remove = [];
		for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
			let edge = vertex.edgesOut[e];
			if (tosca.isTosca(edge, 'Relationship'))
				remove.push(edge);
		}
		for (let e = 0, l = remove.length; e < l; e++)
			remove[e].remove();
	}
}

// For consistent results, we will sort the node templates by name
nodeTemplateVertexes.sort(function(a, b) {
	return a.properties.name < b.properties.name ? -1 : 1;
});

traversal.toCoercibles();

// Resolve all requirements
for (let v = 0, l = nodeTemplateVertexes.length; v < l; v++) {
	let vertex = nodeTemplateVertexes[v];
	let nodeTemplate = vertex.properties;
	let requirements = nodeTemplate.requirements;
	for (let r = 0, ll = requirements.length; r < ll; r++) {
		let requirement = requirements[r];
		resolve(vertex, nodeTemplate, requirement);
	}
}

if (enforceCapabilityOccurrences)
	for (let v = 0, l = nodeTemplateVertexes.length; v < l; v++) {
		let vertex = nodeTemplateVertexes[v];
		let nodeTemplate = vertex.properties;
		let capabilities = nodeTemplate.capabilities;
		for (let capabilityName in capabilities) {
			let capability = capabilities[capabilityName];
			let relationshipCount = countRelationships(vertex, capabilityName);
			let minRelationshipCount = capability.minRelationshipCount;
			if (relationshipCount < minRelationshipCount)
				notEnoughRelationships(capability.location, relationshipCount, minRelationshipCount)
		}
	}

traversal.unwrapCoercibles();

if (puccini.arguments.history !== 'false')
	tosca.addHistory('resolve');
puccini.write(clout)

function resolve(sourceVertex, sourceNodeTemplate, requirement) {
	let location = requirement.location;
	let name = requirement.name;

	if (isSubstituted(sourceNodeTemplate.name, name)) {
		puccini.log.debugf('%s: skipping because in substitution mappings', location.path)
		return;
	}

	let candidates = gatherCandidateNodeTemplates(sourceVertex, requirement);
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
	let priorityCandidates = [];
	for (let c = 0, l = candidates.length; c < l; c++) {
		let candidate = candidates[c];
		if ((candidate.capability.minRelationshipCount !== 0) && (countRelationships(candidate.vertex, candidate.capabilityName) < candidate.capability.minRelationshipCount))
			priorityCandidates.push(candidate);
	}

	let chosen = null;

	if (priorityCandidates.length !== 0)
		// Of the priority candidates, pick the one with the highest minimum relationship count
		// (needs to be fulfilled soonest)
		for (let c = 0, l = priorityCandidates.length; c < l; c++) {
			let candidate = priorityCandidates[c];
			if ((chosen === null) || (candidate.capability.minRelationshipCount > chosen.capability.minRelationshipCount))
				chosen = candidate;
		}
	else
		// Of the candidates, pick the one with highest maximum relationship count
		// (has the most room)
		for (let c = 0, l = candidates.length; c < l; c++) {
			let candidate = candidates[c];
			if ((chosen === null) || isMaxCountGreater(candidate.capability.maxRelationshipCount, chosen.capability.maxRelationshipCount))
				chosen = candidate;
		}

	puccini.log.debugf('%s: satisfied %q with capability %q in node template %q', location.path, name, chosen.capabilityName, chosen.nodeTemplateName);
	addRelationship(sourceVertex, requirement, chosen.vertex, chosen.capabilityName);
}

function gatherCandidateNodeTemplates(sourceVertex, requirement) {
	let path = requirement.location.path;
	let nodeTemplateName = requirement.nodeTemplateName;
	let nodeTypeName = requirement.nodeTypeName;
	let nodeTemplatePropertyValidators = requirement.nodeTemplatePropertyValidators;
	let capabilityPropertyValidatorsMap = requirement.capabilityPropertyValidators;

	let candidates = [];
	for (let v = 0, l = nodeTemplateVertexes.length; v < l; v++) {
		let vertex = nodeTemplateVertexes[v];
		let candidateNodeTemplate = vertex.properties;
		let candidateNodeTemplateName = candidateNodeTemplate.name;

		if ((nodeTemplateName !== '') && (nodeTemplateName !== candidateNodeTemplateName)) {
			puccini.log.debugf('%s: node template %q is not named %q', path, candidateNodeTemplateName, nodeTemplateName);
			continue;
		}

		if ((nodeTypeName !== '') && !(nodeTypeName in candidateNodeTemplate.types)) {
			puccini.log.debugf('%s: node template %q is not of type %q', path, candidateNodeTemplateName, nodeTypeName);
			continue;
		}

		// Node filter
		if ((nodeTemplatePropertyValidators.length !== 0) && !arePropertiesValid(path, sourceVertex, 'node template', candidateNodeTemplateName, candidateNodeTemplate, nodeTemplatePropertyValidators)) {
			puccini.log.debugf('%s: properties of node template %q do not validate', path, candidateNodeTemplateName);
			continue;
		}

		let candidateCapabilities = candidateNodeTemplate.capabilities;

		// Capability filter
		if (capabilityPropertyValidatorsMap.length !== 0) {
			let valid = true;
			for (let candidateCapabilityName in candidateCapabilities) {
				let candidateCapability = candidateCapabilities[candidateCapabilityName];

				// Try by name
				let capabilityPropertyValidators = capabilityPropertyValidatorsMap[candidateCapabilityName];
				if (capabilityPropertyValidators === undefined) {
					// Try by type name
					for (let candidateTypeName in candidateCapability.types) {
						capabilityPropertyValidators = capabilityPropertyValidatorsMap[candidateTypeName];
						if (capabilityPropertyValidators !== undefined) break;
					}
				}

				if ((capabilityPropertyValidators !== undefined) && (capabilityPropertyValidators.length !== 0) && !arePropertiesValid(path, sourceVertex, 'capability', candidateCapabilityName, candidateCapability, capabilityPropertyValidators)) {
					puccini.log.debugf('%s: properties of capability %q in node template %q do not validate', path, candidateCapabilityName, candidateNodeTemplateName);
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
	let path = requirement.location.path;
	let capabilityName = requirement.capabilityName;
	let capabilityTypeName = requirement.capabilityTypeName;

	let candidates = [];
	for (let c = 0, l = candidateNodeTemplates.length; c < l; c++) {
		let candidate = candidateNodeTemplates[c];
		let candidateVertex = candidate.vertex;
		let candidateNodeTemplateName = candidate.nodeTemplateName;

		let candidateCapabilities = [];
		for (let candidateCapabilityName in candidate.capabilities) {
			candidateCapabilities.push({
				name: candidateCapabilityName,
				capability: candidate.capabilities[candidateCapabilityName]
			});
		}

		// For consistent results, we will sort the candidate capabilities by name
		candidateCapabilities.sort(function(a, b) {
			return a.name < b.name ? -1 : 1;
		});

		for (let cc = 0, ll = candidateCapabilities.length; cc < ll; cc++) {
			let candidateCapabilityName = candidateCapabilities[cc].name;

			if ((capabilityName !== '') && (capabilityName !== candidateCapabilityName)) {
				puccini.log.debugf('%s: capability %q in node template %q is not named %q', path, candidateCapabilityName, candidateNodeTemplateName, capabilityName);
				continue;
			}

			let candidateCapability = candidateCapabilities[cc].capability;

			if ((capabilityTypeName !== '') && !(capabilityTypeName in candidateCapability.types)) {
				puccini.log.debugf('%s: capability %q in node template %q is not of type %q', path, candidateCapabilityName, candidateNodeTemplateName, capabilityTypeName);
				continue;
			}

			if (enforceCapabilityOccurrences) {
				let maxRelationshipCount = candidateCapability.maxRelationshipCount;
				if ((maxRelationshipCount !== -1) && (countRelationships(candidateVertex, candidateCapabilityName) === maxRelationshipCount)) {
					puccini.log.debugf('%s: capability %q in node template %q already has %d relationships, the maximum allowed', path, candidateCapabilityName, candidateNodeTemplateName, maxRelationshipCount);
					continue;
				}
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
	let edge = sourceVertex.newEdgeTo(targetVertex);
	edge.metadata['puccini'] = {
		version: '1.0',
		kind: 'Relationship'
	};

	let relationship = requirement.relationship;
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
	let count = 0;
	for (let e = 0, l = vertex.edgesIn.size(); e < l; e++) {
		let edge = vertex.edgesIn[e];
		if (tosca.isTosca(edge, 'Relationship') && (edge.properties.capability === capabilityName))
			count++;
	}
	return count;
}

function arePropertiesValid(path, sourceVertex, kind, name, entity, validatorsMap) {
	let valid = true;

	let properties = entity.properties;
	for (let propertyName in validatorsMap) {
		puccini.log.debugf('%s: applying validators to property %q of %s %q', path, propertyName, kind, name);

		let property = properties[propertyName];
		if (property === undefined) {
			// return false; GOJA: returning from inside for-loop is broken
			valid = false;
			break;
		}

		let validators = validatorsMap[propertyName];
		validators = clout.newValidators(validators, sourceVertex, sourceVertex, entity)
		if (!validators.isValid(property)) {
			// return false; GOJA: returning from inside for-loop is broken
			valid = false;
			break;
		}
	}

	return valid;
}

function isSubstituted(nodeTemplateName, requirementName) {
	for (let vertexId in clout.vertexes) {
		let vertex = clout.vertexes[vertexId];
		if (tosca.isTosca(vertex, 'Substitution')) {
			for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
				let edge = vertex.edgesOut[e];
				if (!tosca.isTosca(edge, 'RequirementPointer'))
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
		throw util.sprintf('%s: could not satisfy %q because %s', location.path, name, message);
	else
		problems.reportFull(11, 'Resolution', location.path, util.sprintf('could not satisfy %q because %s', name, message), location.row, location.column);
}

function notEnoughRelationships(location, relationshipCount, minRelationshipCount) {
	if (typeof problems === 'undefined')
		throw util.sprintf('%s: not enough relationships: %d < %d', location.path, relationshipCount, minRelationshipCount);
	else
		problems.reportFull(11, 'Resolution', location.path, util.sprintf('not enough relationships: %d < %d', relationshipCount, minRelationshipCount), location.row, location.column);
}
