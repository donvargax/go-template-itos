# go-template-itos

A real project template for [itos-template](https://github.com/donvargax/itos-template),
maintained by itos-corp. Answers replace real project values; no template language is used.

## Status

The template is ready. Every combination its manifest allows, `go` and `go + cli`, renders from
the public branches and passes every check the manifest names, the pinned credential scan
among them, so itos-template check reports 2 of 2 combinations and no warning. A project made
from the public template builds, vets and tests, runs its own `count` command, and carries
none of this repository's own itos data.

Ready means a made project starts out building and passing its checks. What a made project
gets, and what it does not get yet:

- [x] Its CI runs the Go checks on its own `main` branch, on Linux, macOS and Windows: see
  [the sample project's run](https://github.com/donvargax/go-template-itos-sample/actions/runs/38081894211).
- [x] Setup gives it its own commit rules and checks, the ones this template is held to.
- [x] A CLI project publishes its binaries when it makes a release: see
  [the sample project's v0.1.0](https://github.com/donvargax/go-template-itos-sample/releases/tag/v0.1.0).
- [ ] Every guide it ships speaks to the made project, not to the template.
- [ ] Automatic dependency updates are proven running on every branch.
- [ ] Its CI scans its own commits for credentials, not only the template's renders.
- [ ] A CLI release is checked against the last release's scenarios, so the `--json` contract
  cannot break unnoticed.
- [ ] A Go project without the command line gets versioned releases too.
- [ ] The template's own releases are cut by CI rather than by hand.

## Making a project

Make a project from the newest release with
[itos-template](https://github.com/donvargax/itos-template) 0.26.0 or later:

```sh
itos-template new https://github.com/donvargax/go-template-itos.git my-project --stack go --feature cli
```

Leave out `--feature cli` for a Go project without the command line. It asks on a terminal for
the project's name, the owner its attribution names and its Go module path, or takes them as
`--answer name=…`, `--answer owner=…` and `--answer module=…`. `--ref v0.1.0` pins a release.
Run the setup steps it prints in the new folder. The made project's `.itos-template.yaml`
records the release it came from.

Each release tags every branch with one version: `main/v0.1.0`, `stack/go/v0.1.0` and
`go/cli/v0.1.0` make v0.1.0, the first.

## Practices

A made project follows these; the guide linked beside each says how. `go` marks what the Go
stack adds, `cli` what the command-line feature adds.

- Design: vertical slices over a domain that owns its behaviour and reaches the outside only
  through its own ports ([docs/engineering.md](docs/engineering.md)); `go` enforces the import
  rules with depguard and proves each rule with a self-test (`docs/architecture.md`).
- Tests: behaviour specified first, failing steps committed before the code; the project's own
  fakes, never a mocking library; `go` adds property tests, an 80% domain coverage floor,
  changed-code mutation proof and runs on Linux, macOS and Windows (`docs/go.md`); `cli` runs
  scenarios against the built binary.
- Code quality: `go` runs gofmt, vet, golangci-lint with gosec and gochecksumtype, and
  govulncheck (`docs/go.md`).
- Dependencies: a library before new code, locked and verified modules, every tool and
  workflow action pinned and checksum-checked ([docs/upgrading.md](docs/upgrading.md)),
  Renovate configured for every branch.
- Security: no credentials in files, fixtures or logs, and provenance kept for copied code
  ([docs/security.md](docs/security.md)).
- Command line: `cli` treats only the exit code and `--json` as a contract with scripts, with
  rules sourced from clig.dev and the GNU and POSIX standards (`docs/CLI.md`).
- Configuration: rules for the first configuration file a project adds
  ([docs/CONFIG.md](docs/CONFIG.md)).
- Releases: `cli` releases on a push to `main` holding a feat or a fix, with notes from
  git-cliff and attested binaries from GoReleaser.
- Work: itos holds every commit to Conventional Commits and its task, in hooks and CI alike;
  decisions are recorded in `docs/decisions/`; coding agents work one item at a time through
  `itos go` and `itos guide work` ([AGENTS.md](AGENTS.md)).

## Branches

- `main` holds the files every stack shares.
- `stack/go` adds a working Go project: a private counting domain and the gates its code is held to.
- `go/cli` adds a command line, with a small count command demonstrating the project's rules.

The combinations are `go` and `go + cli`. Root changes merge down into the branches;
published history is not rebased. A render merges its selected branches, then replaces the
manifest's literals and excludes template-only files.

## Working here

Read [TEMPLATE.md](TEMPLATE.md) first, then [PLAN.md](PLAN.md) and [AGENTS.md](AGENTS.md).
`itos go` is the coordinator's entry point; an implementing agent uses `itos guide work`. The
repository's own notes are [docs/ORCHESTRATING.md](docs/ORCHESTRATING.md). Hooks and CI judge
the repository's rules. Those three files, the itos data and the Template check workflow stay
in the template: a made project's README.md and AGENTS.md speak to it alone.

The root checks configuration, commit rules and documentation caps. Go-specific checks travel
with their code on `stack/go` and `go/cli`; this root has no Go module or binary.

## License and provenance

The template is MIT licensed. `itos-corp` is an illustrative owner literal, replaced with the
made project's owner. [LICENSE](LICENSE) contains the license text.

Infrastructure is harvested from donvargax's
[itos](https://github.com/donvargax/itos) and
[itos-template](https://github.com/donvargax/itos-template), with donvargax's permission to
distribute that copied work under MIT in this template. Those source projects retain their own licenses.
Third-party dependencies and tools retain their licenses and notices.
