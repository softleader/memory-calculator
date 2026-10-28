---
name: test-performance
description: Diagnose slow tests and change test fixtures, environment lifetime, waits, discovery or runner parallelism while preserving behavior coverage. Use for test performance work and CI timing claims.
---

# Test performance

Use `test-design` for the behavior and mutation oracle. This skill owns the cost,
isolation and measurement decisions. A task's explicit coverage or experiment
boundary takes precedence; record excluded behavior when running a reduced scope.

## Choose the relevant reference

- Go tests, envtest, timers or `golang-testing`: [Go boundaries](references/go.md).
- React, Vitest, browser clocks or WebGL: [Browser boundaries](references/browser.md).
- Shell fixtures or before/after timing: [Fixtures and measurement](references/fixtures-and-measurement.md).

## Work from the bottleneck

1. Record the executed cases, entry command, source revision and current gate.
   Use the actual runner's discovery and output: file names alone do not prove
   that CI selects a case.
2. Separate setup/build/download, test execution, waiting and cleanup. Profile
   the expensive part before choosing resource reuse, smaller fixtures, clock
   control or concurrency. A faster cached result is not faster execution.
   Read the runner's CPU use and throttling (cgroup `cpu.stat`) to pick the
   direction. An idle share gains from overlapping independent phases, and work
   charged to another quota (a host Docker daemon, a remote cluster) overlaps
   nearly free. A saturated runner gains only from less work per case or
   frame; added concurrency there stretches every lane.
3. Keep each protected behavior assigned to a test. Reuse expensive resources
   only with independent case state, failure injection and cleanup. Run the
   changed cases alone and in a changed order; exercise parallelism when the
   implementation actually shares concurrent resources. A change that loads
   the runner harder shifts every case's timing: repeat the affected suites
   under the new load and root-cause the races it exposes.
4. Verify the candidate still catches the named break and goes red for
   nothing else (`test-design`).
5. Compare like-for-like samples and report failures with successful runs.
   Use the task's repetition and time budget; routine edits do not inherit
   another task's five-run benchmark or ten-minute target.

Focused runs shorten iteration. Delivery checks remain in the repository's
`AGENTS.md`; report checks not run and distinguish local timing from CI timing.
For a reduced-scope experiment, report its scope and result separately from
the complete gate. Restoring a gate needs a successful run with it enabled.
