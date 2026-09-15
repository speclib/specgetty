---
# specgetty-jdif
title: refacter zoom/scanned projects UI
status: in-progress
type: epic
priority: normal
created_at: 2026-09-15T14:28:51Z
updated_at: 2026-09-15T16:34:57Z
---

The current UI works like this:

- default view at startup: all scanned projects (uncashed)
  - when enter on project open in zoom mode
- when --zoom option given open in single mode
  - when escaping from single all projects are opened

We should change this refactor this UI.

- remove the term ZOOM. replace with "single project view" is needed
- the single project view is the initial view at startup
- the project picker will be a popup called by a menu key option
- when no openspec dir is found ask: open project picker 

The project picker

- a highlighter project is selected with enter from the project picker -> the single view will be switched to this project
- the picker opens in a popup which closes with escape key
- the picker shows all project as table rows with columns showing project stats like: Specs: 27  Changes: 2 active  Archived: 30  Tasks: 0/50
- the picker has its own search filter: search project names, search in project (filename / file contents)
- the project picker should cache its content and offer a refresh action

## OpenSpec change

`openspec/changes/project-picker/` holds proposal.md, design.md, tasks.md
(67 items) and deltas for project-picker, project-cache, project-view (new),
openspec-scanning, fs-watch (modified), zoom-mode, side-panel-layout (retired).
Validates strict.

## Decisions taken while exploring (2026-09-15)

| Topic            | Decision                                                  |
| ---------------- | --------------------------------------------------------- |
| picker key       | `p`, a plain key, not a menu                               |
| `--zoom`         | removed outright, no alias                                 |
| flags            | `--view=single\|all`, default single; `--path` retained    |
| picker from L1   | allowed; switching always lands at the project view        |
| cache            | on disk, paths only, never statistics                      |
| split view       | deleted; git remembers it if it is ever wanted back        |
| stores           | explicitly out of scope, `--view=store` comes later        |
| picker columns   | hard-coded, no `--project-fields` until asked for          |

## Measurements that drove the design

| Step                                      | Cost   |
| ----------------------------------------- | ------ |
| Directory walk, cold                      | ~3.9s  |
| Directory walk, warm                      | 697ms  |
| Listing openspec contents (18 projects)   | 55ms   |
| Reading every artifact (3522 files, 18MB) | 101ms  |
| Resolving the project at the cwd          | ~1ms   |

The walk is the entire cost, so the cache holds paths and nothing else. Parsing
is cheap enough that the picker parses every row eagerly on open.

Careful: `ProjectStatus.ScanTime` and the logged `scanDuration` time only
`ListOpenSpecContents`, not `ParseProjectInfo`. Reading them as "the parse cost"
is misleading.

## Found while exploring

- `filePaths` and `fileCursor` are written but never rendered. Dead since the
  tabs replaced the flat file listing. The matching requirement in
  `openspec-scanning` is equally dead. Both go in this change.
- `openspec/specs/overview-tab/spec.md` describes a tab replaced by the
  persistent header in `overview-as-persistent-header` (archived 2026-04-20).
  Stale, unrelated to this change, left for its own cleanup.
- A store root contains `openspec/`, so the scanner already finds registered
  stores as ordinary projects. Registry lives at
  `~/.local/share/openspec/stores/registry.yaml`.
