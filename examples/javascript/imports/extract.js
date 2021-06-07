// This scriptlet extracts all artifacts to the output directory

traversal = require('tosca.lib.traversal');
tosca = require('tosca.lib.utils');

traversal.coerce();

for (var vertexId in clout.vertexes) {
	var vertex = clout.vertexes[vertexId];
	if (!tosca.isNodeTemplate(vertex))
		continue;
	var nodeTemplate = vertex.properties;

	for (var key in nodeTemplate.artifacts) {
		var artifact = nodeTemplate.artifacts[key];

		// If 'puccini.output' is empty, this will be relative to current directory
		var targetPath = puccini.joinFilePath(puccini.output, artifact.filename);

		puccini.log.noticef('extracting "%s" to "%s"', artifact.sourcePath, targetPath);
		puccini.download(artifact.sourcePath, targetPath);

		//puccini.log.noticef('%s', puccini.exec('cat', targetPath));
	}
}
