// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package collectors

import (
	"reflect"
	"testing"

	ocinfrav1 "github.com/openshift/api/config/v1"
	mcv1 "github.com/stolostron/api/cluster/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kube-state-metrics/pkg/metric"
	pr "sigs.k8s.io/wg-policy-prototypes/policy-report/api/v1alpha2"
)

func Test_getPolicyReportMetricFamilies(t *testing.T) {
	s := scheme.Scheme

	s.AddKnownTypes(pr.SchemeGroupVersion, &pr.PolicyReport{})
	s.AddKnownTypes(ocinfrav1.SchemeGroupVersion, &ocinfrav1.ClusterVersion{})
	s.AddKnownTypes(mcv1.SchemeGroupVersion, &mcv1.ManagedCluster{})
	version := &ocinfrav1.ClusterVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "version",
		},
		Spec: ocinfrav1.ClusterVersionSpec{
			ClusterID: "mycluster_id",
		},
	}

	mc := &mcv1.ManagedCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "managed-cluster",
		},
		Status: mcv1.ManagedClusterStatus{
			ClusterClaims: []mcv1.ManagedClusterClaim{
				{
					Name:  "id.openshift.io",
					Value: "managed-cluster",
				},
			},
		},
	}
	pri := &pr.PolicyReport{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "local-cluster",
			Namespace: "local-cluster",
		},
		Results: []*pr.PolicyReportResult{
			{
				Category: "openshift,configuration,service_availability",
				Policy:   "MASTER_DEFINED_AS_MACHINESET",
				Result:   "fail",
				Properties: map[string]string{
					"total_risk": "4",
				},
			},
		},
	}
	prU := &unstructured.Unstructured{}
	err := scheme.Scheme.Convert(pri, prU, nil)
	if err != nil {
		t.Error(err)
	}

	prm := &pr.PolicyReport{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "managed-cluster",
			Namespace: "managed-cluster",
		},
		Results: []*pr.PolicyReportResult{
			{
				Category: "service_availability",
				Policy:   "MASTER_DEFINED_AS_MACHINESET",
				Result:   "skip",
				Properties: map[string]string{
					"total_risk": "3",
				},
			},
		},
	}
	prUM := &unstructured.Unstructured{}
	err = scheme.Scheme.Convert(prm, prUM, nil)
	if err != nil {
		t.Error(err)
	}

	client := fake.NewSimpleDynamicClient(s, prU, prUM, version, mc)
	tests := []generateMetricsTestCase{
		{
			Obj:  prU,
			Want: `policyreport_info{managed_cluster_id="mycluster_id",category="openshift,configuration,service_availability",policy="MASTER_DEFINED_AS_MACHINESET",result="fail",severity="critical"} 1`,
		},
		{
			Obj:  prUM,
			Want: `policyreport_info{managed_cluster_id="managed-cluster",category="service_availability",policy="MASTER_DEFINED_AS_MACHINESET",result="skip",severity="important"} 1`,
		},
	}
	for i, c := range tests {
		c.Func = metric.ComposeMetricGenFuncs(getPolicyReportMetricFamilies(client))
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %v run:\n%s", i, err)
		}
	}
}

func Test_createPolicyReportListWatchWithClient(t *testing.T) {

	s := runtime.NewScheme()
	s.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.Namespace{})
	s.AddKnownTypes(pr.SchemeGroupVersion, &pr.PolicyReport{})
	s.AddKnownTypes(pr.SchemeGroupVersion, &pr.PolicyReportList{})

	pri := &pr.PolicyReport{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PolicyReport",
			APIVersion: "wgpolicyk8s.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "local-cluster",
			Namespace: "local-cluster",
		},
		Results: []*pr.PolicyReportResult{
			{
				Category:    "openshift,configuration,service_availability",
				Policy:      "MASTER_DEFINED_AS_MACHINESET",
				Result:      "fail",
				Timestamp:   metav1.Timestamp{},
				Scored:      false,
				Description: "test",
				Properties: map[string]string{
					"total_risk": "3",
				},
			},
		},
	}
	prU := &unstructured.Unstructured{}
	err := scheme.Scheme.Convert(pri, prU, nil)
	if err != nil {
		t.Error(err)
	}

	client := fake.NewSimpleDynamicClient(s, pri)
	type args struct {
		client dynamic.Interface
		ns     string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "succeed",
			args: args{
				client: client,
				ns:     "local-cluster",
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createPolicyReportListWatchWithClient(tt.args.client, tt.args.ns)
			l, err := got.ListFunc(metav1.ListOptions{})
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
			lU := l.(*unstructured.UnstructuredList)

			if len(lU.Items) != tt.want {
				t.Errorf("expected a list of %d elements got %d", tt.want, len(lU.Items))
			}
			if !reflect.DeepEqual(lU.Items[0], *prU) {
				t.Errorf("expected of %v got %v", *prU, lU.Items[0])
			}
			w, err := got.WatchFunc(metav1.ListOptions{})
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
			if w == nil {
				t.Errorf("expected the watch to be not nil")
			}
		})
	}
}
