
// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.4.2

const tosca = require('tosca.lib.utils');

exports.evaluate = function() {
	return tosca.getNestedValue.call(this, 'property', 'properties', arguments);
};
