# Working rules

If `TEMPLATE.md` exists, read it first: it holds what only this repository's own development
needs.

This project is a Go module with a private counting example: a domain counting lines and words
through a Files port, a read-only disk adapter and an in-memory fake. It has no command line.

Where to start in each part of the code is in `MAP.md`. Always read it with the architecture
document's sections 1, 2 and 9: `docs/architecture/01-introduction-and-goals.md`,
`docs/architecture/02-architecture-constraints.md` and
`docs/architecture/09-architecture-decisions.md`, which indexes the decision records.
Shared design, testing and dependency rules are in `docs/engineering.md`. They apply to every
stack and interface; stack and feature guides add only what is specific to them.
Go's packages, boundary lists, tests and quality checks are in `docs/go.md`.
Rules for any configuration file a project owns are in `docs/CONFIG.md`.

Do the work yourself and start no other agents.

## Changes and commits

- Specify behavior before building it. When a public entry point and its scenario runner exist,
  features and fixes follow their scenarios, with the failing steps committed first and the
  scenarios kept pending until then. Never add unspecified behavior or bend a scenario to fit
  the code. Until there is an entry point, work is a checked task.
- Write each commit header as a Conventional Commit, `type: subject`, because the release
  reads them. The body says what changed and why.
- Read `git status --short` before staging. Stage only explicit file paths: never a whole
  folder, `.`, `git add -A` or `git add -u`. Untracked files you did not make are someone
  else's work in progress.
- Never bypass hooks, force-push or rewrite published history. Never pipe a commit or push.
- Run slow gates through their configured hooks and CI, not by hand before committing. Fix a
  valid check's failure in the work, never by weakening the check. Stop on a wrong gate.

## Code boundaries

- Keep thin entry points over domain-owned behavior through ports. The domain knows no UI,
  filesystem implementation, JSON or exit codes.
- Test with fakes, never mocks, and give each port a contract suite that holds its fake and
  its adapters to the same behavior.
- Code differing by OS takes the OS as a value. Build tags hold only a call unavailable on one
  OS, with no mutation site; never branch on `runtime.GOOS` inside logic.
- Look for an existing library before building what one likely solves; settle the choice in
  the spec. The count example uses Go's standard-library whitespace tokenization, not a custom
  word parser.
- The additive counting property uses newline-terminated pieces. Arbitrary concatenation can
  join words and is not additive.
- Treat every input as hostile: a file, a flag, an environment variable, a path or a checkout.
- No credentials in fixtures, logs or committed files. Go vulnerability checks are not
  credential scans; the pinned gitleaks is, `tools/bin/pinned gitleaks dir . --redact`.
- Preserve third-party license notices, `CONTRIBUTORS.md` and the comments naming where
  harvested code came from. The project's license does not change its dependencies' licenses.

## Finishing

A change is done when its pushed head is green in CI. Mutation runs in CI on each
push, fresh over the commits it brings, and a surviving mutant or a line no test reaches fails
it. Kill each surviving mutant with a test; only a truly equivalent mutant may be excepted, with a reason for the person to
review. Code a test does not reach needs a test, never an exception. Never reshape code merely
to remove mutation sites, hand-edit generated evidence or copy another project's mutation
cache.
