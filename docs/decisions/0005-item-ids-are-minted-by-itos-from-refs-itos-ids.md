---
status: accepted
date: 2026-10-09
---

# Item ids are minted by itos from refs/itos/ids

## Context and Problem Statement

Decision 4 numbers a new task one above the highest id on any branch, found by hand with git grep, because rc.1's task next-id reads only the local tree. itos 7.0.0-rc.4 mints every task and numbered item id from the counter refs/itos/ids on the remote, shared by all branches (itos decision 45), one past the higher of the counter and the highest id the branch holds, and task add, numbered work add and work promote refuse an id given. Recommended, as T-24 specifies: supersede decision 4 with ids minted by itos, and remove its command from the notes. Or: keep decision 4 beside the counter as a hand check.

Asked as q-11, about T-24.

## Considered Options

The options are those the question names.

## Decision Outcome

Supersede decision 4, as the person approved with T-24: itos mints every task and numbered item id from refs/itos/ids, shared by all branches, so a new id is unique without a hand rule, and the notes lose decision 4's git grep command. Questions are not counter-minted, so decision 3 stands.

### Consequences

None recorded.

## More Information

Supersedes ADR-0004.
