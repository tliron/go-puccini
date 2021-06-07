// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/hot/1.0/js/constraints/allowed_values.js"] = `

exports.validate = function(v) {
	var values = Array.prototype.slice.call(arguments, 1);
	return values.indexOf(v) !== -1;
};
`
}
