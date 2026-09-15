## Context

`specgetty-g4pa` wants the change list at full width, with a drill-down into a
single change.

Today the UI tracks position with a `zoomed` boolean, an `activeView` int, a
`detailTab` int and four independent cursors. Adding a third depth as a fourth
boolean would make every key handler a longer conditional. An explicit depth
value replaces the boolean.

`specgetty-jdif` separately proposes dropping the projects list as a persistent
view and making the project picker a popup over a single project view. That epic
decides what sits above the change list. So this change owns the change list and
the change below it, and asserts nothing about how a user reaches a project.
`specgetty-edee` and `specgetty-oarh` describe escape behaviour and a `--view=`
flag for the model jdif replaces, and are left for jdif to resolve.

## Goals / Non-Goals

**Goals**

- Separate the project tab bar and a change's artifact sub-nav onto different
  levels so they stop sharing left/right
- One list for active and archived changes
- Find a change by name or by what it says, inside one project

**Non-Goals**

- Anything above the project view. The projects list, the picker popup and the
  startup view belong to `specgetty-jdif`.
- Cross-project search. `specgetty-opwv` proposes a generic query language over
  all projects; that decision is deliberately left open and nothing here
  forecloses it.
- Columns needing new scanner data. Tracked in `specgetty-vru8`.
- A runtime column picker. Field selection is config and CLI only.

## Decisions

### Navigation depth as an explicit level

```
  project        tab bar: changes | specs | config
                 the changes tab is the full-width change table
     | enter                ^ esc
     v                      |
  change         artifact sub-tabs own left/right
```

`enter` descends, `esc` ascends. The model carries an explicit depth instead of
a `zoomed` boolean, so a third depth does not become a fourth flag.

This also removes `enter` meaning "zoom in" at the project list and "zoom out"
one level in. Ascending is `esc` only.

What sits above the project view is left exactly as it behaves today. This
change does not define a floor, because `specgetty-jdif` is about to redefine
that end of the stack.

**Alternative rejected**: keeping `zoomed bool` and adding `changeOpen bool`.
Two booleans encode four states, one of which (not zoomed, change open) is
meaningless and has to be guarded everywhere.

### One seam for merge, mode filter and query filter

Every consumer of the change list already routes through `currentChanges()` and
`currentArchivedChanges()`: the renderer, the `a`/`d`/`e` handlers, the confirm
modals, the nav bar hints, and the sub-tab counters. None of them reads
`Info.Changes` directly.

So the merge, the open/archived/both filter and the search query all apply in
one place, and every downstream caller keeps working without knowing a filter
exists. Actions keep operating on the row under the cursor because the cursor
indexes the same slice the renderer drew.

### Fuzzy for names, substring for bodies

`ChangeInfo` already holds `ArtifactContents` and `SpecContents` in memory after
a scan, so searching inside changes costs no extra I/O. But fuzzy subsequence
matching over a 3KB proposal matches nearly everything, which makes it useless
at document length. The two targets need different matchers:

| Query form | Target                  | Matcher           |
| ---------- | ----------------------- | ----------------- |
| `export`   | change name             | fuzzy subsequence |
| `'export`  | change name             | literal substring |
| `:export`  | artifact and spec text  | literal substring |

`github.com/sahilm/fuzzy` is already in the module graph through `bubbles`, and
is the matcher `bubbles/list` itself uses. It offers `Find` (ranked by score)
and `FindNoSort` (original order preserved). Score-ranked is chosen, so the best
match sits under the cursor, matching what `fzf` does.

Case sensitivity is smart-case: an all-lowercase query matches
case-insensitively, any uppercase character makes the whole query
case-sensitive. Same rule as ripgrep and fzf.

### Search input mode borrows fzf's split focus

Type-to-filter is not available because `a`, `d`, `e`, `s`, `l` and `q` are
actions. So `/` opens an input mode where letter keys belong to the text input,
while `up`/`down` and `ctrl+p`/`ctrl+n` still move the list cursor. That lets a
user type, arrow down and press `enter` to open a change without leaving the
prompt.

`enter` from input mode descends into the highlighted change rather than
committing the filter and returning to list navigation. A committed-filter state
is a mode the user can forget they are in; making `esc` the only exit means the
prompt is always visible while a filter is active.

**Cursor stability**: the selected change is remembered by name, not by index.
After a keystroke narrows the list, or after the filesystem watcher triggers a
rescan, the cursor returns to the same change if it survived and clamps to the
last row otherwise. Index-only clamping would silently move the selection onto a
different change, which matters when `a`, `d` and `e` act on it.

### Body matches need to explain themselves

A row that matched on `proposal.md` text shows a change name the user does not
recognise. A table row has no space for an excerpt, so a dim trailing hint names
which artifacts matched:

```
> export-change-as-zip    17/17   1   proposal, tasks
  full-change-list         0/12   2   design
```

## Risks / Trade-offs

- **Removing the archive tab changes muscle memory.** Anyone pressing `3` for
  archive lands somewhere else. Mitigated by the mode filter defaulting to
  open-only, so the changes tab behaves as it does today until the filter is
  touched.
- **Body search over every artifact of every change is linear in text size per
  keystroke.** Acceptable for one project's changes at TUI scale; it would not
  be acceptable if the same matcher were later pointed at all scanned projects,
  which is worth remembering if `specgetty-opwv` is revived.
- **One more level to escape from.** Reading a proposal now costs an `enter`
  that it did not before. Full-width rendering is what buys that keystroke back.

## Open Questions

- The delete action from the original bean still needs a key. `d` is taken by
  discard, and `specgetty-tyri` covers the same ground. Not resolved here.
- The `detail-tabs` spec is stale relative to the code: it describes overview
  and search tabs and `h`/`l` keys that no longer exist. This change touches
  only the tab set, and leaves the rest of that drift alone. The `h`/`l`
  scenario is carried forward verbatim because a MODIFIED requirement replaces
  its whole block, so dropping it here would have been a silent spec change
  rather than a decision. Worth its own cleanup.
