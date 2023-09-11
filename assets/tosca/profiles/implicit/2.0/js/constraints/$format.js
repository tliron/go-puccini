
exports.validate = function(v, format) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	if (!util.isType(v, 'ard.string'))
		return 'not a string';
	try {
		ard.validateFormat(v, format);
	} catch (x) {
		if (x.value && x.value.error)
			// Unwrap Go error
			return x.value.error();
		else
			throw x;
	}
	return true;
};
