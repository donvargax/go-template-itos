---
status: superseded by ADR-0005
date: 2026-10-09
---

# A new task id is one above the highest on any branch

## Context and Problem Statement

Task ids are per branch: each ledger numbers its own tasks, ledgers merge only downward, and main cannot see T-9 to T-11, so its next id would reuse T-9. Recommended: specify a task on the branch that implements it, with an id one above the highest in any branch's ledger on origin (main, stack/go, go/cli). Or: specify every task on main and merge it down before work starts. Or: give each branch its own id prefix.

Asked as q-4.

## Considered Options

The options are those the question names.

## Decision Outcome

Specify a task on the branch that implements it, with an id one above the highest in any branch's ledger on origin: main, stack/go and go/cli. One coordinator per repository keeps the choice free of races.

### Consequences

None recorded.
