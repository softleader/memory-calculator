# Go test boundaries

This repository policy explicitly supersedes conflicting defaults in
`samber/cc-skills-golang@golang-testing`. Keep the upstream skill unchanged so
updates remain attributable. `test-design` owns assertion quality;
`use-modern-go` owns version-specific idioms.

- Keep existing envtest cases in the ordinary race gate. Introducing an
  `integration` build tag requires preserving their CI execution, not silently
  removing them from `go test`. Existing browser tags keep their own entrypoint.
- Preserve the repository's test-file organization. A filename or a target
  runtime under 1 ms is not a reason to rewrite a behavioral test.
- Add gotests, goleak, testify or other dependencies only when the task needs
  them; an upstream dependency list does not authorize installing tools.
- Choose parallelism from independent state and available resources. A suite
  audit does not automatically authorize upstream fan-out or a higher model
  reasoning setting; follow the active task's delegation and resource limits.

## Match the Kubernetes boundary

Use pure tests for transformations and decisions. Keep API validation, status,
resourceVersion, watch, patch and persistence contracts on a real API server.
The controller-runtime fake client is not a faster equivalent for those contracts.
See the [runtime-version fake client documentation](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/client/fake).

When startup is the measured cost, a suite may lazily share envtest's API server
and etcd. Keep namespace, client configuration, HTTP server, cancellation and
fault injection local to each case. A pure case must still run without envtest
assets. Clean up resources the test owns: envtest has no controller manager or
automatic garbage collection. Deleting a namespace is not proof that its
contents were collected. [Kubebuilder envtest](https://book.kubebuilder.io/reference/envtest)

Measure whole-suite wall time separately from Go benchmarks. `ns/op` measures
the benchmarked operation, not envtest startup, compilation or CI elapsed time.

## Red only for the break

`test-design` owns the rule; these are this repository's means of meeting it in Go.

- **Controlled clock.** `testing/synctest` controls self-contained goroutines
  and timers. The fake pi, the Operator binary, envtest and real network I/O
  stay on real time. Check the module's Go directive as well as the toolchain
  before using newer APIs. [Go synctest](https://pkg.go.dev/testing/synctest)
- **Conditions on real time.** `internal/testwait.Until` retries a check and
  gives up `Margin` before `go test -timeout` with the check's last error. In
  `cmd/operator`, `interruptOnceHarnessRuns` signals the Operator once the fake
  pi has written its output. The completion timeout counts from the Operator's
  start and covers the steps before the Agent Harness: a case that needs the
  Agent Harness running skips when the summary shows it never started, as
  `TestCompletionTimeoutInterruptsHarnessThatKeepsOutputting` does.
- **Load evidence.** Run the case beside 12 to 30 `yes >/dev/null` busy loops
  with `go test -race -count=40 -run '^TestName$' ./<package>` and report the
  tally before and after the fix. CI-only envtest flakes reproduce with a
  `go test -race -c` binary run from the package directory with
  `KUBEBUILDER_ASSETS` set, about 8 concurrent processes at `GOMAXPROCS=1`
  plus busy loops; they hit a few percent of runs, so one local pass is a
  sample, not evidence.
