## Why

Bean `specgetty-0ia0`, a regression from `properties-tab`.

The config tab was one document, so `docActive()` returned true for it whatever
held the keyboard, and `pgup`, `pgdown`, `ctrl+f`, `ctrl+b`, `ctrl+d`, `ctrl+u`,
`gg` and `G` all scrolled it at once. Making it a split gave it a list that takes
the keyboard first, and those keys now do nothing until `tab` is pressed.

```
  before                          after properties-tab
  config tab                      properties tab
  focus: detail                   focus: the row list
  pgdown -> scrolls               pgdown -> nothing
```

`document-viewer` already says which of these is right:

> **Move one page**: WHEN a document viewer **is displayed** and the user
> presses `pgdown`, `ctrl+f`, `pgup` or `ctrl+b`, THEN the document SHALL scroll

Displayed, not focused. The specs tab has under-delivered against that same
requirement since it shipped; nobody noticed because its list is the half you
came to use, where the properties tab's list is three fixed rows and the
document is the point.

## What Changes

- The paging and jump keys act on the document whenever one is displayed,
  whichever half of a split tab holds the keyboard
- `j` and `k` are untouched. They have to choose between moving a list and
  scrolling a document, and which half holds the keyboard is what decides that
- The reading position in the panel title is untouched. `specs-tab` says it is
  reported only while the content holds the keyboard, which is deliberate and
  which nobody reported a problem with

Not in scope: where a split tab puts the keyboard when it opens. Both still open
on their list, which is what chooses what the content shows.

## Capabilities

### Modified Capabilities
- `document-viewer`: the paging keys are stated to be independent of focus,
  rather than leaving it to be inferred from the word "displayed"
- `config-tab-display`: the properties tab pages from either half
- `specs-tab`: so does the specs tab, which is the same fix

## Impact

- `src/ui/docview.go`: a `docDisplayed` predicate beside `docActive`
- `src/ui/ui.go`: six key handlers ask the new one

## Rollback

Its own commit, and it restores nothing anyone wants: reverting puts the
regression back.
