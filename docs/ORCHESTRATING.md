# This repository's own notes

`itos go` prints these notes after its coordinator guide. Keep them specific to this template.

## Start of a session

- Read `PLAN.md` and this branch's actual `itos.yaml`. The Go foundation has private example
  code and its Go/code-proof gates, but no public CLI, godog runner or binary release job yet.
- Root is `main`; stacks and features remain real projects. Merge common changes downward,
  never rebase published history. Made-project binary releases are not template release tags.
- Infrastructure is harvested from the owner's itos-template and itos repositories. Preserve
  provenance and real upstream identities; do not globally substitute the owner's handle.
- A Renovate configuration does not prove App installation or bot activity. Verify before
  claiming the automation runs.
- Ask, answer and record every question on `main`, even one specific to a stack or feature
  (decision 3). Question and record numbers are per branch; asked anywhere else, main reuses them.
- itos mints each new task and numbered item id from refs/itos/ids, which every branch shares
  (decision 5); `itos task add` and `itos work add --kind` refuse an id given by hand.
- Main's registry is the template-wide index of every idea. When the work lands on a stack or
  feature branch, leave main's copy open with an `itos work edit --note` naming where; never
  drop it, since dropped reads as unwanted.
- A need only this template has (its branch layout, renders, the decisions it ships) is solved
  here or asked of itos-template; itos gets only what any project adopting a policy needs.
  itos#31 was withdrawn for this; #33 and #39 stay, as general adoption features.

## Briefs and proof

- Name the item, spec, exact harvest files and neighboring checks. The coordinator writes the
  user-facing docs; each branch's claims must match its implemented stage.
- The count example uses standard-library Unicode whitespace and logical lines. Its additive
  property uses newline-terminated pieces, not arbitrary concatenation.
- The proof base is the push's range start, the nearest ancestor with a green CI run: `itos ci
  run` judges the proof from there, and `itos work done` runs none. Adding tests can stale
  imported-code snapshots: refresh the records the gate requires, not the gate.
- This branch's proof covers `internal/`; add `cmd/` with the public CLI, not as an empty placeholder.
- Pick a representative smoke scenario when the CLI runner makes each feature live. Keep live
  scenarios as explained wip on a revert; never delete them to force the revert through.
- Never inherit another project's mutation cache or generator-domain fixtures.

## Lessons

New lessons name a recorded date, last-seen date and exit item (or permanent); keep at most ten.
Report an upstream-tool gap upstream after checking releases and existing reports, rather than
working around it in this template.

- After merging a parent branch down, compare the items both registries hold. The merge keeps the
  child's status, so close an item finished upstream with `itos work done` on the child too.
  Recorded 2026-10-09; last seen 2026-10-09; permanent while registries share items.
- A task's checks run in CI on every push that names it, so a check that passes only once the
  work is whole turns an early push red. Push all of the task's commits together. Recorded
  2026-10-09; last seen 2026-10-09 (T-11); permanent.
- An agent never answers a question, not even one a decision command needs: it stops and the
  coordinator asks the person. Recorded 2026-10-10; last seen 2026-10-10 (T-24's q-11, which the
  person then confirmed); permanent.
