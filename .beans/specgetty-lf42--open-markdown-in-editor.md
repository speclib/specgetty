---
# specgetty-lf42
title: open markdown in $EDITOR
status: completed
type: task
priority: normal
created_at: 2026-09-21T16:37:21Z
updated_at: 2026-09-21T18:40:08Z
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


## Summary of Changes

Shipped as OpenSpec change `open-the-file-in-an-editor`, commit `f24aed3`.

`E` hands the file a pane is showing to your own editor. It is bound wherever a
pane shows exactly one file and nowhere else: a change's proposal, design or
tasks, a spec on the specs tab, and the project row of the properties tab. A
change's specs sub-tab shows every delta at once and the schema and store rows
are assembled reports, so neither carries the key. The nav bar lists `E` exactly
where it applies, driven by the same fact the key itself asks.

The editor is `$VISUAL` when set, `$EDITOR` otherwise, and a value carrying
arguments works, so `EDITOR="code -w"` runs. No shell is involved, so a quoted
path inside the variable does not work; that is stated in the spec and in the
README rather than left to be discovered. With neither variable set the key
reports on the nav bar and does nothing: there is no fallback to `vi`, because
guessing is wrong on a machine that has no `vi`. The interface yields the
terminal while the editor runs and reads the project again when it exits, so an
edit is on screen without pressing `s`. The file opens at its top; passing a
position would work for some editors and be silently ignored by others.

Two defects the same work uncovered, both shipped fixed:

- Toggling a task checkbox in a store-backed project failed with `could not
  save: no such file or directory`, and `Y` copied a path that looked right and
  did not exist. Both built their path from the repository you started in rather
  than from the store the content was read from.
- `store-resolution` stated its rule as a list of four operations. Two writers
  added later never joined that list, and both were wrong for as long as it
  stood. The requirement now states the rule over every operation, so the next
  one is covered without the spec being amended to name it.

`document.path` also changed meaning, from "the file a toggle writes back to" to
"the file this document came from". The toggle's guard is `docHasCursor()`, which
also wants mapped source lines, and only `tasks.md` has those: a test compares
the cursor table before and after and it is identical.

Coverage: ui 87.8% (floor raised 86.8 to 87.6), total 89.5% (88.7 to 89.2). The
gate ran green, and the revert check bit on both halves: reverting the two call
sites builds cleanly and fails the store tests again.
