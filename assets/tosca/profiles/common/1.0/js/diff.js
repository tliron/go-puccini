
const traversal = require('tosca.lib.traversal');
const tosca = require('tosca.lib.utils');

if (!env.arguments.base) {
    throw 'must provide "base" argument';
}

let base = clout.load(env.arguments.base);

traversal.coerce();
traversal.coerce(base);

transcribe.output(diff(clout, base));

function diff(clout, base) {
    let nodes = gatherNodeTemplates(clout);
    let baseNodes = gatherNodeTemplates(base);
    
    let diff = {
        added: [],
        removed: []
    };
    
    for (let n = 0, l = nodes.length; n < l; n++) {
        let node = nodes[n];
        if (!hasNode(baseNodes, node))
            diff.added.push(node.name);
    }
    
    for (let n = 0, l = baseNodes.length; n < l; n++) {
        let node = baseNodes[n];
        if (!hasNode(nodes, node))
            diff.removed.push(node.name);
    }

    return diff;
}

function gatherNodeTemplates(clout) {
    let nodeTemplates = [];
    for (let vertexId in clout.vertexes) {
        let vertex = clout.vertexes[vertexId];
        if (tosca.isNodeTemplate(vertex)) {
            let nodeTemplate = vertex.properties;
            nodeTemplates.push(nodeTemplate);
        }
    }
    return nodeTemplates;
}

function hasNode(nodes, node) {
    for (let n = 0, l = nodes.length; n < l; n++)
        if (nodes[n].name === node.name)
            return true;
    return false;
}
