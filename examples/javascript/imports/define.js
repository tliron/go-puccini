
const traversal = require('tosca.lib.traversal');

// The "clout.define" API accepts the scriptlet source code as text
clout.define('tosca.function.get_input', "\n\
const tosca = require('tosca.lib.utils');\n\
\n\
// This is a copy of the built-in get_input function source\n\
// Except that we added a '* 2' to the returned result\n\
exports.evaluate = function(input) {\n\
	if (arguments.length !== 1)\n\
		throw 'must have 1 argument';\n\
	if (!tosca.isTosca(clout))\n\
		throw 'Clout is not TOSCA';\n\
	var inputs = clout.properties.tosca.inputs;\n\
	if (!(input in inputs))\n\
		throw puccini.sprintf('input \"%s\" not found', input);\n\
	var r = inputs[input];\n\
	r = clout.coerce(r);\n\
	return r * 2;\n\
};\n\
");

traversal.coerce();

puccini.write(clout);
