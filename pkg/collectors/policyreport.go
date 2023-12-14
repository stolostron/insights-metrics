package collectors

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/kube-state-metrics/pkg/metric"
	"sigs.k8s.io/wg-policy-prototypes/policy-report/pkg/api/wgpolicyk8s.io/v1alpha2"
)

var (
	descPolicyReportLabelsName    = "policyreport_info"
	descPolicyReportLabelsHelp    = "Open Cluster Management PolicyReport Info."
	descPolicyReportDefaultLabels = []string{"managed_cluster_id", "category", "policy", "result", "severity"}

	policyReportGvr = schema.GroupVersionResource{
		Group:    "wgpolicyk8s.io",
		Version:  "v1alpha2",
		Resource: "policyreports",
	}
)

func getPolicyReportMetricFamilies(client dynamic.Interface) []metric.FamilyGenerator {
	return []metric.FamilyGenerator{
		{
			Name: descPolicyReportLabelsName,
			Type: metric.Gauge,
			Help: descPolicyReportLabelsHelp,
			GenerateFunc: wrapPolicyReportFunc(func(prObj *unstructured.Unstructured) metric.Family {
				klog.Infof("Getting PolicyReport Info for Cluster Name %s with name %s", prObj.GetNamespace(), prObj.GetName())
				pr := &v1alpha2.PolicyReport{}
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(prObj.UnstructuredContent(), &pr)
				if err != nil {
					klog.Infof("Error unstructuring PolicyReport ")
					return metric.Family{Metrics: []*metric.Metric{}}
				}
				_, errPR := client.Resource(policyReportGvr).Namespace(pr.GetNamespace()).Get(context.TODO(), pr.GetName(), metav1.GetOptions{})
				if errPR != nil {
					klog.Infof("PolicyReport %s not found, err: %s", pr.GetName(), errPR)
					return metric.Family{Metrics: []*metric.Metric{}}
				}
				clusterName := pr.GetNamespace()
				clusterId := getClusterID(client, clusterName)
				metrics := getReports(clusterId, pr)

				f := metric.Family{}

				for _, result := range metrics {
					f.Metrics = append(f.Metrics, &metric.Metric{
						LabelKeys:   descPolicyReportDefaultLabels,
						LabelValues: result,
						Value:       1,
					})
				}

				klog.Infof("Returning %v", string(f.ByteSlice()))
				return f
			}),
		},
	}
}

func wrapPolicyReportFunc(f func(*unstructured.Unstructured) metric.Family) func(interface{}) *metric.Family {
	return func(obj interface{}) *metric.Family {
		PolicyReport := obj.(*unstructured.Unstructured)

		metricFamily := f(PolicyReport)

		for _, m := range metricFamily.Metrics {
			m.LabelKeys = append([]string{}, m.LabelKeys...)
			m.LabelValues = append([]string{}, m.LabelValues...)
		}

		return &metricFamily
	}
}

func createPolicyReportListWatchWithClient(client dynamic.Interface, ns string) cache.ListWatch {
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return client.Resource(policyReportGvr).Namespace(ns).List(context.TODO(), opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return client.Resource(policyReportGvr).Namespace(ns).Watch(context.TODO(), opts)
		},
	}
}

func getReports(clusterID string, pr *v1alpha2.PolicyReport) [][]string {
	var metrics [][]string

	if clusterID == "" {
		return metrics
	}

	for _, reportResult := range pr.Results {
		var severity string

		result := "fail"

		if reportResult.Result != "" {
			result = string(reportResult.Result)
		}

		switch risk := reportResult.Properties["total_risk"]; risk {
		case "4":
			severity = "critical"
		case "3":
			severity = "important"
		case "2":
			severity = "moderate"
		case "1":
			severity = "low"
		default:
			severity = "unknown"
		}

		if reportResult.Policy != "" {
			metrics = append(metrics, []string{
				clusterID,
				reportResult.Category,
				reportResult.Policy,
				result,
				severity,
			})
		}
	}

	return metrics
}
