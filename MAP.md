# Map

Where to start in each part of the code: the files to read first, the fake for its port and the
test command. Add an entry for each slice or domain; say where to look, not what the code does.

## Counting

- Start: `internal/engine/counting/domain/count.go`, the domain, and `port/files.go`, the `Files`
  port it reads through.
- Adapter: `disk/files.go`, the read-only disk implementation.
- Fake: `port/porttest/files.go`, the in-memory `Files`; `files_contract.go`, the contract the
  fake and the disk adapter both pass.
- Test: `go test ./internal/engine/counting/...`

## Import rules

- Start: `.golangci.yml`, each package's role by its path, and `tools/architecture`, the graph
  check for arrows depguard cannot hold; fixtures: `tools/testdata/architecture/cases.json`.
- Test: `tools/bin/architecture-check`.

## Count command

- Start: `cmd/go-template-itos/main.go`, which assembles the command; `internal/slices/count/
  command/count.go`, which carries its values to the domain, with `classify.go` and `input.go`;
  `internal/cli/cli.go`, the generic UI.
- Scenarios: `features/*.feature`, steps in `features/steps.go`; run by `scenarios_test.go` here
  and by `features/features_test.go` through the binary.
- Fake: counting's `porttest/files.go`.
- Test: `go test ./cmd/... ./internal/... ./features -count=1`

## Release

- Start: `release.yml`, `.goreleaser.yaml` (kos: the image, off); `tools/bin/`
  `release-version`, `release-notes`, `notices` (`-check` in CI); library: `internal/release`.
- Test: `go test ./tools/bin/release-version/`, likewise `release-notes`, `notices`

## Pinned tools

- Start: `tools/bin/pinned`, the launcher gates run scanners and linters through.
- Test: `go test ./tools/pinned/...`, against a local release server.

## Nightly

- Start: `.github/workflows/nightly.yml` and the `tools/bin/nightly-*` and `e2e-coverage` it
  runs; `nightly-report` writes the issue. Fuzz targets are `Fuzz*` tests.
- Test: `go test ./tools/bin/nightly-report/`; the rest runs in the nightly.
