# Working rules

If `TEMPLATE.md` exists, read it first: this repository is then the template projects are made
from, and that file says how to work on it.

This project is a Go module with a private counting example: a domain counting lines and words
through a Files port, a read-only disk adapter and an in-memory fake. It has no command line.

Shared design, testing and dependency rules are in `docs/engineering.md`. They apply to every
stack and interface; stack and feature guides add only what is specific to them.
Go's packages, boundary lists, tests and quality checks are in `docs/go.md`.
Rules for any configuration file a project owns are in `docs/CONFIG.md`.

The person's session coordinates with `itos go`. An implementing agent runs `itos guide work`,
takes its item with `itos work take`, does the work itself and starts no agents. Read the registry
with `itos work`; never edit owners or statuses by hand.

## Changes and commits

- Specify behavior before building it. A public entry point's features and fixes follow
  black-box scenarios, with red steps committed first; until there is one, work is a checked
  task. Never add unspecified behavior.
- Other work names a ledger task. Decide each commit's type and paths before editing.
- Commit with `itos commit -F <file>`, never `git commit`. Bodies explain what changed and why;
  no body line starts with a word and a colon, or with `with #,`.
- Read `git status --short` before staging; never stage a whole folder, `.`, or with `git add -A`
  or `-u`.
- Never bypass hooks, force-push or rewrite published history. Push with `itos push` in the
  background, reading its output and watching CI without polling. Never pipe a commit or push.
- Gates run through their configured hooks and CI, not manually before committing. Fix a valid
  check's failure in the work, not by weakening the check; stop on a wrong gate.

## Code boundaries

- Keep thin entry points over domain-owned behavior through ports. The domain knows no UI,
  filesystem implementation, JSON or exit codes. Use our fakes in its tests, never mocks.
- Look for a library before building what one likely solves; settle the choice in the spec. The
  count example uses Go's standard-library whitespace tokenization, not a custom word parser.
- Code differing by OS takes the OS as a value. Build tags hold only a call unavailable on one
  OS, with no mutation site; never branch on runtime.GOOS inside logic.
- The additive counting property uses newline-terminated pieces. Arbitrary concatenation can
  join words and is not additive.
- No credentials in fixtures, logs or committed files. Go vulnerability checks are not
  credential scans; the pinned gitleaks is, `tools/bin/pinned gitleaks dir . --redact`.
- Preserve third-party notices, `CONTRIBUTORS.md` and the comments naming where harvested code
  came from. Owner permission for copied infrastructure does not relicense dependencies or the
  original source projects.

## Finishing

Close with `itos work done` only after the actual pushed head is green: its `itos ci run` judged the
code proof over the push's range, and work done runs none. Before pushing code, record changed-code
proof with the pinned itos-cc, all-tests, fail-uncovered and no-annotate, from the range's start,
the nearest ancestor with a green CI run (`itos ci range --head HEAD` prints it, given
GITHUB_REPOSITORY and GITHUB_TOKEN); commit `.metrics/mutate/`. Kill each survivor with a test; only
truly equivalent mutants may be excepted with a reason for the person to review. An uncovered mutant
needs a test, never an exception. Never reshape code merely to remove mutation sites or copy another
project's cache. The block below is itos's own, written by `itos init --agent-rules`: it points
agents at `itos go`, and `itos guide commit` prints the rules `itos.yaml` sets. Rewrite it with
that command when itos moves, never by hand.

<!-- itos:begin -->

## itos

Start with `itos go` unless your brief names another guide.

Stage only your own files, by explicit path: untracked files you did not make are someone else's work in progress.

<!-- itos:end -->
