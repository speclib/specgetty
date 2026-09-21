## Context

See proposal.md for why. What matters here is what the code already has and what
it is missing.

`docview.go` already carries the concept this change needs:

```go
type document struct {
    key     string
    content string
    lines   []sourceLine // for the cursor
    path    string       // the file on disk, for a toggle to write back
}
```

"A pane that shows one file" is `document.path != ""`. One document sets it
today, the tasks artifact, because a checkbox toggle had to write somewhere. The
field was built for a narrower purpose than it describes.

`currentDoc` has four cases. Two of them are a file, one is several files
concatenated, and one is a report assembled from the project's state:

```
  levelChange, artifact tab      one .md                 a file
  levelChange, specs sub-tab     every delta together    not one file
  levelProject, specs tab        one spec.md             a file
  levelProject, properties       config: a file
                                 schema, store: a report assembled in memory
```

The root question is settled elsewhere in the code. `currentRoot()` returns the
directory content was resolved to, `currentKey()` returns the directory the user
started in, and the scanner is explicit that the project map is keyed by the
second: "Keyed by the path that was asked for, not by the root: two repos
sharing a store are two rows."

## Goals / Non-Goals

Goals: one key over every pane that shows a file; one function that answers
which file; the root rule stated so the next writer cannot miss it.

Non-Goals: editing inside specgetty. Creating or renaming files. Opening an
editor at a position in the file. Any behaviour that depends on the editor being
terminal-based.

## Decisions

### The key is a property of the document, not of the view

`E` asks `document.path` rather than asking which level and tab are active. The
panes that carry the key and the panes that do not then follow from one fact
instead of from a list that each new pane must be added to, which is the same
mistake `store-resolution` made with its four operations.

`document.path` stops meaning "the file a toggle writes back to" and starts
meaning "the file this document came from". The toggle keeps working because
its own guard is `docHasCursor()`, which also requires `docLines`, and only the
tasks artifact has those.

### Every path is built from the resolved root

`renderChangeArtifact` takes the project path as an argument and joins it with
`openspec/changes/...`. Today `currentDoc` passes `m.repoPaths[m.cursor]`, which
is the starting directory. It passes `m.currentRoot()` instead, and
`changeDirPath` is called with the same.

An alternative was to have `renderChangeArtifact` look the root up itself. It
takes the path as a parameter so it stays a function of its arguments, which is
what makes it testable without a model, so the call site is where the fix goes.

### VISUAL before EDITOR

The convention is older than either of us: `EDITOR` names a line editor that
works without a full terminal, `VISUAL` names the full-screen one. A user with
both set has said which they want for an interactive edit. Reading only
`EDITOR`, which is the shorter thing to write, would open `ed` for someone whose
environment is correct.

### The value is split on whitespace, with no shell

`EDITOR="emacsclient -nw"` and `EDITOR="code -w"` are both ordinary. Passing the
whole string to `exec.Command` looks for a binary with a space in its name.

Two ways out. `sh -c` handles quoting and word splitting exactly as a shell
would, at the cost of requiring a shell and of running the file path through
shell parsing, where a path with a space or a quote in it becomes an injection
question. Splitting on whitespace handles the cases people actually have and
passes the path as one argument that nothing parses.

Splitting wins. The case it does not handle, a quoted path inside the variable,
is stated in the spec rather than left to be discovered.

### Nothing is guessed when neither variable is set

A fallback to `vi` is a guess that is wrong on a container with no `vi`, and
when it is wrong the failure arrives as a process that could not start rather
than as an answer. Reporting on the nav bar through `statusMsg` costs one line
and says the true thing. That is the mechanism the copy keys already use, and
the model comment on the field argues the case: a modal is the wrong weight for
something instant and harmless.

### tea.ExecProcess, and a rescan in its callback

`tea.ExecProcess` is documented for this exact case: "useful for spawning other
interactive applications such as editors and shells from within a Program." It
pauses the program, runs the command against the real terminal, and resumes.

Its callback returns a message, which is where the rescan goes. The filesystem
watcher would usually notice the edit on its own, but it is allowed to fail to
start, and then nothing on screen would move. The task toggle already makes
exactly this argument and rescans for exactly this reason.

The terminal may have been resized while the editor had it. Bubbletea reports
the size on resume, so the existing `WindowSizeMsg` path handles it and nothing
here needs to.

### The file opens at its top

Jumping to the line under the cursor is worth something in the spec detail view,
where a node maps to a line. It is `+N` for vim, nano and emacs and
`--goto file:N` for code, so a single implementation works for some editors and
silently ignores the argument, or errors, for the rest. A feature that works for
some people and fails quietly for others is worse than one that does the same
thing for everyone.

## Risks / Trade-offs

- An editor that is not terminal-based, such as one that opens a window and
  returns at once, resumes the interface immediately and the edit lands later. →
  The watcher picks it up. Nothing here assumes the editor holds the terminal,
  and the rescan on resume is additional rather than required.
- Splitting on whitespace mangles `EDITOR='/path/with space/ed'`. → Stated in
  the spec. `$VISUAL` gives the user a second place to put a value that works.
- The root fix touches three call sites at once, one of which is a rendering
  function used by every artifact. → The change is a different argument at the
  call site, and two of the three have a failing test written before the fix.

## Migration Plan

None. The key is new and the two fixes make operations succeed that fail today.
