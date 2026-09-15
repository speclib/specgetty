## Why

`scroll-markdown-documents` made the change artifact pane and the config tab
scrollable and left the specs tab out, because its two halves compete for `j`
and `k` and settling that needs a focus model. The specs tab is where the
longest documents in the tool are read: the largest `spec.md` in this repository
is 6299 bytes, and the pane shows whatever fits and no more.

The wrapping half of that change already applies here, so nothing is lost inside
the visible height any more. What is still unreachable is everything below the
fold.

## What Changes

- The specs tab gains a focus: either the spec list or the spec content has the
  keyboard, and `tab` moves between them.
- With the content focused, the spec is a document viewer with the scrolling,
  position reporting and position retention that capability already describes.
- With the list focused, `j` and `k` move the spec cursor exactly as they do
  today.
- Which half has the keyboard is visible, so the keys are never a guess.
- On the specs tab, `tab` cycles spec list, spec content, and the log panel when
  it is open, rather than toggling the log alone.

## Capabilities

### Modified Capabilities
- `specs-tab`: gains a focus model, and its content pane becomes a document
  viewer rather than a truncated block

## Impact

- `src/ui/ui.go`: `renderSpecsTab` renders its content half through the
  viewport; the `tab` handler cycles focus on this tab; the vertical keys ask
  which half is focused
- `src/ui/docview.go`: `docActive` and `currentDocument` gain the specs case,
  and `docRegion` gains its narrower content width
- No new dependencies
