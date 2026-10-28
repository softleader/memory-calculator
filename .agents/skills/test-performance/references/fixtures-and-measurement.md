# Fixtures and measurement

## Minimum complete fixture

List what the real consumer reads or executes, including transitive package and
configuration dependencies. Copy those inputs, not a convenient parent tree
containing node_modules, build output or unrelated data.

Run the pristine fixture through the real checker first. Give each mutation an
independent copy, verify that it changes its intended input, and require the
checker to reject that break. A missing prerequisite that makes every mutant
fail invalidates the result.

## Comparable samples

Record source SHA, command, scope, tool versions, host/runner resources, worker
limits, cache state and exit status. Define cold/warm explicitly. Run before and
after sequentially under equivalent conditions; use the task's sample count and
retain all samples, including failures and cancellations.

For Go execution measurements, disable result caching with `-count=1`; build
cache is a separate choice. Report setup/build and execution separately. For CI,
use a fixed job and build number, its checkout SHA, stage wall time and artifacts.
Parallel stage durations are not additive build wall time.

Report medians and spread when repeated samples exist. A single successful run
is a sample, not proof of stability. Different revisions, hardware or case lists
are observational comparisons; state what changed rather than assigning the
entire time difference to one optimization.

Acceptance compares the complete scope requested by the task. A fast package,
source-only run or disabled 3D lane cannot establish a full source/browser/Kind
time target. Preserve the smallest evidence needed to reproduce a failure;
follow the repository's secret-scanning rules before publishing logs.
