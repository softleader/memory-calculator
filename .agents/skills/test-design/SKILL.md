---
name: test-design
description: "Load before editing any test, even a one-line change to an expected value, wait or timeout, and when judging whether a test can fail at all: writing a test, fixing a failing, flaky or CI-red one, updating one after an intended code change, adding a mock, stub, fake, helper, wait or cleanup, or asserting a config value, constant or label. Designs tests that catch a real break."
---

# test-design

A test exists to catch a specific break. These principles govern every rule here:

1. Every test names **the break** it catches.
2. Every test exercises the real thing.
3. Every test goes red for its break and for nothing else.

A test written first and watched failing against real code has shown it can fail. A test that stays green under load and in a reordered suite, with nothing broken, has shown, for the conditions it ran under, that its red means a break. These rules are what those checks look for, whichever order the test and the code were written in.

## Name the break

Before writing the test body, answer: **what production change makes this test fail, and is that change a bug or a decision?** A test earns its place by catching a wrong branch, a missing side effect, a wrong argument, a boundary case or a broken contract.

**The oracle is independent.** The oracle is where the expected value comes from. Use literals and hand-checked fixtures; table-driven tests with literal expected values are the preferred shape. An expectation computed by the code under test, or by its helpers, passes whatever that code does.

**Decisions are not breaks.** A constant's value, exact message wording and private structure are decisions: a test that fails only when one of them changes fires on redesign and sleeps through bugs. Test the behaviour that depends on the decision. A value the spec or configuration states is covered by a test of what it changes; reading it back and asserting it covers nothing. For a 30-second timeout, the test shows a call still waiting at 30 seconds is abandoned, not that the setting equals 30.

**Rendered text is a contract only when the reader acts on it.** An error message the user must read to recover, a legal notice, a status the user decides on: assert the rendered text. A label or heading asserted only because the spec wrote it is a decision; assert the action it leads to.

**Behaviour, not text.** A script, configuration file or component is tested by running it against controlled inputs and asserting outputs, side effects or exit codes. Asserting that its source contains a line proves only that the source is the source. Instructions written for agents are tested by the consuming agent's behaviour.

**Your contract, not the framework's.** Test what your code promises at its boundaries: the route it registers, the query it emits, the payload it produces. Upstream mechanics are upstream's tests to write; asserting that the framework invokes a registered handler tests the framework. When upstream behaviour genuinely surprised you, write one narrow characterization test that names the assumption. Inside your code the same boundary holds: constructors, getters, constants and trivial forwarding earn a test only when they validate, normalize, default, derive, enforce or cause a side effect; otherwise assert the first consumer-visible result that depends on them.

Gate, before the test body:

- Name the production change that fails this test. None → redesign around an observable behaviour.
- "The source text changed" → run the artifact and assert its effects.
- Only a decision fails it → test the behaviour that depends on the decision.
- Confirm the oracle is independent of the code under test; if it reuses the code's logic or helpers, replace it with a literal or a hand-checked fixture.

## Exercise the real thing

A **double** is anything standing in for a real collaborator: a mock, stub, fake, spy or recorded response.

**The double earns no assertion.** An assertion on the double passes when the double is present and fails when it is absent; it says nothing about the component. Assert the component's behaviour. When the double is what the assertion checks, unmock it or delete the assertion. An interaction check keeps the call and adds a meaningful assertion on the captured arguments or the resulting state.

**The double sits at the slow or external level.** Learn every side effect of the real method before replacing it; replace the slow or external operation below those side effects and keep what the test depends on real. When unsure, run the test against the real implementation first and observe what has to happen.

**The double keeps the real failure modes.** The real dependency can time out, reject, throw or return partial data. When the test claims to cover one of those outcomes, the double produces it; a double that always succeeds proves the happy path and nothing else, and the error-handling branch stays unexercised.

**Doubles are specific.** When arguments, call counts or ordering are part of the contract, assert them; a fake that accepts anything verifies nothing. Give each branch (success, error, malformed) its own fixture or spy, so the wrong branch cannot satisfy the expectation.

**Doubles mirror real data completely.** Model the complete structure as it exists in reality, every documented field, not only the fields the test reads. A partial double passes while downstream code that reads an omitted field breaks in integration.

**Test-only code lives in test utilities.** Cleanup or reset that only tests call belongs in a test utility, never as a method on the production class. Called only from tests, or the class does not own the resource's lifecycle → test utility.

**Real components beat complex doubles.** When double setup outgrows the test logic, when the double lacks methods the real component has, or when tests break because the double changed, switch to an integration test with real components.

Gate, before adding a double or a test helper:

- List the real method's side effects; keep the ones the test depends on real and replace the slow or external level below them.
- The double's responses mirror the complete real structure, and it produces the failures the test claims to cover.
- A method only tests call lives in test utilities.
- About to assert on the double itself → unmock it or delete the assertion.

## Red only for the break

The verdict is a function of the production code. Load, scheduling, the order the suite runs in, the wall clock and state another test left behind are not inputs: a test any of them can turn red reports on the machine, and its readers learn to rerun it instead of reading it.

