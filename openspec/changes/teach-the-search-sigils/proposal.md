## Why

Bean `specgetty-wxqm`. The search grammar has three matchers and two of them are
invisible:

```
  query        matches                                      how you find out
  ──────────   ─────────────────────────────────────────    ────────────────
  expzip       fuzzy subsequence on the name                it just works
  'export      literal substring on the name                you read the README
  :export      literal substring in every artifact and      you read the README
               spec file in the change
```

The project picker has the same grammar, where `:` searches file paths and
contents across every project found on the machine, which is the most useful
thing in the application and the least likely to be guessed at.

`group-the-change-list` made this worse and better at once. The change list now
always holds the archive, so `:` went from searching one active change to
searching 39 of them, 167 files and 462 KB in this project. That is where the
answers usually are, and nothing on screen suggests the key that reaches them.

## What Changes

- The search prompt names the three matchers while it is focused and empty,
  which is the moment between pressing `/` and knowing what to type:

```
  /_  fuzzy name  'exact  :inside            38 of 38 shown
```

- The legend is dropped rather than wrapped when the prompt is too narrow for
  it, the way the nav bar already drops hints
- A name search that matched nothing suggests the body search with the same
  term. The legend teaches at the point of intent and is gone on the first
  keystroke; this teaches at the point of frustration, when the user is already
  looking for another way:

```
  No changes match "inotify"
  try :inotify to search inside the changes themselves
```

- A `:` search that matched nothing suggests nothing, because there is nothing
  further to try
- Both surfaces get both, since they share one prompt renderer and one grammar
- The match hint is separated from the table by the same gap that separates
  every other column. It is currently appended flush, which nothing showed until
  `date` became a default column this afternoon:

```
  before   box-the-tab-content   16/16   1   2026-09-18design
  after    box-the-tab-content   16/16   1   2026-09-18  design
```

Not in scope: a matcher toggle. A key that cycles name, literal and contents
would make the sigils unnecessary rather than teaching them, and it would leave
two ways to do one thing. It is the alternative to this change, not a companion.

## Capabilities

### Modified Capabilities
- `change-search`: the prompt names the matchers, and a failed name search
  suggests the body search
- `project-picker`: the same, in the picker's own prompt and empty message
- `change-list-view`: the match hint is separated from the table it sits beside

## Impact

- `src/ui/table.go`: `renderSearchPrompt` takes a width and renders the legend;
  `renderTable` separates the hint
- `src/ui/grouped.go`: the same separation in the change list's own renderer
- `src/ui/ui.go` and `src/ui/picker.go`: the two empty messages gain a suggestion

## Rollback

Its own commit. Reverting removes the teaching and restores the flush hint,
which is what ships today.
