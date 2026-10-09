# Inherited engineering practices

The default template should carry the general practices already used by
`donvargax/itos-template`, without its generator code or work history. The person confirmed
the branch split on 2026-10-09. No implementation agent was running when this work began.

## Ownership

| Branch | Owns | Does not own |
| --- | --- | --- |
| `main` | Language- and interface-independent design/testing guidance; common manifest declarations; shared setup policy | Go tools, Go CI, CLI flags or binary releases |
| `stack/go` | Go design guidance, dependency/tool pins, import restrictions, property/coverage/proof tooling and Go CI | CLI behavior, command completion or executable publishing |
| `go/cli` | CLI contract, command/error handling, binary scenarios, completion and binary release wiring | Generator behavior or template-release tags |

Shared changes merge downward: root to stack, stack to feature. A feature's rules never
merge upward merely because all current examples happen to use that feature. Branches
keep their own configuration and work history. Resolve documentation drift with the
branch's actual code in view, not by copying root bootstrap claims into finished branches.

## Work

- T-8 writes the shared guide and its short agent-rule reference on `main`.
- `go-engineering` describes the actual Go policy on `stack/go`.
- `go-architecture` strengthens existing depguard enforcement on `stack/go`, then inherits
  into `go/cli`. Additions needed only for the CLI stay in the feature's configuration.
- `cli-contract` carries the reusable CLI guide and exhaustive failure handling on `go/cli`.
  q-1 selects rejection of repeated non-cumulative flags with usage exit 2, matching the
  source contract. Replace the existing last-negation-wins specification before building
  the change; explicitly cumulative flags remain repeatable. Do not change behavior in a
  documentation or lint commit.
- `made-project-policy` separates reusable rules from template history. Common policy
  belongs on the root; Go and CLI additions belong on their branches. It depends on a
  supported fresh-project initialization mechanism, not an invented workaround around itos.
- `made-project-ci` makes Go CI work on a generated ordinary `main`, with scenario additions
  only on `go/cli`. The root workflow must remain language-neutral.
- `cli-releases` finishes the generic binary-release helpers, pins and caller on `go/cli`.

Each idea is specified and promoted on the branch that owns its implementation. Split
cross-branch integration into separate items where the checks differ. The manifest is a
root-owned file: declaring a new Go check there is distinct from installing its tooling
and CI on the Go branch. Keep only one implementing agent in this checkout.

## What the evidence must show

Template-branch green CI is necessary but does not establish generated-project CI. The
manifest excludes `itos.yaml` and `tasks`, and plain `itos init` creates a generic starter;
neither fact proves that a new project inherits custom quality gates. A project's own
policy must be installed without restoring the template ledger or mutation cache.

For each supported combination, verify the emitted guidance and configuration, fresh
setup, actual checks on an ordinary generated `main`, and platform jobs. For the CLI,
verify binary-release wiring separately from template tagging. Keep exact-head CI links
with the evidence. Mutation sampling is not a complete changed-code proof.

T-7 remains open for the original end-to-end render/readiness proof. Its existing checks
do not certify this larger engineering-inheritance goal. Specify a separate integration
item once setup, CI and release work are settled; do not claim the goal from the existing
README readiness sentence or a documentation-only CI run.
