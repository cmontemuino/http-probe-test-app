# Security Policy

## Supported Versions

Only the latest release receives security updates.

| Version | Supported |
|---------|-----------|
| latest (0.x) | ✅ |
| older releases | ❌ |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via GitHub Security Advisories:

1. Go to the [Security tab](https://github.com/cmontemuino/http-probe-test-app/security/advisories)
2. Click "Report a vulnerability"
3. Fill in the details

You should receive a response within 48 hours. If the issue is confirmed, we will:
- Release a patch as soon as possible
- Credit you in the release notes (unless you prefer to remain anonymous)

## Security Scanning

This project uses multiple layers of security scanning:

- **CodeQL**: Static analysis for code-level vulnerabilities (weekly + on PR)
- **Trivy**: Container image and filesystem vulnerability scanning (on PR + release)
- **govulncheck**: Go-specific vulnerability database checks (CI + local)
- **Renovate**: Automated dependency updates with vulnerability alerts
- **Dependabot**: Additional security-focused dependency alerts

All security findings are tracked in the [Security tab](https://github.com/cmontemuino/http-probe-test-app/security).

## Security Best Practices

When using this container image:

1. **Always use specific versions**: Avoid `latest` tag in production
2. **Monitor security advisories**: Subscribe to this repository's security advisories
3. **Keep updated**: Regularly pull newer versions with security patches
4. **Network isolation**: Run in isolated network environments when testing probe failures
5. **Resource limits**: Set CPU/memory limits to prevent DoS via artificial latency

## Vulnerability Disclosure Timeline

- **0–48h**: Initial response to report
- **2–7 days**: Vulnerability confirmation and patch development
- **7–14 days**: Release with security fix and public disclosure
