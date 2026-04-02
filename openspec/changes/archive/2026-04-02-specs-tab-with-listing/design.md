## Context

The detail panel has 5 tabs. Overview and config are functional. The specs tab needs to show a browsable list of specs with their content. Each spec lives in `openspec/specs/<name>/spec.md`.

The specs tab reuses the same horizontal split pattern as the main layout: narrow list on the left, content on the right.

```
┌─ Detail ──────────────────────────────────────────┐
│ [overview] [specs] [changes] [config] [search]     │
│────────────────────────────────────────────────────│
│ openspec-scanning │ ## ADDED Requirements          │
│ side-panel-layout │                                │
│                   │ ### Requirement: Detect ...     │
│                   │ The scanner SHALL identify...   │
│                   │                                │
│                   │ #### Scenario: Directory with.. │
└───────────────────┴────────────────────────────────┘
```

## Goals / Non-Goals

**Goals:**
- Read spec names and contents at scan time, store in ProjectInfo
- Split the specs tab area: ~30% for spec list, ~70% for spec content
- j/k navigation moves the spec cursor when specs tab is active and detail panel is focused
- Selected spec's `spec.md` is rendered with the existing markdown renderer
- Scrollable spec content for long specs

**Non-Goals:**
- Parsing spec structure (requirements, scenarios) — just render as markdown for now
- Editing specs from the TUI
- Searching within specs

## Decisions

### 1. Read all spec.md contents at scan time
Store a sorted `[]string` of spec names and a `map[string]string` of name→content in ProjectInfo. Spec files are small, so reading them all upfront is fine.

### 2. Spec cursor as separate model field
Add `specCursor int` to the model. When `detailTab == tabSpecs` and `activeView == viewDetail`, j/k moves `specCursor` instead of `fileCursor`. This keeps navigation context separate per tab.

### 3. Reuse renderMarkdown for content
The existing `renderMarkdown` function handles headers, lists, bold, italic — good enough for spec.md files.

### 4. Spec list width: 30% of detail panel inner width
Consistent with the main layout's project list ratio.

## Risks / Trade-offs

- **Trade-off**: Reading all spec contents at scan time increases memory for projects with many large specs. Acceptable — spec files are typically small (1-5KB).
