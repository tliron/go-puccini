// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/hot/1.0/js/functions/list_join.js"] = `

// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#list_join]

exports.evaluate = function() {
	var length = arguments.length;
	if (length < 1)
		throw 'must have at least 1 arguments';
	var a = [];
	for (var i = 1; i < length; i++) {
		var argument = arguments[i];
		if (argument.$string !== undefined)
			argument = argument.$string;
		a.push(argument);
	}
	return a.join(arguments[0]);
};
`
}
