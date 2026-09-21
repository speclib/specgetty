## 1. Removing the filter modes

- [x] 1.1 Remove the three modes, their names and `emptyListMessage` from `changerow.go`, and take the mode parameter off `buildRows`. Proven by the package compiling with no reference to them left
- [x] 1.2 Remove `listMode` and `defaultMode` from the model, the `f` key handler, and the mode entry from the nav bar. Proven by asserting `f` does nothing on the change list and the nav bar names no mode
- [x] 1.3 Remove `ResolveListMode` and `Config.ChangeMode`. Proven by the packages compiling and by the field being absent from the parsed config struct
- [x] 1.4 Remove the `--change-mode` flag, and assert `spg --change-mode=active` reports an unrecognised flag, the way `--zoom` already does
- [x] 1.5 Remove the documented `change_mode` block from `src/config.yml`
- [x] 1.6 Report a configuration that still carries `change_mode` as no longer having an effect, and start normally. Proven by a config fixture with the key and one without, asserting the report appears only for the first

## 2. Grouping

- [x] 2.1 Produce the list as two groups, active first, each carrying a label and the changes in it. Proven by a project fixture with both, asserting order and membership
- [x] 2.2 Order the active group by name. Proven by a fixture whose directories are created out of alphabetical order
- [x] 2.3 Order the archived group by archive date, most recent first. Proven by a fixture spanning several dates, and against this project, asserting `2026-09-21-list-the-repos-not-the-stores` is the first archived row rather than the last
- [x] 2.4 Order archived changes sharing a date stably, so two renders agree. Proven by building the list twice and comparing
- [x] 2.5 Place an archived change whose directory carries no parsable date after those that do, rather than dropping it or letting it float. Proven by a fixture with an undated archive directory
- [x] 2.6 Keep a group and its header when it holds nothing, with a count of zero. Proven for an empty active group, an empty archived group, and a project with neither

## 3. Drawing it

- [x] 3.1 Draw one table with one column header row and a group header between the groups, reusing `layoutFields` and `fitCell` without changing `renderTable`. Proven by asserting the project picker's rendering is byte-identical to before this change
- [x] 3.2 Style group headers with the existing section header style, so the tab introduces no new visual vocabulary
- [x] 3.3 Count the rows a group currently holds in its header, which under a filter is what the filter left. Proven filtered and unfiltered
- [x] 3.4 Map a change's index to the line it is drawn on, counting headers. Proven by a table test over a list long enough to scroll, asserting the line of the first and last change in each group
- [x] 3.5 Scroll so the selected change stays visible with headers counted as the rows they occupy. Proven with the cursor immediately before and immediately after a group boundary, at a panel height that forces scrolling. This is where every off-by-one in this change would live
- [x] 3.6 Keep the frame exactly its terminal's rows and columns at three widths including the 60-column minimum

## 4. Moving through it

- [x] 4.1 Move the cursor between changes only, never resting on a header. Proven by stepping from the last active change to the first archived one and asserting the selection, not the line
- [x] 4.2 Confirm stepping to either end of the list stops on a change rather than a header. The task originally said `g` and `G` should reach the first and last change; they are document keys and have never been bound on the change list, so there was no jump for a header to catch and binding them would be behaviour this change was not asked for
- [x] 4.3 Keep the selected change across a filter, a rescan and a project switch, as it is kept today. Proven by the existing selection tests passing unchanged
- [x] 4.4 Open the change under the cursor with `enter` from either group

## 5. Columns

- [x] 5.1 Add the archive date to the default columns, giving `name,tasks,specs,date`. Proven by asserting the default rendering carries a date column
- [x] 5.2 Leave the date blank on an active row
- [x] 5.3 Remove the state column from the defaults while keeping it configurable. Proven by asserting it is absent by default and present when named in `change_fields`
- [x] 5.4 Confirm the name column stays above the width at which columns are dropped, with four defaults at the 60-column minimum

## 6. Actions

- [x] 6.1 Keep archive, discard and export acting on the change under the cursor, in either group. Proven by triggering each with the cursor in the archived group
- [x] 6.2 Keep archive a no-op on an already archived change
- [x] 6.3 Show a change in the archived group, at the top of it, after archiving it. Proven end to end against a fixture

## 7. Verification

- [x] 7.1 Open this project with the built binary and confirm one table, `ACTIVE (1)` above `ARCHIVED (35)`, dates on archived rows only, and `list-the-repos-not-the-stores` as the first archived row
- [x] 7.2 Confirm `f` does nothing and the nav bar names no mode
- [x] 7.3 Confirm a config carrying `change_mode` reports it once and starts
- [x] 7.4 Revert task 2.3 and confirm the archive order test fails. Check the build succeeds first
- [x] 7.5 Revert task 3.5 and confirm the scroll test fails
- [x] 7.6 Confirm the project picker is unchanged: same rows, same names, same columns, same scrolling
- [x] 7.7 `nix flake check` passes, coverage floors included

## 8. Notes

- [x] 8.0 Fix a crash found while covering the configuration paths, unrelated to grouping but not shippable alongside it: `spg --config <broken.yml> <dir>` dereferenced a nil config, because a parse failure leaves it nil while directory arguments take the branch that writes into it. Two lines, with a test each for the crashing case and the working one

- [x] 8.1 Update the README where it describes the change list: the `f` key, the three modes, and the archive sharing the list are all documented there
- [x] 8.2 Note on `specgetty-vru8` that sorting was considered here and deliberately left to it, together with the reason a sort needs a key per field rather than the rendered value
- [x] 8.3 `f` is now unbound. Left so rather than finding it a job in this change
