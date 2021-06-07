
exports.evaluate = function(v) {
    if (arguments.length !== 1)
        throw 'must have 1 argument';
	return v * 2;
};
