# Map

Where to start in each part of the code. Each entry names the files to read first, the fake
that stands in for its port in tests and the command that tests it. Add an entry when a slice
or domain is added; keep each to a few lines and say where to look, not what the code does.

## Counting

- Start: `internal/counting/count.go`, the domain that counts lines and words, and
  `internal/counting/port/files.go`, the `Files` port it reads through.
- Adapter: `internal/disk/files.go`, the read-only disk implementation of the port.
- Fake: `internal/counting/port/porttest/files.go`, the in-memory `Files` the domain's tests use.
- Test: `go test ./internal/counting/... ./internal/disk/...`
