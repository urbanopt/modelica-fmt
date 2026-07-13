# Releases
Modelica formatter uses [GoReleaser](https://goreleaser.com/) for building assets and documenting tags for release. This is run automatically by GitHub Actions whenever a new tag is pushed to remote. See [.goreleaser.yml](/.goreleaser.yml) and [publish.yml](/.github/workflows/publish.yml) for the configuration.

## How to make a release
Pick the release version and update `CHANGELOG.md` so `# Unreleased` is moved to a versioned heading for the release.

Run the release preflight:

```bash
make release-check VERSION=<version>
```

Where `<version>` is a semantic version without the leading `v`, for example `1.2.3` or `1.2.3-rc.1`. This runs tests, builds the binary with the requested release metadata, verifies that `modelica-fmt -v` reports the requested version id, checks the GoReleaser configuration, and runs a GoReleaser snapshot build.

Then create and push the annotated tag:

```bash
make release VERSION=<version> MESSAGE="<message>"
```

This creates and pushes tag `v<version>` to `origin`, triggering the GitHub Actions GoReleaser workflow.

To tag a specific commit instead of `HEAD`, pass `SHA=<commit>`:

```bash
make release VERSION=<version> MESSAGE="<message>" SHA=<commit>
```

If you need to tag manually, run the same preflight first and then create and push the tag:

```bash
make release-check VERSION=<version>
git tag -a v<version> -m "<message>" [SHA]
git push origin v<version>
```

After the build successfully finishes, go to the releases page on Github, check that it looks good, and publish it.
