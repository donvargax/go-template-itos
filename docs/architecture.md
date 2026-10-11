# Architecture

The code is a shared domain behind ports, adapters apart from it, and thin slices named for
their commands (decision 6 in `docs/decisions/`). `docs/go.md` lists the packages and the
import rules that hold them.

The private counting domain, `internal/counting`, owns line and word semantics and its result.
It reads through its `Files` port, whose sealed failures describe the file input it needs
without exposing an OS's own errors. The `internal/disk` adapter implements real read-only
access; an in-memory fake holds the domain's test inputs and failures. Unit tests assert
outcomes, and the rapid property checks addition for newline-terminated pieces.

The domain knows no command-line flags, JSON or exit codes. Those belong to a slice and the UI,
which a command line adds; a second command reuses the domain through its port, never another
command's package. A context reaches the port from the caller, carrying cancellation only
(decision 7). This stack has no public executable or scenario runner.

Line separators are LF, with CRLF counted once; a nonempty final unterminated line counts.
Words follow the standard-library Unicode-whitespace rule. No encoding detection or linguistic
word segmentation is performed. File data is not executed or changed.
