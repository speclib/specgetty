# Opening a file in your editor

`E` hands the file a pane is showing to your own editor. It is bound wherever a
pane shows exactly one file, and nowhere else:

| Pane                                  | What `E` opens         |
|:--------------------------------------|:-----------------------|
| a change's proposal, design or tasks  | that artifact's file   |
| a spec on the specs tab               | that spec's `spec.md`  |
| the project row of the properties tab | the configuration file |
| a change's specs sub-tab              | nothing: several files |
| the schema and store rows             | nothing: no file       |

The nav bar lists `E` exactly where it would do something, so the pane itself
tells you whether the key applies.

The editor comes from `$VISUAL` when that is set, and from `$EDITOR` otherwise.
That is the convention: `EDITOR` names a line editor that works without a full
terminal and `VISUAL` the full-screen one, so with both set the first one here
is what you meant for an interactive edit.

A value carrying arguments works, because the value is split on whitespace:

```sh
export EDITOR="code -w"
export VISUAL="emacsclient -nw"
```

What that does not handle is a quoted path inside the variable, such as
`EDITOR='/opt/my editor/bin/ed'`. No shell is involved, so the space is read as
a separator. `$VISUAL` gives you a second place to put a value that works.

With neither set, `E` says so on the nav bar and does nothing. There is no
fallback to `vi`: guessing is worse than saying, and it guesses wrong on a
machine that has no `vi`.

The interface stops drawing while the editor has the terminal and resumes when
it exits, then reads the project again, so your edit is on screen without your
asking for it. The file opens at its top; no position within it is passed,
because the syntax for that differs per editor and would work for some people
and silently do nothing for the rest.
