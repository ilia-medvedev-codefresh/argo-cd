package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	appclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	descAppsetDefaultLabels = []string{"namespace", "name"}

	//descAppLabels *prometheus.Desc

	descAppsetInfo = prometheus.NewDesc(
		"argocd_appset_info",
		"Information about applicationset",
		descAppsetDefaultLabels,
		nil,
	)
)

type ApplicationsetMetricsCollector struct {
	appsClientSet *appclientset.Clientset
}

func NewAppsetMetricsCollector(clientset *appclientset.Clientset) (*ApplicationsetMetricsCollector) {
	return (&ApplicationsetMetricsCollector{
		appsClientSet: clientset,
	})
}

// Describe implements the prometheus.Collector interface
func (c *ApplicationsetMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
		ch <- descAppsetInfo
}

// Collect implements the prometheus.Collector interface
func (c *ApplicationsetMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	appsets, _ := c.appsClientSet.ArgoprojV1alpha1().ApplicationSets("").List(context.Background(),v1.ListOptions{})

	for _, appset := range appsets.Items {
		ch <- prometheus.MustNewConstMetric(descAppsetInfo, prometheus.GaugeValue, 1 , appset.Namespace, appset.Name)
	}
}
