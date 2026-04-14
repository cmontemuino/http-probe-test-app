# Contributing

## Development workflow

1. Fork the repository
2. Create a branch from `main`
3. Make your change (keep it focused)
4. Run validators locally:
   - `make test`
   - `make lint`
5. Open a pull request

## Commit messages

This repository expects **Conventional Commits** so release automation can generate changelogs.

Examples:

- `feat: add new endpoint`
- `fix: handle missing env var`
- `chore: bump dependency`

## Local development

- Build: `make build`
- Test (race): `make test`
- Lint: `make lint`

## Dependency Updates

This project uses [Renovate](https://github.com/renovatebot/renovate) to automatically update dependencies.

- **Schedule**: Weekly on Tuesday mornings (UTC)
- **Grouping**: Related updates are combined into single PRs
- **Auto-merge**: Minor/patch updates auto-merge when CI passes
- **Manual review**: Major updates and base image changes require review

### Reviewing Renovate PRs

1. **Check the changelog**: Renovate includes release notes in PR descriptions
2. **Verify CI passes**: All tests must pass before merge
3. **Major updates**: Review breaking changes, run local tests if needed
4. **Docker base image**: Test image builds locally before merging

### Security Alerts

Renovate creates immediate PRs for security vulnerabilities. These are:
- Labeled with `security`
- Auto-merged if CI passes
- Prioritized over regular updates

## Security

### Running Security Scans Locally

Before submitting a PR, run these security checks:

```bash
# Go vulnerability check
make govulncheck

# Container vulnerability scan (requires Docker + Trivy installed)
make docker-build
trivy image --severity CRITICAL,HIGH test:local
```

### Security Workflow

All PRs are automatically scanned for:
1. **Code vulnerabilities** (CodeQL)
2. **Container vulnerabilities** (Trivy)
3. **Dependency vulnerabilities** (Renovate + Dependabot)

If security issues are detected:
- **CRITICAL/HIGH**: CI fails, must be fixed before merge
- **MEDIUM/LOW**: Reported in Security tab, can be addressed later

### Supply Chain Security

Released images include SLSA provenance and SBOM attestations, generated automatically
during the release workflow. See [SECURITY.md](SECURITY.md) for verification instructions.

### OpenSSF Scorecard

This project is monitored by [OpenSSF Scorecard](https://scorecard.dev/viewer/?uri=github.com/cmontemuino/http-probe-test-app),
which evaluates security practices. When contributing, keep in mind:
- Pin GitHub Actions to commit SHAs (Renovate handles this automatically)
- Use conventional commits
- Keep dependencies up to date

### Reporting Security Issues

**Do not open public issues for security vulnerabilities.**

Use GitHub Security Advisories: https://github.com/cmontemuino/http-probe-test-app/security/advisories

See [SECURITY.md](SECURITY.md) for details.
