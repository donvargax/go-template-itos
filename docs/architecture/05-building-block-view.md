# 5. Building block view

The code is engines and slices (decision 11 in `docs/decisions/`): a shared domain is an
engine, with its port and adapters beneath it, and a command is a slice that uses engines
through their ports. `docs/go.md` lists the packages and the import rules that hold them.
This stack has one engine and no slice.

| Block | Path | What it holds | May import | Held by |
| --- | --- | --- | --- | --- |
| Domain | `internal/engine/counting/domain` | Counting of lines and words and its result | the standard library packages its rule lists, other engines | `domain`, `domain-tests` |
| Port | `internal/engine/counting/port` | `Files` and its sealed failures, `Missing` and `Unreadable` | the same short list; no `os`, `net` or `syscall` | `port`, `store-ports-io` |
| Fake and suite | `internal/engine/counting/port/porttest` | The in-memory `Files` and `FilesContract` | the standard library, engines | `porttest` |
| Adapter | `internal/engine/counting/disk` | Read-only disk access, OS errors turned into the port's failures | its engine's port, other engines' APIs | `adapter` |

The domain reads through `Files` and takes a `context.Context` first where it calls the port,
carrying cancellation only (decision 7). Pure logic, such as counting bytes, takes none. Unit
tests assert outcomes against the fake, and the rapid property checks that counts add for
newline-terminated pieces.

The domain knows no command-line flags, JSON or exit codes. Those belong to a slice and the
UI, which a command line adds; a second command reuses the engine through its port, never
another command's package. This stack has no public executable or scenario runner.

Line separators are LF, with CRLF counted once; a nonempty final unterminated line counts.
Words follow the standard-library Unicode-whitespace rule. No encoding detection or linguistic
word segmentation is performed. File data is not executed or changed.
