# Go design and quality

This guide adds what is specific to Go to `docs/engineering.md`, which it does not repeat.
It describes what this project actually holds and checks: the Go foundation and the command
line built on it. The command line's interface contract belongs in its own guide. The
policy is harvested from [itos-template](https://github.com/donvargax/itos-template) (its
decisions on thin layers over a domain, unit tests with fakes, rapid properties and the
OS as a value), adapted to this module with its owner's permission (`CONTRIBUTORS.md`).

## Packages

The code follows decision 11 (`docs/decisions/`): a shared domain is an engine, and a command
is a slice that uses engines through their ports and never another slice.

- An engine is `internal/engine/<name>/{domain,port,port/porttest,<adapter>}`.
- A slice is `internal/slices/<cmd>/{command,domain,port,<adapter>}`.

The counting example is one engine, `internal/engine/counting`:

- `domain` owns counting and its result, and reads input only through its port. A function that
  calls the port takes a `context.Context` first (decision 7); pure logic, such as counting
  bytes, takes none.
- `port` declares `Files` and its sealed set of failures: `Missing`, `Unreadable` and
  `Cancelled`. `Files.Read` takes the context first.
- `port/porttest` holds our in-memory fake of `Files` for the tests.
- `disk` is the adapter that reads real files and translates OS errors into the port's failures.

The `count` command is one slice, `internal/slices/count`:

- `command` is the command, named for it. Its kong struct and `Run` are the UI; its `Query`,
  carrying domain values, and `Handler.Handle(ctx, q)`, which only calls the engine through its
  port, are the application. `Run` assembling the `Query` from kong's flags is the one copy the
  layout accepts. `Classify` is the slice's classifier: it gives the port's failures their
  problems, so `internal/cli` need not know them. `Input` is the checked file argument.
- `stdin` is the slice's adapter of the engine's port: it reads standard input for `-`, stops
  with `Cancelled` once the context is done, and hands every other path to a fallback it is
  given. It imports nothing of ours but the port.

Around them:

- `internal/cli` is the UI every slice shares. It imports nothing of the project's: it turns
  the problems of the classifiers `cmd` hands it into the exit code, by `Precedence`, and prints
  the plain or JSON answer (decision 13). It builds on the kong fork (decision 12), tells a bad
  flag from an unknown command by the fork's typed parse errors, and stops a command on Ctrl-C
  (decision 8).
- `cmd/go-template-itos` assembles everything: it makes the context, feeds the interrupts to
  `internal/cli`, hands it each slice's classifier, gives `stdin` the disk as its fallback, and
  hands the context and the handler to `Run` through kong's bindings. It is the only place that knows every package.
- `internal/version` says what the build is, and `internal/release` is the library of the
  release tools that `tools/bin/release-version` and `tools/bin/release-notes` run. Neither is
  a command.

The domain is shared, not part of any command, so a second command reuses counting through its
port without reaching into another command's package. A context carries cancellation and
deadlines only, always as the first parameter and never in a struct or as a value.

The rules are taken from itos-container's `.golangci.yml` and its `tools/architecture`, and are
default-deny: Go code under `internal` that no role names is refused.

- `unclassified` refuses a package under `internal` or `pkg` that is not a role below, so a new
  package needs its role before lint accepts it.
- `domain` and `port` import only the standard library packages their exact lists name and the
  engines' and slices' packages; `domain-tests` and `port-tests` add the rest of the standard
  library, and the domain's tests rapid.
- `store-ports-io` and `command-io`: a port and a command do no I/O, neither `os`, `net` nor
  `syscall`. Their tests may read fixtures.
- `command`: a command imports the standard library, kong, `internal/cli`, and engines and
  slices; the graph check keeps it from another slice.
- `adapter`: an adapter imports its owner's port and domain and other engines' APIs, no other
  package of ours. Its tests may use `porttest`.
- `porttest`: the fakes import the standard library and engines only.
- `ui`, `version` and `release`: `internal/cli` imports the standard library and kong,
  `internal/version` the standard library, and `internal/release` itself and the standard
  library.
- `mocks`: no mocking library anywhere. Fakes live in `porttest`.
- `rapid`: the property library only in a domain's own tests.
- `kong`: kong only in `cmd`, `internal/cli` and a slice's `command` package.
- forbidigo, taken from itos-template: `runtime.GOOS` is read only where an edge hands the OS in
  as a value, nothing prints with `fmt.Print*`, `print` or `println` outside the writers the UI
  hands down, `os.Exit`, `syscall.Exit` and `log.Fatal*` are for cmd alone, and nothing panics.
  Tests may ask the OS, exit and panic; `tools/` holds programs that print and exit.
- nolintlint: each `//nolint` names its linter and says why, and one that silences nothing is
  reported.
- gochecksumtype: see below.

