
exports.isTosca = function(o, kind) {
	if (o.metadata === undefined)
		return false;
	o = o.metadata['puccini'];
	if (o === undefined)
		return false;
	if (o.version !== '1.0')
		return false;
	if (kind !== undefined)
		return kind === o.kind;
	return true;
};

exports.isNodeTemplate = function(vertex, typeName) {
	if (exports.isTosca(vertex, 'NodeTemplate')) {
		if (typeName !== undefined)
			return typeName in vertex.properties.types;
		return true;
	}
	return false;
};

exports.setOutputValue = function(name, value) {
	if (clout.properties.tosca === undefined)
		return false;
	let output = clout.properties.tosca.outputs[name];
	if (output === undefined)
		return false;

	if (output.$type && output.$type.type)
		switch (output.$type.type.name) {
		case 'boolean':
			value = (value === 'true');
			break;
		case 'integer':
			value = parseInt(value);
			break;
		case 'float':
			value = parseFloat(value);
			break;
		}

	output.$value = value;
	return true;
};

exports.getPolicyTargets = function(vertex) {
	let targets = [];

	function addTarget(target) {
		for (let t = 0, l = targets.length; t < l; t++)
			if (targets[t].name === target.name)
				return;
		targets.push(target);
	}

	for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
		let edge = vertex.edgesOut[e];
		if (exports.isTosca(edge, 'NodeTemplateTarget'))
			targets.push(clout.vertexes[edge.targetID].properties);
		else if (toexportssca.isTosca(edge, 'GroupTarget')) {
			let members = exports.getGroupMembers(clout.vertexes[edge.targetID]);
			for (let m = 0, ll = members.length; m < ll; m++)
				addTarget(members[m])
		}
	}
	return targets;
};

exports.getGroupMembers = function(vertex) {
	let members = [];
	for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
		let edge = vertex.edgesOut[e];
		if (exports.isTosca(edge, 'Member'))
			members.push(clout.vertexes[edge.targetID].properties);
	}
	return members;
};

exports.addHistory = function(description) {
	let metadata = clout.metadata;
	if (metadata === undefined)
		metadata = clout.metadata = {};
	let history = metadata.history;
	if (history === undefined)
		history = [];
	else
		history = history.slice(0);
	history.push({
		timestamp: util.nowString(),
		description: description
	});
	metadata.history = history;
};

exports.getNestedValue = function(singular, plural, args) {
	args = Array.prototype.slice.call(args);
	let length = args.length;
	if (length < 2)
		throw 'must have at least 2 arguments';
	let vertex = exports.getModelableEntity.call(this, args[0]);
	let nodeTemplate = vertex.properties;
	let value = nodeTemplate[plural];
	let a = 1;
	let arg = args[a];
	let nextArg = args[a+1];
	let count = 0;
	if (arg in nodeTemplate.capabilities) {
		value = nodeTemplate.capabilities[arg][plural];
		singular = util.sprintf('capability %q %s', arg, singular);
		arg = args[++a];
	} else for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
		let edge = vertex.edgesOut[e];
		if (!exports.isTosca(edge, 'Relationship'))
			continue;
		let relationship = edge.properties;
		if (relationship.name === arg)
			if (count++ === nextArg) {
				value = relationship[plural];
				singular = util.sprintf('relationship %q %s', arg, singular);
				a += 2;
				arg = args[a];
				break;
			}
	}
	if ((typeof value === 'object') && (value !== null) && (arg in value))
		value = value[arg];
	else
		throw util.sprintf('%s %q not found in %q', singular, arg, nodeTemplate.name);
	value = clout.coerce(value);
	for (let i = a + 1; i < length; i++) {
		arg = args[i];
		if ((typeof value === 'object') && (value !== null) && (arg in value))
			value = value[arg];
		else
			throw util.sprintf('nested %s %q not found in %q', singular, args.slice(a, i+1).join('.'), nodeTemplate.name);
	}
	return value;
};

exports.getModelableEntity = function(entity) {
	let vertex;
	switch (entity) {
	case 'SELF':
		if (!this || !this.site)
			throw util.sprintf('%q cannot be used in this context', entity);
		vertex = this.site;
		break;
	case 'SOURCE':
		if (!this || !this.source)
			throw util.sprintf('%q cannot be used in this context', entity);
		vertex = this.source;
		break;
	case 'TARGET':
		if (!this || !this.target)
			throw util.sprintf('%q cannot be used in this context', entity);
		vertex = this.target;
		break;
	case 'HOST':
		if (!this || !this.site)
			throw util.sprintf('%q cannot be used in this context', entity);
		vertex = exports.getHost(this.site);
		break;
	default:
		for (let vertexId in clout.vertexes) {
			let vertex = clout.vertexes[vertexId];
			if (exports.isNodeTemplate(vertex) && (vertex.properties.name === entity))
				return vertex;
		}
		vertex = {};
	}
	if (exports.isNodeTemplate(vertex))
		return vertex;
	else
		throw util.sprintf('%q node template not found', entity);
};

exports.getHost = function(vertex) {
	for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
		let edge = vertex.edgesOut[e];
		if (exports.isTosca(edge, 'Relationship')) {
			for (let typeName in edge.properties.types) {
				let type = edge.properties.types[typeName];
				if (type.metadata.role === 'host')
					return edge.target;
			}
		}
	}
	if (exports.isNodeTemplate(vertex))
		throw util.sprintf('"HOST" not found for node template %q', vertex.properties.name);
	else
		throw '"HOST" not found';
};

exports.getComparable = function(v) {
	if ((v === undefined) || (v === null))
		return null;
	let c = v.$number;
	if (c !== undefined)
		return c;
	c = v.$string;
	if (c !== undefined)
		return c;
	return v;
};

exports.getLength = function(v) {
	if (v.$string !== undefined)
		v = v.$string;
	let length = v.length;
	if (length === undefined)
		length = Object.keys(v).length;
	return length;
};

exports.compare = function(v1, v2) {
	let c = v1.$comparer;
	if (c === undefined)
		c = v2.$comparer;
	if (c !== undefined)
		return clout.call(c, 'compare', [v1, v2]);
	v1 = exports.getComparable(v1);
	v2 = exports.getComparable(v2);
	if (v1 == v2)
		return 0;
	else if (v1 < v2)
		return -1;
	else
		return 1;
};

// See: https://stackoverflow.com/a/45683145
exports.deepEqual = function(v1, v2) {
	if (v1 === v2)
		return true;

	if (exports.isPrimitive(v1) && exports.isPrimitive(v2))
		return v1 === v2;

	if (Object.keys(v1).length !== Object.keys(v2).length)
		return false;

	for (let key in v1) {
		if (!(key in v2)) return false;
		if (!exports.deepEqual(v1[key], v2[key])) return false;
	}

	return true;
};

exports.isPrimitive = function(obj) {
	return obj !== Object(obj);
};
