## 1. Navigation levels

- [x] 1.1 Add `level int` to `model` (0 projects, 1 project, 2 change), replacing the `zoomed bool`; every existing `zoomed` truth test becomes `level >= 1`
- [x] 1.2 Rewrite the `enter` handler: descend one level, clamped at 2; entering a change records the selected change
- [x] 1.3 Rewrite the `esc` handler: ascend one level, clamped at 0, leaving behaviour above the project view exactly as it is today (`specgetty-jdif` owns that end)
- [x] 1.4 Remove `enter` as an unzoom key (ascending is `esc` only)
- [x] 1.5 Remove the left/right fall-through at `src/ui/ui.go:428` and its mirror in the `left` handler, so artifact sub-tabs cannot spill into the project tab bar
- [x] 1.6 Replace `m.zoomed` reads in `View()`, `renderNavBar()`, the `s` rescan handler and the `tab` handler with `m.level` comparisons
- [x] 1.7 Keep `activeView` for focus between projects, detail and log at level 0; at level 1 and above focus is implicit

## 2. Merged change list

- [x] 2.1 Add a `changeRow` type in `ui`: the `scanner.ChangeInfo`, an `archived bool`, and a `matchedFiles []string` for search hints
- [x] 2.2 Add `listMode` constants (`modeOpen`, `modeArchived`, `modeBoth`) and a `listMode` model field defaulting to `modeOpen`
- [x] 2.3 Rewrite `currentChanges()` to build `[]changeRow` from `Info.Changes` and `Info.ArchivedChanges`, tagging `archived` from which slice a row came from, then apply the mode filter and the query filter
- [x] 2.4 Delete `currentArchivedChanges()` and repoint its six call sites at `currentChanges()`
- [x] 2.5 Collapse `archiveCursor` and `archiveArtifactTab` into `changeCursor` and `changeArtifactTab`; delete `tabArchive` and its entry in `tabNames`
- [x] 2.6 Choose and document a key that cycles the list mode, and render the current mode in the tab area (`tab` is unbound at the change list today; `a d e s l q /` are taken)
- [x] 2.7 Make the archive action a no-op when the row under the cursor is archived
- [x] 2.8 Render per-mode empty states: "No active changes", "No archived changes", and a both-mode equivalent

## 3. Field selection

- [x] 3.1 Define the field id set and a `defaultFields` slice of `name`, `tasks`, `specs`
- [x] 3.2 Add `change_fields` to `scanner.Config` and to the commented example in `src/config.yml`
- [x] 3.3 Add a `--change-fields` string flag in `src/main.go` and pass it through `ui.Run`
- [x] 3.4 Resolve the field list as flag over config over default, and report unknown field names with the list of valid ones rather than failing silently
- [x] 3.5 Add the `archived` indicator as a field so both-mode rows are distinguishable

## 4. Table renderer

- [x] 4.1 Extract a `renderChangeTable(rows []changeRow, fields []string, cursor int, width, height int) string` from the left half of `renderChangeList`
- [x] 4.2 Render a column header row above the rows
- [x] 4.3 Allocate widths: fixed-width columns first, `name` takes the remainder with ellipsis truncation; when the remainder drops below a minimum, drop columns from the right
- [x] 4.4 Highlight the cursor row across the full width
- [x] 4.5 Keep the existing scroll offset behaviour so the cursor stays visible
- [x] 4.6 Append the dim `matchedFiles` hint to a row when search matched it on body text

## 5. Single change view

- [x] 5.1 Extract a `renderChangeDetail(row changeRow, artifactTab int, width, height int) string` from the right half of `renderChangeList`, rendering at full panel width
- [x] 5.2 Delete `renderChangeList`, `renderChangesTab` and `renderArchiveTab`, and dispatch on `m.level` instead
- [x] 5.3 Show the change name and its archived state in the change view header
- [x] 5.4 Reset `changeArtifactTab` to 0 when a different change is opened

## 6. Query matching

