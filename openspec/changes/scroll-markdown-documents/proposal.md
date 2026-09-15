## Why

Markdown is not scrolled in specgetty, it is thrown away. Every content pane
ends in `truncateContent`, which cuts the line slice at the panel height and
pads (`src/ui/ui.go:1563`). No offset exists anywhere in the model. Open a
change whose `proposal.md` is longer than the panel and the rest of the document
is unreachable.

At the change level the vertical axis is inert on purpose. Up, down, `j` and `k`
are excluded by `m.level != levelChange` (`src/ui/ui.go:535`, `src/ui/ui.go:562`)
and `pgdown`, `pgup`, `gg` and `G` move `m.fileCursor` (`src/ui/ui.go:467-496`),
a leftover from a file listing that no longer has a renderer.

A second defect makes the first one worse. `renderMarkdown` takes a width and
never uses it (`src/ui/ui.go:1447`), so nothing wraps. The panel box then wraps
the already-truncated string at render time (lipgloss `style.go:368`) and cuts
the overflow with `MaxHeight`. One long paragraph pushes its neighbours off
screen, so content is lost inside the visible height as well. This has to be
fixed for scrolling to be truthful: a position computed over source lines does
not match rows on screen, and the reported percentage would lie.

`detailViewport` sits in the model (`src/ui/ui.go:151`) and is resized on every
layout pass (`src/ui/ui.go:728`). Its `View()` is never called. Only
`logViewport` scrolls today, and it is the working pattern this change follows.

## What Changes

- `renderMarkdown` wraps at the pane width and returns rendered rows, so the row
  count the renderer produces is the row count the terminal shows. This alone
  stops content being lost silently in every markdown pane, the specs tab
  included.
- The artifact pane of an open change becomes a scrollable document. Up, down,
  `j` and `k` move one row; `pgup`, `pgdown`, `ctrl+f` and `ctrl+b` move a full
  page; `ctrl+d` and `ctrl+u` move half; `gg` and `G` jump to the ends.
- The config tab content becomes a scrollable document under the same keys. It
  has no competing cursor, so nothing has to be given up for it.
- The panel title shows the scroll position as a percentage while the document
  is taller than the pane, and shows nothing while it fits.
- Scroll position resets to the top when the document under the pane changes,
  and survives a rescan of the same document, clamped to the new end.
- `pgdown`, `pgup`, `ctrl+f` and `ctrl+b` keep their current half-page step in
  the project list and the change list. Only documents get full pages.
- `m.detailViewport` is either wired up or deleted. It does not stay dead.

Non-goal: the specs tab content pane. Its left and right halves compete for
`j`/`k` and settling that needs a focus model. It gets the wrapping fix, not the
scrolling, and is tracked separately.

Non-goal: the mouse wheel. Enabling it in bubbletea captures click and drag as
well, which takes native text selection away from a tool people read markdown
in.

## Capabilities

### New Capabilities
- `document-viewer`: a markdown pane that renders a document wrapped to its
  width, scrolls it with the keyboard, reports its position, and decides when
  that position resets

### Modified Capabilities
- `change-list-view`: the artifact pane of an open change is a document viewer,
  so the vertical keys have meaning at that level
- `config-tab-display`: the config content is a document viewer

## Impact

- `src/ui/ui.go`: `renderMarkdown` gains real wrapping; `renderPanel` takes a
  scroll position for its title; the `up`/`down`/`j`/`k`, `pgdown`/`pgup`,
  `ctrl+f`/`ctrl+b` and `gg`/`G` handlers route to the document viewer when one
  is focused; `ctrl+d`/`ctrl+u` are new bindings; `m.detailViewport` is resolved
- `src/ui/changelist.go`: `renderChangeDetail` renders through the viewport
  instead of `truncateContent`
- `src/ui/` new file for the document viewer wrapper around `bubbles/viewport`
- `go.mod`: unchanged, `bubbles/viewport` is already a direct import
- Overlaps with the active `project-picker` change, which deletes `filePaths`
  and `fileCursor` and renumbers the level constants. Whichever lands second
  rebases onto the other; the seams are the same handlers
