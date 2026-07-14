# Releases

Modelica formatter uses [GoReleaser](https://goreleaser.com/) for building assets and documenting tags for release. This is run automatically by GitHub Actions whenever a new tag is pushed to remote. See [.goreleaser.yml](/.goreleaser.yml) and [publish.yml](/.github/workflows/publish.yml) for the configuration.

## CHANGELOG.md is the source of truth

[CHANGELOG.md](/CHANGELOG.md) follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/):
notable changes are added as bullets under `## [Unreleased]` as they land
(e.g. in the same PR that makes the change). When a release is cut,
`make release` moves that section into a dated `## [vX.Y.Z] - YYYY-MM-DD`
heading, and the GitHub Actions release workflow extracts that exact section
and publishes it verbatim as the GitHub release body (via
`scripts/changelog.sh` and goreleaser's `--release-notes` flag). This keeps
CHANGELOG.md and the GitHub releases page in sync automatically — there is no
separate release notes to write.

## How to make a release

Make sure `CHANGELOG.md` has one or more entries under `## [Unreleased]`
describing the release; `make release-check` fails otherwise.

Run the release preflight:

```bash
make release-check VERSION=<version>
```

Where `<version>` is a semantic version without the leading `v`, for example `1.2.3` or `1.2.3-rc.1`. This runs tests, verifies `CHANGELOG.md` has unreleased notes, builds the binary with the requested release metadata, verifies that `modelica-fmt -v` reports the requested version id, checks the GoReleaser configuration, and runs a GoReleaser snapshot build using the unreleased notes.

Then create and push the annotated tag:

```bash
make release VERSION=<version> MESSAGE="<message>"
```

This moves the `## [Unreleased]` section of `CHANGELOG.md` into a new
`## [v<version>] - <date>` heading, commits that change, creates and pushes
tag `v<version>` (and the commit) to `origin`, triggering the GitHub Actions
GoReleaser workflow. Automatic CHANGELOG finalization requires releasing from
`HEAD` (the default).

To tag a specific commit instead of `HEAD`, pass `SHA=<commit>`. In that case
you must finalize `CHANGELOG.md` for the release yourself first (rename
`## [Unreleased]` to `## [v<version>] - <date>` and leave a fresh, empty
`## [Unreleased]` above it), since it isn't safe to automatically rewrite
history at an arbitrary commit:

```bash
make release-check VERSION=<version>
git tag -a v<version> -m "<message>" <SHA>
git push origin v<version>
```

After the build successfully finishes, go to the releases page on GitHub,
check that it looks good, and publish it (releases are created as drafts).
