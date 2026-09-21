## Why

Bean `specgetty-ve3m`. The properties sidepanel is a flat list of four labels in
fifteen columns, with two thirds of the panel empty and nothing saying which
label belongs with which:

```
 ╭─────────────╮ ╭────────────────────────────────────────────────╮
 │ project     │ │ # /stores/nivis-tunnel/openspec/project.md     │
 │ spec-driven │ │                                                │
 │ tinychange  │ │ # Tunnel                                       │
 │ store       │ │                                                │
 │             │ │ A project.                                     │
 │             │ │   ... seven more empty rows ...                │
 ╰─────────────╯ ╰────────────────────────────────────────────────╯
```

`project` and `store` both say where the project's content comes from; the rows
between them are schemas. Nothing on screen groups them, so the list reads as
four unrelated words.

### The rule that would block this does not protect what it claims

`config-tab-display` says the list SHALL take the width its labels need "so that
the content keeps the room its paths require". The longest line the content pane
ever shows is a store root of 60 columns. Measured against both policies:

| terminal | list 17 (labels) | list 20 (floor) | does a 60-column path fit? |
| -------- | ---------------- | --------------- | -------------------------- |
|       60 |       content 38 |      content 37 | no, under either           |
|       80 |       content 58 |      content 55 | no, under either           |
|      100 |       content 78 |      content 75 | yes, under either          |
|      140 |      content 118 |     content 115 | yes, under either          |

The policy never decides it: the path wraps at 80 whatever the list does, and
fits at 100 whatever the list does. The three columns the sections want are free,
and the stated reason for keeping the list minimal does not hold at the sizes
where it would matter.

### The third list of its kind

A list whose cursor indexes items while its pane counts drawn lines now exists
twice, and this is the third:

```
  the change list    cursor: changes    drawn: rows, headers, spacers   lineOfRow()
  the spec outline   cursor: nodes      drawn: wrapped label rows       nodesInRows()
  the properties     cursor: sections   drawn: rows and headers         this change
```

Two instances is a coincidence. Three is where the arithmetic should be one
thing, and it is the same thing: one drawn line per item plus chrome, with the
scroll offset computed in lines rather than in items.

## What Changes

- The properties list is grouped into sections with headers that the cursor
  cannot land on: `PROJECT` holding the configuration and the store, and
  `SCHEMAS` holding one row per schema. Headers are uppercase and carry no
  count, rows are indented under them, and a blank line separates the groups.

- The row labelled `project` becomes `config`. It shows the configuration file
  and the path it came from, which is what the new name says and the old one did
  not.

- The list takes a floor of twenty columns rather than only the width its labels
  need, keeping the one-third cap it already has. A proportional split like the
  specs tab uses is the wrong instrument here: at 200 columns it would hand 58
  columns to labels that are thirteen characters long.

- The line arithmetic behind all three lists becomes one helper over a slice of
  item indices, one entry per drawn line. The change list and the spec outline
  move onto it first, where the tests they already have prove the move changes
  nothing, and the new list is built on it rather than duplicating it a third
  time.

Not in scope:

- **The problems row.** `specgetty-p44t` will report the specs of a project that
  openspec cannot read, and that row belongs under `PROJECT` beside `config` and
  `store`. This change leaves it a seat and does not build it.
- **Counts on the headers.** `ACTIVE (3)` earns its count on a list of 47
  changes; `PROJECT (2)` does not.

## Capabilities

### Modified Capabilities
- `config-tab-display`: the list is grouped into sections with headers the
  cursor skips, and it is no longer sized only to its labels

## Impact

- `src/ui/properties.go`: sections gain a group, and the rows a label change
- `src/ui/ui.go`: `renderPropertiesTab` draws headers, `propertiesSplit` gains
  its floor
- `src/ui/lines.go`: new, the shared line arithmetic
- `src/ui/grouped.go`, `src/ui/listpaging.go`, `src/ui/specdetail.go`: moved onto
  the helper, with no change in behaviour
- `README.md`: the properties tab picture

## Rollback

Its own commit. Reverting returns the flat list and puts the arithmetic back in
three places, one of which would not exist to be put back.
