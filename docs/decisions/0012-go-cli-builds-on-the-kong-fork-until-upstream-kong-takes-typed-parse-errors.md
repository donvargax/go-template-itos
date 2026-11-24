---
status: accepted
date: 2026-11-23
---

# go/cli builds on the kong fork until upstream kong takes typed parse errors

## Context and Problem Statement

Which kong does go/cli build on: upstream alecthomas/kong, or the person's fork github.com/donvargax/kong?

Asked as q-38, about T-69.

## Considered Options

- Upstream alecthomas/kong, classifying parse errors by their messages
- The fork github.com/donvargax/kong, required and imported by its own path

## Decision Outcome

The fork, github.com/donvargax/kong at v1.17.1-itos.1 (the person, 2026-11-23). It adds typed parse errors and an exported Suggest rule, so the UI classifies a bad flag or an unknown command by its type rather than by its message, and suggests the command meant (itos-cc's decisions 35 and 62). The fork declares its own module path, so go.mod requires github.com/donvargax/kong and the code imports that path, as itos-cc does. The fork retires when upstream kong takes the change.

### Consequences

The fork's tag is pinned in go.mod and Renovate follows it. A kong upgrade waits for the fork's rebase onto the upstream release. When upstream takes typed parse errors and the exported Suggest, go/cli moves back to the upstream path and this record is superseded.
