
// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.3.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.3.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.3.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.3.2

exports.evaluate = function(v, separators, index) {
	if (arguments.length !== 3)
		throw 'must have 3 arguments';
	if (v.$string !== undefined)
		v = v.$string;
	var s = v.split(new RegExp('[' + escape(separators) + ']'));
	return s[index];
};

function escape(s) {
	return s.replace(/[\-\[\]\/\{\}\(\)\*\+\?\.\\\^\$\|]/g, '\\$&');
}
