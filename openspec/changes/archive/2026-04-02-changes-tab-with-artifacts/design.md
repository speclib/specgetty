## Context

The changes tab needs to show active changes and let users drill into their artifacts. A change directory typically contains:

```
openspec/changes/my-change/
├── .openspec.yaml
├── proposal.md
├── design.md
├── tasks.md
└── specs/
    └── some-capability/
        └── spec.md
```

But custom schemas may have different/additional `.md` files. The approach is schema-agnostic: just discover what's there.

## Goals / Non-Goals

**Goals:**
- List active changes with task completion stats
- Sub-tab navigation for artifacts within a selected change
- Discover `.md` files dynamically (no hardcoded artifact list)
- Parse tasks.md for `- [x]` / `- [ ]` checkbox counting
- Show specs/ subdirectory as a navigable sub-tab (list + content, same as main specs tab)

**Non-Goals:**
- Schema awareness or validation
- Editing artifacts from the TUI
- Showing change status from `.openspec.yaml` (beyond what's discoverable from files)

## Decisions

### 1. ChangeInfo struct in scanner
```
ChangeInfo {
    Name          string
    ArtifactFiles []string          // sorted .md filenames (e.g. "design.md", "proposal.md", "tasks.md")
    ArtifactContents map[string]string // filename → content
    TasksTotal    int
    TasksDone     int
    SpecNames     []string
    SpecContents  map[string]string
}
```
Read at scan time. `ArtifactFiles` is populated by listing `.md` files in the change dir (excluding hidden files). Task stats parsed from `tasks.md` if it exists.

### 2. Three-level navigation with two cursor fields
- `changeCursor int` — which change is selected in the left list
- `changeArtifactTab int` — which artifact sub-tab is selected (index into the discovered artifact list + "specs" appended at the end)
- j/k moves `changeCursor` when changes tab is active
- left/right (or number keys) moves `changeArtifactTab` within the sub-tabs

### 3. Layout mirrors specs tab pattern
Left 30%: change list with task stats. Right 70%: sub-tab header + content. The "specs" sub-tab renders as a nested list+content split (reusing the same pattern).

### 4. Task stats display
- Change list: `my-change (3/5)` — compact format
- Change header: `Tasks: 3/5 complete` — when change is selected

### 5. Sub-tab labels are filename stems
`proposal.md` → `proposal`, `design.md` → `design`, etc. The specs sub-tab is always last and labeled `specs`.

## Risks / Trade-offs

- **Trade-off**: Reading all artifact contents at scan time increases memory for projects with many changes. Acceptable — change artifacts are typically small.
- **Risk**: The specs-within-a-change sub-tab adds a fourth navigation level (spec list + spec content). Keep it simple — no separate cursor, just render the list with content of the first spec.
