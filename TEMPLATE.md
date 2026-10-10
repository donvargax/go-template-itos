# Working on the template

This repository is the template [itos-template](https://github.com/donvargax/itos-template)
makes Go projects from. This file holds what only the template needs. It, `PLAN.md` and
`docs/ORCHESTRATING.md` are listed in `itos-template.yaml`'s `template_only`, with `itos.yaml`,
`tasks` and the Template check workflow, so no made project has them. `README.md` and
`AGENTS.md` on the rendering branches speak to a made project; keep template-facing text here.

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

## Template boundaries

- Templates are real projects. Never introduce a template language or runtime switches merely
  to select a stack or provider: selection is a branch choice before rendering.
- The manifest's literals are real values the code uses, so write them as the code does; a
  render replaces them with the answers.
- A file no made project needs joins `template_only`, with a comment saying why. A file a made
  project keeps makes no claim that it is a template and names no branch it does not have.
- A decision record that binds made projects ships; one about running the template joins
  `template_only`, as does the decisions index, which a made project's itos writes anew.
- Every render is scanned for credentials by gitleaks, pinned in `tools/bin/pinned` and marked
  as the manifest's root check, so `itos-template check` proves it on every combination.

## Provenance and harvest

- Infrastructure is harvested from the owner's itos and itos-template projects: take proven
  machinery whole, with a comment naming its source, rather than rebuilding it differently.
- The owner's MIT permission for their harvested work does not change the source projects'
  licenses or third-party dependencies' and tools'. `CONTRIBUTORS.md`, license notices and the
  harvest comments ship with every made project; keep them true.
- Do not replace real upstream tool and preset identities when substituting project identity.

## The Go stack

`stack/go` adds the Go module: the private counting example under `internal/count`, its tests,
the Go gates and the tools they pin. It has no command line; `go/cli` adds one. The manifest's
`go` checks are the plan this branch's `itos.yaml` runs, less its itos commands. `README.md` and
`AGENTS.md` here describe the project a `go` render makes, with the manifest's literals written
as the code has them.

A render carries none of the template's itos data. It ships `itos-policy.yaml`, which the
manifest's setup steps make the project's own `itos.yaml` through `itos init --policy` and then
remove: this branch's `itos.yaml` less the differences `tools/bin/policy-drift` lists as the
template's alone, each with its reason. Change the two files together; the drift check in the
plan fails when they part. Once set up, a made project's CI runs these gates on its own main: the
ci job sets up Go for the plan and the platform jobs test on Linux, macOS and Windows, with no
condition on branch names, which only the template has.

## The CLI feature

`go/cli` adds the command line on top of the stack: `cmd/go-template-itos`, `internal/cli`, the
scenarios in `features/` and their godog runner, and the release workflow and its tools. The
manifest's `cli` checks are the plan this branch's `itos.yaml` runs, less its itos commands.
`README.md` and `AGENTS.md` here describe the project a `go + cli` render makes. The release job
runs only on a push to main, which no template branch is. On a made project's main it runs after
the ci job and the platform jobs pass, as the stack's gates do there, and cuts a release when the
commits since the last tag call for one. A made project has published one: go-template-itos-sample
cut v0.1.0 (<https://github.com/donvargax/go-template-itos-sample/releases/tag/v0.1.0>) from its
version --json feat, with its five archives, `checksums.txt` and the archives' attestations, in
run <https://github.com/donvargax/go-template-itos-sample/actions/runs/38088380189> (T-33).
