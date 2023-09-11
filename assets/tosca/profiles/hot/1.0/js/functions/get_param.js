
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#get_param]

const tosca = require('tosca.lib.utils');

exports.evaluate = function(input) {
	if (arguments.length !== 1)
		throw 'must have 1 argument';
	if (!tosca.isTosca(clout))
		throw 'Clout is not TOSCA';
	let inputs = clout.properties.tosca.inputs;
	if (!(input in inputs))
		throw util.sprintf('parameter %q not found', input);
	let r = inputs[input];
	r = clout.coerce(r);
	return r;
};
