---
# specgetty-1q1c
title: extra view layer for Spec
status: draft
type: task
priority: normal
created_at: 2026-09-21T15:56:05Z
updated_at: 2026-09-21T16:35:21Z
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
