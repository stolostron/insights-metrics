module github.com/open-cluster-management/insights-metrics

go 1.16

require (
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20210224082022-3d97a244fca7
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
	k8s.io/klog/v2 v2.8.0
	k8s.io/kube-state-metrics v0.0.0-20190129120824-7bfed92869b6
	sigs.k8s.io/wg-policy-prototypes v0.0.0-20210430040757-037274b4aefe
)
