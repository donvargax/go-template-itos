# go-template-itos

go-template-itos is a Go module, `github.com/itos-corp/go-template-itos`, maintained by
itos-corp. It holds a private counting example and the rules and checks its code is held to.
It has no command-line executable.

## What it does

The counting domain, `internal/count`, counts a text's lines and words through a Files port,
with a read-only disk adapter and an in-memory fake that its tests use. Lines are logical text
lines separated by LF; CRLF is one separator, and a nonempty final line without a newline
counts. Empty input has zero lines. Words are runs separated by Unicode whitespace, using Go's
standard library; punctuation does not split a word. Thus `hello-world` is one word. The domain
has no JSON, flags or exit codes, and file access is read-only through its port. Its additive
property uses newline-terminated pieces to preserve word boundaries.

## Check it

```sh
go build ./...
go test ./...
```

Its other gates are formatting, vet, golangci-lint, an import-boundary self-test,
govulncheck, domain coverage of at least 80 percent and changed-code mutation proof;
[docs/go.md](docs/go.md) describes each. The rules `AGENTS.md` generates from the project's
itos configuration say which of them its hooks and CI run.

## Working here

Read [AGENTS.md](AGENTS.md) for the project's rules. [docs/engineering.md](docs/engineering.md)
holds the shared engineering rules, [docs/go.md](docs/go.md) Go's,
[docs/architecture.md](docs/architecture.md) the layer boundaries,
[docs/security.md](docs/security.md) the trust boundaries,
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
