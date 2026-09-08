---
kind: cli
forge: github
tracker: github-issues
test: go test ./...
lint: golangci-lint run ./...
deploy.mode: release-on-tag
deploy.target: github-releases
---

# today

A Go command-line tool for working with the text of Christian scripture, plus the
packages behind [openscripture.today](https://openscripture.today). `pkg/ref` (the
scripture reference model), `pkg/text`, and `pkg/photo` are importable libraries as
well as CLI internals, so a change to them is not only a change to this binary.

## Releasing

Nothing ships on merge. A release happens only when a `vX.Y.Z` tag is pushed, at
which point `release.yaml` builds the Linux, Apple Silicon, and Apple Intel binaries
and publishes a GitHub release. Pushing a `release/*` branch first runs `prepare.yaml`,
which is the same steps minus publication — a dry run.

Both workflows refuse to run unless `cmd/version.txt` and the top heading of
`Changes.md` agree, and that heading must carry the release date **in US Central on
the day of the push**. A release that straddles midnight Central fails on the second
push even though it passed on the first. The `/release` skill in `.claude/skills/`
drives the whole sequence and documents the rest of the failure modes; use it rather
than doing the steps by hand.

## Working here

- CI fails under 80% coverage, measured with `cmd` and `tools/gen/verses` excluded.
  The exact command is in CLAUDE.md and needs `bash`, not `zsh`.
- `pkg/ref/canonical.go` and `pkg/ref/abbr.go` are generated. Edit the sources in
  `tools/gen/verses/` and regenerate; do not hand-edit the output.
- `depguard` restricts imports to an allowlist in `.golangci.yaml`. Adding a
  dependency is a deliberate edit to that list, not an incidental one.
- `prepare.yaml` and `release.yaml` share their steps and CI diffs them. Edit both
  together.
- Tests need no credentials — every external service (ESV, Unsplash) is faked.

## Maintenance

Upkeep here has been dependency sweeps: `.claude/skills/maintenance-deps` runs
`zed:dependabot-sweep`, and `maintenance-weekly` runs that. Nothing about this
project needs a human mid-run except cutting a release.
