---
# specgetty-1q1c
title: extra view layer for Spec
status: completed
type: task
priority: normal
created_at: 2026-09-21T15:56:05Z
updated_at: 2026-09-21T18:22:47Z
blocked_by:
    - specgetty-9xi9
---

Currently the Specs are rendered as markdown but they always have the same
structure and reading them would be more focussed it a new layer
shows one spec in detail:

- left panel:
  - proposal
  - requirements:
    - my first requirement
      - scenario's
        - my first scenario
        - my second scenario
    - my second requirement
      - scenario's
        - my first scenario


- right panel body with content optimized in focus view, like a card:


```
--------------------------------------------------------------------------
|                                                                        |
|                                                                        |
|        Scenario: Project with tasks across multiple                    |
|                                                                        |
|        WHEN                                                            |
|           a project has 2 active changes, one with 3/5 tasks           |
|           done and another with 2/4 tasks done                         |
|                                                                        |
|        THEN                                                            |
|           `ProjectInfo.TasksTotal` SHALL be 9 and                      |
|           `ProjectInfo.TasksDone` SHALL be 5                           |
|                                                                        |
|                                                                        |
--------------------------------------------------------------------------
```

The Scenario title should be rendered bold
The WHEN THEN AND SHALL keywords should be rendered in highlighted colors.
text between backticks should be highlighted
normal text should be written in a clear color


---

This change is currently for the live specs. When this works we also want to have a look at the specs in active changes, but this need more attention as these as formed as a delta

## Design

Explored and specified as OpenSpec change `open-a-spec-in-detail`.

Settled during exploration:

- A third navigation level. `enter` on the specs tab descends, `esc` returns.
- The left panel lists Purpose (not proposal; live specs have no proposal),
  the requirements and their scenarios.
- Purpose shows its prose. A requirement shows the text above its first
  scenario, nothing more. A scenario shows the card.
- The spec is parsed on demand, held while open.
- A spec that does not fit the structure does not open. The nav bar says why
  and the markdown view on the tab still shows the whole file.
- The outline wraps its labels rather than truncating, and the split goes from
  three tenths to four. At three tenths on an 80-column terminal only 22% of
  requirement titles and 31% of scenario titles would fit.

Deferred, as the bean already noted: the specs inside an open change, which are
deltas and need an outline level the live specs do not have.

## Ordering

Blocked by `specgetty-9xi9` (remove the log panel), which is in progress.

Not a spec conflict: the two changes touch `specs-tab` and `document-viewer`
but never the same requirement, so archive merges them in either order.
`9xi9` removes `The specs tab has a focus` and `The paging keys act on whatever
holds the keyboard`; this one adds three requirements to `specs-tab` and
modifies `A page is what the surface can show`. Disjoint names in both files.

The order is about the code. `9xi9` deletes `focusLog` from 66 references in
`ui.go`, which it says simplifies every `tab` case, every border decision and
`halfPage`. This change adds a third arm to `docRegion`, `docActive`,
`splitTab`, `listPage` and the nav bar. Going second means writing the log
panel into the new level and deleting it again, across seven test files that
reference the panel.

Three lines to edit here once `9xi9` lands:

- `specs/spec-detail-view/spec.md:36`, drop "and on through the log panel when
  it is open" from the requirement text
- `specs/spec-detail-view/spec.md:48`, delete the scenario "Cycling through the
  log panel"
- `tasks.md:74`, rewrite 6.2 as a two-stop cycle between the outline and the
  card

Nothing in `9xi9` refers to this change.


## Summary of Changes

Shipped as OpenSpec change `open-a-spec-in-detail`, commit `af367e7`.

`enter` on the specs tab now opens the selected spec at a third navigation
level: an outline of its requirements and scenarios on the left, a card showing
whichever node the cursor is on at the right. `tab` moves the keyboard between
the halves, `esc` returns to the list with the same spec selected.

What was built:

- `src/ui/specparse.go` parses a spec into Purpose, requirements and scenarios,
  with each scenario's clauses split into keyword and text. It accepts all 24
  live specs in this repository (148 requirements, 464 scenarios) and refuses a
  file that yields no requirements or a requirement with no scenarios, so a spec
  it cannot navigate reports why rather than opening empty.
- `src/ui/specdetail.go` holds the split, the outline and the card. A label too
  long for the outline wraps rather than being cut, and the cursor still moves
  one node per keystroke however many rows that node occupies. Every row of the
  selected node is highlighted, and the outline scrolls a whole node into view.
- The card lays a clause out with its keyword on its own row and the text under
  a hanging indent, which is the reading view the bean asked for. Backticked
  spans are styled and their marks dropped.
- Each node is remembered by path rather than by index, so a rescan that inserts
  a scenario above the cursor leaves the cursor where it was.
- Paging keys (`^f`, `^b`, `^d`, `^u`, `gg`, `G`) follow the keyboard between the
  two halves, with the outline's page walking node heights until one pane of
  rows is covered.

Coverage: ui 87.0% (floor raised 85.8 to 86.8), total 88.9% (floor 88.0 to 88.7).
The gate ran green; the revert check bit on the split agreement, on which half
owns the vertical axis, and on the path-based cursor memory.
