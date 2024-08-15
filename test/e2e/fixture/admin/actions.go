package admin

import (
	"github.com/argoproj/argo-cd/v2/test/e2e/fixture"
)

// this implements the "when" part of given/when/then
//
// none of the func implement error checks, and that is complete intended, you should check for errors
// using the Then()
type Actions struct {
	context      *Context
	ignoreErrors bool
	lastOutput   string
	lastError    error
}

func (a *Actions) prepareExportArgs() []string {
	return []string{
		"admin", "export",
	}
}

func (a *Actions) Export() *Actions {
	a.lastOutput, a.lastError = fixture.RunAdminCli(a.prepareExportArgs()...)
	return a
}

func (a *Actions) Then() *Consequences {
	return &Consequences{
		context: a.context,
		actions: a,
	}
}