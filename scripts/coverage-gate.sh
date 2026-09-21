#!/usr/bin/env bash
#
# Coverage ratchet.
#
# The floors below are the measured coverage at the moment the gate was
# introduced (2026-09-15) and raised as coverage improved. They may only ever
# be raised, never lowered. A change
# that drops a package below its floor fails the gate, so coverage cannot erode
# quietly.
#
# The target is 70% overall and 80% for the core packages (scanner and ui).
# Every package now clears it. These floors track what is actually achieved, so
# they can only rise.
#
# Run it by hand with: bash scripts/coverage-gate.sh
#
# WARNING when raising a floor: coverage is not identical in every environment.
# src/watcher measures 84.6% on a normal machine but 87.2% inside the nix
# sandbox, because its inotify paths behave differently there. A floor must
# therefore track the LOWEST observed value, not the one your last run printed,
# or `nix flake check` and a local run will disagree about whether the gate
# passes. Check both before raising anything.

set -euo pipefail

floor_for() {
  case "$1" in
    # Raised from 42.3 when the CLI action moved out of main() into runApp,
    # which made everything up to the terminal takeover reachable from a test.
    github.com/mipmip/specgetty/src)         echo "87.7" ;;
    github.com/mipmip/specgetty/src/scanner) echo "95.4" ;;
    github.com/mipmip/specgetty/src/ui)      echo "87.6" ;;
    # watcher measures anywhere from 84.6% to 90.2% depending on how its
    # inotify paths fall, so its floor tracks the lowest reading seen.
    github.com/mipmip/specgetty/src/watcher) echo "84.6" ;;
    # Measured at 88.2 in both environments. Held a fraction lower because
    # watcher's flutter moves the total by a few tenths and a zero-margin
    # total would fail at random rather than for a real regression.
    TOTAL)                                   echo "89.2" ;;
    *)                                       echo "" ;;
  esac
}

profile="$(mktemp -t specgetty-cover.XXXXXX)"
trap 'rm -f "$profile"' EXIT

echo "==> go vet"
go vet ./...

echo "==> go test with coverage"
if ! test_output="$(go test -coverprofile="$profile" -covermode=set ./... 2>&1)"; then
  echo "$test_output"
  echo
  echo "Tests FAILED. Fix them before the coverage ratchet is even considered."
  exit 1
fi
echo "$test_output"

status=0

check() {
  local label="$1" actual="$2" floor="$3"
  if [ -z "$floor" ]; then
    printf '  %-45s %6s%%   (no floor set, new package)\n' "$label" "$actual"
    return
  fi
  if awk "BEGIN{exit !($actual >= $floor)}"; then
    printf '  %-45s %6s%%   >= %s%% ok\n' "$label" "$actual" "$floor"
  else
    printf '  %-45s %6s%%   <  %s%% FAILED\n' "$label" "$actual" "$floor"
    status=1
  fi
}

echo
echo "==> coverage ratchet"

while read -r pkg cov; do
  [ -n "$pkg" ] || continue
  check "$pkg" "$cov" "$(floor_for "$pkg")"
done < <(echo "$test_output" \
  | sed -n 's|^ok[[:space:]]\{1,\}\([^[:space:]]\{1,\}\).*coverage: \([0-9.]\{1,\}\)% of statements.*|\1 \2|p')

total="$(go tool cover -func="$profile" | awk '$1=="total:"{sub(/%/,"",$3); print $3}')"
check "TOTAL" "$total" "$(floor_for TOTAL)"

echo
if [ "$status" -ne 0 ]; then
  echo "Coverage gate FAILED. Add tests, or raise nothing: the floors only go up."
  exit 1
fi
echo "Coverage gate passed."
