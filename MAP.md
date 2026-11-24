# Map

Where to start in each part of the code. Each entry names the files to read first, the fake
that stands in for its port in tests and the command that tests it. Add an entry when a slice
or domain is added; keep each to a few lines and say where to look, not what the code does.

## Counting

- Start: `internal/engine/counting/domain/count.go`, the domain that counts lines and words, and
  `internal/engine/counting/port/files.go`, the `Files` port it reads through.
- Adapter: `internal/engine/counting/disk/files.go`, the read-only disk implementation of the port.
- Fake: `internal/engine/counting/port/porttest/files.go`, the in-memory `Files` the domain's
  tests use, and `files_contract.go`, `FilesContract`, which the fake's and the disk adapter's tests
  both run.
- Test: `go test ./internal/engine/counting/...`

## Import rules

- Start: `.golangci.yml`, the role of each package by its path, and `tools/architecture`, the
  graph check for the arrows depguard cannot hold. Their fixtures are
  `tools/testdata/architecture/cases.json`, a case per arrow.
- Test: `tools/bin/architecture-check`, the repository's inspection, then the fixtures.

## Count command

- Start: `cmd/go-template-itos/main.go`, which assembles the command and the context an
  interrupt cancels; `internal/slices/count/command/count.go`, the slice that carries the
  command's values to the domain through the port, with `classify.go` (failures to problems)
  and `input.go` (the checked argument); and `internal/cli/cli.go`, the generic UI.
- Scenarios: `features/*.feature`, with the shared steps in `features/steps.go`. Two drivers run
  them: `cmd/go-template-itos/scenarios_test.go` in process, the default, and
  `features/features_test.go` through the built binary where
  `GO_TEMPLATE_ITOS_FEATURES_BINARY_DRIVER=1`. `tools/bin/scenario-list` lists them.
- Fake: the counting domain's `internal/engine/counting/port/porttest/files.go`; the command's own tests
  use it too.
- Test: `go test ./cmd/... ./internal/... ./features -count=1`

## Pinned tools

- Start: `tools/bin/pinned`, the launcher every gate runs scanners and linters through; its cache is
  in the git common directory, shared by all worktrees.
- Test: `go test ./tools/pinned/...`, which runs the launcher against a local release server.
