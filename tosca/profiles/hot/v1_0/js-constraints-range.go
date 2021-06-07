// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/hot/1.0/js/constraints/range.js"] = `

const tosca = require('tosca.lib.utils');

exports.validate = function(v, bounds) {
	if (arguments.length !== 2)
		throw 'must have 1 arguments';
	if ((bounds.min === undefined) && (bounds.max === undefined))
		throw 'must provide "min" and/or "max"';
	v = tosca.getComparable(v);
	if (bounds.min !== undefined)
		if (tosca.compare(v, bounds.min) < 0)
			return false;
	if (bounds.max !== undefined)
		if (tosca.compare(v, bounds.max) > 0)
			return false;
	return true;
};
`
}
