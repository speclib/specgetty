## Why

The changes tab crams a change list into 30% of the panel and the artifact
viewer into the remaining 70%. Neither gets enough room: the list cannot show
the fields that make a change identifiable at a glance, and the markdown reads
in a narrow column.

Worse, the two navigations share one axis. Pressing right past the last artifact
sub-tab does not stop; it falls through into the parent tab bar, switches you to
a different top-level tab, and resets your position (`src/ui/ui.go:428`). The
sub-nav of a change and the tab bar of a project are different things and should
not share left/right.

Separately, active changes and archived changes are two tabs showing the same
kind of list, so finding "that change about zip export, active or archived" means
looking in two places and reading both by eye.

## What Changes

- Give the change list its own full-width level. Selecting a change and pressing
  `enter` opens that change one level deeper, where its artifact sub-tabs
  (`proposal`, `design`, `tasks`, `specs`) own left/right without colliding with
  the project tab bar.
- **BREAKING (UI)**: merge the changes tab and the archive tab into one list with
  an open/archived/both filter. The archive tab is removed.
- Render the change list as a table of columns instead of a bare name list.
  Which columns appear is selectable through the config file and a
  `--change-fields=` option.
- Add a `/` search filter over the change list: fuzzy matching on change names,
  and a `:` sigil for literal substring search inside artifact and spec text.

Columns in this change are limited to what a scan already computes (name, task
progress, spec count). Columns that need new scanner capability (schema,
completeness against the schema, spec delta counts, created and updated dates)
are tracked separately and land behind the same `--change-fields` mechanism.

## Capabilities

### New Capabilities
- `change-list-view`: the full-width change list level, its columns, the
  open/archived/both filter, the drill-down into a single change, and the
  keys that move between levels
- `change-search`: the `/` filter over the change list, its matching rules,
  ordering, lifetime, and cursor behaviour

### Modified Capabilities
- `changes-tab`: the split-pane layout requirements are replaced by the
  full-width list; artifact sub-tab requirements move down a level
- `archive-tab`: removed as a separate tab, its behaviour folded into the
  open/archived/both filter
- `detail-tabs`: the tab set loses `archive`

## Impact

- `src/ui/ui.go`: navigation level state replaces the `zoomed` boolean plus
  ad-hoc cursors; `currentChanges()` becomes the single seam where the
  active/archived merge, the mode filter and the query filter apply;
  `renderChangeList` splits into a table renderer and an artifact renderer;
  the left/right fall-through at line 428 is removed
- `src/scanner/scan.go`: unchanged
- `src/config.yml` and CLI flags: new field selection setting
- `go.mod`: `github.com/sahilm/fuzzy` promoted from transitive (already present
  via `bubbles`) to a direct dependency; `bubbles/textinput` newly imported
