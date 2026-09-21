#!/usr/bin/env bash
#
# Prints one section of CHANGELOG.md: the lines under a `## [<label>]` heading,
# up to the next heading that begins `## [`.
#
# The label is a version, `0.8.0`, or `Unreleased`. Both are the same shape in
# Keep a Changelog, which is why one reader serves the release script's check
# that something was written and the workflow's need for the text to publish.
#
# awk rather than sed, for the reason release.sh already gives about its own
# changelog editing: the substitutions this would need are GNU extensions, and
# this project ships darwin builds.
set -euo pipefail

label="${1:-}"
file="${2:-CHANGELOG.md}"

if [ -z "$label" ]; then
    echo "usage: changelog-entry.sh <version|Unreleased> [changelog]" >&2
    exit 2
fi

# Matched as a prefix including the closing bracket, so `0.7.2` cannot match
# the heading of `0.7.20`, and so a version named in an entry's prose cannot
# start a section: only a line beginning `## [` can.
awk -v want="## [$label]" '
    index($0, want) == 1 { found = 1; next }
    found && index($0, "## [") == 1 { exit }
    found { print }
    END { if (!found) exit 1 }
' "$file"
