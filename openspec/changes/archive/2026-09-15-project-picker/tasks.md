## 1. Remove the split layout

- [x] 1.1 Delete `renderProjectList`, `leftPanelWidth` and `rightPanelWidth`
- [x] 1.2 Delete the two-panel branch of `View()` so the detail panel always spans the full width, and drop the `viewProjects` case from `renderPanel`
- [x] 1.3 Delete the `viewProjects` constant; `activeView` becomes a two-state focus between the main view and the log panel
- [x] 1.4 Delete `filePaths` and `fileCursor`, `updateFileList`, and every branch that moves `fileCursor` (`G`, `pgup`/`pgdown`, the `default:` arm of `j`/`k`): nothing has rendered them since the tabs were introduced
- [x] 1.5 Rebind `tab`: with one panel there is nothing to switch to unless the log panel is visible
- [x] 1.6 Renumber the level constants to `levelProject = 0` and `levelChange = 1`, and make `esc` at `levelProject` do nothing

## 2. Startup and flags

- [x] 2.1 Replace `--zoom` with `--view`, accepting `single` and `all`, defaulting to `single`; reject other values with a message naming the valid ones
- [x] 2.2 Keep `--path` and treat it as implying `--view=single`
- [x] 2.3 Resolve the startup project from the working directory with `findOpenSpecProject`, and open it without walking the configured scan directories
- [x] 2.4 When no project resolves, prompt "No OpenSpec project here. Open the project picker?" and open the picker on yes
- [x] 2.5 When the prompt is declined, show the project view empty state rather than exiting
- [x] 2.6 Start the filesystem watcher on the resolved project at startup
- [x] 2.7 Rename `initialZoomPath` on the model and in `ui.Run` to match the new vocabulary

## 3. Project view

- [x] 3.1 Render the empty state naming the picker key when no project is selected
- [x] 3.2 Keep the persistent header (path plus statistics) above the tab bar
- [x] 3.3 Keep `s` scoped to rescanning the open project
- [x] 3.4 Keep the minimum terminal size message at 60x20

## 4. Project cache

- [x] 4.1 Define the cache file: `version`, `scanned_at`, the `scandirs` include and exclude lists it was built from, and the discovered `paths`
- [x] 4.2 Resolve its location under `os.UserCacheDir()` following the XDG Base Directory Specification
- [x] 4.3 Write the cache after any full walk
- [x] 4.4 Load the cache instead of walking when it exists and its recorded `scandirs` match the current config
- [x] 4.5 Treat an unreadable, malformed or version-mismatched cache as absent rather than failing
- [x] 4.6 Drop cached paths that no longer exist, using a `stat` per path rather than a walk
- [x] 4.7 Parse statistics from disk on every load; never serve them from the cache

## 5. Generalise the table and filter

- [x] 5.1 Introduce a row abstraction covering what both lists need: a stable key, searchable text sources, and field accessors
- [x] 5.2 Make `filterRows`, `indexOfKey` and `renderChangeTable` operate on that abstraction instead of `changeRow`
- [x] 5.3 Make the field table generic over the row type, keeping `changeRow`'s existing fields working unchanged
- [x] 5.4 Leave `parseQuery`, smart case, `fitCell` and `layoutFields` as they are: they are already row-agnostic
- [x] 5.5 Confirm the change list behaves identically afterwards, by the tests that already cover it

## 6. Picker overlay

- [x] 6.1 Add a `projectRow` type carrying the path, display name and parsed `ProjectInfo`
- [x] 6.2 Add picker state to the model: open flag, cursor, `selectedKey`, its own search input
- [x] 6.3 Bind `p` to open the picker from both levels, and `esc` to dismiss it without changing the current project
- [x] 6.4 Order the key interceptors: confirm modals, then the picker, then the change list search prompt, then ordinary keys
- [x] 6.5 Render the picker as a centred overlay over the current view
- [x] 6.6 Render its rows with name, specs, active changes, archived changes and aggregate task progress
- [x] 6.7 Reuse `projectDisplayNames` so duplicate basenames are disambiguated by parent directory
- [x] 6.8 Render an empty state naming the refresh key when no projects have been discovered

## 7. Picker selection and filter

- [x] 7.1 On `enter`, switch the current project, close the picker, and land at `levelProject` with the changes tab active
- [x] 7.2 Stop the watcher for the previous project before starting it for the new one
- [x] 7.3 Reset the change list state on switch: cursor, list mode and any active filter
- [x] 7.4 Bind `/` to the picker's own filter, using the shared query grammar
- [x] 7.5 Extend `:` matching for projects to cover file paths as well as file contents
- [x] 7.6 Render `No projects match "<query>"` echoing the query when nothing matches
- [x] 7.7 Keep the picker cursor on the same project across a re-filter, by key

## 8. Refresh

- [x] 8.1 Bind `r` in the picker to a full walk that rewrites the cache
- [x] 8.2 Show progress while a refresh runs, reusing the existing spinner
- [x] 8.3 Keep the picker open and restore the cursor to the same project when the refresh completes

## 9. Tests

- [x] 9.1 Cache round trip: write, read back, and assert the paths match
- [x] 9.2 Cache invalidated when the recorded `scandirs` differ from the config
- [x] 9.3 Malformed and version-mismatched cache files are treated as absent, not as errors
- [x] 9.4 A cached path that no longer exists is dropped on load
- [x] 9.5 Startup resolves the project at the working directory without walking
- [x] 9.6 `--view` rejects an unknown value and names the valid ones
- [x] 9.7 `--zoom` is no longer accepted
- [x] 9.8 The picker opens from both levels and dismisses without changing the current project
- [x] 9.9 Selecting a project while a change is open lands at the project view, not in a change
- [x] 9.10 Key precedence: a confirm modal outranks the picker, and the picker outranks the change list keys
- [x] 9.11 Picker filter covers each sigil, smart case, and the no-match message
- [x] 9.12 Picker cursor follows the same project across a re-filter
- [x] 9.13 The empty state renders when no project is selected and names the picker key
- [x] 9.14 `esc` at the project view neither exits nor changes the level
- [x] 9.15 The generalised table renders both row types

## 10. Docs

- [x] 10.1 Rewrite the README navigation section for two levels plus the picker overlay
- [x] 10.2 Document `--view`, note that `--zoom` is gone, and document `--path`
- [x] 10.3 Document the picker keys, its filter grammar and the refresh action
- [x] 10.4 Document the cache: where it lives, what invalidates it, and that new projects need a refresh
- [x] 10.5 Add a CHANGELOG entry covering the breaking changes
