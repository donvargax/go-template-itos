# 9. Architecture decisions

One line per decision record in `docs/decisions/`; a superseded record leaves the list.

- [Non-cumulative CLI flags are accepted once](../decisions/0001-non-cumulative-cli-flags-are-accepted-once.md)
- [Go code is a shared domain behind ports, separate adapters and thin slices named for their commands](../decisions/0006-go-code-is-a-shared-domain-behind-ports-separate-adapters-and-thin-slices-named-for-their-commands.md)
- [context.Context carries cancellation from cmd through Handle to the ports, never a value](../decisions/0007-context-context-carries-cancellation-from-cmd-through-handle-to-the-ports-never-a-value.md)
- [An interrupted command exits 130, says so once, and a second interrupt ends it at once](../decisions/0008-an-interrupted-command-exits-130-says-so-once-and-a-second-interrupt-ends-it-at-once.md)
