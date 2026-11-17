# Proposal

## Why

OpenSpec allows a capability path with more than one segment, and its own
artifact instructions use `identity/user-auth` as the example. specgetty reads
one level under `specs/`, so a project organized that way shows its grouping
directories instead of its specs. In `chatcms-app`, 52 specs across 14
containers read as 14 capabilities, every one of them empty.

The same one-level read collapses change deltas. `hosted-until-connected` has
three, at `onboarding/import`, `tenancy/projects` and `publish/act`; specgetty
names each of them after its container, so they overwrite one another and only
one survives.

There is a quieter consequence behind the listing: a delta is compared against
the live spec of the same name. A delta named `tenancy` finds no live spec
named `tenancy`, so no MODIFIED requirement in a nested project is ever marked
as changed.

## What Changes

- The scanner walks `specs/` recursively, for a project's main specs and for
  each change's deltas. A capability is a directory that directly holds a
  `spec.md`, named by its path relative to `specs/`, with `/` as the separator.
  This is the identity `openspec list --specs` already reports.
- A directory that holds no `spec.md` is no longer listed as a capability, which
  removes the container rows. "No spec.md found" stays, for a capability that is
  listed but whose file cannot be read.
- The specs tab renders capabilities under a header per first path segment,
  with the count, and the remaining path as the row. Capabilities with no
  parent segment stay bare rows at their sort position, so a project with no
  nesting looks as it does today.
- Paging the spec list counts drawn rows rather than capabilities, because a
  page of the pane can now include headers.
- A change delta and its live spec now agree on the capability name, so the
  MODIFIED marks in the change detail view appear in nested projects.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `openspec-scanning`: capability discovery becomes recursive, and the count of
  specs becomes the count of directories holding a `spec.md`.
- `specs-tab`: the spec list is grouped by first path segment instead of being
  a flat list of directory names, and a directory without `spec.md` is no
  longer listed.
- `change-spec-detail-view`: a delta's MODIFIED marks compare against the live
  spec of the same nested path.

## Impact

- `src/scanner/scan.go`: `parseChangeDir` and `ParseProjectInfo` share a
  recursive discovery helper; `SpecCount` counts what the helper returns.
- `src/ui/ui.go`: `renderSpecsTab` draws the grouped list; the spec list's page
  size comes from the drawn lines.
- `src/ui/listpaging.go`: the specs tab's page moves a pane of rows.
- `src/ui/docview.go`, `src/ui/changespec.go`, `src/ui/picker.go`: unchanged in
  behavior. They build paths and keys from the capability string, and
  `filepath.Join` already accepts a path with separators in it.
- Tests: `TestParseProjectInfo` currently asserts that a directory without
  `spec.md` is counted.
- No dependency changes. The scanner stays filesystem-based and does not shell
  out to `openspec`.
- Tracked by bean `specgetty-9xiu`.

## Non-goals

- Grouping the change's delta outline, which would need a group node kind in the
  spec tree. The deltas read as `tenancy/projects` there and are not confused
  with one another; only the main spec list has the container problem.
- Supporting a flat `specs/<name>.md` layout. Nothing writes one.
- Reading specs through the `openspec` CLI.