- [x] 6.1 Promote `github.com/sahilm/fuzzy` to a direct dependency in `go.mod` (already present transitively through `bubbles`)
- [x] 6.2 Write `parseQuery(q string)` returning the matcher kind and the bare term: `'` prefix is literal on name, `:` prefix is literal on body, otherwise fuzzy on name
- [x] 6.3 Implement smart case: fold both sides when the term has no uppercase, compare as-is otherwise
- [x] 6.4 Implement fuzzy name matching with `fuzzy.Find` so results come back ranked, strongest first
- [x] 6.5 Implement literal name matching, preserving alphabetical order (no score to rank by)
- [x] 6.6 Implement body matching over `ArtifactContents` and `SpecContents` plus the change name, recording which filenames hit into `matchedFiles`
- [x] 6.7 Return the unmodified list in normal order when the query is empty

## 7. Search prompt

- [x] 7.1 Add `bubbles/textinput` to the model as `searchInput`, plus a `searchActive bool` and the committed `query string`
- [x] 7.2 Handle `/` at the change list to focus the input; while focused, route all rune keys to it and re-filter on every change
- [x] 7.3 Route `up`, `down`, `ctrl+p` and `ctrl+n` to the list cursor while the input is focused, leaving the query untouched
- [x] 7.4 Handle `enter` while focused: open the change under the cursor, keeping the filter applied underneath
- [x] 7.5 Handle `esc` while focused: clear the query, unfocus, restore the full list
- [x] 7.6 Render the prompt whenever the query is non-empty, even after opening a change and returning
- [x] 7.7 Render `No changes match "<query>"` when a non-empty query matches nothing in the current mode

## 8. Filter lifetime and cursor stability

- [x] 8.1 Add `selectedKey` to the model, keyed on the archived flag plus the change name, since an archived change keeps the name it had while active
- [x] 8.2 After every re-filter and every rescan, restore the cursor to `selectedKey` if that row survives, otherwise clamp to the nearest valid row
- [x] 8.3 Update `selectedKey` whenever the user moves the cursor
- [x] 8.4 Preserve the query across opening a change and returning, and across `fsChangeMsg` rescans
- [x] 8.5 Clear the query, reset `listMode` and reset `selectedKey` in `updateFileList()`, which already runs on project switch

## 9. Keys and nav bar

- [x] 9.1 Remap the number keys to 1-3 for changes, specs and config
- [x] 9.2 Make number keys inert while a change is open so they cannot change the project tab bar
- [x] 9.3 Build the nav bar per level: the projects level keeps its current hints, the change list adds `/ search`, the mode-cycle key and the change actions, an open change shows only sub-tab and `esc` hints
- [x] 9.4 Suppress the whole nav bar key set while the search input is focused, showing the prompt hints instead

## 10. Tests

- [x] 10.1 `parseQuery` table test covering both sigils, a bare term and an empty query
- [x] 10.2 Smart case test: lowercase term matches mixed-case names, term with an uppercase does not
- [x] 10.3 Fuzzy ranking test asserting the strongest match sorts first
- [x] 10.4 Body match test asserting `matchedFiles` names the artifacts that hit
- [x] 10.5 Mode filter test for open, archived and both
- [x] 10.6 Cursor stability test: filter a list so the selected change is removed, assert the cursor clamps; filter so it survives, assert the cursor follows it
- [x] 10.7 Level navigation test through `Update`: `enter` descends, `esc` ascends, and neither runs past the ends of the range
- [x] 10.8 Regression test asserting right at the last artifact sub-tab leaves `detailTab` unchanged
- [x] 10.9 Field resolution test: flag beats config beats default, unknown name reports an error
- [x] 10.10 Column width test asserting columns drop from the right on a narrow panel

## 11. Docs

- [x] 11.1 Update the key table in `README.md` for the level model, `/` search and the mode-cycle key
- [x] 11.2 Document `change_fields` and `--change-fields` in `README.md`
- [x] 11.3 Note in `README.md` that the archive tab is gone and archived changes live behind the list filter
