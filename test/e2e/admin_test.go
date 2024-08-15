package e2e

import (
	"testing"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"

	//appsv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	. "github.com/argoproj/argo-cd/v2/test/e2e/fixture/admin"
	fixtureutils "github.com/argoproj/argo-cd/v2/test/e2e/fixture/admin/utils"

	//appfixture "github.com/argoproj/argo-cd/v2/test/e2e/fixture/app"
	"github.com/stretchr/testify/assert"
)

func TestBackupExportImport(t *testing.T) {
	ctx := Given(t)

	//exportRawOutput := ""
	ctx.
	  When().
	  RunExport().
	  Then().
	  AndCLIOutput(func(output string, err error) {
		assert.NoError(t, err, "export finished with error")
		exportResources, err := fixtureutils.GetExportedResourcesFromOutput(output)
		assert.NoError(t, err, "export format not valid")
		assert.True(t, exportResources.HasResource(kube.NewResourceKey("v1", "ConfigMap", "", "argocd-cm")), "argocd-cm not found in export")
	})
}
