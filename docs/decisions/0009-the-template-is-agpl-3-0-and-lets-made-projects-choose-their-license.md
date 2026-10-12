---
status: accepted
date: 2026-10-12
---

# The template is AGPL-3.0 and lets made projects choose their license

## Context and Problem Statement

The itos family (itos, itos-template, itos-cc) is AGPL-3.0, this template MIT. Should the template move to the AGPL, and what license does a project made from it get? Recommended: the template is AGPL-3.0 with an additional permission under section 7 letting a project made from it, through itos-template or by copying its files, use and distribute that project, the template's parts in it included, under terms of its choice; a copy distributed as a template stays AGPL. A render starts MIT with its owner, the MIT text shipped as LICENSE.project and moved to LICENSE by a go stack setup step, the template's AGPL LICENSE in template_only. Or: the exception covering any use of the files; made projects starting AGPL or with no license; waiting for an itos-template way to render a file at another path; or COPYING holding the AGPL beside an MIT LICENSE.

Asked as q-35.

## Considered Options

- AGPL-3.0 with a section 7 permission for made projects only
- AGPL-3.0 with a permission for any use of its files
- Made projects start AGPL, or with no license
- Wait for itos-template to render a file at another path
- COPYING holds the AGPL beside an MIT LICENSE

## Decision Outcome

The template becomes AGPL-3.0 like the rest of the itos family, with a section 7 additional permission for made projects only: a project made from the template, through itos-template or by copying its files, may use, modify and distribute that project, the template's parts in it included, under terms of its choice; distributing the template itself, or a modified copy, as a template stays under the AGPL. A render starts MIT with its owner: the MIT text ships as LICENSE.project and the go stack's setup moves it to LICENSE, the template's AGPL LICENSE joining template_only. itos-template is asked for a way to render a file at another path, replacing the setup step once it ships.

### Consequences

LICENSE holds the AGPL text and joins template_only with the permission's file; renders get LICENSE.project, the MIT text with the owner, and the go stack's setup moves it to LICENSE until itos-template#8 lets the manifest render it there. Releases from v0.2.0 back stay MIT. The docs naming the owner's MIT permission for harvested work name the AGPL and the permission instead. This record is about the template, so it joins template_only.
