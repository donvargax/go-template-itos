# 8. Crosscutting concepts

This page holds what every package shares. `docs/engineering.md` has the rules for every
stack and `docs/go.md` the Go packages, tests and checks; neither is repeated here.
`docs/CLI.md` is the product contract the code serves.

## Ports, sealed failures and contract suites

A port is an interface in its engine's `port` package, and each method takes a
`context.Context` first, carrying cancellation and nothing else (decision 7). Its failures are
a sealed set: an interface with an unexported marker method, annotated `//sumtype:decl`, its
kinds plain types that name what went wrong without an OS's own errors (`Missing`,
`Unreadable` and `Cancelled` for `Files`). Each slice has one classifier, `Classify` in
`internal/slices/count/command`, that gives every kind its problem in one switch;
gochecksumtype refuses the switch when a kind is left out. `cmd` hands the classifiers to
`internal/cli`, which imports no domain, and an error none claims is a bug, exit 100
(decision 13, and the opening of `docs/CLI.md`).

The fake is ours, never a mock: an in-memory implementation in `port/porttest` that keeps
state. `porttest` also exports one contract suite per port, `FilesContract`, which the fake's
tests and the disk adapter's tests both run. A difference the suite finds is fixed on the wrong
side, never by dropping the case. An adapter's own tests keep what only a real disk can do.
Production code imports no fake, and the `mocks` rule refuses every mocking library.

Every port reaches the command through `cmd`'s `adapters` struct, bound on disk for `main` and
to the fakes for the in-process driver. `TestEveryPortIsBoundOnDiskAndInProcess` fails until a
new port is bound in both.

## The OS as a value

Code that differs by OS takes the OS as a value, detected once at the boundary, so a test can
give any other. Logic never branches on `runtime.GOOS`; forbidigo refuses it outside an edge.
A build-tagged file holds only a call one OS lacks, with no mutation site: the scenarios'
interrupt step is the example, `os.Interrupt` on Unix and `CTRL_BREAK_EVENT` on Windows.

## Testing model

- **The scenarios are the spec.** Every `feat` and `fix` is driven by Gherkin scenarios in
  `features/`, black-box: they judge only what the command prints, its exit code and the files
  it touches (`features/README.md`). Their failing steps are committed first, with the
  scenarios still `@wip`.
- **Two drivers, one set of steps** (decision 14). `cmd/go-template-itos`'s `TestScenarios`
  runs every file in process, with `run` called on buffers and the ports' fakes. `features`'
  `TestFeatures` builds the binary and runs the same files against it, which proves `main`'s
  wiring, the real adapters, exit codes and signals. It runs only where
  `GO_TEMPLATE_ITOS_FEATURES_BINARY_DRIVER=1` is set, which CI sets. A scenario only one driver
  can hold carries `@binary-only`, never a copy.
- **Unit tests** exercise the domain with the `porttest` fake, assert outcomes and typed
  failures, and keep what the scenarios cannot reach. Properties use `pgregory.net/rapid` in
  the domain's own tests; `FuzzCount` and `FuzzNewInput` are the fuzz targets.
- Each domain reaches 80% statement coverage from its own tests (`tools/bin/domain-coverage`).
- Mutation runs fresh in CI on each push. A surviving mutant is a missing test or scenario,
  never an exception.

## Gates by tier

| Tier | Where | What runs | Why there |
| --- | --- | --- | --- |
| Commit | `tools/hooks/pre-commit` | gofmt over the staged Go files | a fraction of a second, no network |
| Push | `tools/hooks/pre-push` | `go vet ./...` over the working tree | compiles every package and its tests in seconds |
| CI, per push | `.github/workflows/ci.yml` | lint, the architecture check, vulnerability and credential scans, tests, coverage, mutation, the scenarios under both drivers; macOS and Windows run `cmd` and `internal` and the scenarios in process | what a green push means (decision 15) |
| Nightly | `.github/workflows/nightly.yml` | every test and the binary driver on three systems, the race detector, fuzz, properties at a thousand checks, a coverage floor, a whole-tree mutation run | slow checks that nothing waits on; a red night opens one issue (decision 15) |
| Release | `.github/workflows/release.yml` | the release, from a push to `main` once CI passes there | the only workflow that cuts one |

The lint and vulnerability checks need a download, so they stay in CI. A valid check's failure
is fixed in the work, never by weakening the check. Tools are pinned (`docs/upgrading.md`).

## Recipe: adding a command

1. Specify it: `features/<command>.feature` and its steps in `features/steps.go`; commit the
   steps red, scenarios `@wip`.
2. Add `internal/slices/<command>/command`, laid out as `internal/slices/count/command`:
   its kong struct, `Query`, `Handler` and `Classify`. Use engines through their ports.
3. Bind it in `cmd/go-template-itos/main.go` and hand its classifier to the UI.
4. Add a port to `adapters` on disk and in process if the command needs a new one.
5. Add the cases the rules must accept and refuse to `tools/testdata/architecture/cases.json`,
   and an entry to `MAP.md`.

## Recipe: moving code into a slice

Move a package in the task that makes it part of a slice. Move its scenarios first, so they
hold the move, then the code with the same bytes and exit codes. Failures are classified by
type, never by text. Anything a second slice needs goes to an engine first. Add the fixtures
for each rule the move touches to `tools/testdata/architecture/cases.json`, never loosening a
rule. `docs/go.md`'s "Adding a package" has the detail.
