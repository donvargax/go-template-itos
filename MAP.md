# Map

Where to start in each part of the code. Each entry names the files to read first, the fake
that stands in for its port in tests and the command that tests it. Add an entry when a slice
or domain is added; keep each to a few lines and say where to look, not what the code does.

## Counting

- Start: `internal/engine/counting/domain/count.go`, the domain, and `port/files.go`, the `Files`
  port it reads through.
- Adapter: `disk/files.go`, the read-only disk implementation of the port.
- Fake: `port/porttest/files.go`, the in-memory `Files`; `files_contract.go`, the contract the
  fake and the disk adapter both pass.
- Test: `go test ./internal/engine/counting/...`

## Import rules

- Start: `.golangci.yml`, the role of each package by its path, and `tools/architecture`, the
  graph check for the arrows depguard cannot hold. Their fixtures are
  `tools/testdata/architecture/cases.json`, a case per arrow.
- Test: `tools/bin/architecture-check`.

## Count command

- Start: `cmd/go-template-itos/main.go`, which assembles the command and the context an
  interrupt cancels; `internal/slices/count/command/count.go`, the slice that carries the
  command's values to the domain through the port, with `classify.go` and
  `input.go`; and `internal/cli/cli.go`, the generic UI.
- Scenarios: `features/*.feature`, steps in `features/steps.go`; run in process by
  `scenarios_test.go` here, through the binary by `features/features_test.go`.
- Fake: counting's `porttest/files.go`.
- Test: `go test ./cmd/... ./internal/... ./features -count=1`

## Pinned tools

- Start: `tools/bin/pinned`, the launcher every gate runs scanners and linters through.
- Test: `go test ./tools/pinned/...`, which runs it against a local release server.

## Nightly

- Start: `.github/workflows/nightly.yml` and the `tools/bin/nightly-*` and `e2e-coverage` it runs;
  `nightly-report` writes the issue. Fuzz targets are `Fuzz*` tests, as `FuzzNewInput`.
- Test: `go test ./tools/bin/nightly-report/`; the rest runs in the nightly.
