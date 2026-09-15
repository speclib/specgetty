#!/usr/bin/env bash
set -e

# =============================================================================
# specgetty Release Script
# Automates semantic versioning releases with changelog updates
# =============================================================================

# --- Colors and Styles ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# --- Helper Functions ---
info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
    exit 1
}

# --- Portable File Editing ---
# BSD sed requires an argument to -i and GNU sed must not have one, so no single
# `sed -i` invocation works on both. Write to a temp file and move it instead.
sed_inplace() {
    local expr="$1"
    shift
    local f tmp
    for f in "$@"; do
        tmp="$(mktemp)"
        # Write back through the original file rather than moving the temp over
        # it: mktemp creates 0600, and a move would carry that mode onto a file
        # that was 0644. git only tracks the executable bit, so the tightened
        # permissions would not show up in a diff.
        sed "$expr" "$f" > "$tmp" && cat "$tmp" > "$f"
        rm -f "$tmp"
    done
}

# insert_version_heading adds "## [X.Y.Z] - DATE" below the [Unreleased] heading.
# This is awk rather than sed because a "\n" in a sed replacement is a GNU
# extension: BSD sed inserts a literal "n", and this project ships darwin builds.
insert_version_heading() {
    local version="$1" today="$2" tmp
    tmp="$(mktemp)"
    awk -v v="$version" -v d="$today" '
        !inserted && /^## \[Unreleased\]/ {
            print
            print ""
            print "## [" v "] - " d
            inserted = 1
            next
        }
        { print }
    ' CHANGELOG.md > "$tmp" && cat "$tmp" > CHANGELOG.md
    rm -f "$tmp"
}

# --- Dependency Check ---
if ! command -v gum &> /dev/null; then
    error "gum is not installed. Install it with:

    brew install gum          # macOS
    nix-env -iA nixpkgs.gum   # Nix
    go install github.com/charmbracelet/gum@latest  # Go

See https://github.com/charmbracelet/gum for more options."
fi

# --- VCS Detection ---
# jj owns the working copy wherever it is present, including a colocated
# workspace that also has a .git. Git HEAD is detached in that case, so pushing
# HEAD would push the empty working-copy commit rather than the release.
if command -v jj &> /dev/null && jj root &> /dev/null; then
    VCS="jj"
else
    VCS="git"
fi
info "Version control: ${BOLD}$VCS${NC}"

# --- Safety Checks ---
info "Running safety checks..."

# An untracked file that belongs in the release must stop the release, so both
# checks below cover untracked files as well as modified ones. `git diff` alone
# does not.
if [[ "$VCS" == "jj" ]]; then
    # jj snapshots the whole working copy into @, so its diff already includes
    # files git would call untracked.
    if [[ -n "$(jj diff --summary 2>/dev/null)" ]]; then
        error "Working copy is not clean. Commit or abandon changes first."
    fi
else
    if [[ -n "$(git status --porcelain)" ]]; then
        error "Working tree is not clean. Commit, stash or ignore changes first."
    fi
fi
success "Working directory is clean"

# Check CHANGELOG.md exists and contains [Unreleased]
if [[ ! -f "CHANGELOG.md" ]]; then
    error "CHANGELOG.md not found"
fi

if ! grep -q "\[Unreleased\]" CHANGELOG.md; then
    error "CHANGELOG.md does not contain [Unreleased] section"
fi
success "CHANGELOG.md has [Unreleased] section"

# --- Read Current Version ---
VERSION_FILE="src/VERSION"
if [[ ! -f "$VERSION_FILE" ]]; then
    error "VERSION file not found at $VERSION_FILE"
fi

CURRENT_VERSION=$(cat "$VERSION_FILE" | tr -d '\n')
info "Current version: ${BOLD}$CURRENT_VERSION${NC}"

# Parse version components
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# --- Version Selection ---
echo ""
echo -e "${BOLD}Select version bump type:${NC}"
BUMP_TYPE=$(gum choose "patch" "minor" "major")

# Calculate new version
case $BUMP_TYPE in
    major)
        NEW_VERSION="$((MAJOR + 1)).0.0"
        ;;
    minor)
        NEW_VERSION="$MAJOR.$((MINOR + 1)).0"
        ;;
    patch)
        NEW_VERSION="$MAJOR.$MINOR.$((PATCH + 1))"
        ;;
esac

info "New version will be: ${BOLD}$NEW_VERSION${NC}"

# Check if tag already exists
if [[ "$VCS" == "jj" ]]; then
    if jj tag list 2>/dev/null | grep -q "^v$NEW_VERSION:"; then
        error "Tag v$NEW_VERSION already exists"
    fi
else
    if git tag -l "v$NEW_VERSION" | grep -q "v$NEW_VERSION"; then
        error "Tag v$NEW_VERSION already exists"
    fi
fi
success "Tag v$NEW_VERSION is available"

