# Architecture

The root carries project policy and shared documentation. It has no executable or Go module.
Stack and feature branches add the code and checks they need while retaining these shared files.

Rendering and template updates belong to itos-template, not to this repository. Templates use
real literals and git branches; their generated projects do not contain a template interpreter.
