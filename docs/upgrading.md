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

The owner's tools can be pinned on release. Other updates normally wait seven days, with
security fixes following the existing security-update policy. In the template, main's Renovate
configuration names this branch as a base branch and merges this branch's
`.github/renovate.json5` over it, so its requests follow this branch's managers and are judged by
its own CI. The Renovate app is installed, but no request has landed yet, so the automation stays
unverified until one lands on this branch.

Fresh mutation records are required when code or relevant test imports change. Record them
from the push's range start, the nearest ancestor with a green CI run, since `itos ci run`
judges the proof from there.

git-cliff publishes a `.sha512` file per archive rather than a checksums file: move its pin by
checking each archive against that file first, then pinning the archive's SHA-256.

The release job in `ci.yml` runs only on a push to `main`, so this branch cuts no CLI binary
release, and a made project's `main` cuts none until `made-project-ci` makes its CI run there.
Template releases and the generator's adoption/update workflow are separate from made-project
binary release tooling and are not claimed here.
