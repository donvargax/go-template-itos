# Go design and quality

This guide adds what is specific to Go to `docs/engineering.md`, which it does not repeat.
It describes what this project actually holds and checks. The policy is harvested from
[itos-template](https://github.com/donvargax/itos-template) (its decisions on thin layers
over a domain, unit tests with fakes, rapid properties and the OS as a value), adapted to
this module with its owner's permission (`CONTRIBUTORS.md`).

## Packages

The code follows decision 6 (`docs/decisions/`): a shared domain behind ports, adapters apart
from it, and thin slices named for their commands. This stack has the domain and its adapter;
the command-line feature adds the slices and the UI.

- `internal/counting` is the domain. It owns counting and its result, and reads input only
  through its port. A function that calls the port takes a `context.Context` first (decision
  7); pure logic, such as counting bytes, takes none.
- `internal/counting/port` declares `Files` and its sealed set of failures, `Missing` and
  `Unreadable`. `Files.Read` takes the context first.
- `internal/counting/port/porttest` holds our in-memory fake of `Files` for the domain's tests.
- `internal/disk` is the adapter. It reads real files and translates OS errors into the port's
  failures, importing nothing of ours but the port.

The domain is shared, not part of any command, so a second command reuses counting through its
port without reaching into another command's package. A context carries cancellation and
deadlines only, always as the first parameter and never in a struct or as a value.

`.golangci.yml` holds these arrows with depguard:

- `domain`: the domain, its port and their tests import nothing of ours but the domain and its
  port subpackages, and never kong.
- `domain-io`: the domain's own code does no I/O, neither `os` nor `net`. Its tests may read
  fixtures.
- `infra`: an adapter imports nothing of ours but the port it implements; only its tests may
  use `porttest`.
- `mocks`: no mocking library anywhere. Fakes live in `porttest`.
- `rapid`: the property library only in the domain's own tests.

### Adding a package

A rule names its packages, so a package it does not name is checked by none of them. Lint
refusing an import after you add one means a rule has not been told; extend the rule as below,
never loosen it.

- A domain package: add its folder as a `**/internal/<name>/**` glob to `domain` and
  `domain-io`, and its import path to `domain`'s `allow`. Give it a `port` subpackage for what
  it reaches outside. `tools/bin/domain-coverage` reads its list from `domain`'s `files`, so the
  new package is held to 80% from then on. If its tests use rapid, add a
  `!**/internal/<name>/*_test.go` exclusion to `rapid`.
- An adapter: add its folder to `infra`'s `files`, and the port it implements to `infra`'s
  `allow`.
- Then add compile-valid fixtures to `tools/testdata/architecture` that the rule must accept and
  refuse. `tools/bin/architecture-check` runs the pinned linter over them, so a rule that stops
  matching fails there instead of passing silently.

Code that differs by OS takes the OS as a value. Logic never branches on `runtime.GOOS`. A
build-tagged file holds only a call that one OS lacks, with no mutation site. The counting
example has no OS-specific logic.

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

The project starts with no mutation records: the template's are taken over its own names and
stay behind. Setup makes the render commit the point judging starts from, so the first push
checks no code. A change to a function needs that function's results: run `itos-cc mutation run`
as the code proof asks of any change, and commit `.metrics/mutate/` with the change.
