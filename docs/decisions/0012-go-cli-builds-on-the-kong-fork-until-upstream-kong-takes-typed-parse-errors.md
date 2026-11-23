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
- The fork github.com/donvargax/kong through a go.mod replace

## Decision Outcome

The fork, github.com/donvargax/kong at v1.17.1-itos.1 (the person, 2026-11-23). It adds typed parse errors and an exported Suggest rule, so the UI classifies a bad flag or an unknown command by its type rather than by its message, and suggests the command meant (itos-cc's decisions 35 and 62). go.mod requires alecthomas/kong and replaces it with the fork, so imports keep the upstream path. The replace retires when upstream kong takes the change.

### Consequences

The fork's tag is pinned in go.mod and Renovate follows it. A kong upgrade waits for the fork's rebase onto the upstream release. When upstream takes typed parse errors and the exported Suggest, the replace goes and this record is superseded.
