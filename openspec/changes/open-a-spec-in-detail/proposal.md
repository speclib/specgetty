## Why

Bean `specgetty-1q1c`. A spec is read as one markdown document, top to bottom.
`change-list-view` renders to 425 rows in the 49 columns the specs tab gives it
on an 80-column terminal, which is thirty screenfuls in a 14-row pane, with no
landmark to aim at and no way to ask for one requirement.

The structure is in the file and the renderer throws it away. Across the 24 live
specs in this project there are 140 requirements and 428 scenarios, and every one
of them sits under a heading that `renderMarkdown` turns into another bold line:

```
  ## Purpose
  ### Requirement: <title>              140 of these
  #### Scenario: <title>                428 of these
  - **GIVEN** **WHEN** **THEN** **AND**
```

Nothing else is in the corpus. No spec has a heading outside those four shapes,
no requirement lacks its text or its scenarios, no scenario lacks its keyword
bullets. The shape is reliable enough to navigate by.

## What Changes

- `enter` on the specs tab descends into a new view of the selected spec, and
  `esc` returns to the tab. It is the third navigation level, alongside the
  change view, and it follows the same rule: `enter` descends, `esc` ascends.

- The view is an outline beside a card. The outline lists the spec's Purpose,
  its requirements and their scenarios. The card shows whichever node the cursor
  is on:

```
╭──────────────────────────────╮ ╭───────────────────────────────────────────╮
│ Purpose                      │ │                                           │
│ Specs tab displays a list of │ │     Scenario: Project with specs          │
│   specs                      │ │                                           │
│   Project with specs         │ │     WHEN                                  │
│   Project with no specs      │ │        the user switches to the specs     │
│ Specs tab displays selected  │ │        tab for a project that has specs   │
│   spec content               │ │                                           │
│   Spec selected              │ │     THEN                                  │
│   Spec without spec.md       │ │        the left side SHALL list all       │
│ Spec list navigation         │ │        spec directory names sorted        │
│   Navigate specs with j/k    │ │        alphabetically                     │
╰──────────────────────────────╯ ╰───────────────────────────────────────────╯
```

- A scenario's card is a re-layout, not a restyle. The source is one bullet,
  `- **WHEN** the user switches to...`, and the card puts the keyword on its own
  row with the clause indented under it. The title is bold, the keywords are
  coloured, backticked spans are highlighted, and the prose is left legible.

- Purpose shows its prose. A requirement shows the text above its first
  scenario, and nothing more. A scenario shows its card.

- The outline wraps its labels rather than truncating them. Requirement titles
  run to a median of 36 columns and a maximum of 63, so in the 22 columns
  today's split would give them on an 80-column terminal, four rows in five
  would be a stub. The outline half is widened from three tenths to four tenths
  for the same reason. Colour carries the level once labels wrap, because a
  wrapped requirement line sits at the same indent as a scenario title.

- A spec that does not fit the shape does not open. `enter` leaves the cursor
  where it is and the nav bar says why, in the transient one-line form that `c`
  and `C` already use. The markdown view on the tab is unchanged and still shows
  the whole file, so nothing becomes unreadable.

- Backticked spans are highlighted in the markdown renderer too, which is where
  the card gets the behaviour from.

Not in scope: the specs inside an open change. Those are deltas, with `ADDED`,
`MODIFIED`, `REMOVED` and `RENAMED` sections above the requirements, and the
outline would need a level the live specs do not have. The parser is written so
that level can be added later, but nothing in this change reads a change's
specs.

## Capabilities

### New Capabilities
- `spec-detail-view`: one spec read as an outline beside a card, with the
  parse that makes it possible and the report when the parse fails

### Modified Capabilities
- `specs-tab`: `enter` descends into the detail view, and a spec that cannot be
  parsed reports instead of opening
- `document-viewer`: what a page means over a surface whose rows are not all one
  row tall

## Impact

- `src/ui/specparse.go`: new, the parser and its node tree
- `src/ui/specdetail.go`: new, the outline, the card and the split
- `src/ui/ui.go`: a third level, and the arms it needs in `docRegion`,
  `docActive`, `splitTab` and `renderNavBar`
- `src/ui/listpaging.go`: `listPage`, `moveListCursor` and `gotoListEnd` gain the
  new surface, and the page walk learns rows that are not all one row tall
- `src/ui/docview.go`: `currentDoc` gains the card, keyed by node
- `README.md`: the new level and its keys

## Rollback

Its own commit. Reverting removes the level and leaves the specs tab exactly as
it ships today, since nothing about the markdown view changes except the
backtick highlight.
