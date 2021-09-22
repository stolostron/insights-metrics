// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package collectors

import (
	"context"

	clusterv1 "open-cluster-management.io/api/cluster/v1"
	ocinfrav1 "github.com/openshift/api/config/v1"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
)

var (
	ScrapeErrorTotalMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ksm_scrape_error_total",
			Help: "Total scrape errors encountered when scraping a resource",
		},
		[]string{"resource"},
	)

	ResourcesPerScrapeMetric = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "ksm_resources_per_scrape",
			Help: "Number of resources returned per scrape",
		},
		[]string{"resource"},
	)

	cvGVR = schema.GroupVersionResource{
		Group:    "config.openshift.io",
		Version:  "v1",
		Resource: "clusterversions",
	}

	mcGVR = schema.GroupVersionResource{
		Group:    "cluster.open-cluster-management.io",
		Version:  "v1",
		Resource: "managedclusters",
	}
)

func getClusterID(c dynamic.Interface, clusterName string) string {
	clusterId := ""
	if clusterName == "local-cluster" {
		cvObj, errCv := c.Resource(cvGVR).Get(context.TODO(), "version", metav1.GetOptions{})
		if errCv != nil {
			klog.Warningf("Error getting cluster version %v \n", errCv)
			return clusterId
		}
		cv := &ocinfrav1.ClusterVersion{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(cvObj.UnstructuredContent(), &cv)
		if err != nil {
			klog.Warningf("Error unmarshal cluster version object%v \n", err)
			return clusterId
		}
		return string(cv.Spec.ClusterID)
	} else {
		mcObj, errMc := c.Resource(mcGVR).Get(context.TODO(), clusterName, metav1.GetOptions{})
		if errMc != nil {
			klog.Warningf("Error getting ManagedCluster %v \n", errMc)
			return clusterId
		}
		mc := &clusterv1.ManagedCluster{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(mcObj.UnstructuredContent(), &mc)
		if err != nil {
			klog.Warningf("Error unmarshal ManagedCluster object%v \n", err)
			return clusterId
		}
		for _, claimInfo := range mc.Status.ClusterClaims {
			if claimInfo.Name == "id.openshift.io" {
				return string(claimInfo.Value)
			}
		}
	}
	return clusterId

}
