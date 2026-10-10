# Architecture

The root carries project policy and shared documentation. It has no executable or Go module.
Stack and feature branches add the code and checks they need while retaining these shared files.

`stack/go` supplies the Go foundation: a private counting domain behind a Files port, with real
disk access at the boundary and an in-memory fake for domain tests. `go/cli` adds a thin command
entry point over that domain. Each branch's copy of this file describes its own code.

Rendering and template updates belong to itos-template, not to this repository. Templates use
real literals and git branches; their generated projects do not contain a template interpreter.
