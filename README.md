# http-probe-test-app

[![CI](https://github.com/cmontemuino/http-probe-test-app/actions/workflows/ci.yml/badge.svg)](https://github.com/cmontemuino/http-probe-test-app/actions/workflows/ci.yml)
[![CodeQL](https://github.com/cmontemuino/http-probe-test-app/workflows/CodeQL/badge.svg)](https://github.com/cmontemuino/http-probe-test-app/actions/workflows/codeql.yml)
[![Trivy](https://github.com/cmontemuino/http-probe-test-app/workflows/Trivy%20Container%20Scan/badge.svg)](https://github.com/cmontemuino/http-probe-test-app/security/code-scanning)
[![Renovate](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://renovatebot.com)

A small HTTP service to **exercise Kubernetes liveness/readiness probes** and related monitoring:
artificial latency, configurable failure modes, and Prometheus metrics.

## Endpoints

| Endpoint | Method | Description |
|---|---:|---|
| `/` | GET | Returns a simple text response; increments metrics. |
| `/healthz` | GET | Liveness endpoint (200/500). Supports `?fail=1` override. |
| `/readyz` | GET | Readiness endpoint (200/503). Supports delay/toggle/thresholds. |
| `/toggle-ready` | POST | Toggles readiness state (affects `/readyz` and `/info`). |
| `/info` | GET | JSON info (version, commit, uptime, readiness, env). |
| `/metrics` | GET | Prometheus metrics. |

## Configuration (environment variables)

| Name | Default | Notes |
|---|---|---|
| `PORT` | `8080` | HTTP listen port. |
| `PREFIX` | `dummy` | Metrics prefix (e.g. `probe_test_requests_total`). |
| `CLUSTER_LABEL` | `unknown` | Exported as a metrics label and returned by `/info`. |
| `POD_NAME` | `unknown` | Metrics label + `/info`. |
| `NAMESPACE` | `unknown` | `/info`. |
| `NODE_NAME` | `unknown` | Metrics label + `/info`. |
| `ENVIRONMENT` | `development` | `/info`. |
| `EXTRA_LATENCY_MS` | `0` | Adds fixed latency to `/` responses. |
| `LATENCY_JITTER_MS` | `0` | Adds random extra latency (0..N) on top of `EXTRA_LATENCY_MS`. |
| `FAIL_LIVENESS_AFTER_N_REQUESTS` | `0` | If >0, `/healthz` returns 500 after N total requests. |
| `FAIL_READINESS_AFTER_N_REQUESTS` | `0` | If >0, `/readyz` returns 503 after N total requests. |
| `READY_DELAY_SECONDS` | `0` | If >0, `/readyz` returns 503 until delay elapses after start. |

## Quick start (container)

```bash
docker run --rm -p 8080:8080 ghcr.io/cmontemuino/http-probe-test-app:latest
curl -sSf http://localhost:8080/healthz
```

## Local development

### Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) (for `make govulncheck`): `go install golang.org/x/vuln/cmd/govulncheck@latest`
- [golangci-lint](https://golangci-lint.run/welcome/install/) (for `make lint`)
- [Docker](https://docs.docker.com/get-docker/) (for `make docker-build`)

### Build, test, lint

```bash
make build
make test
make lint
```

## Docker

```bash
docker build -t http-probe-test-app:local .
docker run --rm -p 8080:8080 http-probe-test-app:local
```

## Kubernetes

See `deploy/` for a minimal `Deployment` and `Service` you can apply to a cluster.

## Supply Chain Security

Every released container image includes:

- **SLSA Build Provenance** (Level 3): Cryptographically signed attestation of how the image was built
- **SBOM** (Software Bill of Materials): Complete list of packages in the image (SPDX format)
- **Sigstore Signatures**: All attestations are signed using Sigstore via GitHub's OIDC provider

### Verifying Image Provenance

Install the [GitHub CLI](https://cli.github.com/) and verify any image:

```bash
gh attestation verify oci://ghcr.io/cmontemuino/http-probe-test-app:latest \
  --owner cmontemuino
```

### Downloading Attestations

```bash
# Download provenance attestation
gh attestation download oci://ghcr.io/cmontemuino/http-probe-test-app:latest \
  --owner cmontemuino \
  --predicate-type https://slsa.dev/provenance/v1

# Download SBOM attestation
gh attestation download oci://ghcr.io/cmontemuino/http-probe-test-app:latest \
  --owner cmontemuino \
  --predicate-type https://spdx.dev/Document
```

### Kubernetes Policy Enforcement

Use [Sigstore Policy Controller](https://docs.sigstore.dev/policy-controller/overview/) or
[Kyverno](https://kyverno.io/) to enforce provenance verification in your cluster before
admitting images.

## Security

This project implements comprehensive security scanning:

- **CodeQL**: Automated code security analysis
- **Trivy**: Container vulnerability scanning
- **govulncheck**: Go vulnerability database checks
- **Renovate**: Automated dependency updates
- **Dependabot**: Security-focused alerts

Security findings are tracked in the [Security tab](https://github.com/cmontemuino/http-probe-test-app/security).

See [SECURITY.md](SECURITY.md) for vulnerability reporting and security policies.

## License

Apache-2.0. See `LICENSE`.
