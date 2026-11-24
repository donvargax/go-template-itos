---
status: accepted
date: 2026-11-23
---

# A push judges and a nightly explores, reporting through one Nightly red issue

## Context and Problem Statement

What runs per push and what waits for a nightly, and how does a red night reach a person?

Asked as q-41, about T-74.

## Considered Options

- Every test on every platform per push
- Per push a judging subset, the rest nightly, reported through one issue

## Decision Outcome

A push judges, a nightly explores (the migration plan, approved by the person 2026-11-23, after itos's nightly and its testing tiers). Per push: Linux runs every test and both scenario drivers; macOS and Windows run the unit tests and the in-process scenarios. Nightly, on Linux, macOS and Windows with Windows sharded: every test and the binary driver, then on Linux the race detector, every fuzz target for a bounded time, the property tests with a seed drawn and printed, a coverage floor of 80% for each package outside the domain, and a whole-tree mutation run. A report job that needs every other job keeps one Nightly red issue per branch: it opens it, comments on it while the nights stay red, naming each red job and its failing tests from go test -json, and closes it on a green night. Nothing waits on the nightly.

### Consequences

macOS and Windows pushes are faster and see the binary driver's real processes and signals a day late. A failure the nightly finds is fixed by a regression test or a committed fuzz or property input, which every push then runs. GitHub runs scheduled workflows only from the default branch, so a template whose default branch has no code starts its stack branches' nightlies from there.
