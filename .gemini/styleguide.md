# Style Guide for http-probe-test-app

## Go Code

- Use stdlib over third-party libraries when possible. The only external dependency is `prometheus/client_golang`.
- All HTTP handlers follow the pattern: `func xyzHandler(cfg Config) http.HandlerFunc` returning a closure.
- Configuration is loaded from environment variables via `getEnv`/`getEnvInt` helpers. Never hardcode defaults inline -- add them to `loadConfig()`.
- Use `atomic` operations for shared counters (`reqCount`) and `sync.RWMutex` for shared state (`readyToggle`).
- Tests use `httptest.NewRequest` + `httptest.NewRecorder` for handler tests. Call `resetGlobalState()` at the start of each test function.
- Test files mirror the source: `main.go` -> `main_test.go`, `config_test.go` for config helpers.
- No `t.Parallel()` -- tests mutate global state by design.
- CGO_ENABLED=0 for builds, CGO_ENABLED=1 for tests (race detector).
- Prometheus metrics use the naming convention `{prefix}_test_*`.

## Container Image

- Runtime stage uses `FROM scratch` (minimal attack surface).
- Build args: `VERSION` and `GIT_COMMIT` are injected via `-ldflags`.
- OCI labels (`org.opencontainers.image.*`) go in the runtime stage only.

## GitHub Actions Workflows

- Pin runners to `ubuntu-24.04` (not `ubuntu-latest`).
- Pin actions to major version tags (e.g., `@v6`, `@v4`), not SHAs -- Renovate handles updates.
- Use least-privilege permissions: per-job `permissions` blocks, never workflow-level write access.
- Schedule cron jobs staggered across the week (Tue/Wed/Sat pattern).

## Documentation

- README.md: user-facing (endpoints, config, usage).
- SECURITY.md: security policy, verification commands.
- CONTRIBUTING.md: developer-facing, references SECURITY.md for details.
- Each doc should be self-contained for its audience.

## Commits

- Use conventional commits: `feat:`, `fix:`, `chore:`, `ci:`, `docs:`.
- Only `feat:` and `fix:` trigger releases via release-please.
- Scopes are optional: `feat(docker):`, `fix(ci):`.
