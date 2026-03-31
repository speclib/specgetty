## Context

Current layout is vertical:
```
┌─ Repositories ──────────────────────────┐
│ /home/pim/project-a                      │
│ /home/pim/project-b                      │
└──────────────────────────────────────────┘
┌─ Contents ──────────────────────────────┐
│ d  changes                               │
│ f  config.yaml                           │
└──────────────────────────────────────────┘
 q quit  s scan  tab switch ...
```

Target layout is horizontal:
```
┌─ Projects ──┐┌─ Detail ─────────────────┐
│              ││ /home/pim/project-a      │
│  project-a   ││                          │
│  project-b   ││ d  changes               │
│              ││ d  specs                  │
│              ││ f  config.yaml            │
│              ││                           │
└──────────────┘└──────────────────────────┘
 q quit  s scan  tab switch ...
```

## Goals / Non-Goals

**Goals:**
- Horizontal two-panel layout with fixed-width left panel
- Left panel shows basenames only (e.g. `specgetty` not `/home/pim/cVibeCoding/specgetty`)
- Right panel shows full path as header line, then file listing below
- Existing navigation (j/k, tab, gg/G, pgup/dn) works in both panels
- Log panel at bottom spans full width when toggled

**Non-Goals:**
- Tab buttons in the detail header (that's the next change)
- Parsing openspec contents for stats
- Any changes to the scanner or data model

## Decisions

### 1. Left panel width: 30% of terminal, min 20, max 40 chars
A fixed ratio keeps it readable on both small and large terminals. The project list only shows basenames so it doesn't need much space.

### 2. Right panel gets remaining width
Simple: `rightWidth = totalWidth - leftWidth - 3` (accounting for borders and separator).

### 3. Project names are basenames with dedup
Show `filepath.Base(path)`. If two projects have the same basename, append a disambiguator (the parent dir name). This keeps the list scannable.

**Alternative considered**: Always show full path in the left panel. Rejected — defeats the purpose of the layout change, and the full path is shown in the detail panel.

### 4. Detail panel structure
Line 1: full project path (styled as a header). Then a blank line. Then the file listing (reusing existing renderFileList logic but at the new width).

### 5. Layout calculation replaces height-based with width+height
`recalcLayout()` currently only deals with heights. It needs to also compute panel widths. The left panel height spans from top to nav bar (minus log panel). The right panel matches.

## Risks / Trade-offs

- **Risk**: Small terminals (< 60 cols) may not have enough space for both panels → Mitigation: Fall back to full-width single panel if terminal is too narrow.
- **Trade-off**: Basenames lose path context at a glance, but full path is always visible in the detail panel header.
