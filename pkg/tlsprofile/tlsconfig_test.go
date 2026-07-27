// Copyright Contributors to the Open Cluster Management project

package tlsprofile

import (
	"context"
	"crypto/tls"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

func newFakeAPIServer(tlsProfile map[string]interface{}) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "config.openshift.io/v1",
			"kind":       "APIServer",
			"metadata": map[string]interface{}{
				"name": "cluster",
			},
			"spec": map[string]interface{}{},
		},
	}
	if tlsProfile != nil {
		obj.Object["spec"].(map[string]interface{})["tlsSecurityProfile"] = tlsProfile
	}
	return obj
}

func newFakeDynamicClient(objects ...runtime.Object) *dynamicfake.FakeDynamicClient {
	s := runtime.NewScheme()
	s.AddKnownTypeWithName(
		schema.GroupVersionKind{Group: "config.openshift.io", Version: "v1", Kind: "APIServer"},
		&unstructured.Unstructured{},
	)
	s.AddKnownTypeWithName(
		schema.GroupVersionKind{Group: "config.openshift.io", Version: "v1", Kind: "APIServerList"},
		&unstructured.UnstructuredList{},
	)
	return dynamicfake.NewSimpleDynamicClient(s, objects...)
}

func TestGetTLSConfig_NoAPIServer(t *testing.T) {
	client := newFakeDynamicClient()

	cfg, _, ok := GetTLSConfig(context.TODO(), client)

	if cfg == nil {
		t.Fatal("expected non-nil config on Intermediate fallback")
	}
	if ok {
		t.Error("snapshot should be invalid on fallback")
	}
	if cfg.MinVersion != tls.VersionTLS12 {
		t.Errorf("expected TLS 1.2 fallback, got %d", cfg.MinVersion)
	}
}

func TestGetTLSConfig_NoProfile(t *testing.T) {
	apiServer := newFakeAPIServer(nil)
	client := newFakeDynamicClient(apiServer)

	cfg, _, ok := GetTLSConfig(context.TODO(), client)

	if cfg == nil {
		t.Fatal("expected non-nil config for default profile")
	}
	if !ok {
		t.Error("snapshot should be valid")
	}
	if cfg.MinVersion != tls.VersionTLS12 {
		t.Errorf("expected TLS 1.2, got %d", cfg.MinVersion)
	}
	if len(cfg.CipherSuites) == 0 {
		t.Error("expected cipher suites to be set")
	}
}

func TestGetTLSConfig_IntermediateProfile(t *testing.T) {
	apiServer := newFakeAPIServer(map[string]interface{}{
		"type": "Intermediate",
	})
	client := newFakeDynamicClient(apiServer)

	cfg, snapshot, ok := GetTLSConfig(context.TODO(), client)

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if !ok {
		t.Error("snapshot should be valid")
	}
	if snapshot == nil {
		t.Error("expected non-nil snapshot for explicit profile")
	}
	if cfg.MinVersion != tls.VersionTLS12 {
		t.Errorf("expected TLS 1.2, got %d", cfg.MinVersion)
	}
}

func TestGetTLSConfig_OldProfile(t *testing.T) {
	apiServer := newFakeAPIServer(map[string]interface{}{
		"type": "Old",
	})
	client := newFakeDynamicClient(apiServer)

	cfg, _, ok := GetTLSConfig(context.TODO(), client)

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if !ok {
		t.Error("snapshot should be valid")
	}
	if cfg.MinVersion != tls.VersionTLS10 {
		t.Errorf("expected TLS 1.0, got %d", cfg.MinVersion)
	}
}

func TestGetTLSConfig_ModernProfile(t *testing.T) {
	apiServer := newFakeAPIServer(map[string]interface{}{
		"type": "Modern",
	})
	client := newFakeDynamicClient(apiServer)

	cfg, _, ok := GetTLSConfig(context.TODO(), client)

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if !ok {
		t.Error("snapshot should be valid")
	}
	if cfg.MinVersion != tls.VersionTLS13 {
		t.Errorf("expected TLS 1.3, got %d", cfg.MinVersion)
	}
}

func TestGetTLSConfig_CustomProfile(t *testing.T) {
	apiServer := newFakeAPIServer(map[string]interface{}{
		"type": "Custom",
		"custom": map[string]interface{}{
			"ciphers":       []interface{}{"ECDHE-RSA-AES256-GCM-SHA384"},
			"minTLSVersion": "VersionTLS13",
		},
	})
	client := newFakeDynamicClient(apiServer)

	cfg, _, ok := GetTLSConfig(context.TODO(), client)

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if !ok {
		t.Error("snapshot should be valid")
	}
	if cfg.MinVersion != tls.VersionTLS13 {
		t.Errorf("expected TLS 1.3, got %d", cfg.MinVersion)
	}
}

func TestIntermediateProfileTLSConfig(t *testing.T) {
	cfg := IntermediateProfileTLSConfig()

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.MinVersion != tls.VersionTLS12 {
		t.Errorf("expected TLS 1.2, got %d", cfg.MinVersion)
	}
	if len(cfg.CipherSuites) == 0 {
		t.Error("expected cipher suites to be set")
	}
}
