# Issue #67 source acceptance evidence

## Candidate

- Baseline commit: `c33af5700052687a8a87b06375f8268888c5289d`
- Host: `darwin/arm64`
- Toolchain: `go1.25.9`
- Language and toolchain directives remain `go 1.25.0` and `toolchain go1.25.9`.

| Module | Before | Candidate |
|---|---:|---:|
| `golang.org/x/crypto` | `v0.45.0` | `v0.52.0` |
| `golang.org/x/mod` | `v0.31.0` | `v0.35.0` |
| `golang.org/x/net` | `v0.47.0` | `v0.55.0` |
| `golang.org/x/sync` | `v0.19.0` | `v0.20.0` |
| `golang.org/x/sys` | `v0.41.0` | `v0.45.0` |
| `golang.org/x/text` | `v0.32.0` | `v0.37.0` |
| `golang.org/x/tools` | `v0.39.0` | `v0.44.0` |

`go.mod` and `go.sum` change only these seven pinned modules. MVS also selects
`x/term v0.43.0` and `x/telemetry v0.0.0-20260409153401-be6f6cb8b1fa`
transitively; neither is added to `go.mod`.

## Production-path contracts

- `TestWebApplicationType` executes the existing None, Reactive, Servlet, and
  resolver-error specs and restores the package resolver.
- `TestRun_ExportAndLoadedClassCountPrecedence` runs
  `prep -> boot -> calc -> out` in an isolated child process.
- `TestRun_LinuxResolverContract` ran on real `linux/amd64` and pins the ordered
  `JAVA_TOOL_OPTIONS` value after `/etc/resolv.conf` is read.
- `TestCalculator_Execute_CertificateLoaderContract` runs the public calculator
  path with isolated PKCS#12 and PEM inputs and proves the copied truststore was
  rewritten.

The certificate contract was checked against the exact upstream
`paketo-buildpacks/libjvm v1.46.0` tag at commit
`d0895b1355131c76a1ef2d998ea1cfcda19c1cce`. That runtime uses a passwordless
PKCS#12 truststore; `changeit` applies only to its JKS path.

## Verification

| Gate | Result |
|---|---|
| `go test -count=1 ./...` | pass |
| `go test -count=1 -race ./...` | pass |
| `go vet ./...` | pass |
| `go mod tidy -diff` | empty |
| Cross-package coverage before upgrade | `56.0%` |
| Cross-package coverage after upgrade | `56.0%` |
| Linux amd64 build | pass |
| Linux arm64 build | pass |
| macOS amd64 build | pass |
| macOS arm64 build | pass |
| Linux amd64 orchestration and certificate contracts | pass |
| `gofmt -l .` | empty |
| `git diff --check` | pass |

## Vulnerability comparison

Raw pinned outputs:

- [`govulncheck-before.txt`](govulncheck-before.txt)
- [`govulncheck-after.txt`](govulncheck-after.txt)

| `govulncheck v1.6.0` level | Before | Candidate |
|---|---:|---:|
| Symbol-level | 6 | 6 |
| Imported package only | 3 | 3 |
| Required module only | 20 | 6 |

No symbol-level finding was added. The six reachable findings are unchanged
and remain outside Issue #67: five affect Go `1.25.9`, and one affects
`go-pkcs12 v0.6.0`.

Verbose output still reports `GO-2026-5932` for
`golang.org/x/crypto/openpgp` at module level with `Fixed in: N/A`; it is not an
imported-package or symbol-level result. Harbor may still report that known
residual risk.

This is source acceptance only. GitHub Release `1.2.6` and downstream Harbor
rebuild/rescan evidence are separate acceptance stages.
