---
name: release
description: Cut a software release for this project — set the version and changelog date, push a release branch, open and merge a release PR, tag it, and monitor the GitHub release through to publication.
---

# Release

Run the full release process for `today`, from changelog to published GitHub release.

Take an optional version argument (`/release 1.2.0`). Without one, work the version
out from the changelog and the commits since the last tag.

## How releasing works here

Three workflows are involved. Knowing what each checks saves a lot of failed runs.

| Workflow | Trigger | What it does |
| --- | --- | --- |
| `test.yaml` ("Test and Sanity") | push to `master`, any PR | lint, tests, 80% coverage floor, workflow-sync check |
| `prepare.yaml` ("Prepare for Release") | push to `release/*` | version + changelog checks, then builds binaries — a dry run of the release |
| `release.yaml` ("Release") | push of a `v*` tag | the same checks and builds, then creates and publishes the GitHub release |

`prepare.yaml` and `release.yaml` share their first ten steps, and `test.yaml`
diffs the two to make sure they stay in sync. **Never edit one without the other.**

Both release workflows derive the version from the ref name with
`grep -Eo '[0-9]+\.[0-9]+\.[0-9]+.*$'`, so branch `release/v1.2.0` and tag `v1.2.0`
both yield `1.2.0`. Then they enforce two things:

1. `grep -q "$RELEASE_VERSION" cmd/version.txt` — a substring match, so `1.2.0`
   in the file satisfies it.
2. The first line of `Changes.md` equals `## $RELEASE_VERSION  $date` **exactly**
   — two spaces between version and date, and `$date` is *the day the workflow
   runs*, in `America/Chicago`.

That second check is the one that bites. Read "The midnight problem" below before
starting a release late in the US Central day.

## Steps

### 1. Preflight

```shell
git checkout master && git pull
git status --short                 # must be clean
gh run list --branch master --limit 1   # must be green
git tag --sort=-v:refname | head -3
```

Stop and ask if the tree is dirty or `master` is red. Do not release from a red master.

### 2. Decide the version

If the user named a version, use it. Otherwise read the top (unreleased) section of
`Changes.md` and `git log $(git describe --tags --abbrev=0)..master --oneline`, then
apply semver. This changelog marks intent with emoji:

- `:boom: Breaking Change :boom:` → **major**
- `:computer:` (new or changed CLI behavior) or `:sparkles:` → **minor**
- `:hammer:` (fixes) or dependency bumps alone → **patch**

Now that the project is past 1.0, any breaking change to an exported package API or
to command behavior is a major bump. State the version you picked and why before
proceeding.

### 3. Set the changelog heading

The top section of `Changes.md` must be the release being cut. Between releases this
repo leaves the heading as `## WIP  TBD`, or as `## X.Y.Z  TBD` once a version has been
decided. **Neither form may survive into the release** — replace it wholesale with the
version and date:

```
## 1.2.0  2026-08-05
```

Two spaces. Get the date from the release workflow's timezone, not the local one:

```shell
TZ=America/Chicago date +%Y-%m-%d
```

The heading must begin with a digit. `release.yaml` extracts the release notes by
reading from the first `^## [0-9]` line to the next one, so a heading left as `## WIP`
would silently publish the *previous* version's notes.

Preview exactly what will be published (the workflow's `sed` is GNU-only and fails on
macOS BSD `sed`, so use `awk` locally):

```shell
awk '/^## [0-9]/{if(n++)exit;next} n' Changes.md
```

If the section is thin — dependency bumps only, or entries that do not describe what
changed for a user — improve it now. This text becomes the public release notes.

### 4. Match `cmd/version.txt`

```shell
echo 1.2.0 > cmd/version.txt
go build ./... && go run . version    # should print "today v1.2.0"
```

### 5. Push the release branch

```shell
git checkout -b release/v1.2.0
git add Changes.md cmd/version.txt
git commit -m "chore: releng"        # the conventional message in this repo
git push -u origin release/v1.2.0
```

This fires `prepare.yaml`.

