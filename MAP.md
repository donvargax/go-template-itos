# Map

Where to start in each part of the code. Each entry names the files to read first, the fake
that stands in for its port in tests and the command that tests it. Add an entry when a slice
or domain is added; keep each to a few lines and say where to look, not what the code does.

## Counting

- Start: `internal/engine/counting/domain/count.go`, the domain that counts lines and words, and
  `internal/engine/counting/port/files.go`, the `Files` port it reads through.
- Adapter: `internal/engine/counting/disk/files.go`, the read-only disk implementation of the port.
- Fake: `internal/engine/counting/port/porttest/files.go`, the in-memory `Files` the domain's
  tests use.
- Test: `go test ./internal/engine/counting/...`

## Import rules

- Start: `.golangci.yml`, which gives each package a role by its path and refuses any Go code
  under `internal` that has none, and `tools/architecture`, the graph check that resolves the
  owner and layer arrows depguard cannot. The fixtures that prove both are
  `tools/testdata/architecture/cases.json`, a case per arrow.
- Test: `tools/bin/architecture-check`, which inspects the repository and then runs the fixtures.

## Count command

- Start: `cmd/go-template-itos/main.go`, which assembles the command and the context an
  interrupt cancels; `internal/count/count.go`, the slice that carries the command's values to
  the domain through the port; and `internal/cli/cli.go`, the shared UI that maps each failure to
  its exit code and prints the answer.
- Scenarios: `features/count.feature` and `features/cli.feature`, run by `features/features_test.go`
  against the built binary.
- Fake: the counting domain's `internal/counting/port/porttest/files.go`; the command's own tests
  use it too.
- Test: `go test ./cmd/... ./internal/... ./features -count=1`

## Pinned tools

- Start: `tools/bin/pinned`, the launcher every gate runs scanners and linters through; its cache is
  in the git common directory, shared by all worktrees.
- Test: `go test ./tools/pinned/...`, which runs the launcher against a local release server.
