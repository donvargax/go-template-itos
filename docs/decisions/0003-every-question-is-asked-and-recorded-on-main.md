---
status: accepted
date: 2026-10-09
---

# Every question is asked and recorded on main

## Context and Problem Statement

Question and decision numbers are per branch, and changes merge only downward, so a question asked on stack/go or go/cli takes a number main will reuse and the next merge down holds two different records under it. Recommended: ask, answer and record every question on main, even one specific to a stack or feature, so records merge down like other shared docs and numbers stay unique. Or: ask only template-wide questions on main and renumber branch-specific collisions by hand at each merge. Or: keep per-branch numbering and resolve collisions at each merge.

Asked as q-3.

## Considered Options

The options are those the question names.

## Decision Outcome

Ask, answer and record every question on main, even one specific to a stack or feature; records merge down like other shared docs and numbers stay unique. q-2, first asked on stack/go, is asked again on main with the same text and answer, so the merge down meets an identical record.

### Consequences

None recorded.
