## 1. The line map

- [x] 1.1 Have `renderMarkdown` return, alongside its rows, an entry per source line recording the source text and the range of screen rows it produced
- [x] 1.2 Keep the existing string-returning behaviour for every caller that does not need the map, so the specs tab and the config tab are untouched
- [x] 1.3 Store the map on the model beside the document, and rebuild it whenever the document is rebuilt, so the two cannot drift
- [x] 1.4 Check the map against a wrapped line, a line that does not wrap, and a blank line

## 2. Rendering the boxes

- [x] 2.1 In `styleMarkdownLine`, draw `- [ ]` as `▢` and `- [x]` as `▣`, matching the same shape the scanner counts (`^- \[ \] ` and `^- \[x\] `, at column 0)
- [x] 2.2 Leave any other list item alone, including an indented or `*` prefixed one, since the scanner does not count those either and disagreeing with it would be worse than ignoring them
- [x] 2.3 Confirm both glyphs measure one cell with `ansi.StringWidth`, which is what the wrapper counts

## 3. The cursor

- [x] 3.1 Add a cursor holding a source line index, active only for a change's tasks artifact
- [x] 3.2 Highlight every screen row the selected line occupies, and no row belonging to another line
- [x] 3.3 Move the cursor on `j`, `k`, up and down when it is active, leaving those keys scrolling by a row everywhere else
- [x] 3.4 Scroll the viewport so the whole selected line is visible after a move, including a line that is taller than one row
- [x] 3.5 Clamp at the first and last source line
- [x] 3.6 Reset the cursor when the document changes, by the same rule that already resets the scroll position
- [x] 3.7 Leave `pgup`, `pgdown`, `ctrl+f`, `ctrl+b`, `ctrl+d`, `ctrl+u`, `gg` and `G` moving the view, not the cursor

## 4. The toggle and the save

- [x] 4.1 Bind `space`, which is unbound today; note that bubbletea v2 reports it from `String()` as `"space"`
- [x] 4.2 Do nothing when the selected line carries no checkbox
- [x] 4.3 Read `tasks.md` from disk at the moment of the toggle rather than using the copy held in memory
- [x] 4.4 Find the line by matching the remembered source text exactly, not by counting checkboxes, so an insertion above the cursor cannot redirect the edit
- [x] 4.5 Refuse and report when the line is absent from the file as it now stands
- [x] 4.6 Refuse and report when the line occurs more than once, since the intended one cannot be identified
- [x] 4.7 Flip `[ ]` and `[x]` on that line and write the whole freshly read buffer back
- [x] 4.8 Write atomically: a temporary file beside `tasks.md`, then rename over it
- [x] 4.9 Carry the original file mode onto the replacement. `mktemp` creates `0600`, and git tracks only the executable bit, so tightening `0644` to `0600` would never appear in a diff. This exact bug shipped in `scripts/release.sh` and was caught by testing, not by reading
- [x] 4.10 Report a failed write on the status line and leave the file untouched

## 5. Tests

- [x] 5.1 The glyphs appear for checked and unchecked tasks, and both measure one cell
- [x] 5.2 The highlight covers every row of a wrapped task and none of its neighbour's
- [x] 5.3 `j` and `k` move the cursor a source line at a time and clamp at both ends
- [x] 5.4 A move to a line below the fold scrolls it fully into view
- [x] 5.5 `space` on an unchecked task checks it on disk, and on a checked task unchecks it
- [x] 5.6 `space` on a line with no checkbox writes nothing
- [x] 5.7 The save survives an edit made after the scan: write extra tasks to the file behind specgetty's back, toggle, and assert the extra tasks are still there
- [x] 5.8 A toggle whose line has been removed refuses and reports, leaving the file byte-identical
- [x] 5.9 A toggle whose line appears twice refuses and reports
- [x] 5.10 The file mode is the same before and after a toggle
- [x] 5.11 A read-only file reports a failure and leaves the contents unchanged
- [x] 5.12 No cursor and no `space` behaviour in the specs tab, the config tab or a non-tasks artifact

## 6. Verification

- [x] 6.1 `go build ./...` and `go vet ./...` clean
- [x] 6.2 Grep the test output for FAIL rather than counting PASS lines, which cannot see a failure
- [x] 6.3 `bash scripts/coverage-gate.sh` passes and the floors are raised
- [x] 6.4 `nix flake check` passes
- [x] 6.5 Toggle a task in this repository's own `tasks.md`, then `git diff` it: exactly one line changed, and only within the brackets

      Done against `archive/2026-09-16-upgrade-bubbletea-v2/tasks.md`:
      `1 file changed, 1 insertion(+), 1 deletion(-)`, mode still 644, reverted
      afterwards.
- [x] 6.6 Run `spg`, move the cursor through a group of wrapped tasks and tick one. The highlight covering several rows, and the glyphs in a real font, are things no test here can see

      Handed over, not done: this one needs a terminal. What was checked
      instead is the rendered output in a test, where the highlight does span
      both rows of a wrapped task and both glyphs measure one cell.
