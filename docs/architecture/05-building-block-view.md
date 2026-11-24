# 5. Building block view

The code is engines and slices (decision 11 in `docs/decisions/`): a shared domain is an
engine, with its port and adapters beneath it, and a command is a slice that uses engines
through their ports. `docs/go.md` lists the packages and the import rules that hold them.

The private counting engine, `internal/engine/counting`, has its domain in `domain`, which owns
line and word semantics and its result. It reads through the engine's `Files` port, whose sealed
failures describe the input it needs without exposing an OS's own errors. The engine's `disk`
adapter and the count slice's `stdin` adapter implement it; an in-memory fake holds the domain's test inputs and failures. Unit tests assert outcomes, and the
rapid property checks addition for newline-terminated pieces.

The domain knows no command-line flags, JSON or exit codes. The `count` slice, `internal/slices/count`,
carries the command's domain values to the domain through the port and counts nothing itself;
`internal/cli`, the UI every slice shares, maps each failure to its exit code and prints the
JSON or the text answer. A second command is a new slice beside `count` that reuses the domain,
never `count`'s package. `cmd/go-template-itos` assembles them, making the context that a Ctrl-C
cancels (decisions 7 and 8). `features/` holds the public behaviour as scenarios, which godog
runs in process against the ports' fakes by default and against the binary the harness builds
where CI opts in (decision 14).

Line separators are LF, with CRLF counted once; a nonempty final unterminated line counts.
Words follow the standard-library Unicode-whitespace rule. No encoding detection or linguistic
word segmentation is performed. File data is not executed or changed.
