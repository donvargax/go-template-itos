---
status: accepted
date: 2026-11-23
---

# Go code is laid out as engines and slices, held by default-deny import rules

## Context and Problem Statement

How should a made project's Go code be laid out and held: decision 6's shared domain with lax import lists, or itos-container's engines and slices under default-deny rules?

Asked as q-37, about T-68.

## Considered Options

- Keep decision 6's layout and its lax depguard lists
- itos-container's engines and slices with default-deny rules and an architecture check

## Decision Outcome

itos-container's layout under default-deny rules (the migration plan, approved by the person 2026-11-23). A shared domain is an engine, internal/engine/<name>/{domain,port,port/porttest,<adapter>}; a command is a slice, internal/slices/<cmd>/{command,domain,port,<adapter>}, using engines through their ports and never another slice. Unclassified Go code is refused, domain and port import only what their exact allowlists name, port and command do no I/O, and the architecture check proves the rules refuse what they should.

### Consequences

Every new package needs a role before lint accepts it. A second command reuses a domain only through an engine's port. The rules come from itos-container's .golangci.yml and its tools/architecture check.

## More Information

Supersedes ADR-0006.
