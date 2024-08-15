package e2e

import (
	application "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	adminFixture "github.com/argoproj/argo-cd/v2/test/e2e/fixture/admin"
	"github.com/argoproj/argo-cd/v2/test/e2e/fixture/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdminExport(t *testing.T) {
	ctx := adminFixture.Given(t)

	app.GivenWithSameState(t).
		Name("export-cmd").
		Path("guestbook-logs").
		When().
		CreateApp().
		Then().
		And(func(app *application.Application) {
			assert.NotNil(t, app)
		})
	
	ctx.
		When().
		Export().
		Then().
		AndCLIOutput(func(output string, err error) {
			assert.Contains(t, output, "name: export-cmd")
			
			// DO import here
			
			
		})
}
