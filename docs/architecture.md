# Architecture

The private counting domain owns line/word semantics and its result. It reads through a Files
port whose outcomes describe the required file input without exposing OS-specific models. A
disk adapter implements real read-only access; an in-memory fake holds domain test inputs and
failures. Unit tests assert outcomes, and the rapid property checks addition for
newline-terminated pieces.

The domain does not know command-line flags, JSON or exit codes. Those belong to the `go/cli`
feature's future thin handler and UI. This Go foundation has no public executable or godog
command-line runner yet.

Line separators are LF, with CRLF counted once; a nonempty final unterminated line counts.
Words follow the standard-library Unicode-whitespace rule. No encoding detection or linguistic
word segmentation is performed. File data is not executed or changed.

Root policy and documentation are inherited from `main`. Language-specific gates travel with
their code; stack and feature branches remain buildable projects. Rendering and updating belong
to itos-template, not to this example's domain.
