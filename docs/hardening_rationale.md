# Container & Runtime Hardening Rationale

As-of: `2026-06-01`. Scope: all 23 Go microservices (`services.yaml`).

This document records the least-privilege hardening decisions applied to the
Dockerfiles and Kubernetes base manifests, the threats they address, the residual
risks that remain, and the explicit non-goals. It closes the documentation half of
issue #1 ("Harden Dockerfile and runtime defaults for least privilege").

## 1. Image hardening (Dockerfiles)

Every service uses a two-stage build that produces a static binary
(`CGO_ENABLED=0`, `-trimpath -ldflags="-s -w"`) and ships it on
`gcr.io/distroless/static-debian12:nonroot`.

| Decision | Rationale |
|---|---|
| Distroless `static-debian12` runtime | No shell, package manager, or libc userland → drastically smaller attack/CVE surface than a general-purpose base. Previously 14/23 services ran on `ubuntu:18.04` (EOL since 2023; a recurring source of fixable High/Critical CVEs at the Grype gate). |
| `:nonroot` tag + `USER 65532:65532` | The image defaults to an unprivileged uid even before Kubernetes applies its own `securityContext`. Previously 9/23 "distroless" images had no `USER` and ran as root. |
| `CGO_ENABLED=0` static binary | Self-contained; no dynamic loader or system libraries needed at runtime, which is what makes a distroless `static` base viable. |
| CA certificates | Provided by the distroless base, so outbound TLS works without an `apt-get install ca-certificates` layer. |

`user-service` keeps its baked `config.yaml` and dual ports (18080/19090) but now
also runs on the same distroless `:nonroot` base.

## 2. Runtime hardening (Kubernetes base manifests)

Each `services/<svc>/deploy/kubernetes/base/deployment.yaml` sets:

- **Pod-level:** `runAsNonRoot: true`, `fsGroup: 65532`, `seccompProfile: RuntimeDefault`, `automountServiceAccountToken: false`.
- **Container-level:** `runAsUser/runAsGroup: 65532`, `allowPrivilegeEscalation: false`, `readOnlyRootFilesystem: true`, `capabilities.drop: [ALL]`, with a writable `emptyDir` mounted at `/tmp`.
- **Reliability defaults:** CPU/memory requests+limits and `readiness`/`liveness`/`startup` probes.
- **Admission contract:** the `security.grype.io/high_critical` and `security.stock-trading.dev/sbom-digest` annotations are intentionally left empty in the base and must be populated by the CI overlay to pass the Kyverno policies (see [devsecops_ci_admission.md](devsecops_ci_admission.md)).

## 3. Residual risks (accepted)

- **Base-image CVEs over time.** Distroless reduces but does not eliminate CVE exposure; the Grype gate + nightly matrix are the compensating control for newly disclosed issues.
- **Application-layer risk.** Injection, authz flaws, and business-logic abuse are out of scope for this control set (see Non-goals) and are not addressed by image/runtime hardening.
- **Secrets at rest.** Secrets are referenced via Kubernetes `Secret` objects; envelope encryption / external KMS is not configured in the local validation environment.
- **`readOnlyRootFilesystem` + `/tmp`.** A writable `emptyDir` at `/tmp` is allowed for ephemeral scratch; services must not persist state to the container filesystem.

## 4. Non-goals

- Full enterprise compliance maturity (e.g., SLSA L4, FIPS-validated crypto).
- Network policy / service-mesh mTLS (separate concern from workload hardening).
- Multi-tenant isolation hardening (namespaces are single-tenant in this thesis scope).
- Performance tuning of resource requests/limits (values are conservative defaults, not load-tested).
