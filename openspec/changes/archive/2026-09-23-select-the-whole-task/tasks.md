## 1. The task item

- [x] 1.1 Add the boundary rule: a run starts at a `- [ ] ` or `- [x] ` line at
      column zero and continues while the next line begins with whitespace and
      is not itself a checkbox, ending at the first unindented line, blank line,
      heading or next checkbox
- [x] 1.2 Test the boundary against this repository's own tasks files, and
      against the awkward cases: a task at the end of the file with no trailing
      blank, two checkboxes with nothing between them, an indented `- [ ]` which
      is a continuation rather than a task, and a file whose only checkbox is on
      its last line
- [x] 1.3 Build the `itemLines` slice for the rendered tasks document, marking
      headings, blank lines and the `Tasks: n/m complete` prefix as chrome, the
      way `ownersOfRows` does it for the spec outline
- [x] 1.4 Verify the slice accounts for the prefix rows the artifact renderer
      adds, so a mapping is not off by two the moment a tasks file has any task

## 2. The cursor selects a task

- [x] 2.1 Make the cursor index items rather than source lines, so `j` and `k`
      step task to task and pass over headings and blank lines
- [x] 2.2 Highlight every row of the selected item as one unbroken band,
      through `span`, and verify no row of a neighbouring task, heading or blank
      line is caught in it
- [x] 2.3 Scroll the whole item into view through `offsetFor`, and verify an
      item taller than the pane shows its first row rather than its last
- [x] 2.4 Verify a `tasks.md` with no checkbox at column zero has no cursor and
      scrolls by rows, which should need no special case
- [x] 2.5 Verify the toggle still rewrites the checkbox line and leaves every
      continuation line byte for byte as it was

## 3. The keys reach the cursor

- [x] 3.1 Add `docPage()` beside `listPage()`, returning how many items fill one
      pane of rows, through `fit`
- [x] 3.2 Route `pgdown`, `ctrl+f`, `pgup`, `ctrl+b`, `ctrl+d` and `ctrl+u`
      through the same three-way split `j` and `k` already use, so a cursored
      document moves its cursor and a cursorless one moves its rows
- [x] 3.3 Route `gg` and `G` the same way, landing on the first and last item
- [x] 3.4 Verify the cursor is never left off screen after any of these keys,
      which is the fault being fixed and the one assertion that must not be
      omitted
- [x] 3.5 Verify a page moves further in a taller terminal, and that no document
      without a cursor changed behaviour at all

## 4. The refresh is silent

- [x] 4.1 Split the scanning flag in two: one for a first scan with nothing on
      screen, one for a refresh of a project already shown
- [x] 4.2 Raise the scan indicator on the first only, and verify a refresh
      leaves what is on screen where it is
- [x] 4.3 Let every key through during a refresh, and verify a burst of toggles
      loses none of them
- [x] 4.4 Route the toggle's own rescan through the refresh flag, so a watcher
      event arriving during it coalesces through the existing `scanPending` path
      rather than queueing a third read
- [x] 4.5 Verify the startup path is untouched: an empty screen still says what
      it is doing and still answers only `q` and `ctrl+c`
- [x] 4.6 Verify a key handled during a refresh cannot cause a wrong write, the
      toggle re-reading the file and refusing a line it cannot identify

## 5. The nav bar

- [x] 5.1 Advertise `space` on a pane that has a task cursor, and not elsewhere
- [x] 5.2 Describe `jk/↑↓` as navigating on that pane and as scrolling on the
      others
- [x] 5.3 Place `space` ahead of the paging hints, so a narrow terminal drops
      `gg/G` before it drops the toggle
- [x] 5.4 Verify the hints against a narrow terminal, the nav bar dropping from
      the end rather than colliding with the version on the right

## 6. Closing

- [x] 6.1 Run the gate, `nix flake check`, including the coverage ratchet
- [x] 6.2 Record the change in `CHANGELOG.md`
- [x] 6.3 Re-record the demo if the tasks pane appears in it
