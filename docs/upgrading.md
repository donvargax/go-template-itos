# Upgrading

Read the current branch's configuration and tool comments before moving a pin. Move itos's
repository pin with `itos pin <version>` and the workflow action with it, to that release's
commit with a `# v<version>` comment and no `version` input, so the action installs the launcher
of the release its ref names. Then apply each step of the release's Upgrading section that the
branch's configuration meets, and run `itos init --agent-rules`. Preserve verified checksums
when upgrading the tools fetched by `tools/bin/pinned`; never update only a version beside
stale hashes.

Go dependencies are locked in go.mod/go.sum. Use the project's gates to check module
verification and tidy-diff, and review code and vulnerability results before an update lands.
The Go toolchain is pinned too: its standard library is part of the resulting code, so security
patches matter even when module versions stay unchanged.

Tools by itos's owner can be pinned on release. Other updates normally wait seven days, with
security fixes following the existing security-update policy. Renovate reads
`.github/renovate.json5` from `main`, so its requests follow that file's managers and are judged
by the project's own CI. The file alone does not make Renovate run: the automation stays
unverified until the Renovate app is installed and a request has landed.

Fresh mutation records are required when code or relevant test imports change. Record them
from the push's range start, the nearest ancestor with a green CI run, since `itos ci run`
judges the proof from there.

This project has no CLI binary release yet.
