# Go design and quality

This guide adds what is specific to Go to `docs/engineering.md`, which it does not repeat.
It describes what this project actually holds and checks. The policy is harvested from
[itos-template](https://github.com/donvargax/itos-template) (its decisions on thin layers
over a domain, unit tests with fakes, rapid properties and the OS as a value), adapted to
this module with its owner's permission (`CONTRIBUTORS.md`).

## Packages

The code follows decision 11 (`docs/decisions/`): a shared domain is an engine, and a command
is a slice that uses engines through their ports and never another slice.

- An engine is `internal/engine/<name>/{domain,port,port/porttest,<adapter>}`.
- A slice is `internal/slices/<cmd>/{command,domain,port,<adapter>}`. This stack has no command
  line, so it has engines only; the command-line feature adds the slices and the UI.

The counting example is one engine, `internal/engine/counting`:

- `domain` owns counting and its result, and reads input only through its port. A function that
  calls the port takes a `context.Context` first (decision 7); pure logic, such as counting
  bytes, takes none.
- `port` declares `Files` and its sealed set of failures, `Missing` and `Unreadable`.
  `Files.Read` takes the context first.
- `port/porttest` holds our in-memory fake of `Files` for the domain's tests.
- `disk` is the adapter. It reads real files and translates OS errors into the port's failures.

The domain is shared, not part of any command, so a second command reuses counting through its
port without reaching into another command's package. A context carries cancellation and
deadlines only, always as the first parameter and never in a struct or as a value.

The rules are taken from itos-container's `.golangci.yml` and its `tools/architecture`, and are
default-deny: Go code under `internal` that no role names is refused.

- `unclassified` refuses a package under `internal` or `pkg` that is not a role below, so a new
  package needs its role before lint accepts it.
- `domain` and `port` import only the standard library packages their exact lists name and the
  engines' packages; `domain-tests` and `port-tests` add the rest of the standard library, and
  the domain's tests rapid.
- `store-ports-io`: a port does no I/O, neither `os`, `net` nor `syscall`. Its tests may read
  fixtures.
- `adapter`: an adapter imports its engine's port and other engines' APIs, no other package of
  ours. Its tests may use `porttest`.
- `porttest`: the fakes import the standard library and engines only.
- `mocks`: no mocking library anywhere. Fakes live in `porttest`.
- `rapid`: the property library only in a domain's own tests.
- forbidigo, taken from itos-template: `runtime.GOOS` is read only where an edge hands the OS in
  as a value, nothing prints with `fmt.Print*`, `print` or `println` outside the writers the UI
  hands down, `os.Exit`, `syscall.Exit` and `log.Fatal*` are for cmd alone, and nothing panics.
  Tests may ask the OS, exit and panic; `tools/` holds programs that print and exit.
- nolintlint: each `//nolint` names its linter and says why, and one that silences nothing is
  reported.

`tools/bin/architecture-check` runs the graph check, which the import prefixes of depguard
cannot express: an engine's domain imports no adapter, an adapter imports no other engine's
adapter, test support is never linked into a binary, only cmd imports a slice from outside the
slices, and a `//nolint` of depguard is refused.

### Adding a package

A package under `internal/engine/<name>/` has the role its folder names, `domain`, `port`,
`port/porttest` or an adapter, and no other folder is accepted until the rules name it. Lint or
the check refusing an import after you add one means the package is in the wrong role; move it,
never loosen a rule.

- A new engine needs no edit to the rules. Its `domain` is measured by `tools/bin/domain-coverage`
  from then on, held to 80%.
- A new adapter is a folder beside `domain` and `port`; it imports its engine's port.
- Then add the cases the rules must accept and refuse to `tools/testdata/architecture/cases.json`:
  `Kind` `lint` names the depguard rules that must refuse an import, `graph` the check's
  violation, and `allow` an arrow both must accept. The check runs the pinned linter and itself
  over each case in a module of its own, so a rule that stops matching fails there instead of
  passing silently.

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

Tests run on Linux, macOS and Windows in CI's platform jobs. The architecture check's
fixtures run there as well, because they check paths and diagnostics that differ by platform.

## Quality checks

CI's `static` and `test` jobs run these checks in order:

- `gofmt`, `go build`, `go vet`, `go mod verify` against `go.sum`, and `go mod tidy -diff`
  so the module files match the code.
- golangci-lint: the standard linters plus depguard, forbidigo, nolintlint, gosec on non-test
  code, and revive's file-length limit.
- `tools/bin/architecture-check`, which inspects the packages against the rules above and runs
  the fixtures that prove the rules refuse what they should.
- `go tool actionlint`, which checks the workflows.
- `tools/bin/govulncheck`, which runs `go tool govulncheck -test ./...` with the toolchain
  `go.mod` names, covering the dependencies in use, test-only ones and the standard library.
- `tools/bin/pinned gitleaks git --redact`, a credential scan of the commits a push or pull
  request brings.
- The tests with `RAPID_CHECKS=1`.
- `tools/bin/domain-coverage`, which requires each domain package to reach 80% statement
  coverage from its own tests.
- `tools/bin/pinned itos-cc mutation run --count 10 --since <base> --seed <base>
  --fail-uncovered=lines internal`, run fresh on each push. The base is the pull request's base
  or the push's before; a branch's first push has no base and runs the packages with no
  `--since` and no `--fail-uncovered`, seeded by the head, where only a surviving mutant fails
  the step. With a base, a mutant that survives or a changed line no test reaches fails it. No mutation result is committed: `.metrics/` is ignored.

Tools are pinned:

- The Go toolchain is pinned in `go.mod`, and govulncheck is a `tool` dependency there.
- `tools/bin/pinned` fetches golangci-lint, itos-cc and gitleaks at fixed versions and
  checks each against its release's own SHA-256.
- Workflow actions are pinned by digest.

`docs/upgrading.md` says how to move any of these pins.

## Setup

A new project needs no itos. `itos-template new` prints one setup step, the move of
`LICENSE.project` to `LICENSE`, which gives the project its license. From then on the workflow
runs the Go steps and the Linux, macOS and Windows tests on the project's own `main`.

The project starts with no mutation results and commits none. A branch's first push mutates the
whole packages, judging survivors only, and each later push the functions it changes; a surviving mutant is a test to write.
