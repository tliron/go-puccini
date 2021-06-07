
exports.validate = function(v) {
	var values = Array.prototype.slice.call(arguments, 1);
	return values.indexOf(v) !== -1;
};
