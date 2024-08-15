package admin

// this implements the "then" part of given/when/then
type Consequences struct {
	context *Context
	actions *Actions
}

func (c *Consequences) And() *Consequences {
	c.context.t.Helper()
	return c
}

func (c *Consequences) AndCLIOutput(block func(output string, err error)) *Consequences {
	block(c.actions.lastOutput, c.actions.lastError)
	return c
}

func (c *Consequences) Given() *Context {
	return c.context
}

func (c *Consequences) When() *Actions {
	return c.actions
}
