---
# specgetty-pimg
title: left and right leak out of the spec detail view
status: completed
type: bug
priority: normal
created_at: 2026-09-21T20:37:21Z
updated_at: 2026-09-21T21:02:33Z
---

Pressing `left` or `right` in the spec detail view yanks the keyboard from
the card back to the outline, and moves a tab you cannot see.

`levelSpec` has no arm in the arrow-key handlers, so both fall through to the
project branch:

```go
case "right":
    if m.level == levelChange { ... }
    else if m.detailTab < len(tabNames)-1 {
        m.detailTab++                 // a tab that is not on screen
        m.focus = m.defaultFocus()    // focusListPane, the outline
        cmds = append(cmds, m.enterTab()...)
    }
```

## What happens

- `m.focus` is reset to `focusListPane`, so a reader with the keyboard on the
  card loses it mid-read. This is the visible symptom.
- `m.detailTab` advances underneath. `esc` from `levelSpec` forces
  `detailTab = tabSpecs`, so the damage is masked on the way out and the bug
  reads as a focus glitch rather than as state corruption.
- `enterTab()` is guarded on `levelProject`, so the schema subprocess does not
  fire. That part is safe.

## Expected

`left` and `right` do nothing at `levelSpec`. There are no sibling views to
move between: the view is one spec, and the two halves are reached with `tab`.

The comment on the handler already states the rule the code does not follow:

> Sub-tabs belong to an open change, the tab bar to the project. Neither spills
> into the other.

A third level spills into the second.

## Notes

- Shipped in `open-a-spec-in-detail` (commit `af367e7`). No test presses either
  key at `levelSpec`.
- `specgetty-zqhs` wants `left`/`right` to drive an old/new/diff row on the card
  of a change-spec detail view, so what these keys mean at a spec level has to
  be settled before that lands.
- The `else if` shape is what let this through: every arm added to `level` has
  to remember to guard, rather than the handler switching on the level. Worth
  considering a `switch m.level` here instead.

## OpenSpec change

`keep-the-arrows-in-the-spec-view`, tinychange schema, validated strict.

Spec delta adds one requirement to `spec-detail-view`: the arrow keys do
nothing in the view. It mirrors `change-list-view`'s existing requirement
`Artifact sub-navigation is scoped to the open change`, one level down.

Task 1.2 replaces the `else if` chain with a `switch m.level`, so a level added
later has to state what these keys do rather than inheriting the project tab
bar by omission. That shape is what let this through.

Schema source: https://github.com/speclib/openspec-tinychange-schema


## Summary of Changes

Shipped as OpenSpec change `keep-the-arrows-in-the-spec-view`, commit `f4700a7`.

`left` and `right` now do nothing at `levelSpec`. The view shows one spec with
no sibling to move to, and its two halves are reached with `tab`.

Both handlers switch on the level rather than chaining, which is task 1.2 and
the more important half. The `else if` shape gave every level it did not name
the project tab bar by omission, so a level added later inherited behaviour
nobody chose for it. A `switch m.level` makes the next level state its answer.

The revert check bit on exactly the two tests the tasks predicted: the
unchanged-state assertion and the card keeping the keyboard. The build was
confirmed clean before reverting, as task 2.5 asks.

Task 2.4's second half nearly shipped as a skip: `taskModel` builds a change
with one artifact, so there was nowhere for the keys to go and the subtest
skipped rather than asserted. It now uses the fixture carrying `proposal.md` and
`tasks.md`, and fails rather than skips if a fixture ever loses its second tab.

Coverage: ui 87.6%, total 89.3%, floors unchanged. Gate green.

Next: `follow-the-spec-grammar` gives this level a second state for a file that
does not parse, and it is now written on top of arrow keys that behave.
