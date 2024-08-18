package metrics

import (
	// "context"
	// "log"
	"net/http"
	"net/http/httptest"
	// "strings"
	"testing"

	argoappv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	appclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"sigs.k8s.io/yaml"
)

const fakeAppset1 = `
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: test
  namespace: argocd
spec:
  generators:
  - git:
      directories:
      - path: test/*
      repoURL: https://github.com/test/test.git
      revision: HEAD
  template:
    metadata:
      name: '{{.path.basename}}'
    spec:
      destination:
        namespace: '{{.path.basename}}'
        server: https://kubernetes.default.svc
      project: default
      source:
        path: '{{.path.path}}'
        repoURL: https://github.com/test/test.git
        targetRevision: HEAD
status:
  conditions:
  - lastTransitionTime: "2024-01-01T00:00:00Z"
    message: Successfully generated parameters for all Applications
    reason: ApplicationSetUpToDate
    status: "False"
    type: ErrorOccurred
  - lastTransitionTime: "2024-01-01T00:00:00Z"
    message: Successfully generated parameters for all Applications
    reason: ParametersGenerated
    status: "True"
    type: ParametersGenerated
  - lastTransitionTime: "2024-01-01T00:00:00Z"
    message: ApplicationSet up to date
    reason: ApplicationSetUpToDate
    status: "True"
    type: ResourcesUpToDate
`

func newFakeAppset(fakeAppYAML string) *argoappv1.ApplicationSet {
	var appset argoappv1.ApplicationSet
	err := yaml.Unmarshal([]byte(fakeAppYAML), &appset)
	if err != nil {
		panic(err)
	}
	return &appset
}

func TestApplicationsetCollector(t *testing.T) {
	appset := newFakeAppset(fakeAppset1)
	appsets := []runtime.Object{appset}

	fakeClientSet := appclientset.NewSimpleClientset(appsets...)

	appsetCollector := newAppsetCollector(fakeClientSet, []string{})

	metrics.Registry.MustRegister(appsetCollector)

	req, err := http.NewRequest("GET", "/metrics", nil)
    assert.NoError(t, err)

	rr := httptest.NewRecorder()
    handler := promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{})
    handler.ServeHTTP(rr, req)

    // Check the response
    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Contains(t, rr.Body.String(), `
# TYPE argocd_appset_info gauge
argocd_appset_info{name="test",namespace="argocd"} 1
`,
)
}
