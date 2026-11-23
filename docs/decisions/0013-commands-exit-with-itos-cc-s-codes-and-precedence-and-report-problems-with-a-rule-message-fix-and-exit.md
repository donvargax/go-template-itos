---
status: accepted
date: 2026-11-23
---

# Commands exit with itos-cc's codes and precedence and report problems with a rule, message, fix and exit

## Context and Problem Statement

Which exit codes, precedence and problem shape do go/cli's commands use: the template's 0, 1, 2, 3, 70, 75 and 130, or itos-cc's 0, 1, 75, 80, 81, 82, 100 and 130 with an explicit precedence?

Asked as q-39, about T-69.

## Considered Options

- Keep 0, 1, 2, 3, 70, 75 and 130
- itos-cc's 0, 1, 75, 80, 81, 82, 100 and 130 with an explicit precedence

## Decision Outcome

itos-cc's (the migration plan, approved by the person 2026-11-23). 0 success, 1 a check said no, 75 temporary, 80 a usage or config error, 81 an unknown command, 82 a requirement not met, 100 an error nothing classified, 130 interrupted. With problems of several codes, the run exits with the first of 130, 80, 81, 100, 82, 75, 1. A problem carries a stable dotted lower-case rule, a message, a fix and its exit, with its subject as keys of their own in --json. Codes stay clear of 2, 3 and 70, which shells, Go and sysexits already use.

### Consequences

Scripts that read 2, 3 or 70 from a made project read 80, 81, 82 or 100 instead, a breaking change of the contract. Rule IDs such as COUNT_INPUT_MISSING become count.input-missing. Decision 8's 130 stays.
