# go-template-itos

go-template-itos is a command line, built from the Go module
`github.com/itos-corp/go-template-itos` and maintained by itos-corp.
`go-template-itos count <file>` counts a file's lines and words, `-` for standard input, and
prints the counts for a person or one JSON object with `--json`.

## What it does

`internal/count` is the private counting domain beneath the command, with a Files port, a disk
adapter and an in-memory fake, and `internal/cli` is the UI that turns a typed failure into an
exit code. Lines are logical text lines separated by LF; CRLF is one separator, and a nonempty
final line without a newline counts. Empty input has zero lines. Words are runs separated by
Unicode whitespace, using Go's standard library; punctuation does not split a word. Thus
`hello-world` is one word. The domain has no JSON, flags or exit codes, and file access is
read-only through its port. Its additive property uses newline-terminated pieces to preserve
word boundaries.

The public behaviour lives in [features/](features/README.md) as scenarios, which godog runs
against the binary the harness builds for the platform it is on. The release machinery a
command line needs, its workflow and its version and notes tools, is here as well.

## Check it

```sh
go build ./...
go test ./...
go test ./features -count=1
```

Its other gates are formatting, vet, golangci-lint, an import-boundary self-test,
govulncheck, domain coverage of at least 80 percent and changed-code mutation proof;
[docs/go.md](docs/go.md) describes each. The rules `AGENTS.md` generates from the project's
itos configuration say which of them its hooks and CI run.

## Working here

Read [AGENTS.md](AGENTS.md) for the project's rules. [docs/CLI.md](docs/CLI.md) holds the
command line's contract and rules, [docs/engineering.md](docs/engineering.md) the shared
engineering rules, [docs/go.md](docs/go.md) Go's, [docs/architecture.md](docs/architecture.md)
the layer boundaries, [docs/security.md](docs/security.md) the trust boundaries,
[docs/upgrading.md](docs/upgrading.md) how to move dependency pins,
[docs/CONFIG.md](docs/CONFIG.md) the rules for a configuration file the project owns, and
[docs/decisions/](docs/decisions/) its decision records. `itos go` coordinates the
work; an implementing agent uses `itos guide work`.

## License and provenance

The project is MIT licensed; see [LICENSE](LICENSE). Its infrastructure comes from the
copyright owner's [itos](https://github.com/donvargax/itos) and
[itos-template](https://github.com/donvargax/itos-template), with permission to distribute that
owner's copied work under MIT. The original projects and third-party dependencies retain their
own licenses and notices; [CONTRIBUTORS.md](CONTRIBUTORS.md) records the provenance.
