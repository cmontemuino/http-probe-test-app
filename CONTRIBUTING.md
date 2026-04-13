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
