---
status: accepted
date: 2026-10-09
---

# Made projects inherit quality policy through upstream itos initialization

## Context and Problem Statement

Automatic quality-policy inheritance is not supported by current itos init: it creates a generic starter, while an existing supplied config is diagnosed rather than getting fresh ledger/registry data. Recommended: pursue a focused upstream initialization capability that accepts reusable policy and creates project-local fresh work data, then keep template automatic setup dependent on it. Alternative: explicitly manual project adoption using a reviewed complete policy and newly authored initial data, without claiming automatic inheritance. Which should the default template target? In either case the generator stays unaware of itos, and policy is split between main, stack/go and go/cli.

Asked as q-2, about made-project-policy.

## Considered Options

The options are those the question names.

## Decision Outcome

Pursue upstream itos support for initializing reusable policy with fresh project-local work data. File the focused consumer request after reproducing the current limitation and checking releases and existing issues. Automatic template policy setup depends on that capability; do not build a workaround or copy template history. Continue independent template guidance, enforcement, CI and release work. Keep the generator unaware of itos and preserve the main/stack/go/go/cli split.

### Consequences

None recorded.
