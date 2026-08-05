# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```shell
go build ./...
go test ./...
go test ./pkg/ref/ -run TestCanon_Resolve_Proper              # one test
go test ./pkg/ref/ -run 'TestCalculateRefStats/single_chapter' # one subtest (spaces become underscores)
golangci-lint run ./...                                        # CI pins v2.8.0
go run . ref "2 John 1:1-4" --stat                             # run the CLI without installing
```

Coverage gate (CI fails under 80%; `cmd` and the generator are excluded):

```shell
coveredFiles=$(go list ./... | grep -v 'github.com/zostay/today/\(cmd\|tools/gen/verses\)')
go test $coveredFiles -coverprofile=coverage.out && go tool cover -func=coverage.out
```

Note this needs `bash`, not `zsh` — zsh will not word-split `$coveredFiles`.

Regenerate the canon data (see "Generated data" below):

```shell
cd tools/gen/verses && go generate .
```

## Architecture

### The reference model (`pkg/ref`)

This is the heart of the project, and it is a three-layer pipeline. Understanding
the split is the difference between a one-line change and a scattered one.

1. **Parse** (`parse.go`) — a hand-written recursive-descent parser over runes.
   Produces an AST. Accepts `:` or `.` as the chapter/verse separator and several
   dash characters.
2. **Validate** (`Validate()` methods in `ref.go`) — canon-*independent*
   structural plausibility only. `Philemon 12:4` and `Sterling 2:2` both validate;
   validation does not know what books exist.
3. **Resolve** (`Canon.Resolve` in `books.go`) — maps the AST onto a `Canon` and
   yields `[]Resolved`. This is the only layer that knows real books and verses.

The AST types:

- `Verse` leaves: `N{Number}` is a bare number whose meaning is contextual — a
  *chapter* in a book with chapters, a *verse* in a book without. `CV{Chapter, Verse}`
  is explicit.
- `Relative` (no book): `Single`, `Range`, `AndFollowing` (`ff`/`ffc`/`ffb`),
  `Related` (comma list).
- `Absolute` (has a book): `Proper` (book + relative), `Multiple` (semicolon list),
  and `Resolved`.

**`Resolved` is the currency of the whole codebase.** It is a normalized, single,
contiguous range: `{Book *Book, First, Last Verse}`. Formatting (`format.go`),
statistics (`stats.go`), text lookup (`pkg/text`), pericopes, and canon filtering
all consume `Resolved` and never the parsed AST. Consequently:

- Normalization decisions belong in `Resolve`, where they propagate to every
  consumer for free. `ensureVerseMatchesBook` is the choke point — it widens `N`
  into `CV` for books with chapters and folds chapter-1 `CV` down to `N` for books
  without.
- Output formatting is a property of `Resolved`, not of the input string.
  `Resolved.compactRef` decides whether to print `Genesis`, `Genesis 12`,
  `Genesis 12:4`, or `Genesis 12:4-6`, and every formatter delegates to it.
  So `today ref 2jn1.1-4` prints `2 John 1-4`: the input form is discarded at
  resolution.

**`Book.JustVerse`** marks the five books with no chapters (Obadiah, Philemon,
2 John, 3 John, Jude). Their `Book.Verses` hold `N`; every other book's hold `CV`.
Both forms of a chapterless reference are accepted — `2 John 1-4` and
`2 John 1:1-4` resolve identically — but any chapter other than 1 is an error.

**`Book.Verses` is a flat ordered list of every verse in the book.** `Contains`,
`Resolved.Verses()`, and filtering all walk it. This is also why a chapter cannot
be assumed to start at verse 1: `Canon.Filtered(exclude...)` clones a canon and
physically deletes verses from these slices, so a filtered canon can have chapters
beginning mid-way.

### Generated data

`pkg/ref/canonical.go` (every verse of the Protestant canon, plus categories) and
`pkg/ref/abbr.go` (book names and accepted abbreviations) are **generated — do not
hand-edit**. Their sources live in `tools/gen/verses/`: `esv.json` (verse
structure), `categories.yaml`, `abbr.yaml`. Change the YAML/JSON, then regenerate.

`AbbrTree` (`tree.go`) is a prefix tree over every accepted abbreviation; it drives
the liberal book-name matching that lets `2jn`, `2 Jn.`, and `Second John` all
resolve, and returns `MultipleMatchError` on ambiguity.

### Service layers

Each external service sits behind a narrow interface with a swappable
implementation, so tests use fakes rather than network calls:

- `pkg/text` — `Resolver` interface (operates on `*ref.Resolved`), wrapped by
  `Service`, which is the string-accepting front door (`svc.VerseText(ctx, "John 3:16")`).
  `pkg/text/esv` implements `Resolver` against the ESV API.
- `pkg/photo` — `Source` interface + `Service`; `pkg/photo/unsplash` implements it.
- `pkg/ost` — client for openscripture.today, composing a text and a photo service.
- `cmd/` — cobra commands, all registered in `cmd/root.go`.

### Credentials

ESV: `ESV_API_TOKEN`, else `~/.esv.yaml`. Unsplash: `UNSPLASH_API_TOKEN`, else
`~/.unsplash.yaml`. Both files use the key `access_key`. The test suite needs
neither — every external service is faked.

## Constraints that bite

- **depguard** restricts imports by an allowlist in `.golangci.yaml`. Non-test code
  may import only stdlib, `github.com/zostay/*`, cobra, wrap, go-unsplash, resize,
  go-dateparser, levenshtein, yaml.v3, and oauth2. Test code may add testify.
  Adding any other dependency means editing the allowlist deliberately.
- **paralleltest / tparallel** — every test and subtest needs `t.Parallel()`.
- **godot** — top-level comments must end in a period.
- Existing tests use the `tt := tt` loop-variable copy. Modern Go does not need it
  and gopls flags it, but the enabled linters do not; match the surrounding style.

## Release process

`Changes.md` and `cmd/version.txt` must agree, and the workflows check it. The top
heading of `Changes.md` must read exactly `## X.Y.Z  YYYY-MM-DD` (two spaces, the
date being the release day in US Central) and `cmd/version.txt` must contain the
same version. During development the heading carries `TBD` in place of the date; a
`chore: releng` commit fills it in at release time.

Pushing a `release/*` branch runs the prepare workflow (a dry run); pushing a `vX.Y.Z`
tag runs the release workflow, which builds binaries and creates the GitHub release.
CI also diffs the two workflow files to make sure their shared steps stay in sync —
edit both together.

The `/release` skill (`.claude/skills/release/SKILL.md`) drives the whole process and
documents the failure modes, including the one where the changelog date must match the
US Central date on *both* the branch push and the tag push.
