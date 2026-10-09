# Working rules

Read `PLAN.md` first. This repository is itself the template: its branches must remain real
projects. The root currently has no Go module or executable. Add each language-specific gate
with the code it checks; never claim a check runs before it exists.

The person's session coordinates with `itos go`. An implementing agent runs `itos guide work`,
takes its named item with `itos work take`, does the work itself and starts no agents. Read work
through `itos work`; never edit owners or statuses by hand. Bootstrap initialization is owned
by the source project's T-25 until this new repository has its own usable registry and history.

## Changes and commits

- Specify behavior before building it. When a public entry point and its scenario runner exist,
  features and fixes follow their scenarios, with red steps committed first. Before that, the
  bootstrap is a checked task, not an excuse to add unspecified behavior.
- Everything else names a task in the ledger. Split commit types and their paths before editing.
- Commit with `itos commit -F <file>`, never `git commit`. Bodies explain what changed and why;
  no body line starts with a word and a colon, or with `with #,`.
- Read `git status --short` before staging, and stage only explicit file paths. Never stage a
  whole folder, `.` or use `git add -A` or `-u`.
- Never bypass hooks, force-push or rewrite published history. Push with `itos push` in the
  background; read its output and watch CI without polling. Do not pipe a commit or push.
- Run slow gates through their configured hooks and CI, not by hand before committing. A valid
  check's failure is fixed in the work, not by weakening the check. Stop on a wrong gate.

## Template boundaries

- Templates are real projects. Never introduce a template language or runtime switches merely
  to select a stack or provider: selection is a branch choice before rendering.
- Look for an existing library before building what one likely solves; the person settles that
  choice in the spec. Preserve thin entry points and domain-owned behavior through ports.
- No credentials belong in answers, fixtures, logs or committed files. No credential scan is
  installed at this root-bootstrap stage; the template-check piece adds and proves it.
- Preserve third-party license notices. MIT permission for the owner's harvested infrastructure
  does not change dependencies' licenses or the original source projects' licenses.

## Finishing

An item closes through `itos work done` only after its actual pushed head is green. Add the
complete code proof with the first Go code, and its mutation results with the code they prove.
Never copy another project's mutation cache. The rules `itos init --agent-rules` generates below
must describe this branch's actual configuration.
