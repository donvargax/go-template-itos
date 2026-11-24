# go-template-itos

go-template-itos is a command line, built from the Go module
`github.com/itos-corp/go-template-itos` and maintained by itos-corp.
`go-template-itos count <file>` counts a file's lines and words, `-` for standard input, and
prints the counts for a person or one JSON object with `--json`.

## What it does

`internal/engine/counting` is the private counting engine beneath the command, with a domain, a
Files port, a disk adapter and an in-memory fake; `internal/slices/count` is the `count`
command's slice, and
`internal/cli` is the UI that turns a typed failure into an exit code and stops a command on
Ctrl-C, exit 130. Lines are logical text lines separated by LF; CRLF is one separator, and a
nonempty final line without a newline counts. Empty input has zero lines. Words are runs
separated by Unicode whitespace, using Go's standard library; punctuation does not split a word.
Thus `hello-world` is one word. The domain has no JSON, flags or exit codes, and file access is
read-only through its port. Its additive property uses newline-terminated pieces to preserve
word boundaries.

The public behaviour lives in [features/](features/README.md) as scenarios, which godog runs
in process against the ports' fakes, and against the binary the harness builds for the platform
it is on where `GO_TEMPLATE_ITOS_FEATURES_BINARY_DRIVER=1` is set. The release machinery a
command line needs, its workflow and its version and notes tools, is here as well.

## Check it

```sh
go build ./...
go test ./...
GO_TEMPLATE_ITOS_FEATURES_BINARY_DRIVER=1 go test ./features -count=1
```

Its other gates are formatting, vet, golangci-lint, the architecture check,
govulncheck, domain coverage of at least 80 percent and a fresh mutation run over the pushed
commits; [docs/go.md](docs/go.md) describes each. `.github/workflows/ci.yml` runs them.

## Hooks

A clone runs the hooks it declares in its own `.git/config`, so a fresh clone has none until it
declares them, and every worktree of the clone runs what the clone declares. Git 2.54 or later
runs hooks declared in config. The project has two, each a few lines of `sh`:

- [`tools/hooks/pre-commit`](tools/hooks/pre-commit) runs gofmt over the staged version of each
  Go file a commit adds, copies, modifies or renames, and names each one gofmt would change. It
  reads the index, not the working tree, so it refuses an unformatted staged file whose working
  copy was fixed later and accepts a formatted one with unstaged edits beside it. It formats
  nothing, stages nothing and skips deleted files; a commit with no staged Go file runs no gofmt.
- [`tools/hooks/pre-push`](tools/hooks/pre-push) refuses a push, with vet's output, when
  `go vet ./...` fails: vet compiles every package and its tests. A plain `git push` pushes
  HEAD, and the hook vets the checked-out tree, which is HEAD's only when the tree is clean; when
  it is not, the hook says so.

Declare both once in each clone, from the top of the working tree:

```sh
git config hook.go-template-itos-pre-commit.event pre-commit
git config hook.go-template-itos-pre-commit.command tools/hooks/pre-commit
git config hook.go-template-itos-pre-push.event pre-push
git config hook.go-template-itos-pre-push.command tools/hooks/pre-push
```

`git hook list pre-commit` and `git hook list pre-push` name what a clone runs.
`tools/bin/pre-commit-check` and `go run ./tools/bin/pre-push-check` prove the hooks: each builds
a scratch repository (the second with a local bare remote), declares the hook as above and makes
real commits or pushes, accepted and refused. CI runs both. Neither proves another clone has
declared the hooks; that clone's `git config` does. golangci-lint, govulncheck and the tests
stay in CI only.

## Working here

Read [AGENTS.md](AGENTS.md) for the project's rules. [docs/CLI.md](docs/CLI.md) holds the
command line's contract and rules, [docs/engineering.md](docs/engineering.md) the shared
engineering rules, [docs/go.md](docs/go.md) Go's, [docs/architecture/05-building-block-view.md](docs/architecture/05-building-block-view.md)
the layer boundaries, [docs/security.md](docs/security.md) the trust boundaries,
[docs/upgrading.md](docs/upgrading.md) how to move dependency pins,
[docs/CONFIG.md](docs/CONFIG.md) the rules for a configuration file the project owns, and
[docs/decisions/](docs/decisions/) its decision records. `itos go` coordinates the
work; an implementing agent uses `itos guide work`.

## License and provenance

The project is MIT licensed; see [LICENSE](LICENSE). Its infrastructure comes from
donvargax's [itos](https://github.com/donvargax/itos) and
[itos-template](https://github.com/donvargax/itos-template), both AGPL, and reaches it under the
template's permission letting a project made from it choose its terms. The original projects and
third-party dependencies retain their own licenses and notices; [CONTRIBUTORS.md](CONTRIBUTORS.md) records the provenance.
