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

Infrastructure is harvested from the copyright owner's
[itos](https://github.com/donvargax/itos) and
[itos-template](https://github.com/donvargax/itos-template), with permission to distribute that
owner's copied work under MIT in this template. Those source projects retain their own licenses.
Third-party dependencies and tools retain their licenses and notices.
