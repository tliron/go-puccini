
const traversal = require('tosca.lib.traversal');
const tosca = require('tosca.lib.utils');

traversal.coerce();

// From: cdnjs.com
let jQueryVersion = '3.6.0';
let jQueryUiVersion = '1.12.1';

// From: jsdeliver.net
let jQueryLayoutVersion = '1.8.5';
let visJsVersion = '9.0.4';

let colorNode = 'rgb(100,200,255)';
let colorPolicy = 'rgb(255,165,0)';
let colorSubstitution = 'rgb(150,200,255)';
let colorWorkflow = 'rgb(100,255,100)';

let nodes = [];
let edges = [];
let header;

if (tosca.isTosca(clout)) {
	let templateName = clout.properties.tosca.metadata.template_name;
	let templateAuthor = clout.properties.tosca.metadata.template_author;
	let templateVersion = clout.properties.tosca.metadata.template_version;
	let description = clout.properties.tosca.description;

	header = '<h1>Clout from TOSCA Service Template</h1>';
	if (templateName)
		header += '<h2>' + escapeHtml(templateName) + '</h2>';
	if (templateVersion)
		header += '<p><b>Version</b>: ' + escapeHtml(templateVersion) + '</p>';
	if (templateAuthor)
		header += '<p><b>Author</b>: ' + escapeHtml(templateAuthor) + '</p>';
	if (description)
		header += '<p>' + formatDescription(description) + '</p>';
} else {
	header = '<h1>Clout</h1>';
}

for (let id in clout.vertexes) {
	let vertex = clout.vertexes[id];
	addVertex(id, vertex);
}

function formatDescription(description) {
	let r = '';
	let paragraphs = description.split('\n');
	for (let p in paragraphs) {
		let paragraph = paragraphs[p];
		if (paragraph)
			r += '<p>' + escapeHtml(paragraph) + '</p>';
	}
	return r;
}

