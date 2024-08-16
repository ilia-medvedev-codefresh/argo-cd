package e2e

import (
	"testing"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/stretchr/testify/assert"

	"github.com/argoproj/argo-cd/v2/test/e2e/fixture"
	fixtureutils "github.com/argoproj/argo-cd/v2/test/e2e/fixture/admin/utils"
	appfixture "github.com/argoproj/argo-cd/v2/test/e2e/fixture/app"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	. "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	. "github.com/argoproj/argo-cd/v2/test/e2e/fixture/admin"
)

func TestBackupExportImport(t *testing.T) {
	ctx := Given(t)

	// Create application in argocd namespace
	appctx := appfixture.GivenWithSameState(t)

	//var exportRawOutput string

	// Create application in test namespace
	appctx.
		Path(guestbookPath).
		Name("exported-app1").
		When().
		CreateApp().
		Then().
		And(func(app *Application) {
			assert.Equal(t, app.Name, "exported-app1")
			assert.Equal(t, app.Namespace, fixture.TestNamespace())
		})

	// Create app in other namespace
	appctx.
		Path(guestbookPath).
		Name("exported-app-oter-namespace").
		SetAppNamespace(fixture.AppNamespace()).
		When().
		CreateApp().
		Then().
		And(func(app *Application) {
			assert.Equal(t, app.Name, "exported-app-oter-namespace")
			assert.Equal(t, app.Namespace, fixture.AppNamespace())
		})

	ctx.
		When().
	  	RunExport().
	  	Then().
	  	AndCLIOutput(func(output string, err error) {
		//exportRawOutput = output
		assert.NoError(t, err, "export finished with error")
		exportResources, err := fixtureutils.GetExportedResourcesFromOutput(output)
		assert.NoError(t, err, "export format not valid")
		assert.True(t, exportResources.HasResource(kube.NewResourceKey("", "ConfigMap", "", "argocd-cm")), "argocd-cm not found in export")
		assert.True(t, exportResources.HasResource(kube.NewResourceKey(v1alpha1.ApplicationSchemaGroupVersionKind.Group, v1alpha1.ApplicationSchemaGroupVersionKind.Kind, "", "exported-app1")), "test namespace application not in export")
		assert.True(t, exportResources.HasResource(kube.NewResourceKey(v1alpha1.ApplicationSchemaGroupVersionKind.Group, v1alpha1.ApplicationSchemaGroupVersionKind.Kind, fixture.AppNamespace(), "exported-app-oter-namespace")), "app namespace application not in export")
	})
}
