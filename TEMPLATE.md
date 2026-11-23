# Working on the template

This repository is the template [itos-template](https://github.com/donvargax/itos-template)
makes Go projects from. This file holds what only the template needs. It, `PLAN.md` and
`docs/ORCHESTRATING.md` are listed in `itos-template.yaml`'s `template_only`, with `itos.yaml`,
`tasks` and the Template check workflow, so no made project has them. `README.md` and
`AGENTS.md` on the rendering branches speak to a made project; keep template-facing text here.

## Working with itos here

`AGENTS.md` is what a made project keeps, so it names no itos command. The template's own
development uses itos, and these rules are the ones it left out.

- The person's session coordinates with `itos go`. An implementing agent runs `itos guide work`,
  takes its item with `itos work take`, does the work itself and starts no agents. Read work
  through `itos work`; never edit owners or statuses by hand.
- Every change not driven by scenarios names a task in the ledger. Split commit types and their
  paths before editing: a docs commit cannot touch `itos.yaml`, `itos-template.yaml` or
  `tools/bin`, so that part goes in a chore commit naming the task.
- Commit with `itos commit -F <file>`, never `git commit`. No body line starts with a word and
  a colon, or with `with #,`.
- Push with `itos push` in the background, read its output and watch CI without polling.
- Close an item with `itos work done` only after its actual pushed head is green. Its
  `itos ci run` judged the pushed range's checks, and `work done` runs none. Mutation runs fresh in
  CI on each push (`ci.yml`'s test job), so no mutation result is committed: `.metrics/` is ignored.
- The block `itos init --agent-rules` writes into `AGENTS.md` points at `itos go` and says to
  stage only your own files. A made project's `AGENTS.md` carries no such block, and no setup
  step writes it back (decision 10).
- The template's own itos checks run in `.github/workflows/itos.yml`, which `template_only`
  keeps out of renders; `ci.yml` is what a made project keeps and calls plain commands. Each
  branch's `itos.yaml` `ci.steps` hold only what itos or the template needs, never a step
  `ci.yml` runs.

The template's AGPL and its permission for made projects are under Provenance and harvest below.

## Reading order

Read `PLAN.md` first: the template's shape and its order of work. Then `AGENTS.md`, the rules
every branch and every made project keeps, and the guides it names. A coordinator's session also
reads `docs/ORCHESTRATING.md`, the notes `itos go` prints after its guide.

## Branches

`main` is the root and holds the files every stack shares; it never renders alone, so its
`README.md` stays the template's front page. `stack/go` branches off it and adds a working Go
project; `go/cli` branches off the stack and adds a command line. The combinations are `go` and
`go + cli`. A render merges the selected branches, replaces the manifest's literals with the
answers and leaves out the `template_only` paths.

Root changes merge down, `main` into `stack/go`, then `stack/go` into `go/cli`, with normal
merge commits; published history is never rebased or forced. A merge keeps the receiving
branch's own `itos.yaml`, CI workflow, Renovate configuration and work registry, takes the
questions and decision records whole, and merges the ledger as whole entries in id order. Each
branch adds to this file a section of its own.

`stack/go` supplies the Go foundation: a private counting domain behind a Files port, with real
disk access at the boundary and an in-memory fake for domain tests. `go/cli` adds a thin command
entry point over that domain. Each branch's `docs/architecture/05-building-block-view.md` describes its own code.
Root policy and documentation are inherited from `main`; language-specific gates travel with
their code, and stack and feature branches remain buildable projects.

Renovate reads its configuration from main, whose `.github/renovate.json5` names main, stack/go
and go/cli as base branches and merges each branch's own file over it, so each branch's requests
follow that branch's managers and are judged by its own CI. The Renovate app is installed, but no
request has landed yet, so the automation stays unverified until one lands on each branch.

## Template boundaries

- Templates are real projects. Never introduce a template language or runtime switches merely
  to select a stack or provider: selection is a branch choice before rendering.
- The manifest's literals are real values the code uses, so write them as the code does; a
  render replaces them with the answers.
- A file no made project needs joins `template_only`, with a comment saying why. A file a made
  project keeps makes no claim that it is a template and names no branch it does not have. It
  never names the template by name or URL either: a render replaces those literals, so the
  sentence would name the made project. What only the template needs is said here instead.
- A decision record that binds made projects ships; one about running the template joins
  `template_only`, as does the decisions index, which a made project's itos writes anew.
- Every render is scanned for credentials by gitleaks, pinned in `tools/bin/pinned` and marked
  as the manifest's root check, so `itos-template check` proves it on every combination and
  stops warning about a render nothing scans.
- Template checks and made-project checks are different evidence. The template excludes its own
  configuration and work history from renders; a made project adopts the selected branch's
  policy at setup, through the steps `itos-template new` prints. Claim a guarantee for made
  projects only once an emitted project, after setup, has shown it on its own ordinary branch
  and CI.
- Rendering and template updates belong to itos-template, not to this repository or its example
  domain. Template releases and the generator's adoption/update workflow are separate from
  made-project binary release tooling, and no guide claims them.

## Provenance and harvest

- Infrastructure is harvested from the owner's itos and itos-template projects: take proven
  machinery whole, with a comment naming its source, rather than rebuilding it differently.
- The template's AGPL and its permission for made projects (decision 9, `LICENSE-EXCEPTION`)
  do not change the source projects' licenses or third-party dependencies' and tools'. `CONTRIBUTORS.md`, license notices and the
  harvest comments ship with every made project; keep them true.
- `itos-corp` is the template's illustrative project owner, replaced in attribution by the answer
  the person supplies when making a project. `CONTRIBUTORS.md` calls it the project owner the
  MIT license names, `LICENSE.project` here and `LICENSE` once setup moves it, which stays true
  after the replacement.
- Do not replace real upstream tool and preset identities when substituting project identity.

## The Go stack

`stack/go` adds the Go module: the private counting domain `internal/counting` with its port,
the `internal/disk` adapter, their tests, the Go gates and the tools they pin, laid out as
decision 6 says. It has no command line; `go/cli` adds one, with its own interface contract and
checks, the `count` slice and the `internal/stdin` adapter over this domain. The manifest's `go`
checks are the commands this branch's `ci.yml` runs, the import-boundary self-test included;
`itos-template check` runs them after the root's gitleaks scan on each supported render. That is
evidence about the template, not about a made project's own CI, which scans only the commits a
push or pull request brings. `README.md` and the shipped guides here describe the project a `go`
render makes, with the manifest's literals written as the code has them.

A render carries none of the template's itos data and needs no itos (decision 10): its
`ci.yml` calls plain commands, `.github/workflows/itos.yml` and `itos.yaml` stay here, and its
setup is the license move alone. Each job sets the toolchain with `.github/actions/go-toolchain`,
harvested from itos, so `tools/bin/govulncheck` and CI judge the standard library `go.mod` names.
The platform jobs run on pull requests and on pushes to any branch but main, so they run on the
template's own direct pushes and, in a made project, on the pull request before its main push.

## The CLI feature

`go/cli` adds the command line on top of the stack: `cmd/go-template-itos`, `internal/cli`, the
scenarios in `features/` and their godog runner, and the release workflow and its tools. The
manifest's `cli` checks are the commands this branch's `ci.yml` runs: the scenarios and the
import-boundary self-test included, but not the release helpers' tests, and not the smoke-set
check in `itos.yaml`, an itos command a made project has no use for. `itos-template check` runs them after the root's
gitleaks scan on each supported render, which is evidence about the template, not about a made
project's own CI. `README.md`, `AGENTS.md` and the shipped guides here describe the project a
`go + cli` render makes. The release job runs only on a push to main, which no template branch
is, so this branch cuts no CLI binary release. On a made project's main it runs after the static and
test jobs pass and the platform jobs pass or skip, as the stack's gates do there, and cuts a release when the commits
since the last tag call for one. A made project has published one: go-template-itos-sample
cut v0.1.0 (<https://github.com/donvargax/go-template-itos-sample/releases/tag/v0.1.0>) from its
version --json feat, with its five archives, `checksums.txt` and the archives' attestations, in
run <https://github.com/donvargax/go-template-itos-sample/actions/runs/38088380189> (T-33).

## Architecture of the root

The root carries project policy and shared documentation. It has no executable or Go module.
Stack and feature branches add the code and checks they need while retaining these shared files.
A made project's architecture document is `docs/architecture/`, a lean arc42: one file per
section, only sections with content. Sections 1, 2 and 9 are what an agent always reads, so
`tools/bin/doc-caps` holds each to 2048 bytes; section 9 indexes the decision records a made
project ships.

Rendering and template updates belong to itos-template, not to this repository. Templates use
real literals and git branches; their generated projects do not contain a template interpreter.

## Security and upgrading, for the template

itos's hooks and commit rules are workflow controls, not a security sandbox.

The root uses a pinned itos. Move itos's repository pin with `itos pin <version>` and the
workflow action with it, to that release's commit with a `# v<version>` comment and no `version`
input, so the action installs the launcher of the release its ref names. Then apply each step of
the release's Upgrading section that the branch's configuration meets, and run
`itos init --agent-rules`. Do not change only the launcher or action version.
