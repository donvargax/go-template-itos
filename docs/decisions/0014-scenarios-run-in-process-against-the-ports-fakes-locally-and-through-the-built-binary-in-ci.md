---
status: accepted
date: 2026-11-23
---

# Scenarios run in process against the ports' fakes locally and through the built binary in CI

## Context and Problem Statement

How do go/cli's scenarios run: only through the built binary, or in process against the ports' fakes locally with the binary in CI?

Asked as q-40, about T-71.

## Considered Options

- The built binary only
- In process against the fakes by default, the binary behind an opt-in variable CI sets

## Decision Outcome

In process locally, the binary in CI (the migration plan, approved by the person 2026-11-23, after itos-cc's decisions 44, 52 and 53). One set of feature files runs under two drivers. The in-process driver calls cmd's run with buffers for the UI and each port's porttest fake, so go test ./... builds no binary and starts no process; it is the default. The binary driver builds the program and runs it with the real adapters, proving main's wiring, real exit codes and signals; it runs only where an opt-in variable is set, which CI sets. A scenario only one driver can hold carries @binary-only, never a copy. A port's fake may stand in for its adapter only because both pass the port's one contract suite in porttest.

### Consequences

Local runs cost milliseconds per scenario and spare the machine; mutation kills slice and domain mutants without a rebuild. Every port needs a contract suite its fake and its adapters pass, and a new port is bound in both drivers, a test failing until it is. What only @binary-only scenarios reach is proven in CI alone.
