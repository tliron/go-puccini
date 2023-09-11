
const traversal = require('tosca.lib.traversal');
const tosca = require('tosca.lib.utils');

traversal.coerce();

if (tosca.isTosca(clout))
    transcribe.output(clout.properties.tosca.outputs);
