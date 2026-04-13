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
