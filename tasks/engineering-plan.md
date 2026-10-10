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

- T-8 wrote `docs/engineering.md` and its short agent-rule reference on `main`.
  The guide commit is `9a58617892771a6eb9d67b3d04f9ac3c16b1f9cb`; its
  [exact-head CI passed](https://github.com/donvargax/go-template-itos/actions/runs/37980659570).
  The item is closed. The guide reached `stack/go` in merge
  `5b7a38cc26cff141db752c456ea155126e3e1760`, preserving that branch's registry and
  Go configuration. Its [CI passed on all three platforms](https://github.com/donvargax/go-template-itos/actions/runs/37981702016).
  T-11 carried it into `go/cli`.
- T-10 wrote `docs/go.md`, the Go stack's guide, and its agent-rule reference on `stack/go`.
  The guide commit is `a8a8c2093f4ab2b8a3ecb5b51e5a8ede360683a7`; its
  [exact-head CI passed on all three platforms](https://github.com/donvargax/go-template-itos/actions/runs/38010453331).
  The item is closed. T-11 carried it into `go/cli`.
- `go-architecture` strengthens existing depguard enforcement on `stack/go`, then inherits
  into `go/cli`. Additions needed only for the CLI stay in the feature's configuration.
  The stack-owned implementation is T-9, specified in its ledger; reserve that ID across
  branches. T-9 is closed. Its final code head is
  `376037fe8f2cd925c6f4c07b5136a10a90e752ab`; [exact-head CI passed on all three platforms](https://github.com/donvargax/go-template-itos/actions/runs/37983078990).
  The self-test compiles 99 positive/negative fixtures before checking the actual pinned
  linter and configuration. No configured `internal` code-proof path changed. T-11 carried it
  into `go/cli`.
- T-13 added the import-boundary self-test to the manifest's `go` and `cli` checks. itos-template
  check passed both combinations at the merged heads; the main
  [run](https://github.com/donvargax/go-template-itos/actions/runs/38012832848) and the
  [go/cli run](https://github.com/donvargax/go-template-itos/actions/runs/38013111934) passed.
  No CI runs that check yet (`template-check-ci`).
- `cli-contract` writes the reusable CLI guide on `go/cli`, without changing behavior.
- `cli-exhaustive-errors` is a separate checked lint/refactor item on `go/cli`.
- `cli-once-only-flags` implements q-1: reject repeated non-cumulative flags with usage exit
  2, matching the source contract. Specify the replacement of last-negation-wins before
  building the change; explicitly cumulative flags remain repeatable. Do not change this
  behavior in a documentation or lint commit.
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

When merging shared documentation, preserve the child's registry and configuration rather
than copying a parent's owners or statuses. Merge task specifications additively: CI reads
footers from imported commits too, and a whole-ledger snapshot would omit root T-8 and fail
the plan as an unknown task. Register child-owned work through itos after committing the
merge. A copied question remains evidence of the person's answer, not a child work item.

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
