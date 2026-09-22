package ui

import (
	"errors"

	tea "charm.land/bubbletea/v2"

	"os"
	"path/filepath"
	"strings"
)

// The shapes a task line can take.
//
// These are the prefixes the scanner counts (`^- \[ \] ` and `^- \[x\] `, at
// column zero), and matching them exactly is deliberate: an indented or `*`
// prefixed checkbox is not counted in the totals, so drawing it as a box or
// letting it be toggled would make the display disagree with the numbers beside
// it.
const (
	taskUncheckedPrefix = "- [ ] "
	taskCheckedPrefix   = "- [x] "
)

// The glyphs.
//
// U+25A2 and U+25A3, both one cell wide and both from Geometric Shapes, so they
// are drawn as a matched pair. The obvious choice, U+2610 and U+2611, was tried
// first and rejected: the checked one has an emoji presentation, so terminals
// draw it larger and wider than `ansi.StringWidth` reports, and a renderer that
// counts one cell while the terminal draws two drifts a column on every wrapped
// line.
//
// They differ by a filled centre rather than by colour, which keeps the state
// readable on the highlighted row, where the background would swallow a colour
// difference.
const (
	checkboxUnchecked = "▢"
	checkboxChecked   = "▣"
)

// isTaskLine reports whether a source line carries a checkbox this tool owns.
func isTaskLine(line string) bool {
	return strings.HasPrefix(line, taskUncheckedPrefix) ||
		strings.HasPrefix(line, taskCheckedPrefix)
}

// continuesTaskItem reports whether a line belongs to the task above it.
//
// A task is rarely one line. OpenSpec's own generators wrap a task onto as many
// lines as it needs and indent the rest, so the checkbox line and its
// continuations are one thing to a reader and have to be one thing to the
// cursor:
//
//   - [ ] 1.1 Add `.github/workflows/check.yml` running on push and on
//     installing nix and running `nix flake check` and nothing else
//
// A continuation is indented and not blank. Everything else starts something
// new: an unindented line, a blank line, a heading, or the next checkbox.
//
// An indented `- [ ]` is a continuation rather than a task of its own, which
// falls out of the prefixes above matching at column zero and is the answer
// this file already gives everywhere else: the scanner does not count such a
// line, so a cursor that stopped on it would offer to toggle something the
// totals do not know about.
func continuesTaskItem(line string) bool {
	if isTaskLine(line) {
		return false
	}
	if strings.TrimSpace(line) == "" {
		return false
	}
	return line[0] == ' ' || line[0] == '\t'
}

// toggleTaskLine returns the line with its checkbox flipped.
func toggleTaskLine(line string) (string, bool) {
	if rest, ok := strings.CutPrefix(line, taskUncheckedPrefix); ok {
		return taskCheckedPrefix + rest, true
	}
	if rest, ok := strings.CutPrefix(line, taskCheckedPrefix); ok {
		return taskUncheckedPrefix + rest, true
	}
	return line, false
}

// writeFileAtomic replaces path with data in a single step.
//
// The temporary file is created beside the target, because os.Rename is atomic
// only within a filesystem; a temp file in /tmp would be a copy across a
// boundary and lose the guarantee.
//
// The original's permissions are carried onto the replacement. os.CreateTemp
// makes 0600, and moving that over a 0644 file quietly tightens it. Git tracks
// only the executable bit, so the change would never appear in a diff. That bug
// shipped in scripts/release.sh and was found by testing, not by reading.
func writeFileAtomic(path string, data []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	// Any failure from here on leaves the original untouched.
	cleanup := func() { _ = os.Remove(tmpName) }

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return err
	}
	if err := os.Chmod(tmpName, info.Mode().Perm()); err != nil {
		cleanup()
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		cleanup()
		return err
	}
	return nil
}

// Errors a toggle can refuse with. They are distinguishable so the message can
// say which happened, since the two mean different things to the reader.
var (
	errTaskLineGone      = errors.New("that task is no longer in the file")
	errTaskLineAmbiguous = errors.New("that task appears more than once")
)

// toggleTaskInFile flips the checkbox on the line matching want.
//
// The file is read here rather than taken from what was scanned earlier, and
// the line is found by its exact text rather than by counting checkboxes.
// Counting survives a re-read but not an insertion above the cursor: the
// positions shift, the wrong task is ticked, and nothing says so. Matching on
// text either finds the line the reader is looking at or refuses.
func toggleTaskInFile(path, want string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// A trailing newline produces a final empty element that must survive the
	// round trip, so the file does not grow or shrink a byte on every toggle.
	lines := strings.Split(string(raw), "\n")

	found := -1
	for i, line := range lines {
		if line == want {
			if found >= 0 {
				return errTaskLineAmbiguous
			}
			found = i
		}
	}
	if found < 0 {
		return errTaskLineGone
	}

	flipped, ok := toggleTaskLine(lines[found])
	if !ok {
		return errTaskLineGone
	}
	lines[found] = flipped

	return writeFileAtomic(path, []byte(strings.Join(lines, "\n")))
}

// toggleSelectedTask flips the checkbox under the document cursor and reports
// what happened, for the status line.
//
// The refusals matter more than the success: a toggle that silently ticked the
// wrong task, or that quietly did nothing, would be worse than one that says it
// could not.
func (m model) toggleSelectedTask() (string, tea.Cmd) {
	sel, ok := m.selectedTask()
	if !ok || m.docPath == "" {
		// No task under the cursor, which now means no task in the document at
		// all: the cursor selects tasks and nothing else. Saying nothing is
		// right, the key not applying here, and a message for every stray press
		// would be noise.
		return "", nil
	}

	switch err := toggleTaskInFile(m.docPath, sel.text); {
	case err == nil:
		// Success needs no message: the box changes shape, which is the
		// feedback. It does need a rescan, so the box and the task counts come
		// from the file rather than from what was read earlier. The filesystem
		// watcher would usually do this, but it is allowed to fail to start,
		// and then nothing on screen would move.
		if len(m.repoPaths) > 0 && m.cursor < len(m.repoPaths) {
			return "", m.doScanSingle(m.repoPaths[m.cursor])
		}
		return "", nil
	case errors.Is(err, errTaskLineGone):
		return "tasks.md changed on disk, press s to rescan", nil
	case errors.Is(err, errTaskLineAmbiguous):
		return "that task line appears more than once, edit it by hand", nil
	default:
		return "could not save: " + err.Error(), nil
	}
}
