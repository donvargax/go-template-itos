# Upgrading

Read the current branch's configuration and tool comments before moving a pin. Preserve verified
checksums when upgrading the tools fetched by `tools/bin/pinned`; never update only a version
beside stale hashes.

Go dependencies are locked in go.mod/go.sum. Use the project's gates to check module
verification and tidy-diff, and review code and vulnerability results before an update lands.
The Go toolchain is pinned too: its standard library is part of the resulting code, so security
patches matter even when module versions stay unchanged.

Tools by itos's owner can be pinned on release. Other updates normally wait seven days, with
security fixes following the existing security-update policy. Renovate reads
`.github/renovate.json5` from `main`, so its requests follow that file's managers and are judged
by the project's own CI. The file alone does not make Renovate run: the automation stays
unverified until the Renovate app is installed and a request has landed.

Mutation runs fresh in CI on each push, so a move of itos-cc needs no records renewed.

git-cliff publishes a `.sha512` file per archive rather than a checksums file: move its pin by
checking each archive against that file first, then pinning the archive's SHA-256.

The release job in `ci.yml` runs only on a push to `main`, after the ci and platform jobs pass,
and cuts a release when the commits since the last tag call for one: GoReleaser's five archives
and `checksums.txt`, the archives attested, published with the notes `tools/bin/release-notes`
writes.
