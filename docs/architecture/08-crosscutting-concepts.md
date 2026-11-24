# 8. Crosscutting concepts

This page holds what every package shares. `docs/engineering.md` has the rules for every
stack and `docs/go.md` the Go packages, tests and checks; neither is repeated here.

## Ports, sealed failures and contract suites

A port is an interface in its engine's `port` package, and each method takes a
`context.Context` first, carrying cancellation and nothing else (decision 7). Its failures are
a sealed set: plain types that say what went wrong without an operating system's own errors
(`Missing` and `Unreadable` for `Files`).

The fake is ours, never a mock: an in-memory implementation in `port/porttest` that keeps
state. `porttest` also exports one contract suite per port, `FilesContract`, which the fake's
tests and the disk adapter's tests both run. A difference the suite finds is fixed on the wrong
side, never by dropping the case. An adapter's own tests keep what only a real disk can do,
such as permissions and error texts. Production code imports no fake. The `mocks` rule refuses
every mocking library.

## The OS as a value

Code that differs by OS takes the OS as a value, detected once at the boundary, so a test can
give any other. Logic never branches on `runtime.GOOS`; forbidigo refuses it outside an edge.
A build-tagged file holds only a call one OS lacks, with no mutation site. The counting
example has none.

## Testing model

- Unit tests exercise the domain with the `porttest` fake and assert outcomes and typed
  failures, never a script of calls.
- Properties use `pgregory.net/rapid` in the domain's own tests (the `rapid` rule holds that).
  Every run uses one case per property; the nightly runs a thousand.
- `FuzzCount` is the example fuzz target; add one wherever hostile input enters.
- Each domain reaches 80% statement coverage from its own tests (`tools/bin/domain-coverage`).
- Mutation runs fresh in CI on each push. A surviving mutant is a missing test, never an
  exception.

This stack has no public entry point, so it has no scenarios; a command line adds them.

## Gates by tier

| Tier | Where | What runs | Why there |
| --- | --- | --- | --- |
| Commit | `tools/hooks/pre-commit` | gofmt over the staged Go files | a fraction of a second, no network |
| Push | `tools/hooks/pre-push` | `go vet ./...` over the working tree | compiles every package and its tests in seconds |
| CI, per push | `.github/workflows/ci.yml` | lint, the architecture check, vulnerability and credential scans, tests, coverage, mutation; macOS and Windows run `internal` | what a green push means (decision 15) |
| Nightly | `.github/workflows/nightly.yml` | every test on three systems, the race detector, fuzz, properties at a thousand checks, a coverage floor, a whole-tree mutation run | slow checks that nothing waits on; a red night opens one issue (decision 15) |

The lint and vulnerability checks need a download, so they stay in CI. A valid check's failure
is fixed in the work, never by weakening the check. Tools are pinned (`docs/upgrading.md`).

## Recipe: adding an engine

1. Make `internal/engine/<name>/{domain,port,port/porttest}` and an adapter beside them; no
   rule edit is needed, and the domain is measured from then on.
2. Write the port's contract suite in `porttest` and run it from the fake's and the adapter's
   tests.
3. Add the cases the rules must accept and refuse to `tools/testdata/architecture/cases.json`.
4. Add an entry to `MAP.md`.

`docs/go.md`'s "Adding a package" has the detail. When a command line is added, its code goes
in a slice that imports the engine through its port.
