package js

//
// ExecutionContext
//

type ExecutionContext struct {
	CloutContext *CloutContext
	Site         any
	Source       any
	Target       any
}

func (self *CloutContext) NewExecutionContext(site any, source any, target any) *ExecutionContext {
	return &ExecutionContext{
		CloutContext: self,
		Site:         site,
		Source:       source,
		Target:       target,
	}
}

func (self *ExecutionContext) Call(scriptletName string, functionName string, arguments ...any) (any, error) {
	if exports, err := self.CloutContext.JSContext.Environment.Require(scriptletName); err == nil {
		return self.CloutContext.JSContext.Environment.GetAndCall(exports, functionName, self, arguments...)
	} else {
		return nil, err
	}
}
