## What and why

## Clean-room provenance (required for any change under `eval/baselib/`)

Delete this section only for a change that touches no file of `eval/baselib/`.

- **Specification read** (the `.Rd` pages, by name; the Language Definition or R Internals sections; nothing else):
- **R function bodies seen** (any R closure printed, deparsed or read, by anyone who wrote the change, with how: oracle, installed R, source tree, a diff, memory). Write "none" only if true. If one was seen, say which, and who rewrote it blind:
- **Gate score before any reshaping** (`go run ./cmd/cleanroom check`, the first run on the first draft, not the last):
- **Gate score of the submitted files**:
- **Rule**: a baselib function whose first draft fails the gate is implemented as a Go native from the Rd page and the oracle, never reshaped in R; record the failed first-draft score in the PR. Natives written for that reason:
- **Oracle experiments** (the inputs sent; the oracle refuses to print base closures, see README "Clean-room process"):

The gate is a floor, not proof: this section is the record.

Refs/Closes: use `Closes #N` or `Refs #N` for any issue mentioned (a bare `#N` is refused by the merge tool).
