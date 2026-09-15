---
# specgetty-jdif
title: refacter zoom/scanned projects UI
status: draft
type: epic
priority: normal
created_at: 2026-09-15T14:28:51Z
updated_at: 2026-09-15T14:41:02Z
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