### 6. Open the release PR

```shell
gh pr create --base master --title "Release v1.2.0" --body "..."
```

Summarize what is in the release and link the changelog section. Opening the PR fires
`test.yaml` and triggers a Copilot review.

### 7. Wait for CI, and fix what fails

Both workflows must pass — `prepare.yaml` on the branch push and `test.yaml` on the PR:

```shell
gh pr checks <pr> --watch
gh run list --branch release/v1.2.0 --limit 5
```

On failure, read the log (`gh run view <id> --log-failed`), fix, push, and wait again.
Repeat until green or until the failure needs a human. Common failures:

- **"Changes.md is out of date!"** — the heading date does not equal the runner's
  US Central date. Usually the date rolled over. Fix the heading and push again.
- **"cmd/version.txt does not match"** — the file and the branch name disagree.
- **Coverage below 80%** — new uncovered code. Add tests; do not lower the floor.
- **"Prepare and Release workflows are not in sync!"** — someone edited one of the
  two release workflows alone. Both must change together.

Stop and report if you hit something with no clear fix, rather than guessing at
release-time changes.

### 8. Handle the Copilot review

Copilot reviews every PR here automatically. Its findings are usually low value, so
triage rather than obey:

- **Simple, obvious, and release-engineering-shaped** (a changelog typo, a wrong
  version string, a stale doc line) — just fix it in the release PR.
- **Anything larger** — a refactor, a design opinion, a bug in code this release did
  not touch — open a GitHub issue (`gh issue create`) and link it from the PR. Do not
  expand the release PR to cover it.
- **Never let a Copilot finding block the release** unless it is genuinely critical:
  it identifies a real defect that this release would ship to users. In that case stop
  and tell the user; do not decide alone to abandon the release.

Say which findings you fixed, which you filed, and which you dismissed.

### 9. Merge

```shell
gh pr merge <pr> --merge --delete-branch
```

`master` is protected and requires the "Test and Sanity" check, so merge through the
PR — do not push to `master` directly.

### 10. Re-check the date, then tag

**Before tagging, confirm the changelog date still matches the US Central date**, since
`release.yaml` re-runs that check when the tag lands:

```shell
git checkout master && git pull
head -n1 Changes.md
TZ=America/Chicago date +%Y-%m-%d
```

If they no longer agree, see "The midnight problem". If they agree:

```shell
git tag v1.2.0
git push origin v1.2.0
```

### 11. Monitor the release

The tag fires `release.yaml`, which re-runs the checks, builds Linux amd64, macOS arm64,
and macOS amd64 binaries, creates a **draft** release with notes from `Changes.md`,
uploads the three binaries, then un-drafts it.

```shell
gh run watch $(gh run list --workflow release.yaml --limit 1 --json databaseId -q '.[0].databaseId')
gh release view v1.2.0
```

Verify the release is published rather than draft, that all three binaries are attached,
and that the notes match the changelog section. Report the release URL, the version, what
shipped, and anything that needed fixing along the way.

If `release.yaml` fails *after* the tag is pushed, the tag exists but the release does
not. Fix the cause on `master` through a PR, delete and re-push the tag
(`git push --delete origin v1.2.0` then re-tag the new commit), and say clearly that you
did so.

### 12. Open the next cycle (optional)

Offer to add a fresh `## WIP  TBD` heading at the top of `Changes.md` on `master` so the
next change has somewhere to go.

## The midnight problem

Both release workflows require the `Changes.md` date to equal the day *they run*, in
`America/Chicago`. The branch push and the tag push are separate runs, often minutes
apart but sometimes not. If the release crosses midnight Central, the first run passed
and the second will fail.

Prevention: do the whole flow in one sitting, and avoid starting after about 22:00
Central.

Recovery, if the date rolls over before tagging: update the `Changes.md` heading to the
new date on `master` via a small PR (`master` is protected, so it needs one), then tag
the resulting commit. Do not tag a commit whose changelog date is stale — the release
workflow will reject it.
