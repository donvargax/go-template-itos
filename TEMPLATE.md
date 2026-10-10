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

`stack/go` supplies the Go foundation: a private counting domain behind a Files port, with real
disk access at the boundary and an in-memory fake for domain tests. `go/cli` adds a thin command
entry point over that domain. Each branch's copy of `docs/architecture.md` describes its own code.
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
- The owner's MIT permission for their harvested work does not change the source projects'
  licenses or third-party dependencies' and tools'. `CONTRIBUTORS.md`, license notices and the
  harvest comments ship with every made project; keep them true.
- `itos-corp` is the template's illustrative project owner, replaced in attribution by the answer
  the person supplies when making a project. `CONTRIBUTORS.md` calls it the project owner
  `LICENSE` names, which stays true after the replacement.
- Do not replace real upstream tool and preset identities when substituting project identity.
