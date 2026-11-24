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
- [x] It needs no itos: its CI, hooks and release call plain commands, and it keeps no itos
  registry, commit footers or itos files.
- [x] A nightly on Linux, macOS and Windows runs the slow checks and keeps one issue open
  while it is red.
- [x] Its releases verify their attestation before publishing, ship their licence notices and
  can publish a container image.
- [x] A CLI project publishes its binaries when it makes a release: see
  [the sample project's v0.1.0](https://github.com/donvargax/go-template-itos-sample/releases/tag/v0.1.0).
- [x] Every guide it ships speaks to the made project, not to the template.
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
`--answer name=…`, `--answer owner=…` and `--answer module=…`. `--ref v0.3.0` pins a release.
Run the setup steps it prints in the new folder. The made project's `.itos-template.yaml`
records the release it came from.

Each release tags every branch with one version: `main/v0.3.0`, `stack/go/v0.3.0` and
`go/cli/v0.3.0` make v0.3.0, the newest. It makes the template AGPL-3.0 with a permission
letting made projects choose their license, while a made project still starts MIT (the template's decision 9),
and moves to itos 8.3.0. v0.2.0 laid the Go code out as a shared domain, separate adapters and a
slice per command (the template's decisions 6 to 8), so a project made from it adds a command without
restructuring first; v0.1.0 came before that layout.

## Practices

A made project follows these; the guide linked beside each says how. `go` marks what the Go
stack adds, `cli` what the command-line feature adds.

- Design: vertical slices over a domain that owns its behaviour and reaches the outside only
  through its own ports ([docs/engineering.md](docs/engineering.md)); `go` lays the code out as
  engines and slices under default-deny import rules, with a context carrying cancellation to
  the ports, each rule proved by a self-test (`docs/go.md`, `docs/architecture/`).
- Tests: behaviour specified first, failing steps committed before the code; the project's own
  fakes, never a mocking library, each held to its adapters by one contract suite; `go` runs
  on Linux, macOS and Windows, and its nightly adds the race detector, fuzzing, property tests,
  a coverage floor and mutation proof (`docs/go.md`); `cli` runs one set of scenarios in
  process against the fakes, and against the built binary in CI.
- Code quality: `go` runs gofmt, vet, golangci-lint with gosec and gochecksumtype, and
  govulncheck (`docs/go.md`).
- Dependencies: a library before new code, locked and verified modules, every tool and
  workflow action pinned and checksum-checked ([docs/upgrading.md](docs/upgrading.md)),
  Renovate configured for every branch.
- Security: no credentials in files, fixtures or logs, and provenance kept for copied code
  ([docs/security.md](docs/security.md)).
- Command line: `cli` treats only the exit code and `--json` as a contract with scripts, exits
  with codes 0, 1, 75, 80, 81, 82, 100 and 130, reports each problem with a rule, a message,
  a fix and its exit, and sources its rules from clig.dev and the GNU and POSIX standards
  (`docs/CLI.md`).
- Configuration: rules for the first configuration file a project adds
  ([docs/CONFIG.md](docs/CONFIG.md)).
- Releases: `cli` releases on a push to `main` holding a feat or a fix, with notes from
  git-cliff and binaries from GoReleaser. Before 1.0 a breaking change makes a minor release,
  and a commit whose header is not a Conventional Commit counts as a patch and is named in the
  run. The attestation is verified before publishing, the Go proxy is asked for the new
  version, `THIRD_PARTY_NOTICES.md` ships in every archive, and a container image can be
  switched on.
- Work: decisions are MADR records in `docs/decisions/`; commit headers are Conventional
  Commits, which the release reads, and nothing checks them. Coding agents follow
  [AGENTS.md](AGENTS.md); itos is not part of a made project.

## Branches

- `main` holds the files every stack shares.
- `stack/go` adds a working Go project: a private counting domain and the gates its code is held to.
- `go/cli` adds a command line, with a small count command demonstrating the project's rules.

The combinations are `go` and `go + cli`. Root changes merge down into the branches;
published history is not rebased. A render merges its selected branches, then replaces the
manifest's literals and excludes template-only files.

## Working here

This repository's own development keeps using itos; a made project does not. Read
[TEMPLATE.md](TEMPLATE.md) first, then [PLAN.md](PLAN.md) and [AGENTS.md](AGENTS.md). The
repository's own notes are [docs/ORCHESTRATING.md](docs/ORCHESTRATING.md). Hooks and CI judge
the repository's rules. Those three files, the itos data and the Template check workflow stay
in the template: a made project's README.md and AGENTS.md speak to it alone.

The root checks configuration, commit rules and documentation caps. Go-specific checks travel
with their code on `stack/go` and `go/cli`; this root has no Go module or binary.

## License and provenance

The template is licensed under the [GNU AGPL version 3](LICENSE), as the rest of the itos family
is, with an [additional permission](LICENSE-EXCEPTION): a project made from it may use, modify
and distribute that project, the template's parts included, under terms of its choice.
Distributing the template itself, or a modified copy, as a template stays under the AGPL.
Releases up to v0.2.0 were MIT licensed and stay so.

A made project starts MIT licensed: the go stack's setup moves `LICENSE.project` to `LICENSE`,
with `itos-corp`, an illustrative owner literal, replaced by the project's owner (the template's decision 9).

Infrastructure is harvested from donvargax's
[itos](https://github.com/donvargax/itos) and
[itos-template](https://github.com/donvargax/itos-template), both AGPL; as their copyright owner,
donvargax grants the same permission for the copied work. Those source projects retain their own
licenses. Third-party dependencies and tools retain their licenses and notices.
