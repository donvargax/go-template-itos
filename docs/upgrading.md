# Upgrading

Read the branch's configuration and tool comments before changing a pin. The root uses pinned tools and
pinned Renovate presets.

The owner's tools can be pinned when released. Other tools and dependencies normally wait seven
days, with security updates following the established security-update policy. When tools with
download checksums arrive, move their versions and verified checksums together. Dependencies
retain their upstream licenses and notices.

Renovate reads its configuration from main's `.github/renovate.json5`; `TEMPLATE.md` says how
it reaches the other branches and what is verified of it.

Go modules, binary release tooling and template combination checks are not installed at this
root-bootstrap stage. Their upgrade instructions must be completed when the corresponding
pieces are implemented. The root does not yet define template release tags or claim an
adoption/update command is available.
