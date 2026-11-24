---
status: accepted
date: 2026-11-23
---

# A release reads Conventional Commit headers, releases any other as a patch and says so, and is attested, verified and announced to the Go proxy

## Context and Problem Statement

How does a made project's release choose its version and prove what it publishes, with no itos and no check on commit headers?

Asked as q-42, about T-75.

## Considered Options

- Fail the release when a header is not a Conventional Commit
- Release nothing for such a commit, as today
- Release it as a patch and name it in the version step's output

## Decision Outcome

The version comes from the commits' Conventional Commit headers as today, plus two rules. Before 1.0, a breaking change makes a minor release, not a major one (itos-cc's rule, go-template-itos#2, taken as is). A commit whose header is not a Conventional Commit counts as a patch-level change instead of failing the release (the person, 2026-11-23), and the version step's output names each commit it read that way. A merge commit is skipped when its header is Git's or GitHub's default, since its parents are read. The notes quote every BREAKING-CHANGE and Upgrading footer from each commit's last paragraph, read by the notes program itself rather than by itos. The archives are attested, and that attestation is verified on the built linux amd64 archive before the release is published. Afterwards proxy.golang.org is asked for the new version, and no answer only prints a warning. THIRD_PARTY_NOTICES.md, written by tools/bin/notices from the modules the binary links and Go's own license, is checked in CI and shipped in every archive.

### Consequences

A made project that writes free-form headers still releases its fixes, as patches, and the release run names the commits read that way. A feature written without its header ships only as a patch until its commit says feat. A 0.y.z project reaches 1.0.0 only by choice. A release whose attestation cannot be verified fails before it is published. The proxy step never fails a release that is already published. A new dependency fails CI until THIRD_PARTY_NOTICES.md is written again.
