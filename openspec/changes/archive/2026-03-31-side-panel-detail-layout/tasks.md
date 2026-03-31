## 1. Layout restructure

- [x] 1.1 Add `leftPanelWidth()` and `rightPanelWidth()` methods to model (30% ratio, min 20, max 40 for left)
- [x] 1.2 Update `recalcLayout()` to compute widths for both panels and set viewport dimensions accordingly
- [x] 1.3 Update minimum terminal width check from height-only to also require >= 60 columns

## 2. Left panel: project basenames

- [x] 2.1 Add `projectDisplayNames()` method that returns basenames, with parent dir disambiguation for duplicates
- [x] 2.2 Update `renderRepoList()` to use display names instead of full paths, and use left panel width

## 3. Right panel: detail view

- [x] 3.1 Create `renderDetailPanel()` that shows full project path as header + file listing below
- [x] 3.2 Update `renderStatusContent()` to call `renderDetailPanel()` at the right panel width
- [x] 3.3 Show "No project selected." when project list is empty

## 4. Compose the layout

- [x] 4.1 Update `View()` to join left and right panels horizontally with `lipgloss.JoinHorizontal`, log panel below spanning full width
- [x] 4.2 Update `renderPanel()` to accept variable widths (left panel uses leftPanelWidth, right uses rightPanelWidth)
- [x] 4.3 Rename panel title from "Projects" to "Projects" (left) and "Contents" to "Detail" (right)

## 5. Tests and verification

- [x] 5.1 Update `TestRepoPanelHeight` to work with new layout dimensions
- [x] 5.2 Update `TestStatusPanelHeight` for horizontal layout
- [x] 5.3 Add test for `projectDisplayNames` with unique and duplicate basenames
- [x] 5.4 Verify build and all tests pass
