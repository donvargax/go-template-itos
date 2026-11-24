---
status: accepted
date: 2026-11-23
---

# A release can publish a container image through GoReleaser's ko, off by default, pinned, attested and verified

## Context and Problem Statement

How does a made project's release publish a container image, and how is it turned on and proved?

Asked as q-43, about T-76.

## Considered Options

- No image; a project that needs one adds its own Dockerfile
- A Dockerfile and docker buildx in release.yml
- GoReleaser's kos section, off by default, switched on in the tree

## Decision Outcome

Through GoReleaser's kos section in .goreleaser.yaml, so the release that builds the archives also builds the image with ko, from the same static binary (CGO_ENABLED=0), with no Dockerfile (go-template-itos#18). The image is off by default. The switch is the section's own skip line in the tree, so turning it on is one reviewed commit. The base is the distroless static nonroot image pinned by digest, which a Renovate regex manager over .goreleaser.yaml moves. It builds for linux/amd64 and linux/arm64 and pushes to ghcr.io/<owner>/<project> with the workflow's token (packages write). It is tagged with the version only, never latest. ko's SBOM is kept. When the release's artifacts list an image, release.yml attests its digest with actions/attest-build-provenance (subject-name and subject-digest, push-to-registry), then verifies that attestation with gh attestation verify oci://… before Publish, as it does for the archives. CI runs goreleaser check, so the section stays valid while it is off.

### Consequences

A CLI-only project publishes no image and needs no packages write in practice, though the release job is granted it. A project that deploys turns the image on with one commit and gets a pinned base, both Linux architectures, an SBOM and a verified attestation. The image is pushed before the release is published, so a run that fails after GoReleaser leaves the version's image tag in the registry, which the next run's push replaces. Consumers pull the image by digest, as the README says.
