# 9. Architecture decisions

One line per decision record in `docs/decisions/`; a superseded record leaves the list.

- [Non-cumulative CLI flags are accepted once](../decisions/0001-non-cumulative-cli-flags-are-accepted-once.md)
- [context.Context carries cancellation, never a value](../decisions/0007-context-context-carries-cancellation-from-cmd-through-handle-to-the-ports-never-a-value.md)
- [An interrupted command exits 130, says so once, and a second interrupt ends it at once](../decisions/0008-an-interrupted-command-exits-130-says-so-once-and-a-second-interrupt-ends-it-at-once.md)
- [Go code is laid out as engines and slices, held by default-deny import rules](../decisions/0011-go-code-is-laid-out-as-engines-and-slices-held-by-default-deny-import-rules.md)
- [go/cli builds on the kong fork until upstream kong takes typed parse errors](../decisions/0012-go-cli-builds-on-the-kong-fork-until-upstream-kong-takes-typed-parse-errors.md)
- [Commands exit with itos-cc's codes and precedence and report problems with a rule, message, fix and exit](../decisions/0013-commands-exit-with-itos-cc-s-codes-and-precedence-and-report-problems-with-a-rule-message-fix-and-exit.md)
- [Scenarios run in process against the ports' fakes locally and through the built binary in CI](../decisions/0014-scenarios-run-in-process-against-the-ports-fakes-locally-and-through-the-built-binary-in-ci.md)
- [A push judges and a nightly explores, reporting through one Nightly red issue](../decisions/0015-a-push-judges-and-a-nightly-explores-reporting-through-one-nightly-red-issue.md)
- [A release reads Conventional Commit headers and is attested, verified and announced to the Go proxy](../decisions/0016-a-release-reads-conventional-commit-headers-releases-any-other-as-a-patch-and-says-so-and-is-attested-verified-and-announced-to-the-go-proxy.md)
- [A release can publish a container image through ko, off by default](../decisions/0017-a-release-can-publish-a-container-image-through-goreleaser-s-ko-off-by-default-pinned-attested-and-verified.md)
