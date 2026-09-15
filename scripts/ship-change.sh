#!/usr/bin/env bash
#
# Ship one OpenSpec change: stage, gate, archive, commit, push.
#
#   bash scripts/ship-change.sh <change-name> "<commit subject>"
#
# The gate is `nix flake check`, which builds the package, runs go vet, runs the
# full test suite and enforces the coverage ratchet in scripts/coverage-gate.sh.
#
# Staging happens BEFORE the gate on purpose: nix flakes only see files that git
# knows about, so an untracked new file would be invisible to the check and the
# gate would pass on code it never compiled.
#
# The gate is never bypassed. There is no --force and no --no-verify.

set -euo pipefail

change="${1:-}"
subject="${2:-}"

if [ -z "$change" ] || [ -z "$subject" ]; then
  echo "usage: bash scripts/ship-change.sh <change-name> \"<commit subject>\"" >&2
  exit 2
fi

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"

if [ ! -d "openspec/changes/$change" ]; then
  echo "No such active change: openspec/changes/$change" >&2
  echo "Active changes:" >&2
  openspec list >&2
  exit 2
fi

branch="$(git rev-parse --abbrev-ref HEAD)"
if [ "$branch" != "main" ]; then
  echo "On branch '$branch', not main. Refusing to ship." >&2
  exit 2
fi

echo "==> Checking the change is complete"
if grep -q '^- \[ \] ' "openspec/changes/$change/tasks.md" 2>/dev/null; then
  echo "tasks.md still has unchecked items. Refusing to ship a partial change." >&2
  grep -n '^- \[ \] ' "openspec/changes/$change/tasks.md" | head >&2
  exit 1
fi

echo "==> Validating the change"
openspec validate "$change" --strict

echo "==> Staging the working tree"
git add -A

if git diff --cached --quiet; then
  echo "Nothing staged. Refusing to ship an empty commit." >&2
  exit 1
fi

echo "==> Gate: nix flake check"
nix flake check

echo "==> Archiving the change"
openspec archive "$change" -y

echo "==> Staging the archive move"
git add -A

echo "==> Committing"
git -c user.name="Pim Snel" -c user.email="pim@technative.eu" \
  commit -m "$subject"

echo "==> Pushing main"
git push origin main

echo
echo "Shipped '$change' as: $(git rev-parse --short HEAD)"
