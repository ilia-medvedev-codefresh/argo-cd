package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	appclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	metricsutil "github.com/argoproj/argo-cd/v2/util/metrics"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	descAppsetDefaultLabels = []string{"namespace", "name"}

	descAppsetLabels *prometheus.Desc

	descAppsetInfo = prometheus.NewDesc(
		"argocd_appset_info",
		"Information about applicationset",
		descAppsetDefaultLabels,
		nil,
	)
)

type ApplicationsetMetricsCollector struct {
	appsClientSet *appclientset.Clientset
	labels  	  []string
}

func NewAppsetMetricsCollector(clientset *appclientset.Clientset, labels []string) (*ApplicationsetMetricsCollector) {

	if len(labels) > 0 {
		descAppsetLabels = prometheus.NewDesc(
			"argocd_appset_labels",
			"Applicationset labels translated to Prometheus labels",
			append(descAppsetDefaultLabels,metricsutil.NormalizeLabels("label",labels)...),
			nil,
		)
	}

	return (&ApplicationsetMetricsCollector{
		appsClientSet: clientset,
		labels:		   labels,
	})
}

// Describe implements the prometheus.Collector interface
func (c *ApplicationsetMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
		ch <- descAppsetInfo
		if len(c.labels) > 0 {
			ch <- descAppsetLabels
		}
}

// Collect implements the prometheus.Collector interface
func (c *ApplicationsetMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	appsets, _ := c.appsClientSet.ArgoprojV1alpha1().ApplicationSets("").List(context.Background(),v1.ListOptions{})

	for _, appset := range appsets.Items {
		var labelValues =  make([]string,0)
		commonLabelValues := []string{appset.Namespace, appset.Name}

		for _,label := range c.labels {
			labelValues = append(labelValues, appset.GetLabels()[label])
		}
		ch <- prometheus.MustNewConstMetric(descAppsetInfo, prometheus.GaugeValue, 1 , appset.Namespace, appset.Name)
		if len(c.labels) > 0 {
			ch <- prometheus.MustNewConstMetric(descAppsetLabels, prometheus.GaugeValue, 1, append(commonLabelValues, labelValues...)...)
		}
	}
}
