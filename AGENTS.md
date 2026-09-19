# insights-metrics

## Overview

`insights-metrics` is a Go service that watches Kubernetes PolicyReport custom
resources and exposes their metrics for Prometheus. It is based on
`kube-state-metrics` and uses client-go dynamic clients and namespace-scoped
list/watch reflectors.

## Repository Layout

- `main.go`: command-line setup, Kubernetes configuration, TLS profile polling,
  and the metrics and health HTTP servers.
- `pkg/collectors/`: PolicyReport collection, filtering, and metric generation.
- `pkg/options/`: command-line options and default collector configuration.
- `pkg/tlsprofile/`: API-server TLS security profile lookup and polling.
- `build/`: CI build, dependency, deployment, and unit-test scripts.
- `docs/NETWORK_POLICIES.md`: network policy documentation.
- `.tekton/`: Tekton pipeline definitions for ACM releases.

## Development Commands

Run commands from the repository root:

```bash
go test ./...                         # Fast local test run
go build ./...                        # Compile all packages
make deps                             # Tidy Go modules
make test                             # Verbose tests and cover.out
make lint                             # Install/run golangci-lint v2.9.0
make build                            # Build output/insights-metrics
make run                              # Run the service with go run
build/run-unit-tests.sh               # CI sequence: deps, lint, test, coverage
```

The Makefile bootstraps build-harness extensions unless
`USE_VENDORIZED_BUILD_HARNESS` is set. Targets that depend on the harness
require a valid `.build-harness-bootstrap`; use `USE_VENDORIZED_BUILD_HARNESS`
when a vendored harness is available. `make lint` downloads golangci-lint if it
is not already installed.

## Architecture

`main.go` builds Kubernetes REST configuration from either in-cluster settings
or explicit kubeconfig/API-server flags. It creates a dynamic client, polls the
API-server TLS profile, and starts two HTTP servers: one for service telemetry
and one for collected PolicyReport metrics. The collector builder creates a
dynamic list/watch reflector per configured namespace and writes metric
families through kube-state-metrics stores. `/metrics` serves metrics and
`/healthz` reports service health.

## Integrations

- GitHub CLI: `gh` is not installed in this environment. Use the configured
  GitHub MCP tools for GitHub operations when available.
- Jira CLI: no Jira CLI is assumed; use the configured Jira MCP tools for Jira
  operations.
- GitHub tokens, when needed by repository tooling, follow `GH_TOKEN` and
  per-organization `GH_TOKEN_<ORG>` conventions.

## Personal configuration

Read personal config at the start of any task that needs an assignee, email, or project key.
Canonical path: ~/.config/user.local.md (tool-agnostic, global).
If the file does not exist, fall back to agent memory (`user-config`), then placeholders.
Run `make personalize` to generate or update the file (if this repo uses Fleet Engineering tooling).

## Fleet Engineering Skills

Fetch and apply the relevant skill when the task matches its domain.

