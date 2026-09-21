---
# specgetty-1q1c
title: extra view layer for Spec
status: in-progress
type: task
priority: normal
created_at: 2026-09-21T15:56:05Z
updated_at: 2026-09-21T18:08:39Z
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
