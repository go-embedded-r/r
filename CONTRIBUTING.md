# Contributing

This repository is BSD-3-Clause and R is GPL-2-or-later. The rules below keep
it that way; a pull request that breaks one is not merged.

## Clean-room rule

- Work from **specifications**: the *R Language Definition*, the *R Internals*
  manual, the help pages of the functions, and the behaviour of the oracle
  (`go run ./cmd/roracle eval '...'`).
- Do **not** open R's `src/main`, `src/nmath`, `src/library/*/R` or `tests/`
  while writing the Go (or R) code that corresponds to them. Reading them to
  count or to compare is for the roadmap's measurements, not for implementing.
- Do **not** commit anything copied from R: no source, no test file, no
  `.Rout.save`, no help page text. Do not commit a download of any of it, and
  do not make the program download it at run time.
- Golden files come from `roracle`, from inputs **written in this repository**.
  A new input goes in `testdata/snippets.R` or `testdata/values.R`; then run
  `go run ./cmd/roracle regen` and commit the result.
- Base R's R code is rewritten here under BSD. Rewritten means designed from
  the documented behaviour and checked against the oracle, not transcribed.

## Every change

- Tests with it; the suite must pass with and without R
  (`RGO_REQUIRE_ORACLE=1 go test ./...` with R 4.6.2, plain `go test ./...`
  without).
- A parser or decoder gets a fuzz target with its corpus committed.
- Untrusted input is bounded (size, depth, allocation) and refused with an error.
- `go vet`, `gofmt`, `govulncheck` clean.
