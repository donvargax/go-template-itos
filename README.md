# go-template-itos

go-template-itos is a Go module, `github.com/itos-corp/go-template-itos`, maintained by
itos-corp. It holds a private counting example and the rules and checks its code is held to.
It has no command-line executable.

## What it does

The counting domain, `internal/engine/counting/domain`, counts a text's lines and words through a Files port,
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

Read [AGENTS.md](AGENTS.md) for the project's rules. [docs/engineering.md](docs/engineering.md)
holds the shared engineering rules, [docs/go.md](docs/go.md) Go's,
[docs/architecture/05-building-block-view.md](docs/architecture/05-building-block-view.md) the layer boundaries,
[docs/security.md](docs/security.md) the trust boundaries,
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
