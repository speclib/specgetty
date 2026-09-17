## Context

Three measurements taken against the real repository shaped this.

**Most task lines wrap.** At an 86 column pane, across 511 task lines in the
archive:

| screen rows | tasks | share |
| ----------- | ----- | ----- |
| 1           | 169   | 33%   |
| 2           | 308   | 60%   |
| 3           | 31    | 6%    |
| 4           | 3     | ~0%   |

So a task is not a screen row. Selecting and highlighting "the line" means
highlighting one to four rows, and moving the cursor means moving between source
lines, not between rows.

**The current pipeline throws that away on purpose.** `renderMarkdown` returns
already-wrapped final rows so the viewport can slice them and report a position
that matches the screen. Nothing remembers which rows came from which source
line.

**specgetty has never written an artifact's contents.** Its writes today are its
own cache, a directory rename for discard, and a zip written elsewhere. This is
the first time it edits a file the user also edits.

## Goals / Non-Goals

**Goals**

- Tick a task off without leaving the tool
- Never discard an edit made in an editor
- Leave every other document exactly as it behaves now

**Non-Goals**

- Colouring the checkbox separately from its text. Deferred to a wider markdown
  highlighting change.
- Editing anything else: no text entry, no adding or removing tasks, no editing
  other artifacts.
- Undo. `space` is its own inverse, and a second press restores the previous
  state.

## Decisions

### One structure serves the highlight and the save

The renderer keeps, for each source line, the text it started from and the range
of screen rows it produced:

```
  { sourceIndex, sourceText, rowStart, rowEnd }
      │            │           └──────┴── the highlight covers these rows
      │            └── the key the save looks the line up by
      └── a hint for disambiguating, never the thing written to
```

The highlight needs the row range. The save needs the source text. Building one
map rather than two keeps them from disagreeing.

### The save is content-addressed, not index-addressed

The obvious approach is to re-read the file and flip the *n*th checkbox. That
survives a re-read but not an insertion above the cursor:

```
  on screen        meanwhile, saved from an editor
  ▢ 1.1   ←        ▢ 0.9     ← the nth checkbox is now this one
  ▢ 1.2            ▢ 1.1
                   ▢ 1.2
```

The wrong task gets ticked and nothing says so. So the save works like this:

1. read `tasks.md` fresh from disk
2. find the line equal to the remembered source text
3. not found, or found more than once, refuse and say so
4. flip `[ ]` and `[x]` in that line
5. write atomically

The refusal path is the point. It turns "silently ticked the wrong task" into a
message telling the reader to rescan. Duplicate lines are the only ambiguity and
real tasks carry numbers, so a collision needs two byte-identical lines.

**Why writing the whole buffer back is still safe.** The bytes were read a
moment ago and one character changed, so the window for losing someone else's
write is the microseconds between read and rename, not the minutes since the
last scan. Byte-offset patching would reintroduce the stale-offset problem it
was meant to avoid.

The refusal has somewhere to speak: the status line added by
`copy-change-name-and-path` is the right weight for it.

### Atomic means same directory, and the mode must be carried

`os.Rename` is atomic only within a filesystem, so the temporary file goes
beside `tasks.md`, never in `/tmp`.

The mode has to be read from the original and applied to the replacement.
`mktemp` creates `0600`, and moving that over a `0644` file quietly tightens the
permissions. Git tracks only the executable bit, so it never shows in a diff.
This exact bug appeared in `scripts/release.sh` during the session that added
it, and was found by testing rather than by reading.

### The glyphs are `▢` and `▣`

Both are one cell wide, both come from Geometric Shapes, and neither has an
emoji presentation.

The obvious pair, `☐` and `☑`, was tried first and rejected: Ghostty draws
`☑` as a coloured emoji while `☐` stays a text glyph, so the two sizes do not
match and the checked state is wider than `ansi.StringWidth` reports. A wrapper
that counts one cell while the terminal draws two drifts a column on every
wrapped line, which is the same silent misalignment the picker box had.

Dropping from three cells (`- [ ]`) to one also returns two columns to the text,
which moves one-row tasks from 33% to 36%.

Because colour is deferred, the shape carries the state alone. `▢` and `▣`
differ by a filled centre, which reads at a glance and survives the selected
row, where the green background would swallow a colour difference anyway.

### The rescan closes the loop

Writing `tasks.md` makes the filesystem watcher fire, which rescans the project
and re-renders. The document key does not change, so the reading position is
kept, and the counts in the change list update from the file rather than from
anything held in memory. The write is therefore self-verifying, and there is no
second write to loop on.

## Risks / Trade-offs

- **One keystroke edits a tracked file.** There is no confirmation and no undo,
  by decision. `space` is its own inverse, which makes a mistake one keypress to
  correct, but it is still the first key in specgetty that changes a file the
  user owns.
- **`j` and `k` mean two distances.** In `tasks.md` they move a source line,
  which may be four screen rows. In `proposal.md` they still move one row. The
  alternative, a separate key for the cursor, spends a binding on something the
  arrows already mean everywhere else.
- **A wrapped task's highlight is several rows tall.** On a short pane a
  four-row task plus its neighbours may fill most of the view.

## Open Questions

- Whether the cursor stops on every source line or only on checkbox lines.
  Stopping everywhere matches "the lines should be selectable" and is
  predictable, at the cost of pressing `j` through blank lines and headings
  between task groups. Stopping only on checkboxes makes ticking faster and
  makes prose unreachable by cursor. This design assumes every source line, and
  it is the cheapest thing to change after using it once.
