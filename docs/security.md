# Security and trust boundaries

Use templates you trust: a made project contains the selected template's files and workflows.
Do not put credentials in template answers, committed files, fixtures or logs. The owner and
module answers are public project identity, not authentication data.

Pinned tools keep their upstream identities and licenses. Upgrades must preserve checksum and
provenance verification; the source owner's permission to distribute copied infrastructure under
MIT does not grant permission to change third-party licenses.

gitleaks is pinned in tools/bin/pinned as the other tools are; `TEMPLATE.md` says where it
runs. Code-specific checks arrive with their code.

itos's hooks and commit rules are workflow controls, not a security sandbox. A template's own
check commands can execute programs; run them only in an environment appropriate to that trust.
