# Browser and frontend boundaries

Use the installed `vitest` skill for runner APIs and `test-design` for what the
test must catch. Read the current package and config before choosing discovery,
environment or workers. Validate discovery when adding a new extension such as
`.test.tsx`; an existing `.test.ts` include does not select it.

For the Web UI, `cluster/web-ui/frontend/e2e/README.md` decides which breaks
keep a browser case and which go to Vitest; apply it before adding, moving or
merging a browser case. Compare Web UI canvas pixels with `paintedCanvas`
(`e2e/office-frame.ts`): a
plain element screenshot includes the labels and cut-ins above the canvas,
which change with the same data and pass without a WebGL repaint.

## Source and deployed bundles

A Vite-module observer needs the source server; it cannot import source modules
from a production dist server. Check both the runner's grep groups and the
actual execution counts. A source-only success does not prove deployed chunks,
authentication or Kind acceptance.

An explicit reduced-scope switch must list exclusions and remain distinguishable
from a complete run. Check its effective CI value, not only its declared default.
When a group contains only excluded cases, permit an empty selection only in the
explicit reduced mode; a typo in the complete-mode filter must still fail.

## Clocks and rendering

For on-demand rendering, install the observer, request a legitimate repaint,
then wait for the observable result. Reduced motion does not guarantee another
RAF will arrive spontaneously. `runFor` and `fastForward` have different callback
semantics; advancing a clock does not prove an asynchronous response completed.
[Playwright clock](https://playwright.dev/docs/clock)

Long navigation can step the real simulation without rasterizing every step,
but retain small-step semantics and independent real-frame tests for event
wiring, pause/resume, selection repaint and lifecycle behavior. Before accepting
this split, demonstrate that missing the final repaint makes the visual test
fail. DOM coordinates alone cannot prove the scene was drawn.

Read back the renderer, draw-buffer size and DPR for performance claims.
`WEB_UI_SOFTWARE_RENDERING` and `WEB_UI_SOFTWARE_RENDERER` are configuration
names; the backend is whatever the read-back reports. SwiftShader and Mesa
llvmpipe both draw on the CPU, and their cost is mostly fill: it scales with
drawing-buffer pixels (DPR squared) far more than with draw calls. Compare
renderers by per-frame and cold-ready time on the same scene and quota.
Raising workers inside one constrained runner can increase contention.
Compare worker settings under the same CPU quota and test list.
Keep hardware and software screenshot baselines separate.
[Chromium SwiftShader](https://chromium.googlesource.com/chromium/src/+/HEAD/docs/gpu/swiftshader.md)

For a failed simulation, save the initial checkpoint, obstacles, seed and
asserted final state. Reproduce without graphics when possible to distinguish
navigation correctness from rendering cost. Repeatedly selecting a passing seed
does not resolve the defect.
