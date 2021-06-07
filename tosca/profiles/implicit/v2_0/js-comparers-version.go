// This file was auto-generated from a YAML file

package v2_0

func init() {
	Profile["/tosca/implicit/2.0/js/comparers/version.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.2.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.2.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.2.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.2.2

exports.compare = function(a, b) {
	if (a.$comparer !== b.$comparer)
		throw 'both values must be of type "version"';
	if (a.major !== b.major)
		return a.major < b.major ? -1 : 1;
	if (a.minor !== b.minor)
		return a.minor < b.minor ? -1 : 1;
	if (a.fix !== b.fix)
		return a.fix < b.fix ? -1 : 1;
	var aq = a.qualifier.toLowerCase();
	var bq = b.qualifier.toLowerCase();
	if (aq !== bq) // note: the qualifier is compared alphabetically, *not* semantically
		return aq < bq ? -1 : 1;
	if (a.build !== b.build)
		return a.build < b.build ? -1 : 1;
	return 0;
};
`
}
