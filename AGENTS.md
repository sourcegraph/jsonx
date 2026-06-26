# AGENTS.md

Guidance for coding agents working in this repository.

## Overview

`jsonx` is a single-package Go library (`github.com/sourcegraph/jsonx`) providing
an error-tolerant JSON parser, scanner, formatter, and editor that supports
trailing commas and comments. It is ported from Visual Studio Code's
comment-aware JSON APIs.

## Layout

All source lives in the repository root as a flat `package jsonx`:

- `scanner.go` — tokenizer.
- `parser.go` / `tree.go` — parsing and node tree.
- `visitor.go` — streaming visitor API.
- `format.go` — formatting.
- `edit.go` — applying edits to JSON text.
- `json.go` — package doc and `go:generate` directive.
- `json_stringer.go` — generated `String()` methods (do not edit by hand).
- `*_test.go` — tests for each file.

## Setup

Requires Go (see `go.mod` for the module's Go version; CI builds against `1.x`).
There are no third-party dependencies.

```bash
go build ./...
```

## Test, lint, generate

```bash
go test ./...        # run tests
go vet ./...         # vet
go generate ./...    # regenerate json_stringer.go after changing enum types
```

`json_stringer.go` is generated via the `//go:generate stringer` directive in
`json.go`. Regenerate it rather than editing it directly when changing the
`ParseErrorCode`, `ScanErrorCode`, `SyntaxKind`, or `NodeType` types.

## Conventions

- Standard Go formatting (`gofmt`); format code before committing.
- Keep changes within the flat root package layout.
