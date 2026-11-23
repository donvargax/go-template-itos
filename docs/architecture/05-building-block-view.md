# 5. Building block view

The code is engines and slices (decision 11 in `docs/decisions/`): a shared domain is an
engine, with its port and adapters beneath it, and a command is a slice that uses engines
through their ports. `docs/go.md` lists the packages and the import rules that hold them.

The private counting engine, `internal/engine/counting`, has its domain in `domain`, which owns
line and word semantics and its result. The domain reads through the engine's `Files` port,
whose sealed failures describe the file input it needs without exposing an OS's own errors. The
`disk` adapter implements real read-only access; an in-memory fake holds the domain's test inputs and failures. Unit tests assert
outcomes, and the rapid property checks addition for newline-terminated pieces.

The domain knows no command-line flags, JSON or exit codes. Those belong to a slice and the UI,
which a command line adds; a second command reuses the engine through its port, never another
command's package. A context reaches the port from the caller, carrying cancellation only
(decision 7). This stack has no public executable or scenario runner.

Line separators are LF, with CRLF counted once; a nonempty final unterminated line counts.
Words follow the standard-library Unicode-whitespace rule. No encoding detection or linguistic
word segmentation is performed. File data is not executed or changed.
