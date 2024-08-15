package admin

import (
	"strings"

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

func (a *Actions) prepareExportCommand() []string {
	a.context.t.Helper()
	args := []string{"export"}

	if a.context.applicationNamespaces != nil {
		args = append(args, "--application-namespaces", strings.Join(a.context.applicationNamespaces, ","))
	}

	if a.context.applicationsetNamespaces != nil {
		args = append(args, "--applicationset-namespaces", strings.Join(a.context.applicationsetNamespaces, ","))
	}

	return args
}

func (a *Actions) RunExport() *Actions {
	a.context.t.Helper()
	a.runCli(a.prepareExportCommand()...)
	return a
}

func (a *Actions) IgnoreErrors() *Actions {
	a.ignoreErrors = true
	return a
}

func (a *Actions) DoNotIgnoreErrors() *Actions {
	a.ignoreErrors = false
	return a
}

func (a *Actions) prepareSetPasswordArgs(account string) []string {
	a.context.t.Helper()
	return []string{
		"account", "update-password", "--account", account, "--current-password", fixture.AdminPassword, "--new-password", fixture.DefaultTestUserPassword,
	}
}

func (a *Actions) runCli(args ...string) {
	a.context.t.Helper()
	a.lastOutput, a.lastError = RunCli(args...)
}

func (a *Actions) Then() *Consequences {
	a.context.t.Helper()
	return &Consequences{a.context, a}
}
