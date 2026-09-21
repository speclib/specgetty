## Why

Bean `specgetty-9xi9`. The log panel costs a quarter of the focus ring, a key,
and a clause in five capabilities, and it holds nothing worth reading.

One scan of this author's configuration writes 35 lines into it:

```
  30 lines   /home/pim/gh.x/y  3.7ms      per-project scan telemetry
   2 lines   walkDuration / scanDuration  timings
   3 lines   ERROR: ...                   not real errors, see below
```

A picker refresh writes all 35 again. The three real error producers that exist,
a watcher that failed to start, an unreadable scan directory and an unreadable
store registry, are buried under thirty lines of telemetry in a panel nobody
opens. It is not an error surface; it is where errors go to die.

The three `ERROR` lines are the panel's own doing:

```
  ERROR: /home/pim/gh.*: lstat /home/pim/gh.*: no such file or directory
```

`Walk` adds each include to the walk list and expands the glob afterwards, so
the literal pattern `gh.*` is walked as though it were a directory. It always
fails. With `ignore_dir_errors` on, which is the default, the failure is logged
and swallowed. With it off it is not swallowed:

```console
$ spg --debug --ignore_dir_errors=false
lstat /home/pim/tc*: no such file or directory      # the whole scan fails
```

So `spg -i=false` has been broken for every configuration that uses a glob
include, which is the shape the shipped default config has. The log panel is
what made that survivable enough to go unnoticed.

## What Changes

- The log panel goes: the `l` key, the viewport, the focus state, the border
  rules that mention it, and the two-view panel switch
- Log output is discarded while the terminal interface is running. It is
  redirected into the panel today, and simply deleting that redirection would
  send it to stderr, which under a full-screen interface means writing over the
  frame
- `--debug` keeps logging exactly as it does now. It is the one place this
  output is useful and correct, because no interface is running
- A glob include contributes only what it expands to. The literal pattern stops
  being walked, which removes the three spurious errors and makes
  `--ignore_dir_errors=false` work for a configuration that uses globs

Not in scope: giving the three real errors a new home. They are as good as
invisible today, so discarding them changes little in practice, but "auto-rescan
died and nothing said so" deserves better than silence. Recorded as a follow-up
bean rather than folded in, because it is a new surface rather than a removal.

## Capabilities

### Modified Capabilities
- `panel-layout`: the border rules stop enumerating the log panel
- `specs-tab`: the focus ring is the list and the content
- `config-tab-display`: the same for the properties tab
- `document-viewer`: the paging keys no longer have a log panel to reach
- `project-view`: the project view fills the terminal with nothing below it
- `openspec-scanning`: a glob include contributes only its expansions

## Impact

- `src/ui/ui.go`: 66 references across the model, the key handling, the layout
  and the rendering. `focusLog` leaves the focus ring, which simplifies every
  `tab` case, every border decision and `halfPage`
- `src/ui/docview.go`, `src/ui/changelist.go`: the same guards
- `src/scanner/find.go`: the glob expansion stops appending the pattern itself
- Seven test files reference the panel
- `src/config.yml` and the README document the `l` key

## Rollback

Its own commit. Reverting restores the panel and the glob bug together, which is
why they belong in one change: the panel is what has been absorbing the bug.
