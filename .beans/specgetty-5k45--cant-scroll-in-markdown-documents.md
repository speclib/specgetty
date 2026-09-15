---
# specgetty-5k45
title: can't scroll in markdown documents
status: in-progress
type: bug
priority: normal
created_at: 2026-09-15T15:54:19Z
updated_at: 2026-09-15T17:55:18Z
---

cursors and vim keys should work and page down / up

Proposal written: `openspec/changes/scroll-markdown-documents`.

Two defects behind this, not one:

1. No scroll offset anywhere. Every content pane ends in truncateContent (src/ui/ui.go:1563), which cuts the line slice at the panel height. At levelChange the vertical keys are excluded by `m.level != levelChange` (src/ui/ui.go:535, 562) and the page and jump keys move m.fileCursor, a leftover from a file listing with no renderer.
2. renderMarkdown takes a width and never uses it (src/ui/ui.go:1447). The panel box wraps the already-truncated string afterwards (lipgloss style.go:368) and MaxHeight cuts the overflow, so a long paragraph pushes its neighbours off screen. Content is lost inside the visible height too.

Decided: full pages in documents only, lists keep their half page; percentage in the panel title, no scrollbar; no mouse wheel (bubbletea capture would take native text selection away). Specs tab deferred to specgetty-ze7f.

Overlaps with the active project-picker change, which deletes fileCursor and renumbers the level constants. Sequence them rather than running both at once.
