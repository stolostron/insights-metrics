// Copyright Contributors to the Open Cluster Management project

package tlsprofile

import (
	"context"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCurrentTLSProfileData_NoAPIServer(t *testing.T) {
	client := newFakeDynamicClient()

	data, err := currentTLSProfileData(client)

	if err == nil {
		t.Error("expected error when APIServer doesn't exist")
	}
	if data != nil {
		t.Error("expected nil data")
	}
}

func TestCurrentTLSProfileData_NoProfile(t *testing.T) {
	apiServer := newFakeAPIServer(nil)
	client := newFakeDynamicClient(apiServer)

	data, err := currentTLSProfileData(client)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if data != nil {
		t.Error("expected nil data when no profile is set")
	}
}

func TestCurrentTLSProfileData_WithProfile(t *testing.T) {
	apiServer := newFakeAPIServer(map[string]interface{}{
		"type": "Intermediate",
	})
	client := newFakeDynamicClient(apiServer)

	data, err := currentTLSProfileData(client)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil data")
	}
	if data["type"] != "Intermediate" {
		t.Errorf("expected type Intermediate, got %v", data["type"])
	}
}

func TestPollTLSProfile_NoChangeDoesNotExit(t *testing.T) {
	apiServer := newFakeAPIServer(map[string]interface{}{
		"type": "Intermediate",
	})
	client := newFakeDynamicClient(apiServer)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Should complete without exiting when the profile doesn't change.
	pollTLSProfile(ctx, client, 50*time.Millisecond)
}

func TestPollTLSProfile_DetectsChange(t *testing.T) {
	apiServer := newFakeAPIServer(map[string]interface{}{
		"type": "Intermediate",
	})
	client := newFakeDynamicClient(apiServer)

	// Read initial state
	initial, err := currentTLSProfileData(client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if initial["type"] != "Intermediate" {
		t.Fatalf("expected Intermediate, got %v", initial["type"])
	}

	// Simulate a profile change
	updated := newFakeAPIServer(map[string]interface{}{
		"type": "Old",
	})
	_, err = client.Resource(apiServerGVR).Update(
		context.TODO(), updated, metav1.UpdateOptions{})
	if err != nil {
		t.Fatalf("unexpected error updating APIServer: %v", err)
	}

	current, err := currentTLSProfileData(client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if current["type"] != "Old" {
		t.Errorf("expected Old, got %v", current["type"])
	}
}

func TestCurrentTLSProfileData_StableNormalization(t *testing.T) {
	apiServer := newFakeAPIServer(map[string]interface{}{
		"type": "Custom",
		"custom": map[string]interface{}{
			"ciphers":       []interface{}{"ECDHE-RSA-AES256-GCM-SHA384", "ECDHE-RSA-AES128-GCM-SHA256"},
			"minTLSVersion": "VersionTLS12",
		},
	})
	client := newFakeDynamicClient(apiServer)

	data1, _ := currentTLSProfileData(client)
	data2, _ := currentTLSProfileData(client)

	if len(data1) != len(data2) {
		t.Error("expected stable normalization between reads")
	}
}

func TestCurrentTLSProfileData_NoSpec(t *testing.T) {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "config.openshift.io/v1",
			"kind":       "APIServer",
			"metadata": map[string]interface{}{
				"name": "cluster",
			},
		},
	}
	client := newFakeDynamicClient(obj)

	data, err := currentTLSProfileData(client)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if data != nil {
		t.Error("expected nil data when no spec")
	}
}