`tools/bin/architecture-check` runs the graph check, which the import prefixes of depguard
cannot express: a domain imports no adapter, an adapter imports no other owner's adapter, a
slice imports no other slice, only cmd imports a slice from outside the slices, test support is
never linked into a binary, and a `//nolint` of depguard is refused.

### Adding a package

A package under `internal/engine/<name>/` or `internal/slices/<cmd>/` has the role its folder
names, `domain`, `port`, `port/porttest`, `command` or an adapter, and no other folder is
accepted until the rules name it. Lint or the check refusing an import after you add one means
the package is in the wrong role; move it, never loosen a rule.

- A command: add its slice as `internal/slices/<command>/command`, using engines through their
  ports. Add it to `cmd`'s bindings and give it scenarios.
- A new engine needs no edit to the rules. Its `domain` is measured by `tools/bin/domain-coverage`
  from then on, held to 80%.
- A new adapter is a folder beside `domain` and `port`; it imports its owner's port.
- Then add the cases the rules must accept and refuse to `tools/testdata/architecture/cases.json`:
  `Kind` `lint` names the depguard rules that must refuse an import, `graph` the check's
  violation, and `allow` an arrow both must accept. The check runs the pinned linter and itself
  over each case in a module of its own, so a rule that stops matching fails there instead of
  passing silently.

A new sealed failure set, an interface with an unexported marker method like
`port.ReadFailure`, is annotated `//sumtype:decl` and given its problems, with their codes, in one
switch in its slice's classifier that names every kind. gochecksumtype then refuses the switch when a kind is
added and left out; without the annotation nothing checks it.

Code that differs by OS takes the OS as a value. Logic never branches on `runtime.GOOS`. A
build-tagged file holds only a call that one OS lacks, with no mutation site: the scenarios'
interrupt step is the example, `os.Interrupt` on Unix and `CTRL_BREAK_EVENT` on Windows.

## Secure by design: where each input is checked

Every input from outside is checked where it enters, at the outer layer that meets it, into a
checked value: a type with unexported fields and one constructor that checks it, so holding one
proves it was checked. The layers inside trust it and never check it again. Bad input is refused
with a problem naming it, never repaired; only a lossless normalization is allowed. Everything is
under `internal/`, which no other module can import, so nothing outside reaches an inner layer.

| Edge                  | Input                                   | Checked value                                  | Where it is checked                                                                                  |
| --------------------- | --------------------------------------- | ---------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| Command line          | `count`'s file argument                 | `command.Input` (path, or `-` for stdin)       | `command.NewInput`, called by `Run`; an empty argument or one with a NUL byte is `count.input-invalid` |
| Command line          | `--json` and `--no-json`                | kong's bool, in the command's `CLI`          | `internal/cli`'s `Flags` mappers and hook: a value, or a second giving, is `flags.switch-value` or `flags.repeated` |
| Command line          | an unknown flag or command              | none: refused                                  | kong's parse, told by its typed errors in `internal/cli`'s `ParseProblem`: `flags.unknown`, `command.unknown` |
| Command line          | `completion`'s shell                    | kong's enum                                    | kong's `enum` tag, refused as a usage failure                                                        |
| Environment           | `GO_TEMPLATE_ITOS_JSON`                 | the same bool as `--json`                      | kong's `env` tag, through the same mappers; the variable is not the flag, so it is never given twice |
| Standard input        | the bytes behind `-`                    | `[]byte` from the Files port                   | `stdin.Files.Read`: a read that fails or is cancelled is a typed `port.ReadFailure`; the bytes are only counted, never interpreted |
| Files                 | the bytes of the named file             | `[]byte` from the Files port                   | `disk.Files.Read`: a missing or unreadable file is a typed `port.ReadFailure`                        |
| Build                 | the version and commit stamped by the linker | strings                                   | `internal/version`, which only prints them; none is parsed                                           |

## Tests

The domain's unit tests use the `porttest` fake and assert results and typed failures.
The disk adapter's tests use real temporary files, since translating OS errors is what
the adapter does. The `count` slice's tests use the fake too, standard input being its `-`
key; reading standard input, and stopping when the context is done, is `internal/stdin`'s,
held by its own tests. The scenarios in `features/` run under two drivers (decision 14): in process by default, through
`cmd`'s `run` with buffers and the `porttest` fake, and as the built binary where
`GO_TEMPLATE_ITOS_FEATURES_BINARY_DRIVER=1` is set, which CI sets. A scenario only the binary can
hold, an interrupt or the build's stamp, is tagged `@binary-only`. The binary driver interrupts the
program once it waits on an open standard input to show the 130 exit. A second Ctrl-C
ending the process at once is held by `internal/cli`'s unit tests, which feed it fake
interrupts on every system: `count` has no cleanup for a scenario to see skipped.

