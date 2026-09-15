---
# specgetty-ze7f
title: specs tab content pane should scroll
status: in-progress
type: task
priority: normal
created_at: 2026-09-15T16:32:06Z
updated_at: 2026-09-15T18:45:37Z
blocked_by:
    - specgetty-5k45
---

The specs tab shows a spec list on the left and the spec content on the right (src/ui/ui.go:renderSpecsTab). The content pane truncates and does not scroll, because j/k already moves the spec list cursor.

Deferred out of the openspec change `scroll-markdown-documents`, which makes the change artifact pane and the config tab scrollable and gives every markdown pane real wrapping. The specs tab gets the wrapping fix from that change but not the scrolling.

The focus-axis question is already decided: `tab` toggles focus between the spec list and the content pane, and the focused pane gets the active border treatment renderPanel already applies. `tab` is free at that level, the handler at src/ui/ui.go:440 explicitly does nothing once inside a project.

Implementation is then the document-viewer capability applied to a third pane.

- [ ] add a focus field for the specs tab (list or content)
- [ ] bind `tab` to toggle it at levelProject on the specs tab
- [ ] render the focused pane with the active border colour
- [ ] route the vertical keys to the content viewport when it has focus
- [ ] reset scroll when the selected spec changes
- [ ] update the nav bar hints