**Reach the state under test through the condition.** A **condition** is an observable fact the test can check: the process wrote its output, the state became `ready`, the second report landed. When the assertion needs the system to have reached a state first, wait for that condition and act on it. A sleep, a short idle timeout or a fixed number of ticks is a guess at how long the work takes; under load the guess loses, and the test fails with nothing broken.

**Time under test is a controlled input.** When the behaviour under test is itself about time (a timeout, a debounce, a poll interval), drive it from a clock the test controls, so the assertion reads the behaviour and not the machine's speed. Where the clock cannot be controlled, because the time passes in another process or on a real connection, anchor the interval on an observed condition and size it to cover every piece of work that must finish inside it: a wait that must observe a second poll starts when the first poll was observed and lasts at least one interval; a timeout shortened for the test starts only after the setup the assertion depends on has finished.

**A wait gives up on the runner's deadline and says what it waited for.** The deadline bounds a stuck test; it never decides whether an assertion holds. Every wait, cleanup included, draws on the time the runner gives the test, and the waits before cleanup leave it room: a shorter cutoff of a wait's own fails while the runner still had time, and a fresh budget per wait outruns cleanup. On expiry the failure names the condition and the last state observed, so a real red is read, not rerun.

**The test owns its state.** Create what the test needs in its own setup, in a location the test owns, and leave nothing behind: files, environment, globals, registered handlers, started processes. State created at module load or inherited from a neighbouring test makes the verdict depend on the order the suite runs in.

**A red without a break is a test bug.** When a test fails and nothing is broken, find what besides the break it depends on and remove the dependency: a guessed interval becomes a condition, a cutoff of the wait's own becomes the runner's deadline. A longer guess or a retry on the test keeps the dependency; it records that the verdict is not yet a function of the code.

Gate, before the test body:

- Name what besides the break could turn this test red: load, order, the clock, state another test leaves, the environment. Say what pins each; one you cannot pin → redesign around a condition the test can observe or a clock it controls.
- Every wait waits for a condition and gives up on the runner's deadline, with room for cleanup; a measured duration appears only where time is the behaviour under test.
- Everything the test starts, writes or registers is released when the test ends.

## Only the tests the behaviour needs

Ship the tests the behaviour needs and only those. Trivial code and human prose earn none; a test written to satisfy process, or to raise a coverage number, costs maintenance forever and protects nothing.

## Mutation list

Before reporting a test file done, write the mutation list and include it in the completion report. For each plausible mutation of the production code, trace it through the paths that remain (a retry, a fallback, cleanup) and name the test that goes red; if another path still yields the asserted outcome, that test does not catch the mutation. Distinguish removing one path from removing the entire behaviour. Consider:

- a wrong constant or argument
- a wrong branch handler
- a missing state change or side effect
- an empty or default return
- missing validation for zero, empty, null, unauthorized or malformed input

A mutation no test catches marks either an unprotected behaviour or a test that cannot fail: fix one or the other before the report, or list the mutation as `unprotected` with the reason a test is out of reach. The list is the evidence that the check happened: a report that says the tests were checked without the list has not done this step.

For a test that waits, polls, signals a process or shares a resource with other tests, the report also names what besides the break could turn it red and what pins each: the condition it waits for, the clock it controls, the state it owns. When a runner is available, repeated runs with the machine against the test, under load or in another order, are the evidence; give the command, the tally and the seed or trace when the runner keeps one.

## Warning signs

- Setup and assertion share the same object, guaranteeing equality.
- The oracle is hidden behind a loop, builder or helper.
- The test fails on every intentional change and never on accidental breakage.
- The test can fail only through a crash, a panic or a missing selector.
- The test reads source text, or asserts that a removed symbol stays removed.
- The test would still pass if only the framework remained.
- The test exists for coverage and checks no outcome or side effect.
- An assertion targets the double, or fails when the double is removed.
- A method is called only from test files.
- Double setup is more than half the test, or nobody can say why the double is needed.
- A sleep or a short timeout stands between setup and assertion, and nothing under test is about time.
- The test passes alone and fails in the suite, or passes on one machine and fails under load.
- The last fix to the test was a longer timeout or a retry.
- Test state is created at module load, or in a location the test does not own.

## Tool traps

The rules above apply as written in any stack. With these tools the obvious move satisfies a rule's letter and still misses the break:

- **Mockito `verify()` under ArchUnit.** `verify()` alone satisfies `testsShouldIncludeAssertions()`, and satisfying that rule leaves the assertion ineffective: pair the `verify()` with an assertion on the captured arguments or the resulting state.
- **Spring configuration values.** Bind the test's own value with `@SpringBootTest(properties = …)` and exercise what it controls; the production value stays in `application.yml`. A timeout bound down to milliseconds starts only after every piece of setup the assertion depends on has finished, or load, not the code, trips it.
- **Mocked repositories.** A stubbed `findBy…` returns whatever the test says and never runs the query: use a real engine through Testcontainers or the repository's test profile.
- **Module mocks.** `vi.mock('./module')` replaces the whole module and swallows side effects the test depends on, such as a cache write: replace the client below it, or intercept HTTP at the network boundary with MSW and keep the fetching code real.
