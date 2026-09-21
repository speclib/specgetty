---
# specgetty-lf42
title: open markdown in $EDITOR
status: in-progress
type: task
priority: normal
created_at: 2026-09-21T16:37:21Z
updated_at: 2026-09-21T18:31:08Z
---

When a view represents a single markdown (or yaml) a keybinding (E) should allow to open the file in $EDITOR

## OpenSpec change

`open-the-file-in-an-editor`, validated strict.

Settled during exploration:

- `E` everywhere a pane shows exactly one file. The rule is
  `document.path != ""`, a field `docview.go` already has and which only the
  tasks artifact sets today.
- Three panes gain it: a change's proposal/design/tasks, a spec on the specs
  tab, the project row of the properties tab. A change's specs sub-tab shows
  several files and the schema and store rows show assembled reports, so none
  of the three gets the key.
- `$VISUAL` first, then `$EDITOR`. A value with arguments works, split on
  whitespace, no shell. Neither set reports on the nav bar and guesses nothing.
- `tea.ExecProcess`, which the library documents for this case. Its callback
  rescans, the way the task toggle already does rather than trusting the
  watcher alone.
- The file opens at its top. No line jump: the syntax differs per editor, so it
  would work for some people and silently do nothing for the rest.

## Bug found and folded in

`store-resolution` requires every filesystem operation to act on the resolved
root, and states it as a list of four: reading, archive, discard, export. Two
writers were added later and never joined it:

```
  tasks toggle   renderChangeArtifact(m.repoPaths[m.cursor], ...)   origin
  Y copy path    changeDirPath(m.repoPaths[m.cursor], r)            origin
```

In a store-backed project the toggle fails with a path error and `Y` copies a
path that does not exist. `E` would have been the third. The delta restates the
requirement as a rule over every operation, with a scenario saying an operation
added later is covered without the requirement being amended.