# --- Confirmation ---
echo ""
echo -e "${BOLD}Release Summary:${NC}"
echo -e "  Current version: $CURRENT_VERSION"
echo -e "  New version:     $NEW_VERSION"
echo -e "  Bump type:       $BUMP_TYPE"
echo ""

if ! gum confirm "Proceed with release?"; then
    warn "Release cancelled"
    exit 0
fi

# --- Nix vendorHash Update ---
update_nix_vendor_hash() {
    if ! command -v nix &> /dev/null; then
        warn "nix is not installed — skipping vendorHash update in flake.nix"
        return 0
    fi

    info "Updating vendorHash in package.nix and flake.nix..."

    # Save the current vendorHash
    OLD_HASH=$(grep 'vendorHash' package.nix | sed 's/.*"\(.*\)".*/\1/')

    # Temporarily set vendorHash to empty to force nix to compute the correct one
    sed_inplace "s|vendorHash = \".*\"|vendorHash = \"\"|" package.nix

    # Run nix build and capture the expected hash from the error output
    NIX_OUTPUT=$(nix build .#default 2>&1 || true)
    NEW_HASH=$(echo "$NIX_OUTPUT" | grep "got:" | sed 's/.*got: *//')

    if [[ -z "$NEW_HASH" ]]; then
        # Restore old hash if we couldn't determine the new one
        sed_inplace "s|vendorHash = \".*\"|vendorHash = \"$OLD_HASH\"|" package.nix flake.nix
        warn "Could not determine new vendorHash — restored previous hash"
        return 0
    fi

    # Both files declare a vendorHash: package.nix for the binary, flake.nix for
    # the checks derivation that runs the test suite. If they drift,
    # `nix flake check` fails and every ship-change.sh run is blocked.
    sed_inplace "s|vendorHash = \".*\"|vendorHash = \"$NEW_HASH\"|" package.nix flake.nix
    success "Updated vendorHash to $NEW_HASH"
}

# --- Execute Release ---
echo ""
info "Creating release..."

# Update VERSION file
echo "$NEW_VERSION" > "$VERSION_FILE"
success "Updated VERSION file"

# Update CHANGELOG.md
TODAY=$(date +%Y-%m-%d)
insert_version_heading "$NEW_VERSION" "$TODAY"
success "Updated CHANGELOG.md with version $NEW_VERSION"

# Update nix flake vendorHash
update_nix_vendor_hash

# Determine the target branch before committing: under jj it names the bookmark
# to move, and under git it is the push target.
MAIN_BRANCH=""
for candidate in main master; do
    if git rev-parse --verify "refs/heads/$candidate" &>/dev/null; then
        MAIN_BRANCH="$candidate"
        break
    fi
    if git ls-remote --heads origin "$candidate" 2>/dev/null | grep -q "$candidate"; then
        MAIN_BRANCH="$candidate"
        break
    fi
done

if [[ -z "$MAIN_BRANCH" ]]; then
    error "Could not determine main branch (tried 'main' and 'master')"
fi

if [[ "$VCS" == "jj" ]]; then
    # jj commits the whole working copy, and the safety check above proved the
    # only changes present are the ones this script just made.
    jj commit -m "chore: release v$NEW_VERSION"
    success "Created release commit"

    # The release commit is @- now, because jj commit leaves a fresh empty @.
    jj bookmark set "$MAIN_BRANCH" -r @-
    success "Moved bookmark $MAIN_BRANCH"

    jj tag set "v$NEW_VERSION" -r @-
    success "Created tag v$NEW_VERSION"

    info "Pushing to remote..."
    # A bookmark that is not tracking its remote counterpart cannot be pushed.
    # Tracking an already-tracked bookmark warns and changes nothing.
    jj bookmark track "$MAIN_BRANCH" --remote=origin 2>/dev/null || true

    # Push the bookmark FIRST. If it fails, stop before the tag is published,
    # otherwise the remote is left with a tag pointing at a commit no branch
    # reaches, and the release workflow builds from a commit nobody can find.
    jj git push --bookmark "$MAIN_BRANCH"

    # jj git push carries bookmarks and never tags. A colocated workspace
    # exports refs to git automatically, so the tag is already in .git and
    # plain git can push it.
    git push origin "v$NEW_VERSION"
    success "Pushed commit and tag to remote ($MAIN_BRANCH)"
else
    git add "$VERSION_FILE" CHANGELOG.md package.nix flake.nix
    git commit -m "chore: release v$NEW_VERSION"
    success "Created release commit"

    git tag -a "v$NEW_VERSION" -m "Release v$NEW_VERSION"
    success "Created tag v$NEW_VERSION"

    info "Pushing to remote..."
    git push origin HEAD:refs/heads/$MAIN_BRANCH
    git push origin "v$NEW_VERSION"
    success "Pushed commit and tag to remote ($MAIN_BRANCH)"
fi

# --- Done ---
echo ""
echo -e "${GREEN}${BOLD}Release v$NEW_VERSION complete!${NC}"
echo ""
echo "GitHub Actions will now build and publish the release."
echo "Check the progress at: https://github.com/mipmip/specgetty/actions"
