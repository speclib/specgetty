---
# specgetty-0i4n
title: a bad scandir glob kills the process via log.Fatal
status: todo
type: bug
priority: normal
created_at: 2026-09-15T20:16:30Z
updated_at: 2026-09-15T20:16:30Z
---

A scan directory glob whose parent cannot be read kills the process.

`src/scanner/find.go`, inside `Walk`'s glob expansion:

    parent := filepath.Dir(globPath)
    entries, err := os.ReadDir(parent)
    if err != nil {
        log.Fatal(err)
    }

`log.Fatal` writes and then calls `os.Exit(1)`. There is no recovery and no
error returned to the caller.

## Verified

Running `Walk` with an include of `/definitely/not/here/c*` in a subprocess:

    child exit status 1
    open /definitely/not/here: no such file or directory

## Why it matters more than it looks

`--ignore_dir_errors` defaults to **true** and its usage string is "Don't halt
on errors while finding dirs". That flag is only consulted in the `walkone`
loop further down (`find.go:155`). The glob expansion runs before it and
ignores it, so the one setting meant to prevent this exact failure does not
cover it.

Two things make it worse now:

- `log.SetOutput` points the logger at the TUI's log panel, so the message may
  not even reach the terminal before the process exits.
- Since the project picker landed, the scan happens when the user presses `p`
  or `r`, not at startup. The crash arrives mid-session, from an alt-screen
  TUI, which can leave the terminal in a bad state.

## Trigger

Any include glob pointing into a directory that does not exist: an unmounted
drive, a path that moved, an env var that expanded to nothing. A non-glob
include is unaffected, because the expansion branch only runs for entries
ending in `*`.

## Shape of a fix

Treat it like any other directory error: log it, skip that include, carry on.
Returning the error would also be defensible, but silently continuing is what
the flag already promises.

Not fixed in `cover-scanner-and-export` because it is pre-existing and out of
that change's scope, and because a test for it has to spawn a subprocess to
survive.
