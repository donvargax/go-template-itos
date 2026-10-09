# Upgrading

Read the branch's configuration and tool comments before changing a pin. The root uses a pinned
itos and pinned Renovate presets. Move itos's repository pin and the workflow action together
through its supported upgrade command; do not change only the launcher or action version.

The owner's tools can be pinned when released. Other tools and dependencies normally wait seven
days, with security updates following the established security-update policy. When tools with
download checksums arrive, move their versions and verified checksums together. Dependencies
retain their upstream licenses and notices.

Go modules, binary release tooling and template combination checks are not installed at this
root-bootstrap stage. Their upgrade instructions must be completed when the corresponding
pieces are implemented. This template does not yet define template release tags or claim an
adoption/update command is available.
