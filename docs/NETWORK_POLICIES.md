# Network Policies

The `multiclusterhub-operator` deploys a [`NetworkPolicy`](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
for insights-metrics as part of the Insights Helm chart (`pkg/templates/charts/toggle/insights/`),
following the principle of least privilege: **deny by default, allow only the specific traffic
each component needs.** The policy is created alongside the other Insights resources and is
automatically removed if the Insights component is disabled on the MultiClusterHub CR.

## Design principles

- **Pod-scoped selector.** The `NetworkPolicy` selects only insights-metrics pods
  (`app: policyreport`, `component: insights-metrics`, `release: policyreport`), never the
  whole namespace. This is important because insights-metrics shares a namespace
  (`open-cluster-management` in a typical ACM install) with unrelated ACM components — a
  namespace-wide policy would inadvertently restrict traffic for pods that aren't part of
  Insights.
- **`policyTypes: [Ingress]` only.** OVN-Kubernetes (the default CNI on OpenShift) handles
  `kubernetes.default.svc` ClusterIP traffic through the OVN service load balancer *before*
  NetworkPolicy evaluation, so no egress rule type can match kube-API traffic. Applying an
  Egress policyType would silently block insights-metrics from reaching the Kubernetes API.
- **Well-known namespace labels.** Ingress rules that reference OpenShift system namespaces
  use the `kubernetes.io/metadata.name` label, which the API server automatically stamps on
  every namespace (Kubernetes 1.21+). This avoids relying on custom labels that may not exist
  in every cluster.

## Component network flows

### insights-metrics

| Direction | Peer | Port | Rationale |
|---|---|---|---|
| Ingress | `openshift-monitoring` namespace | 8443/TCP | Prometheus (`prometheus-k8s`) scrapes metrics through the kube-rbac-proxy sidecar, which terminates TLS on port 8443 and proxies to the metrics container on `127.0.0.1:8383`. The `ServiceMonitor` resource configures a 60-second scrape interval. |
| Egress | *(not restricted — Ingress-only policy)* | — | insights-metrics requires egress to the Kubernetes API to watch cluster resources for metrics collection (PolicyReports, ManagedClusters, etc.) and to poll the APIServer TLS security profile. OVN-Kubernetes handles `kubernetes.default.svc` ClusterIP traffic via the OVN service load balancer before NetworkPolicy evaluation, so no egress rule can match kube-API traffic. Applying an Egress policyType would silently block the metrics collector from reaching the Kubernetes API. |

### Port details

| Port | Bind address | Purpose |
|---|---|---|
| 8443 | `0.0.0.0` | kube-rbac-proxy sidecar (TLS-terminated, externally reachable). This is the port exposed by the Kubernetes Service and scraped by Prometheus. |
| 8383 | `127.0.0.1` | Metrics container (HTTP, localhost only). Serves `/metrics` with `policyreport_info` gauge data. Only reachable from within the pod via the kube-rbac-proxy. |
| 8444 | `127.0.0.1` | Telemetry port (HTTP, localhost only). Serves `/healthz` for liveness/readiness probes and self-monitoring metrics (`ksm_scrape_error_total`, `ksm_resources_per_scrape`). Not externally reachable; no NetworkPolicy rule needed. |

## Testing

Because these policies are enforced by the cluster's CNI plugin (not the Kubernetes API server),
functional verification — confirming that legitimate traffic still flows and that traffic
outside these rules is blocked — requires testing against a real cluster with a
NetworkPolicy-enforcing CNI (e.g. OVN-Kubernetes on OpenShift).