| Skill | When to use |
|---|---|
| [bug-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/bug-specialist/SKILL.md) | Bug triage, reproduction steps, fix planning |
| [epic-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/epic-specialist/SKILL.md) | Multi-sprint epics with outcomes |
| [feature-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/feature-specialist/SKILL.md) | Large customer-facing capabilities |
| [initiative-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/initiative-specialist/SKILL.md) | Multi-team strategic programs |
| [jira-create](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/jira-create/SKILL.md) | Interactive issue creation with specialist delegation |
| [jira-qe-readiness](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/jira-qe-readiness/SKILL.md) | Check whether a Jira ticket has enough information for QE |
| [jira-report](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/jira-report/SKILL.md) | Jira portfolio reports and issue quality reviews |
| [jira-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/jira-specialist/SKILL.md) | General Jira triage, linking, and transitions |
| [jira-type-audit](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/jira-type-audit/SKILL.md) | Audit and correct Jira issue types across the hierarchy |
| [outcome-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/outcome-specialist/SKILL.md) | Strategic outcomes tied to OKRs |
| [release-dod](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/release-dod/SKILL.md) | Release Definition of Done checklists |
| [risk-report](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/risk-report/SKILL.md) | Detect Jira risk signals and draft status-report risks |
| [risk-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/risk-specialist/SKILL.md) | Risk register and mitigation planning |
| [spike-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/spike-specialist/SKILL.md) | Time-boxed research and proof-of-concept work |
| [story-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/story-specialist/SKILL.md) | User stories and acceptance criteria |
| [supportex-review](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/supportex-review/SKILL.md) | Review SUPPORTEX support exception requests |
| [task-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/task-specialist/SKILL.md) | Internal technical task planning |
| [ticket-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/ticket-specialist/SKILL.md) | Stakeholder request intake and triage |
| [backlog-grooming](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/backlog-grooming/SKILL.md) | Jira backlog grooming and sprint-readiness analysis |
| [breaking-changes](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/breaking-changes/SKILL.md) | Detect breaking API, database, config, or behavior changes |
| [ci-triage](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/ci-triage/SKILL.md) | Diagnose failing PR checks |
| [coderabbit-sync](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/coderabbit-sync/SKILL.md) | Maintain the Fleet reference CodeRabbit configuration |
| [cve-sustaining-handoff](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/cve-sustaining-handoff/SKILL.md) | Resolve CVE trackers and hand off sustaining work |
| [vulnerability-slack-report](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/vulnerability-slack-report/SKILL.md) | Report overdue vulnerability issues |
| [diagnosing-bugs](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/diagnosing-bugs/SKILL.md) | Reproduce, minimize, and fix unclear failures |
| [finish-work](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/finish-work/SKILL.md) | Commit, push, open PR, and update Jira |
| [github-org-access](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/github-org-access/SKILL.md) | Modify GitHub organization access configuration |
| [init-context-docs](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/init-context-docs/SKILL.md) | Assess and bootstrap repository AI-readiness documentation |
| [opencode-setup](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/opencode-setup/SKILL.md) | Install and configure OpenCode |
| [org-repo-audit](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/org-repo-audit/SKILL.md) | Audit organization repositories for SDLC readiness |
| [pr-fix](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/pr-fix/SKILL.md) | Fix merge conflicts, CI failures, and review comments |
| [pr-hygiene](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/pr-hygiene/SKILL.md) | Manage stale and re-review-needed pull requests |
| [pr-review](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/pr-review/SKILL.md) | Review GitHub pull requests with inline comments |
| [pr-review-detailed](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/pr-review-detailed/SKILL.md) | Perform layered checklist-based code analysis |
| [pr-review-fix](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/pr-review-fix/SKILL.md) | Iteratively review and fix local changes before commit |
| [release-notes](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/release-notes/SKILL.md) | Generate categorized release notes |
| [renovate-prs](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/renovate-prs/SKILL.md) | Manage dependency update pull requests |
| [repo-content-audit](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/repo-content-audit/SKILL.md) | Find unlinked or orphaned repository content |
| [repo-setup](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/repo-setup/SKILL.md) | Onboard repositories to the Fleet Engineering SDLC |
| [rhacm-addon-wizard](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/rhacm-addon-wizard/SKILL.md) | Guide RHACM add-on development |
| [scored-code-review](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/scored-code-review/SKILL.md) | Deprecated scored code review workflow |
| [session-summary](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/session-summary/SKILL.md) | Summarize session work with Jira and GitHub context |
| [start-work](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/start-work/SKILL.md) | Create a Jira sub-task |
| [test-coverage-gap](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/test-coverage-gap/SKILL.md) | Analyze risk-prioritized test coverage gaps |
| [f2f-daily-summary](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/f2f-daily-summary/SKILL.md) | Capture daily F2F notes as Jira sub-tasks |
| [f2f-epic-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/f2f-epic-specialist/SKILL.md) | Create and manage F2F meeting epics |
| [presentation-task](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/presentation-task/SKILL.md) | Log delivered presentations as Jira work |
| [scrum-status](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/scrum-status/SKILL.md) | Capture scrum bullets and generate team reports |

The authoritative catalog is [Fleet Engineering skills README](https://github.com/OpenShift-Fleet/agentic-sdlc/blob/main/skills/README.md).
