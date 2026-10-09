# Architecture

The root carries project policy and shared documentation. It has no executable or Go module.
Stack and feature branches add the code and checks they need while retaining these shared files.

`stack/go` will supply the Go foundation. `go/cli` will demonstrate thin entry points over a
domain through ports, with real disk access at the boundary and an in-memory fake for domain
tests. The count example and its detailed interfaces are not implemented at this root stage.

Rendering and template updates belong to itos-template, not to this repository. Templates use
real literals and git branches; their generated projects do not contain a template interpreter.
