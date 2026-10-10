# This repository's own notes

`itos go` prints these notes after its coordinator guide. They hold only what is specific to this
template; do not repeat the generic guide here.

## Start of a session

- Read `PLAN.md` and the branch's actual `itos.yaml` before assigning work. This root-bootstrap
  stage has no Go code, scenario runner, mutation proof or binary release job yet.
- The template's root is `main`; stack and feature branches are real projects. Keep history and
  merge root changes downward. Do not treat a made project's binary release as a template tag.
- Infrastructure came from the owner's itos-template and itos repositories. Harvest proven
  machinery, retaining provenance, rather than rebuilding it differently for this template.
- GitHub App installation and bot operation must be verified before claiming Renovate is active;
  a committed configuration alone proves neither.
- Ask, answer and record every question on `main`, even one specific to a stack or feature
  (decision 3). Question and record numbers are per branch; asked anywhere else, main reuses them.
- A new task's id is one above the highest on any branch (decision 4), never `itos task next-id`
  alone: `git fetch` and `git grep -hoE 'id: T-[0-9]+' origin/main origin/stack/go origin/go/cli
  -- 'tasks/phase-*.yaml' | sort -t- -k2n | tail -1`.

## Briefs

- Name the implementing item, its spec, exact harvest files and neighboring checks.
- Attach code-specific gates with their code. The final Go/CLI template carries the full rules;
  an unfinished root must not pretend to exercise nonexistent code.
- The coordinator writes user-facing project docs in separate docs commits. Keep each claim in
  step with the branch's implementation.
- No mutation cache or generator-domain fixture is inherited from another project.

## Lessons

Record a lesson only when it prevents a repeated mistake. Give its recorded date, last-seen date
and exit item (or permanent); keep at most ten. An upstream-tool gap belongs upstream, after
checking versions and existing reports, not in a workaround here.

- After merging a parent branch down, compare the items both registries hold. The merge keeps the
  child's status, so close an item finished upstream with `itos work done` on the child too.
  Recorded 2026-10-09; last seen 2026-10-09; permanent while registries share items.
