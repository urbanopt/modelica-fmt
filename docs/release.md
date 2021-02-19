# Releases
Modelica formatter uses [GoReleaser](https://goreleaser.com/) for building assets and documenting tags for release. This is run automatically by GitHub Actions whenever a new tag is pushed to remote. See [.goreleaser.yml](/.goreleaser.yml) and [publish.yml](/.github/workflows/publish.yml) for the configuration.

## How to make a release
First, create a tag:
```bash
git tag -a v<version> -m "<message>" [SHA]
```
Where `v<version>` is a valid [semantic version](https://semver.org/) (e.g. `v1.2` or `v1.2-pr.1`) and `<message>` is the tagging message (e.g. "First official release")

Next, push the tag up to remote, triggering a GitHub Actions workflow build which runs GoReleaser:
```bash
git push origin v<version>
```

After the build successfully finishes, go to the releases page on Github, check that it looks good, and publish it.
