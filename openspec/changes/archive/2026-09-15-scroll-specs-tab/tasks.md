## 1. Focus model

- [x] 1.1 Add a `specsFocus` field to `model` with the spec list as its zero value, so switching to the tab focuses the list
- [x] 1.2 Reset it to the list whenever the detail tab changes, so leaving and returning starts from the list
- [x] 1.3 On the specs tab, make `tab` cycle spec list, spec content, and the log panel when the log is open; leave `tab` unchanged everywhere else
- [x] 1.4 Leave the focus on the list when the project has no specs, so the keys cannot land on a pane that is not there

## 2. The specs content becomes a document

- [x] 2.1 Extract the list and content width split from `renderSpecsTab` into one helper, so the renderer and `docRegion` cannot disagree about the content width
- [x] 2.2 Teach `docRegion` the specs case: the content half's width, and the tab's content height
- [x] 2.3 Teach `docActive` that the specs tab owns the vertical axis only when the content holds the keyboard
- [x] 2.4 Teach `currentDocument` the specs case, keyed by project and spec name so selecting another spec starts at the top
- [x] 2.5 Render the content half through `m.docViewport.View()` instead of `truncateContent`
- [x] 2.6 Keep "No spec.md found" working for a spec directory without a file

## 3. Showing the focus

- [x] 3.1 Light the spec list cursor when the list is focused and dim it when the content is focused
- [x] 3.2 Confirm the panel title shows a percentage only while the content is focused, which falls out of `docActive`
- [x] 3.3 Update the nav bar on the specs tab: `tab` to move focus, and the scroll keys when the content holds them

## 4. Tests

- [x] 4.1 Focus starts on the list, and returns to the list after leaving the tab and coming back
- [x] 4.2 `tab` cycles list, content, and list again with the log closed
- [x] 4.3 `tab` cycles list, content, log with the log open
- [x] 4.4 `j` and `k` move the spec cursor with the list focused, and scroll the document with the content focused
- [x] 4.5 Every row of a long spec is reachable once the content is focused
- [x] 4.6 Selecting a different spec shows it from its first row
- [x] 4.7 The title reports a position only while the content is focused
- [x] 4.8 Spec content is wrapped to the content half, not to the full panel width
- [x] 4.9 A project with no specs keeps the focus on the list and renders without error
- [x] 4.10 Run `scripts/coverage-gate.sh` and hold the ratchet
