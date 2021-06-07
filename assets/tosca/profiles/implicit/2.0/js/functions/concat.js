
// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.3.1

exports.evaluate = function() {
	var a = [];
	var length = arguments.length;
	for (var i = 0; i < length; i++) {
		var argument = arguments[i];
		if (argument.$string !== undefined)
			argument = argument.$string;
		a.push(argument);
	}
	return a.join('');
};
