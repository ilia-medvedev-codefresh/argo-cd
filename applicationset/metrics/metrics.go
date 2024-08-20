package metrics

import (
	"context"
	"time"

	argoappv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	appclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	metricsutil "github.com/argoproj/argo-cd/v2/util/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	descAppsetLabels *prometheus.Desc
	descAppsetDefaultLabels = []string{"namespace", "name"}
	descAppsetInfo = prometheus.NewDesc(
		"argocd_appset_info",
		"Information about applicationset",
		append(descAppsetDefaultLabels, "resource_update_status"),
		nil,
	)
)

type ApplicationsetMetrics struct {
	reconcileHistogram *prometheus.HistogramVec
}

type appsetCollector struct {
	appsClientSet appclientset.Interface
	labels  	  []string
}

func NewApplicationsetMetrics(clientset *appclientset.Clientset, appsetLabels []string) (*ApplicationsetMetrics) {

	reconcileHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "argocd_appset_reconcile",
			Help: "Application reconciliation performance in seconds.",
			// Buckets can be set later on after observing median time
		},
		descAppsetDefaultLabels,
	)

	appsetCollector := newAppsetCollector(clientset, appsetLabels)

	// Rgister collectors and metrics
	metrics.Registry.MustRegister(reconcileHistogram)
	metrics.Registry.MustRegister(appsetCollector)

	return &ApplicationsetMetrics{
		reconcileHistogram: reconcileHistogram,
	}
}

func (m *ApplicationsetMetrics) ObserveRconcile(appset *argoappv1.ApplicationSet, duration time.Duration) {
	m.reconcileHistogram.WithLabelValues(appset.Namespace,appset.Name).Observe(duration.Seconds())
}

func newAppsetCollector(clientset appclientset.Interface, labels []string) (*appsetCollector) {

	descAppsetDefaultLabels = []string{"namespace", "name"}

	if len(labels) > 0 {
		descAppsetLabels = prometheus.NewDesc(
			"argocd_appset_labels",
			"Applicationset labels translated to Prometheus labels",
			append(descAppsetDefaultLabels,metricsutil.NormalizeLabels("label",labels)...),
			nil,
		)
	}

	return (&appsetCollector{
		appsClientSet: clientset,
		labels:		   labels,
	})
}

// Describe implements the prometheus.Collector interface
func (c *appsetCollector) Describe(ch chan<- *prometheus.Desc) {
		ch <- descAppsetInfo
		if len(c.labels) > 0 {
			ch <- descAppsetLabels
		}
}

// Collect implements the prometheus.Collector interface
func (c *appsetCollector) Collect(ch chan<- prometheus.Metric) {
	appsets, _ := c.appsClientSet.ArgoprojV1alpha1().ApplicationSets("").List(context.Background(),v1.ListOptions{})

	for _, appset := range appsets.Items {
		var labelValues =  make([]string,0)
		commonLabelValues := []string{appset.Namespace, appset.Name}

		for _,label := range c.labels {
			labelValues = append(labelValues, appset.GetLabels()[label])
		}

		if len(c.labels) > 0 {
			ch <- prometheus.MustNewConstMetric(descAppsetLabels, prometheus.GaugeValue, 1, append(commonLabelValues, labelValues...)...)
		}

		resourceUpdateStatus := "Unknown"

		for _,condition := range appset.Status.Conditions {
			if condition.Type == argoappv1.ApplicationSetConditionResourcesUpToDate {
				resourceUpdateStatus = condition.Reason
			}
		}

		ch <- prometheus.MustNewConstMetric(descAppsetInfo, prometheus.GaugeValue, 1 , appset.Namespace, appset.Name, resourceUpdateStatus)
	}
}
