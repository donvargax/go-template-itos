# go-template-itos

A real Go project template for [itos-template](https://github.com/donvargax/itos-template),
maintained by itos-corp. Answers replace real project values after git merges the selected
branches. No template language is used.

## Status

This branch is the Go foundation. Its private counting example has a domain, a Files port, a
disk adapter and an in-memory fake, with example and property tests. There is no command-line
executable here; `go/cli` adds that public entry point. The manifest declares render checks
for both combinations, including a pinned credential scan. A render excludes the template's
itos configuration and work history. Fresh-project policy and CI still need separate
verification before this branch's gates can be claimed for generated projects.

## Check the Go project

```sh
go build ./...
go test ./...
```

CI runs the Go gates and tests on Linux, macOS and Windows. Domain coverage must stay at least
80 percent. Changed code needs fresh mutation proof, recorded before it is pushed and judged
by CI over each push's range; recorded mutants are sampled in CI. Properties run one random
case per CI push, and more in ordinary local runs.

## The example

Lines are logical text lines separated by LF; CRLF is one separator, and a nonempty final line
without a newline counts. Empty input has zero lines. Words are runs separated by Unicode
whitespace, using Go's standard library; punctuation does not split a word. Thus `hello-world`
is one word. The domain has no JSON, flags or exit codes, and file access is read-only through
its port. Its additive property uses newline-terminated pieces to preserve word boundaries.

Read [docs/engineering.md](docs/engineering.md) for shared engineering rules,
[docs/architecture.md](docs/architecture.md) for the layer boundaries,
[docs/security.md](docs/security.md) for trust boundaries, and
[docs/upgrading.md](docs/upgrading.md) before moving dependency pins.

## Working here

Read [AGENTS.md](AGENTS.md) and [PLAN.md](PLAN.md). `itos go` coordinates; an implementing agent
uses `itos guide work`. Hooks and CI judge the rules configured for this branch. Common changes
merge from root `main` into stacks and then features; published history is not rebased.

## License and provenance

The template is MIT licensed; `itos-corp` is the illustrative attribution owner replaced by a
made project's answer. Infrastructure comes from the copyright owner's
[itos](https://github.com/donvargax/itos) and
[itos-template](https://github.com/donvargax/itos-template), with permission to distribute that
owner's copied work under MIT here. The original projects and third-party dependencies retain
their own licenses and notices. See [LICENSE](LICENSE) and [CONTRIBUTORS.md](CONTRIBUTORS.md).
