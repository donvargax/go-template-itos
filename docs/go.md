# Go design and quality

This guide adds what is specific to Go to `docs/engineering.md`, which it does not repeat.
It describes what this project actually holds and checks. The policy is harvested from
[itos-template](https://github.com/donvargax/itos-template) (its decisions on thin layers
over a domain, unit tests with fakes, rapid properties and the OS as a value), adapted to
this module with its owner's permission (`CONTRIBUTORS.md`).

## Packages

The counting example is one slice under `internal/count`:

- `domain` owns counting and its result. It reads input only through the `Files` port.
- `port` declares `Files` and its typed failures, `Missing` and `Unreadable`.
- `port/porttest` holds our in-memory fake of `Files` for the domain's tests.
- `disk` is the concrete adapter. It reads real files and translates OS errors into the
  port's failures.

There is no `cmd/` package and no public entry point yet. A command line adds its own
handler and input adapter over this domain.

`.golangci.yml` enforces these boundaries with depguard, one rule per layer:

- `domain` and `ports` may import only the domain and its ports from this module.
- `domain-io` denies `os` and `net` to production domain and port code. Their tests may
  still read fixtures.
- `infra` lets a production adapter import only the ports, never the domain, another
  adapter or a fake. It names `disk`, and `input`, the adapter a command line adds.
- `mocks` denies mocking libraries everywhere. Fakes live in `porttest`.
- `rapid` limits the property library to the domain's tests.

The lists name each package explicitly, so a new package is unchecked until a rule names
it. When you add a domain package, a slice or an adapter:

1. Add its folder as a `**/<folder>/**` glob to the matching rule's `files`, and its
   import path to the `allow` list of each rule that may import it.
2. Add compile-valid fixtures that the rule must accept and reject to
   `tools/testdata/architecture`. `tools/bin/architecture-check` builds every fixture, then
   runs the pinned linter and this configuration against them, so a rule that stops
   matching fails there instead of passing silently.
3. A new domain package is measured by `tools/bin/domain-coverage`, which reads its list
   from the `domain` rule's `files`.

Code that differs by OS takes the OS as a value. Logic never branches on `runtime.GOOS`.
A build-tagged file holds only a call that one OS lacks, with no mutation site. The
counting example has no OS-specific logic.

## Tests

The domain's unit tests use the `porttest` fake and assert results and typed failures.
The disk adapter's tests use real temporary files, since translating OS errors is what
the adapter does.

Properties use `pgregory.net/rapid` and live in the domain's tests. The counting property
joins newline-terminated pieces and checks that their counts add. Arbitrary concatenation
can join two words, so it is not promised to be additive. Every run sets
`RAPID_CHECKS=1`: one new random case per property, as fast as an example test, with the
seed printed on failure. This project has no scheduled high-count run.

Tests run on Linux, macOS and Windows in CI's platform jobs. The import-boundary self-test
runs there as well, because it checks paths and diagnostics that differ by platform.

## Quality checks

`itos.yaml`'s CI plan runs these checks in order, and its generated rules in `AGENTS.md`
list them exactly:

- `gofmt`, `go build`, `go vet`, `go mod verify` against `go.sum`, and `go mod tidy -diff`
  so the module files match the code.
- golangci-lint: the standard linters plus depguard, gosec on non-test code, and revive's
  file-length limit.
- `tools/bin/architecture-check`, the self-test above.
- `go tool govulncheck -test ./...`, covering the dependencies in use, test-only ones and
  the standard library.
- The tests with `RAPID_CHECKS=1`.
- `tools/bin/domain-coverage`, which requires each domain package to reach 80% statement
  coverage from its own tests.
- `itos-cc mutation sample --count 10`, a sample of recorded mutation results. It is not
  proof of the changed code.

Changed-code proof is judged in CI. Before pushing code, run
`tools/bin/pinned itos-cc mutation run --since <base> --fail-uncovered --all-tests
--no-annotate internal`, where the base is the push's range start, the nearest ancestor with
a green CI run (`itos ci range --head HEAD` prints it, given GITHUB_REPOSITORY and
GITHUB_TOKEN), and commit `.metrics/mutate/`. When a push's range touches
`internal/**/*.go`, `itos ci run` runs `itos.yaml`'s `proof.code` check from that base as a
static step, before the tests, and fails while a mutant survives, is uncovered or has no
result. `itos work done` runs no proof; it needs that green run.

Tools are pinned:

- The Go toolchain is pinned in `go.mod`, and govulncheck is a `tool` dependency there.
- `tools/bin/pinned` fetches golangci-lint, itos-cc and gitleaks at fixed versions and
  checks each against its release's own SHA-256.
- Workflow actions are pinned by digest.

`docs/upgrading.md` says how to move any of these pins.

## Setup

A new project has no `itos.yaml` or `tasks` yet but has `itos-policy.yaml`: its commit rules,
CI plan and code-proof gate. `itos-template new` prints the setup steps that adopt it.
`itos init --policy itos-policy.yaml` gives the project its own `itos.yaml` and a fresh ledger,
then the policy file is removed. The last step, `itos init --agent-rules`, rewrites
`AGENTS.md`'s rules to describe the project's own config. From then on the workflow runs the Go
steps and the Linux, macOS and Windows tests on the project's own `main`.

The project starts with the mutation records in `.metrics/mutate/` and the reviewed exceptions
in `itos-cc.yaml` that came with its code. Their hashes were taken before its name, owner and
module were filled in, so they are a starting point, not proof: re-prove each function you
change with `itos-cc mutation run`, as the code proof asks of any change, and re-apply an
exception's reviewed reason with `itos-cc mutation except` where the same equivalent mutant
survives.
