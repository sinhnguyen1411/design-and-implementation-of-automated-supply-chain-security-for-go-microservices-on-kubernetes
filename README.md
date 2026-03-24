# Design and Implementation of Automated Supply Chain Security for Go Microservices on Kubernetes

This repository provides a practical DevSecOps baseline for implementing and validating software supply chain security controls for a Go microservice deployed on Kubernetes.

## What This Repository Delivers
- A CI pipeline for build, SBOM generation, vulnerability scanning, signing, and attestation.
- Kubernetes admission policies for signature/provenance and metadata enforcement.
- Reproducible local validation workflow using Kind + Kyverno.
- Thesis-aligned documentation, traceability, and evidence artifacts.

## Architecture Overview
```mermaid
flowchart TD
  subgraph DevLayer["Developer and Source"]
    Dev["Developer"]
    Repo["Git repository"]
    Dev -->|"git push"| Repo
  end

  subgraph CILayer["CICD Pipeline and SCS"]
    CI["CICD pipeline"]
    Integrity["Verify Go dependency integrity"]
    Govuln["Scan Go vulnerabilities govulncheck"]
    Build["Build Go binary and Docker image"]
    SBOM["Generate SBOM Syft"]
    Scan["Scan vulnerabilities Grype"]
    FailBuild["Fail pipeline and block image"]
    Sign["Sign image Cosign"]
    Attest["Create SLSA provenance"]
  end

  Repo -->|"trigger pipeline"| CI
  CI --> Integrity --> Govuln --> Build --> SBOM --> Scan
  Govuln -->|"Fail Go vulnerability gate"| FailBuild
  Scan -->|"Fail CVE high or critical"| FailBuild
  Scan -->|"Pass"| Sign --> Attest

  subgraph DeployLayer["Registry and Kubernetes"]
    Registry["Secure container registry"]
    Deploy["Apply CI rendered deployment overlay"]
    K8s["Kubernetes cluster"]
    AC["Admission controller or Kyverno"]
    Verify["Verify signature provenance SBOM and policy"]
    Decision{"Policy ok"}
    Reject["Reject deployment"]
    Pod["Running pod on Kubernetes"]
  end

  Sign --> Registry
  Attest --> Registry
  Dev --> Deploy --> K8s
  Registry --> AC
  K8s --> AC
  AC --> Verify --> Decision
  Decision -->|"No"| Reject
  Decision -->|"Yes"| Pod

  subgraph LogLayer["Logs and Evidence"]
    LogsCI["CICD logs"]
    LogsAC["Admission controller logs"]
    Evidence["Security evidence for report"]
  end

  CI --> LogsCI --> Evidence
  AC --> LogsAC --> Evidence
```

## Quickstart
### 1) Local service run
```bash
go test ./...
go run main.go server --config cmd/server/config/local.yaml
```

### 2) Trigger secure supply-chain workflow
- Push to branch `Thesis-SCS` or manually run `.github/workflows/secure-supply-chain.yml`.

### 3) Bootstrap local admission demo
```bash
./scripts/devsecops_kind_bootstrap.sh
kubectl get clusterpolicies
```

## Thesis Documentation
- [Thesis specification (English)](docs/thesis_spec_en.md)
- [Interactive architecture diagram (HTML + Mermaid)](docs/scs_architecture_diagram.html)
- [Go dependency integrity baseline](docs/go_dependency_integrity_baseline.md)
- [CI and admission flow](docs/devsecops_ci_admission.md)
- [Implementation roadmap and milestones](docs/implementation_roadmap.md)
- [Demo evidence logs](docs/demo_evidence.md)

## Notes
- Current enforcement baseline is Kyverno-based.
- Sigstore Policy Controller remains an optional future extension.
