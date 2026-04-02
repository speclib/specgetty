## Context

This builds on the side-panel-detail-layout change. The detail panel currently shows raw file listings. This change adds semantic understanding of OpenSpec project contents and a tab system for future extensibility.

An OpenSpec project's `openspec/` directory has a known structure:
```
openspec/
├── config.yaml
├── specs/
│   ├── capability-a/spec.md
│   └── capability-b/spec.md
├── changes/
│   ├── active-change/
│   │   └── .openspec.yaml
│   └── another-change/
│       └── .openspec.yaml
└── archive/
    └── old-change/
        └── .openspec.yaml
```

## Goals / Non-Goals

**Goals:**
- Parse openspec/ to extract: spec count, active change names, archived change names with dates
- Display structured overview in the detail panel
- Add tab header with placeholder tabs for future views
- Tab switching via h/l or 1-5 number keys when detail panel is focused

**Non-Goals:**
- Making specs/changes/config/search tabs functional (placeholder only)
- Deep parsing of spec contents or task completion status (just counts)
- Editing or modifying OpenSpec project contents from the TUI

## Decisions

### 1. Stats extraction at scan time
Extend `ProjectStatus` in the scanner to include a `ProjectInfo` struct with parsed stats. This avoids re-reading the filesystem when switching between projects in the UI.

### 2. Minimal parsing — directory listing only
Count specs by counting subdirectories in `openspec/specs/`. List changes by reading `openspec/changes/` entries. List archived by reading `openspec/archive/` entries. Archive dates come from the directory's modification time (avoid parsing YAML).

**Alternative considered**: Using `openspec status --json` CLI for each project. Rejected — too slow for many projects, and creates a dependency on the openspec CLI being installed.

### 3. Tab system as a simple enum + renderer
Add a `detailTab` int to the model. Each tab has a render function. Only overview is implemented; others return a styled "Not yet implemented" message. This is easy to extend later.

### 4. Tab header styling
Active tab gets a highlighted background. Inactive tabs are dimmed. Use lipgloss styles consistent with the existing nav bar.

## Risks / Trade-offs

- **Risk**: Projects without standard openspec structure may cause stat parsing errors → Mitigation: Default to zero counts on any read error, never crash.
- **Trade-off**: Reading archive mtime for dates is approximate but avoids YAML parsing complexity.
