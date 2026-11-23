# 9. Architecture decisions

One line per decision record in `docs/decisions/`; a superseded record leaves the list.

- [Non-cumulative CLI flags are accepted once](../decisions/0001-non-cumulative-cli-flags-are-accepted-once.md)
- [context.Context carries cancellation from cmd through Handle to the ports, never a value](../decisions/0007-context-context-carries-cancellation-from-cmd-through-handle-to-the-ports-never-a-value.md)
- [An interrupted command exits 130, says so once, and a second interrupt ends it at once](../decisions/0008-an-interrupted-command-exits-130-says-so-once-and-a-second-interrupt-ends-it-at-once.md)
- [Go code is laid out as engines and slices, held by default-deny import rules](../decisions/0011-go-code-is-laid-out-as-engines-and-slices-held-by-default-deny-import-rules.md)
