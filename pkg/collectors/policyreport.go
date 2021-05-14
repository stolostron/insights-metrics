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
	"sigs.k8s.io/wg-policy-prototypes/policy-report/api/v1alpha2"
)

var (
	descPolicyReportLabelsName    = "acm_policyreport_info"
	descPolicyReportLabelsHelp    = "ACM PolicyReport Info."
	descPolicyReportDefaultLabels = []string{"cluster_id", "category", "policy", "result"}

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
				klog.Infof("Cluster Name %s", prObj.GetName())
				pr := &v1alpha2.PolicyReport{}
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(prObj.UnstructuredContent(), &pr)
				if err != nil {
					klog.Infof("Error unstructuring PolicyReport ")
					return metric.Family{Metrics: []*metric.Metric{}}
				}
				_, errPR := client.Resource(policyReportGvr).Namespace(pr.GetName()).Get(context.TODO(), pr.GetName(), metav1.GetOptions{})
				if errPR != nil {
					klog.Infof("PolicyReport %s not found, err: %s", pr.GetName(), errPR)
				}
				clusterName := pr.GetName()
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
	category, policy, result := "", "", ""
	for _, reportResult := range pr.Results {
		var metric []string
		if reportResult.Category != "" {
			category = reportResult.Category
		}
		if reportResult.Result != "" {
			result = string(reportResult.Result)
		}
		if reportResult.Policy != "" && category != "" && result != "" {
			policy = reportResult.Policy
			metric = append(metric, clusterID, category, policy, result)
			metrics = append(metrics, metric)
		}

	}
	return metrics
}
