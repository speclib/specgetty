## Why

Bean `specgetty-bfbo`. A change can only be exported to one place:

```go
func exportDestPath(semanticName string) string {
    home, _ := os.UserHomeDir()
    return filepath.Join(home, semanticName+"-"+dateStr+".zip")
}
```

The home directory is not an oversight, it is written down: `export-change` says
the zip "SHALL be written to the user's home directory". So exporting a change
to hand to someone means exporting it and then moving it.

## What Changes

- A directory can be configured once, and edited at the moment of export. The
  configured value is what the prompt opens on, so someone who exports to the
  same place always types nothing and someone who does not can still redirect a
  single export

```
  ╭────────────────────────────────────────────╮
  │  Export "follow-the-store"                 │
  │                                            │
  │  → ~/Downloads/_                           │
  │    follow-the-store-2026-09-21.zip         │
  │                                            │
  │  tab completes   ⏎ export   esc cancel     │
  ╰────────────────────────────────────────────╯
```

- A directory is typed, never a filename. The name carries a stripped archive
  prefix and an export date, both of which do work, and neither of which
  survives being retyped
- `tab` completes a path against the filesystem. The text input already
  supports suggestions with `tab` bound to accept them, so this is feeding it
  directory entries rather than building anything
- A directory that does not exist is refused with a message, leaving what was
  typed in place to be corrected. Creating it silently would turn a typo into a
  directory
- An existing file is confirmed before it is replaced. Silently overwriting a
  generated name in your own home directory is a redo; silently overwriting a
  path someone just typed is not the same act, and the old rule does not survive
  the destination becoming arbitrary
- `edit_command` is removed. It ships in the configuration file as
  `edit_command: code %WORKING_DIRECTORY`, parses into a field, and nothing has
  ever read that field. A configuration that still sets it is reported, by the
  mechanism `group-the-change-list` built for `change_mode`

Not in scope: remembering a typed directory for the rest of the session. The
configured default already answers "I always export to the same place", and a
third source for one value is worth avoiding.

Not in scope: moving the retired-key requirement out of `change-list-view`,
where it was introduced for `change_mode` and now serves a second key that has
nothing to do with the change list. Generalised in place; its home is a
follow-up.

## Capabilities

### Modified Capabilities
- `export-change`: the destination is configurable and editable, a missing
  directory is refused, and an existing file is confirmed
- `modal-presentation`: a modal may accept typing, which none has until now
- `change-list-view`: the retired-setting rule stops naming one key

## Impact

- `src/ui/ui.go`: the export state machine grows a prompt state and an overwrite
  state; `exportDestPath` takes a directory
- `src/ui/`: the first modal that holds a text input, and the key handling that
  implies
- `src/scanner/scan.go`: `export_dir` replaces `EditCommand` in the config, and
  `edit_command` joins the retired keys
- `src/config.yml`: `edit_command` out, `export_dir` in and documented

## Rollback

Its own commit. Reverting restores the home directory as the only destination
and restores a configuration key nothing reads.
