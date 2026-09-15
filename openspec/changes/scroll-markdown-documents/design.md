## Context

Three facts about the current code shape this design.

`logViewport` already works. It is a `bubbles/viewport`, it is sized in
`recalcLayout`, its `View()` is rendered into a panel, and the key handlers call
`LineUp`, `LineDown`, `GotoTop` and `GotoBottom` on it. Everything this change
needs exists in the file already, applied to one panel.

`renderMarkdown` produces styled text, not plain text. Headers get
`mdHeaderStyle`, inline spans get `mdBoldStyle` and `mdItalicStyle`
(`src/ui/ui.go:1447-1509`). Any wrapping added to it has to count display cells,
not bytes, or an escape sequence gets counted as visible width and the rows come
out short.

Only one document pane is ever on screen. The config tab lives at
`levelProject`, the artifact pane at `levelChange`, and the two cannot both be
visible. One viewport serves both.

## Goals / Non-Goals

**Goals:**
- Every row of a markdown document is reachable in the change artifact pane and
  the config tab
- No row is lost inside the visible height, in any markdown pane, the specs tab
  included
- The reported position matches what is on screen
- No existing key changes meaning

**Non-Goals:**
- Scrolling the specs tab content pane. It needs a focus model first
- Mouse support
- Horizontal scrolling. Long lines wrap, they do not run off the right edge
- Improving the markdown renderer beyond wrapping. Fenced code blocks, tables
  and nested lists render the way they render today

## Decisions

### Wrap inside renderMarkdown, before the viewport

The current order is truncate, then wrap. `truncateContent` cuts at the panel
height while lines are still unwrapped, then the lipgloss box wraps what
survived and `MaxHeight` cuts the overflow. The visible result is fewer rows of
document than the height allows, and which rows disappear depends on paragraph
length.

Reverse it. `renderMarkdown(content, width)` wraps each styled line to `width`
and returns the rows the terminal will show. The viewport then slices rows that
are already final, so `YOffset` and `ScrollPercent` mean what they say.

Wrapping has to be ANSI-aware and it has to happen after styling, not before.
Wrapping the raw source first and styling each row afterwards would break an
inline `**span**` that straddles a boundary: the opening and closing markers
land on different rows and neither is recognised.
`github.com/charmbracelet/x/ansi` is already in `go.mod` as an indirect
dependency through lipgloss and provides `Wordwrap`. It gets promoted to a
direct dependency.

The width passed in must match the width of the content area after border and
padding. If the two disagree the lipgloss box re-wraps and the row
count drifts again. `renderPanel` keeps its `Width()` as a safety net, but in
normal operation it must be a no-op.

### One viewport, keyed by document identity

A single `docViewport` field replaces the dead `detailViewport`, plus a `docKey
string` holding an identity for the document currently loaded: project path, tab,
change name and artifact name joined together.

On every render pass the viewer computes the key for the document that should be
shown. Key changed means `SetContent` and `GotoTop`. Key unchanged means
`SetContent` and clamp the offset to the new end.

This one rule covers every reset case without a handler for each: moving between
artifact sub-tabs, opening a different change, switching tabs, switching
projects. It also gives the right behaviour when the file system watcher
rewrites the file under the reader, which is common in this tool: the document
is the same document, so the position stays where the eye is.

### Full pages in documents, half pages in lists

`halfPage()` currently returns half the panel height and serves `pgdown`,
`pgup`, `ctrl+f` and `ctrl+b` in every pane. Vim treats `ctrl+f` as a full page
and `ctrl+d` as half, and the bean asks for page up and page down.

Documents get the vim split: `pgup`, `pgdown`, `ctrl+f` and `ctrl+b` a full
page, `ctrl+d` and `ctrl+u` half. Lists keep `halfPage()` unchanged, so the
project list and the change list behave exactly as they do today.

The keys mean slightly different distances in a list and in a document. The
alternative was changing how the two lists have always moved, for a consistency
nobody has asked for. Nothing regresses this way.

### Percentage in the title, not a scrollbar

`renderPanel` already composes a title into the top border. It takes an
additional position argument and appends a percentage when the document exceeds
the pane. A scrollbar would cost a column of text width on every row, on a pane
whose whole problem is that documents do not fit.

The indicator is absent when everything fits, which is itself the signal that
there is nothing below.

### The specs tab is deferred, with its answer recorded

When the specs tab content pane becomes scrollable, `tab` toggles focus between
the spec list and the content, and the focused pane gets the active border
treatment `renderPanel` already applies. `tab` is free at that level: the
handler at `src/ui/ui.go:440` explicitly does nothing once inside a project.

This is written down so the next change does not reopen the question. It is not
implemented here.

## Risks / Trade-offs

- **Conflict with `project-picker`.** That change deletes `filePaths` and
  `fileCursor`, collapses `activeView` and renumbers the level constants. This
  change rewrites the same four key handlers. The seams are identical, so
  whichever lands second rebases onto the other. Worth sequencing rather than
  running both at once.
- **Resize moves the reader.** Re-wrapping at a new width changes how many rows
  a paragraph occupies, so an absolute row offset points somewhere slightly
  different afterwards. The offset is clamped, not remapped. Mapping the top
  visible source line back to a row after re-wrap would be exact, and is more
  machinery than a terminal resize deserves.
- **Wrapping changes how existing panes look.** The specs tab and the config tab
  currently let long lines run to the panel edge and get wrapped by lipgloss.
  After this change they wrap one step earlier, at the content width. Long
  unbroken tokens such as URLs and store paths will break differently.
- **Code fences wrap.** `renderMarkdown` has no notion of a fenced block, so a
  long line of code wraps like prose instead of being clipped. That is the same
  loss of alignment the specs tab shows today, made visible rather than cut off.
