---
status: accepted
date: 2026-11-23
---

# A made project needs no itos: its CI, hooks and docs call plain commands

## Context and Problem Statement

Should a project made from this template need itos, or should its CI, hooks and docs call plain commands?

Asked as q-36, about agents-no-itos.

## Considered Options

- Made projects adopt the template's itos policy at setup (decision 2)
- Made projects get plain commands and no itos data

## Decision Outcome

A made project needs no itos (the person, 2026-11-23). Its CI calls plain commands directly: go test, golangci-lint, govulncheck and itos-cc's mutation run. It gets no itos registry, ledger, commit-footer hooks, gate commands, itos.yaml or itos-policy.yaml, and its AGENTS.md names no itos command. The template's own development keeps using itos meanwhile, in files template_only keeps out of renders.

### Consequences

The setup steps that run itos init go, and so do itos-policy.yaml and policy-drift. The template's own itos rules live in TEMPLATE.md, its itos checks in a template-only workflow, and this record is template_only, since a made project has nothing to compare it with.

## More Information

Supersedes ADR-0002.
