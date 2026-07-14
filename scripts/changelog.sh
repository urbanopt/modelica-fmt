#!/usr/bin/env bash
# Keep CHANGELOG.md and GitHub releases in sync.
#
# CHANGELOG.md is the single source of truth for release notes: `make
# release` finalizes the "[Unreleased]" section into a dated version
# section, and the same section is extracted and handed to goreleaser as the
# GitHub release body (see publish.yml and the Makefile).
#
# Usage:
#   changelog.sh extract <heading>            Print the notes under "## [<heading>]"
#   changelog.sh finalize <tag> <date> [file]  Turn "## [Unreleased]" into
#                                               "## [Unreleased]" + "## [<tag>] - <date>"
set -euo pipefail

usage() {
  echo "usage: $0 extract <heading> [file]" >&2
  echo "       $0 finalize <tag> <date> [file]" >&2
  exit 2
}

cmd="${1:-}"
[ -n "$cmd" ] || usage
shift

case "$cmd" in
  extract)
    heading="${1:-}"
    file="${2:-CHANGELOG.md}"
    [ -n "$heading" ] || usage
    awk -v tag="$heading" '
      BEGIN { pat = "^## \\[" tag "\\]"; found = 0; n = 0 }
      $0 ~ pat { found = 1; next }
      found && /^## \[/ { exit }
      found { n++; buf[n] = $0 }
      END {
        start = 1; end = n
        while (start <= end && buf[start] ~ /^[[:space:]]*$/) start++
        while (end >= start && buf[end] ~ /^[[:space:]]*$/) end--
        for (i = start; i <= end; i++) print buf[i]
      }
    ' "$file"
    ;;
  finalize)
    tag="${1:-}"
    date="${2:-}"
    file="${3:-CHANGELOG.md}"
    [ -n "$tag" ] && [ -n "$date" ] || usage
    if ! grep -q '^## \[Unreleased\]' "$file"; then
      echo "error: $file has no '## [Unreleased]' heading" >&2
      exit 2
    fi
    tmp="$(mktemp)"
    awk -v tag="$tag" -v date="$date" '
      { print }
      !inserted && $0 ~ /^## \[Unreleased\]/ {
        print ""
        print "## [" tag "] - " date
        inserted = 1
      }
    ' "$file" > "$tmp"
    mv "$tmp" "$file"
    ;;
  *)
    usage
    ;;
esac
