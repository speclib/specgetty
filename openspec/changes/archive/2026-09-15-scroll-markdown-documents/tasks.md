## 1. Wrapping in the markdown renderer

- [x] 1.1 Promote `github.com/charmbracelet/x/ansi` from an indirect to a direct dependency in `go.mod`
- [x] 1.2 Make `renderMarkdown(content string, width int)` wrap each styled line to `width` with `ansi.Wordwrap`, after styling rather than before, so an inline `**span**` that straddles a boundary keeps its styling
- [x] 1.3 Do the same in `renderYAML`, which has the same unused `width` parameter
- [x] 1.4 Work out the exact content width `renderPanel` leaves after border and padding, and expose it so callers wrap to that number rather than guessing; a disagreement here makes the lipgloss box re-wrap and the row count drift
- [x] 1.5 Leave `Width()` on the panel box as a safety net, and assert in a test that it is a no-op for already-wrapped content
- [x] 1.6 Check what wrapping does to a long unbroken token such as a store path or a URL, and to a line containing double-width runes

## 2. The document viewer

- [x] 2.1 Add a `docViewport viewport.Model` and a `docKey string` to `model`, replacing the dead `detailViewport` (`src/ui/ui.go:151`)
- [x] 2.2 Add a helper that takes a document key and its rendered rows: on a changed key `SetContent` then `GotoTop`, on an unchanged key `SetContent` then clamp the offset to the new end
- [x] 2.3 Build the key from project path, tab, change name and artifact name, so every reset case is decided in one place
- [x] 2.4 Size `docViewport` in `recalcLayout` alongside `logViewport`, and re-wrap on resize before clamping (`src/ui/ui.go:720-735`)
- [x] 2.5 Add a predicate for whether a document viewer currently owns the vertical axis, so the key handlers have one thing to ask

## 3. Change artifact pane

- [x] 3.1 Replace the `truncateContent` call in `renderChangeDetail` with the viewport's `View()` (`src/ui/changelist.go:147`)
- [x] 3.2 Keep the change name line and the sub-tab row outside the scrolling region, and subtract them from the viewport height
- [x] 3.3 Render the tasks header (`Tasks: n/m complete`) as part of the document, so it scrolls away with the content it belongs to
- [x] 3.4 Reset to the top when the artifact sub-tab changes via left or right (`src/ui/ui.go:504-520`), which falls out of the key if it includes the artifact name

## 4. Config tab

- [x] 4.1 Render `renderConfigTab` content through the viewport (`src/ui/ui.go:1401-1425`)
- [x] 4.2 Keep the dimmed file source line fixed above the scrolling region
- [x] 4.3 Confirm the key includes the project path, so switching project resets and switching tab and back does not

## 5. Keys

- [x] 5.1 Route `up`, `k`, `down`, `j` to the viewport when a document viewer owns the axis, removing the `m.level != levelChange` exclusions (`src/ui/ui.go:535`, `src/ui/ui.go:562`)
- [x] 5.2 Route `pgdown`, `ctrl+f`, `pgup`, `ctrl+b` to a full page in a document viewer, leaving `halfPage()` untouched for the project list and the change list (`src/ui/ui.go:474-502`)
- [x] 5.3 Add `ctrl+d` and `ctrl+u` as half page in a document viewer only
- [x] 5.4 Route `gg` and `G` to `GotoTop` and `GotoBottom` in a document viewer, instead of the `fileCursor` branches (`src/ui/ui.go:352-366`, `src/ui/ui.go:459-473`)
- [x] 5.5 Update `renderNavBar` for `levelChange`: the current hint list offers no vertical key at all (`src/ui/ui.go:1640-1650`)
- [x] 5.6 Check that none of the new bindings collide with the search prompt, the confirm modals or the log panel

## 6. Position indicator

- [x] 6.1 Give `renderPanel` a scroll position argument and append a percentage to the title when the document exceeds the pane (`src/ui/ui.go:1586-1613`)
- [x] 6.2 Show nothing when the content fits, so absence of the indicator means there is nothing below
- [x] 6.3 Account for the indicator in the top border width arithmetic, which currently subtracts only `lipgloss.Width(title)`
- [x] 6.4 Decide what happens when the project name plus the percentage exceed the panel width

## 7. Cleanup

- [x] 7.1 Remove the `fileCursor` and `filePaths` branches from the four key handlers this change rewrites, or leave them to `project-picker` if that lands first; do not leave both owners
- [x] 7.2 Check whether `truncateContent` still has callers after this change and remove it if not (`src/ui/ui.go:1563`)

## 8. Tests

- [x] 8.1 `renderMarkdown` wraps a long line to the given width, counted in display cells, not bytes
- [x] 8.2 A wrapped bold span keeps its styling on both rows
- [x] 8.3 A document of n rows in a pane of h rows exposes rows n-h to n after `GotoBottom`, with nothing dropped in between
- [x] 8.4 Key handling: one row, full page, half page and jump each move the offset by the expected amount and stop at the bounds
- [x] 8.5 The project list and the change list still move by a half page on `pgdown`, unchanged
- [x] 8.6 A changed document key resets to the top; an unchanged key with new content keeps the offset, clamped
- [x] 8.7 The title carries a percentage for a tall document and none for a short one
- [x] 8.8 Run `scripts/coverage-gate.sh` and hold the ratchet