function escapeHtml(unsafe) {
	// See: https://stackoverflow.com/a/6234804/849021
	return unsafe
		.replace(/&/g, "&amp;")
		.replace(/</g, "&lt;")
		.replace(/>/g, "&gt;")
		.replace(/"/g, "&quot;")
		.replace(/'/g, "&#039;");
}

function jsonify(data) {
	return JSON.stringify(data, null, '\t').replace(/^/mg, '\t').substr(1);
}

function addVertex(id, vertex) {
	let node = {
		id: id,
		label: id,
		data: tosca.isTosca(vertex) ? vertex.properties : vertex
	};

	if (tosca.isTosca(vertex, 'NodeTemplate'))
		addNodeTemplate(node);
	else if (tosca.isTosca(vertex, 'Group'))
		addGroup(node);
	else if (tosca.isTosca(vertex, 'Policy'))
		addPolicy(node);
	else if (tosca.isTosca(vertex, 'Substitution'))
		addSubstitution(node);
	else if (tosca.isTosca(vertex, 'Workflow'))
		addWorkflow(node);
	else if (tosca.isTosca(vertex, 'WorkflowStep'))
		addWorkflowStep(node);
	else if (tosca.isTosca(vertex, 'WorkflowActivity'))
		addWorkflowActivity(node);
	else
		node.data = vertex;

	nodes.push(node);

	for (let e = 0, l = vertex.edgesOut.size(); e < l; e++)
		addEdge(id, vertex.edgesOut[e]);
}

function addEdge(id, e) {
	let edge = {
		from: id,
		to: e.targetID,
		arrows: {
			to: true
		},
		font: {
			align: 'middle'
		},
		smooth: {type: 'dynamic'},
		length: 300,
		data: tosca.isTosca(e) ? e.properties : e
	};

	if (tosca.isTosca(e, 'Relationship'))
		addRelationship(edge);
	else if (tosca.isTosca(e, 'RequirementMapping'))
		addRequirementMapping(edge);
	else if (tosca.isTosca(e, 'CapabilityMapping'))
		addCapabilityMapping(edge);
	else if (tosca.isTosca(e, 'PropertyMapping'))
		addPropertyMapping(edge);
	else if (tosca.isTosca(e, 'InterfaceMapping'))
		addInterfaceMapping(edge);
	else if (tosca.isTosca(e, 'OnFailure'))
		addOnFailure(edge);
	else
		edge.data = e;

	edges.push(edge);
}

function addNodeTemplate(node) {
	node.label = node.data.name;
	node.shape = 'box';
	node.color = colorNode;
}

function addGroup(node) {
	node.label = node.data.name;
	node.shape = 'circle';
	node.color = colorNode;
}

function addPolicy(node) {
	node.label = node.data.name;
	node.shape = 'circle';
	node.color = colorPolicy;
}

function addRelationship(edge) {
	edge.label = edge.data.name;
	edge.color = {color: colorNode};
}

function addSubstitution(node) {
	node.label = 'substitution';
	node.shape = 'box';
	node.color = colorSubstitution;
	node.shapeProperties = {borderDashes: true};
}

function addRequirementMapping(edge) {
	edge.label = 'requirement\n' + edge.data.requirement;
	edge.color = {color: colorSubstitution};
	edge.dashes = true;
}

function addCapabilityMapping(edge) {
	edge.label = 'capability\n' + edge.data.capability;
	edge.color = {color: colorSubstitution};
	edge.dashes = true;
}

function addPropertyMapping(edge) {
	edge.label = 'property\n' + edge.data.property;
	edge.color = {color: colorSubstitution};
	edge.dashes = true;
}

function addInterfaceMapping(edge) {
	edge.label = 'interface\n' + edge.data.interface;
	edge.color = {color: colorSubstitution};
	edge.dashes = true;
}

function addWorkflow(node) {
	node.label = node.data.name;
	node.shape = 'diamond';
	node.color = colorWorkflow;
}

function addWorkflowStep(node) {
	node.label = node.data.name;
	node.shape ='diamond';
	node.color = colorWorkflow;
}

function addOnFailure(edge) {
	edge.label = 'onFailure';
	edge.color = {color: 'rgb(255,100,100)'};
}

function addWorkflowActivity(node) {
	node.label = node.data.name;
	node.shape = 'triangle';
	node.color = colorWorkflow;
}

let template = '\
<!doctype html>\n\
<html>\n\
<head>\n\
	<title>Clout from TOSCA Service Template</title>\n\
	<meta charset="utf-8"/>\n\
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/%s/jquery.min.js" type="text/javascript"></script>\n\
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jqueryui/%s/jquery-ui.min.js" type="text/javascript"></script>\n\
	<link href="https://cdnjs.cloudflare.com/ajax/libs/jqueryui/%s/jquery-ui.min.css" rel="stylesheet" type="text/css" />\n\
	<script src="https://cdn.jsdelivr.net/npm/layout-jquery3@%s/dist/jquery.layout_and_plugins.min.js" type="text/javascript"></script>\n\
	<link href="https://cdn.jsdelivr.net/npm/layout-jquery3@%s/dist/layout-default.css" rel="stylesheet" type="text/css" />\n\
	<script type="text/javascript" src="https://cdn.jsdelivr.net/npm/vis-network@%s/standalone/umd/vis-network.min.js" type="text/javascript"></script>\n\
	<link href="https://cdn.jsdelivr.net/npm/vis-network@%s/styles/vis-network.min.css" rel="stylesheet" type="text/css" />\n\
	<link href="https://fonts.googleapis.com/css?family=Noto+Sans" rel="stylesheet" />\n\
	<link href="https://fonts.googleapis.com/css?family=Source+Code+Pro" rel="stylesheet" />\n\
	<style type="text/css">\n\
		body {\n\
			font-family: \'Noto Sans\', sans-serif;\n\
		}\n\
		.ui-layout-pane {\n\
			background-color: #F0F0F0;\n\
			padding: 5px 5px 5px 5px;\n\
		}\n\
		h1, h2, h3, p {\n\
			margin-top: 0;\n\
			margin-bottom: 8px;\n\
		}\n\
		#network {\n\
			width: 100%%;\n\
			height: 100%%;\n\
		}\n\
		#info {\n\
			font-family: \'Source Code Pro\', monospace;\n\
			white-space: pre-wrap;\n\
		}\n\
		#corner {\n\
			float: right;\n\
			font-size: small;\n\
			text-align: right;\n\
		}\n\
	</style>\n\
	<script type="text/javascript">\n\
$(document).ready(function () {\n\
	$(\'body\').layout({\n\
		north__resizable: false,\n\
		east__size: \'25%%\',\n\
		livePaneResizing: true\n\
	});\n\
	let nodes = new vis.DataSet(%s);\n\
	let edges = new vis.DataSet(%s);\n\
	let network = new vis.Network(\n\
		document.getElementById(\'network\'),\n\
		{\n\
			nodes: nodes,\n\
			edges: edges\n\
		},\n\
		{\n\
			layout: {\n\
				randomSeed: 1,\n\
				improvedLayout: true\n\
			}\n\
		}\n\
	);\n\
	network.on("click", function (params) {\n\
		if (params.nodes.length === 1) {\n\
			node = nodes.get(params.nodes[0]).data;\n\
			$(\'#info\').text(JSON.stringify(node, null, 4));\n\
		} else if (params.edges.length === 1) {\n\
			edge = edges.get(params.edges[0]).data;\n\
			$(\'#info\').text(JSON.stringify(edge, null, 4));\n\
		}\n\
	});\n\
});\n\
	</script>\n\
</head>\n\
<body>\n\
	<div class="ui-layout-north">\n\
		<div id="corner">Generated by <a href="https://puccini.cloud/">Puccini</a></div>\n\
		%s\n\
	</div>\n\
	<div class="ui-layout-center">\n\
		<div id="network"></div>\n\
	</div>\n\
	<div class="ui-layout-east">\n\
		<div id="info"></div>\n\
	</div>\n\
</body>\n\
</html>';

let html = util.sprintf(
	template,
	jQueryVersion,
	jQueryUiVersion, jQueryUiVersion,
	jQueryLayoutVersion, jQueryLayoutVersion,
	visJsVersion, visJsVersion,
	jsonify(nodes),
	jsonify(edges),
	header
);

transcribe.output(html);
