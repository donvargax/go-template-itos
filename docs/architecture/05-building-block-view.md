# 5. Building block view

The code is engines and slices (decision 11 in `docs/decisions/`): a shared domain is an
engine, with its port and adapters beneath it, and a command is a slice that uses engines
through their ports. `docs/go.md` lists the packages and the import rules that hold them.

| Block | Path | What it holds | May import | Held by |
| --- | --- | --- | --- | --- |
| Assembly | `cmd/go-template-itos` | `main` and `run`: the context, the interrupts, the adapters, kong's bindings | everything of ours | `kong`; `architecture-check` lets it in |
| UI | `internal/cli` | `Precedence`, problems, flags, the JSON and plain answers, completion, interrupts | nothing of ours | `ui` |
| Slice | `internal/slices/count/command` | The `count` command line, its `Query` and `Handler`, `Classify` (its failures into problems) and `Input` (the checked argument) | `internal/cli`, engines, its own adapters | `command`, `architecture-check` |
| Slice adapter | `internal/slices/count/stdin` | Standard input for `-`; stops with `Cancelled`; hands any other path to a fallback | the engine's port | `adapter` |
| Domain | `internal/engine/counting/domain` | Counting of lines and words and its result | the standard library packages its rule lists, other engines | `domain` |
| Port | `internal/engine/counting/port` | `Files` and its sealed failures, `Missing`, `Unreadable` and `Cancelled` | the same short list; no `os` or `net` | `port`, `store-ports-io` |
| Fake and suite | `internal/engine/counting/port/porttest` | The in-memory `Files` and `FilesContract` | the standard library, engines | `porttest` |
| Engine adapter | `internal/engine/counting/disk` | Read-only disk access, OS errors turned into the port's failures | its engine's port | `adapter` |
| Features | `features` | The scenarios and the harness of both drivers | the harness's own packages | the scenarios, not lint |

`internal/version` says what the build is, and `internal/release` is the library of the release
tools in `tools/bin`. Neither is a command.

The domain reads through `Files`. A context reaches the port from the caller, first in each
call that uses it and carrying cancellation only (decision 7). The domain knows no flag, JSON or
exit code: `count` carries a command's values to the domain and counts nothing itself, and
`internal/cli` maps each failure to its exit code (decision 13). A second command reuses the
engine through its port, never `count`'s package. An interrupt cancels the context, and the
command exits 130 (decision 8).

Line separators are LF, with CRLF counted once; a nonempty final unterminated line counts.
Words follow the standard-library Unicode-whitespace rule. No encoding detection or linguistic
word segmentation is performed. File data is not executed or changed.
