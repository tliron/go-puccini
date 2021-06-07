// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/hot/1.0/js/constraints/allowed_pattern.js"] = `

exports.validate = function(v, re) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	if (v.$string !== undefined)
		v = v.$string;
	return new RegExp('^' + re + '$').test(v);
};
`
}
