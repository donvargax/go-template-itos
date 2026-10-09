# go-template-itos

A real project template for [itos-template](https://github.com/donvargax/itos-template),
maintained by itos-corp. Answers replace real project values; no template language is used.

## Status

This is the root bootstrap, not yet a usable Go template. The Go stack, CLI feature and checks
of their rendered combinations are the next pieces. Do not advertise a working count command,
credential scanner or binary release before those pieces pass their checks.

## Branches

- `main` holds the files every stack shares.
- `stack/go` will add a working Go project.
- `go/cli` will add a command line, with a small count command demonstrating the project's rules.

The planned combinations are `go` and `go + cli`. Root changes merge down into the branches;
published history is not rebased. A render merges its selected branches, then replaces the
manifest's literals and excludes template-only files.

## Working here

Read [AGENTS.md](AGENTS.md) and [PLAN.md](PLAN.md). `itos go` is the coordinator's entry point;
an implementing agent uses `itos guide work`. The repository's own notes are
[docs/ORCHESTRATING.md](docs/ORCHESTRATING.md). Hooks and CI judge the repository's rules.

The root checks configuration, commit rules and documentation caps. Go-specific checks will
arrive with their code; this root has no Go module or binary.

## License and provenance

The template is MIT licensed. `itos-corp` is an illustrative owner literal, replaced with the
made project's owner. [LICENSE](LICENSE) contains the license text.

Infrastructure is harvested from the copyright owner's
[itos](https://github.com/donvargax/itos) and
[itos-template](https://github.com/donvargax/itos-template), with permission to distribute that
owner's copied work under MIT in this template. Those source projects retain their own licenses.
Third-party dependencies and tools retain their licenses and notices.
