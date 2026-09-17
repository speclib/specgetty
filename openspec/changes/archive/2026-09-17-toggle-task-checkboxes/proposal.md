## Why

A change's tasks live in `tasks.md`, and specgetty already reads them, counts
them and shows the totals everywhere. Ticking one off means leaving the tool,
opening an editor, finding the line and saving. The tool that knows which task
you are looking at cannot mark it done.

## What Changes

- Checkbox lines render as `▢` and `▣` instead of `- [ ]` and `- [x]`
- The tasks pane gains a cursor. The selected source line is highlighted across
  every screen row it occupies, which matters because most task lines wrap
- `j`, `k` and the arrows move that cursor and bring it into view, instead of
  scrolling the pane by one row
- `space` toggles the checkbox on the selected line and saves immediately
- The save is atomic, and is applied to the file as it is on disk at that
  moment rather than to the copy specgetty last read

Only `tasks.md` gets this. Every other document keeps the behaviour it has.

Not in scope: colouring the checkbox differently from its text. That belongs
with a wider look at markdown highlighting, and this change is about behaviour.

## Capabilities

### New Capabilities
- `task-checkboxes`: how checkbox lines render, how the cursor selects a source
  line, what `space` does, and the rules the save follows so that it cannot
  discard someone else's edit

### Modified Capabilities
- `document-viewer`: moving by one row becomes moving by one source line in a
  document that has a cursor. Every other movement is unchanged

## Impact

- `src/ui/docview.go`: the renderer keeps a map from source line to the screen
  rows it occupies, which is what the highlight and the toggle both need
- `src/ui/ui.go`: `styleMarkdownLine` renders the glyphs; `j`/`k` consult the
  cursor; `space` is bound for the first time
- New file for the toggle and the atomic write
- `src/scanner/scan.go`: unchanged. The counts already come from a rescan, which
  the filesystem watcher triggers on our own write