Properties use `pgregory.net/rapid` and live in the domain's tests. The counting property
joins newline-terminated pieces and checks that their counts add. Arbitrary concatenation
can join two words, so it is not promised to be additive. Every run sets
`RAPID_CHECKS=1`: one new random case per property, as fast as an example test, with the
seed printed on failure. The nightly runs them at a thousand checks (below).

A push judges and a nightly explores (decision 15). On a push, Linux runs every test and the
architecture check's fixtures, and macOS and Windows run the tests of `internal`, with each
platform's paths, shell and line endings.

`.github/workflows/nightly.yml` runs once a day and by hand. Nothing waits on it.

- `platform`: every test and the architecture check's fixtures on Linux, macOS and Windows, each
  `go test -json` kept for the report.
- `explore`, on Linux, every step run whatever an earlier one said: the race detector;
  `tools/bin/nightly-fuzz`, each `Fuzz` target for a minute; `tools/bin/nightly-properties`, the
  packages whose tests import rapid at a thousand checks under a seed drawn and printed;
  `tools/bin/e2e-coverage`, which holds each package outside the domain, the ports and test
  support to 80% of its statements, the domain read from `architecture-check --domains`; and
  `itos-cc mutation run internal`, the whole tree.
- `report` keeps one open issue per branch, "Nightly red on <branch>", under the label
  `nightly-red`: the first red night opens it, each red night after comments on it, naming each
  red job, its failed steps and its failing tests from `go test -json`
  (`tools/bin/nightly-report`), and a green night closes it.

A failure the nightly finds is fixed by a regression test, or by committing the fuzz input or the
property's fail file it printed, which every push then replays. `FuzzCount` in the counting
domain is the example target: add one wherever hostile input enters.

GitHub schedules a workflow from the default branch only. A project whose default branch has the
code gets the nightly from its `schedule`.

## Quality checks

CI's `static` and `test` jobs run these checks in order:

- `gofmt` over `cmd`, `internal` and `features`, `go build`, `go vet`, `go mod verify`
  against `go.sum`, and `go mod tidy -diff` so the module files match the code.
- golangci-lint: the standard linters plus depguard, forbidigo, nolintlint, gosec on non-test
  code, revive's file-length limit and gochecksumtype, which refuses a type switch over a sealed set that
  leaves a kind out, a default arm not counting.
- `tools/bin/architecture-check`, which inspects the packages against the rules above and runs
  the fixtures that prove the rules refuse what they should.
- `go tool actionlint`, which checks the workflows.
- `tools/bin/govulncheck`, which runs `go tool govulncheck -test ./...` with the toolchain
  `go.mod` names, covering the dependencies in use, test-only ones and the standard library.
- `tools/bin/pinned gitleaks git --redact`, a credential scan of the commits a push or pull
  request brings.
- The unit tests of `cmd` and `internal` with `RAPID_CHECKS=1`, then the release helpers'
  tests, over example repositories.
- `tools/bin/domain-coverage`, which requires each domain package to reach 80% statement
  coverage from its own tests.
- The scenarios in process, with the unit tests, and through the binary,
  `go test ./features -count=1` with `GO_TEMPLATE_ITOS_FEATURES_BINARY_DRIVER=1`, which builds
  the binary Go's test cache cannot track; then `tools/bin/scenario-list`, which lists the live
  scenarios with no itos.
- `tools/bin/pinned itos-cc mutation run --count 10 --since <base> --seed <base> --all-tests
  --fail-uncovered=lines cmd internal`, run fresh on each push. The base is the pull request's
  base or the push's before; a branch's first push has no base and runs the packages with no
  `--since` and no `--fail-uncovered`, seeded by the head, where only a surviving mutant fails
  the step. With a base, a mutant that survives or a changed line no test reaches fails it.
  `--all-tests` is there because `cmd` is reached only through the scenarios. No mutation
  result is committed: `.metrics/` is ignored.

Tools are pinned:

- The Go toolchain is pinned in `go.mod`, and govulncheck is a `tool` dependency there.
- `tools/bin/pinned` fetches golangci-lint, itos-cc, gitleaks, GoReleaser and git-cliff at
  fixed versions and checks each against a pinned SHA-256, copied from its release's own
  checksums or, for git-cliff, taken once the archive matched its release's `.sha512`.
- Workflow actions are pinned by digest.

`docs/upgrading.md` says how to move any of these pins.

## Setup

A new project needs no itos. `itos-template new` prints one setup step, the move of
`LICENSE.project` to `LICENSE`, which gives the project its license. From then on the workflow
runs the Go steps, the scenarios and the Linux, macOS and Windows tests on the project's own
`main`, and the release job after them, as `docs/CLI.md` describes.

The project starts with no mutation results and commits none. A branch's first push mutates the
whole packages, judging survivors only, and each later push the functions it changes; a surviving mutant is a test to write.
`itos-cc.yaml` comes with the template's reviewed exceptions, whose hashes cover a function's
body and so hold in a render.
