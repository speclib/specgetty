---
# specgetty-knqz
title: an empty exclude entry panics the scan
status: completed
type: bug
priority: low
created_at: 2026-09-15T20:16:30Z
updated_at: 2026-09-18T09:25:50Z
---

An empty entry in `scandirs.exclude` panics the scan.

`src/scanner/find.go`:

    func skip(needle string, haystack []string) bool {
        for _, f := range haystack {
            if(f[0:1]=="/"){        // panics when f is ""

## Verified

    skip("/a/b/c", []string{""})
    -> runtime error: slice bounds out of range [:1] with length 0

## Trigger

`skip` is called for every directory the walker touches, with
`config.ScanDirs.Exclude` as the haystack. An empty string reaches it from a
config like

    exclude:
      -
      - .terraform

where the dash with nothing after it parses as an empty string. The panic then
happens inside a walk goroutine, which takes the program down.

## Shape of a fix

Skip empty entries, and use `strings.HasPrefix(f, "/")` rather than slicing, so
the comparison cannot go out of range regardless. A config-load-time warning for
an empty exclude entry would be a friendlier addition, but is not required to
stop the crash.

Not fixed in `cover-scanner-and-export`: pre-existing and outside that change's
scope.


## Summary of Changes

Shipped in `c2f6dcf`, OpenSpec change `survive-bad-scan-config`.

`skip` now passes over an empty haystack entry and tests the leading slash with
`strings.HasPrefix` rather than `f[0:1]`, so it cannot go out of range whatever
the entry holds.

The same bug sat on the other half of the config: `Walk` tested an include for a
trailing `*` with `globPath[len(globPath)-1:]`, which panics on an empty include
for the same reason. Fixed alongside it, since a YAML dash with nothing after it
reaches both lists the same way.

Covered by `TestSkipIgnoresEmptyExcludeEntries` and `TestWalkIgnoresAnEmptyInclude`.
Both were confirmed to fail with the old code restored and building.
