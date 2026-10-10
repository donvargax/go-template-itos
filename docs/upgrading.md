# Upgrading

Read the branch's configuration and tool comments before changing a pin. The root uses a pinned
itos and pinned Renovate presets. Move itos's repository pin with `itos pin <version>` and the
workflow action with it, to that release's commit with a `# v<version>` comment and no `version`
input, so the action installs the launcher of the release its ref names. Then apply each step of
the release's Upgrading section that the branch's configuration meets, and run
`itos init --agent-rules`. Do not change only the launcher or action version.

The owner's tools can be pinned when released. Other tools and dependencies normally wait seven
days, with security updates following the established security-update policy. When tools with
download checksums arrive, move their versions and verified checksums together. Dependencies
retain their upstream licenses and notices.

Renovate reads its configuration from main, whose `.github/renovate.json5` names main, stack/go
and go/cli as base branches and merges each branch's own file over it, so each branch's requests
follow that branch's managers and are judged by its own CI. The Renovate app is installed, but no
request has landed yet, so the automation stays unverified until one lands on each branch.

Go modules, binary release tooling and template combination checks are not installed at this
root-bootstrap stage. Their upgrade instructions must be completed when the corresponding
pieces are implemented. This template does not yet define template release tags or claim an
adoption/update command is available.
