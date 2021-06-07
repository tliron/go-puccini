// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/common/1.0/js/coerce.js"] = `

const traversal = require('tosca.lib.traversal');
const tosca = require('tosca.lib.utils');

traversal.coerce();
if (puccini.arguments.history !== 'false')
	tosca.addHistory('coerce');
puccini.write(clout);
`
}
