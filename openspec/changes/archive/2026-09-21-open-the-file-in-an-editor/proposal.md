## Why

Bean `specgetty-lf42`. Every pane in this tool shows a file it will not let you
change. The one exception is a task checkbox, which `space` toggles in place.

The workaround already shipped, and its own specification names where it is
going. `copy-to-clipboard` exists so that `Y` puts a change's directory on the
clipboard, and its Purpose says why:

> either its name, for pasting into a command, or its absolute path, for pasting
> into a shell or an editor.

So the path is already being carried out of the application by hand, to be
pasted into an editor, because there is no key that opens one. That is the whole
case for this change.

The second half is a defect the same work uncovers. `store-resolution` requires
that every filesystem operation act on the resolved root, and it says so as a
list:

> Reading specs and changes, and the archive, discard and export actions, SHALL
> act on the resolved root.

Four operations, closed. Two writers were added afterwards and never joined the
list, and both build their paths from the directory the user stood in:

```
  archive        doArchiveChange(m.currentRoot(), ...)                 root
  discard        doDiscardChange(m.currentRoot(), ...)                 root
  export         doExportChange(m.currentRoot(), ...)                  root
  tasks toggle   renderChangeArtifact(m.repoPaths[m.cursor], ...)      origin
  Y copy path    changeDirPath(m.repoPaths[m.cursor], r)               origin
```

In a store-backed project the toggle fails with `could not save: no such file or
directory`, and `Y` copies a path that looks right and does not exist. Opening a
file in an editor would be the third writer to walk into the same gap, so the
rule is restated here as a rule rather than as a list.

## What Changes

- `E` opens the file the pane is showing in the user's editor. It is bound
  wherever a pane shows exactly one file, and nowhere else:

```
  a change's proposal, design or tasks        one .md            E
  a spec on the specs tab                     one spec.md        E
  the project row of the properties tab       one config file    E
  a change's specs sub-tab                    several files      no key
  the schema and store rows                   a report, no file  no key
```

- The editor is `$VISUAL` if it is set, then `$EDITOR`. That is the convention:
  `VISUAL` names the full-screen editor and `EDITOR` the line editor, and a user
  with both set means the first one here.

- A value carrying arguments works, so `EDITOR="emacsclient -nw"` and
  `EDITOR="code -w"` both run. The value is split on whitespace, which does not
  handle a quoted path inside the variable. No shell is involved.

- With neither set, the key reports on the nav bar and does nothing. There is no
  fallback to `vi`: guessing is worse than saying, and it guesses wrong on a
  machine that has no `vi`.

- The application stops drawing while the editor runs and resumes when it exits,
  then rescans the open project, so an edit is on screen without the user asking
  for it.

- The file opens at its top. Jumping to the line under the cursor was considered
  and dropped: the syntax differs per editor, so it would work for some people
  and silently do nothing for the rest.

- **Fixed**: the task toggle and the `Y` copy act on the resolved root, so both
  work in a store-backed project. `store-resolution` states the rule over every
  operation instead of naming four.

Not in scope: creating a file, renaming one, or any editing inside specgetty
itself. The editor is the editor.

## Capabilities

### New Capabilities
- `edit-in-editor`: the key, which panes carry it, how the editor is chosen and
  what happens while it runs

### Modified Capabilities
- `store-resolution`: the rule is stated over every filesystem operation rather
  than over a list of four
- `copy-to-clipboard`: `Y` copies the path under the resolved root
- `task-checkboxes`: the toggle writes to the file under the resolved root

## Impact

- `src/ui/editor.go`: new, the editor resolution and the command
- `src/ui/docview.go`: `document.path` is set for every file-backed document,
  built from the resolved root
- `src/ui/clipboard.go`: `changeDirPath` takes the root
- `src/ui/ui.go`: the `E` case, and the nav bar hint where it applies
- `README.md`: the key, and what it needs in the environment

## Rollback

Its own commit. Reverting removes the key and returns the toggle and the copy to
building their paths from the origin, which is the bug they ship with today.
