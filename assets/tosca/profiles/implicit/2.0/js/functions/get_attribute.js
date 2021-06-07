
// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.5.1

const tosca = require('tosca.lib.utils');

exports.evaluate = function(entity, first) {
	return tosca.getNestedValue.call(this, 'attribute', 'attributes', arguments);
};
