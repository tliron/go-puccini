
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

const tosca = require('tosca.lib.utils');

exports.validate = function(v, lower, upper) {
	if (arguments.length !== 3)
		throw 'must have 2 arguments';
	return (tosca.compare(v, lower) >= 0) && (tosca.compare(v, upper) <= 0);
};
