## Why

Bean `specgetty-4a37`. The README has no picture in it. A terminal application
is described there in eleven tables of keys, and a reader deciding whether to
install it has to assemble the screen in their head.

Two things were measured before this was written, because both would have
decided the shape of it.

**Recording is cheap.** A real thirteen-second session of `spg` at 1200x700
records to a 133 KB GIF, about ten kilobytes per second. Length was expected to
be the binding constraint and is not: the whole shot list fits in about two
minutes and a megabyte.

**Recording against this machine is not possible.** A trial tape drifted into
the project picker and produced this frame:

```
  project                          store          specs  active  archived
  linny-mcp-server                                  27      0       34
  voorzetramenshop                                  58      4       60
  nivis                            nivis            11      0       64
  registry                         registry         10      1        7
  pim-snel-curriculum-vitae                          0      0        0
  …  27 of 27 shown
```

Twenty-seven real project names, several of them client work. The picker is on
the shot list, so the recording has to run against fixtures with their own
`--config`, and that is most of the work in this change.

There is a third reason for fixtures, and it is the one that decides them. The
most interesting shot is a change's spec deltas: the outline marks which
scenarios a `MODIFIED` requirement leaves alone, and the card offers the
difference, the original and the new text. That comparison is offered only for a
change that has **not** been archived. This repository holds exactly one active
change at a time, whichever happens to be in flight, so recorded here that shot
shows something different every time and some weeks shows nothing.

## What Changes

- A fixture tree under `demo/`, recorded against instead of the machine:

```
  demo/projects/harbour-app      reads its content from the store below
  demo/projects/ferry-times      a plain project
  demo/projects/quay-signage     a plain project
  demo/stores/tideclock          the specs and changes, including one change
                                 that is permanently active
```

  The active change carries a `MODIFIED` requirement whose original is in the
  store's live specs, which is what makes the difference view reproducible.

- Seven tapes under `demo/tapes/`, one per feature rather than one long
  recording, plus a hero. A feature that changes then invalidates fifteen
  seconds of recording rather than two minutes of it.

- `make demo` records them. It generates two things first, because neither can
  be committed: the store registry, whose `local_path` must be an absolute path
  matching wherever the repository is checked out, and the directory the export
  shot writes its zip into.

- The hero leads with the spec detail view. The change list is the part a reader
  can already picture from the README's tables; an outline beside a card is the
  part nothing else does, and it reads as something in a still frame, which is
  what shows before a GIF loads.

- The README gains the hero at the top and the six others beside the sections
  they illustrate.

Not in scope:

- **Recording at release time.** `make demo` is run by hand for now. The hook in
  `release.sh` is a later decision and this change does not take it.
- **Any change to the application.** The demo records what `spg` already does.
  A recording that needed a behaviour change would be a demo of something else.

## Capabilities

No capability changes. This is fixtures, tapes and a make target, so the change
sets `skip_specs: true` rather than inventing a requirement to carry it.

## Impact

- `demo/`: new, the fixture tree, the tapes, the recording config
- `Makefile`: a `demo` target
- `README.md`: the hero and six section images
- `.gitignore`: the generated registry, the export directory and the output

## Rollback

Its own commit, and it touches nothing the application reads. Reverting removes
the pictures and the fixtures, and `spg` behaves identically with or without
them.
