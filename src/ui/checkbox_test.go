package ui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

func tasksFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "tasks.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// --- the glyphs ---

func TestCheckboxGlyphsAreOneCell(t *testing.T) {
	// The renderer wraps by counting display cells. A glyph the terminal draws
	// wider than this reports drifts a column on every wrapped line.
	for _, g := range []string{checkboxUnchecked, checkboxChecked} {
		if w := ansi.StringWidth(g); w != 1 {
			t.Errorf("%q is %d cells, want 1", g, w)
		}
	}
	if checkboxUnchecked == checkboxChecked {
		t.Error("the two states must differ in shape")
	}
}

func TestTaskLineRendersAsABox(t *testing.T) {
	unchecked := renderMarkdown("- [ ] 1.1 do the thing", 80)
	if !strings.Contains(unchecked, checkboxUnchecked) {
		t.Errorf("unchecked task lost its box: %q", unchecked)
	}
	if strings.Contains(unchecked, "[ ]") {
		t.Errorf("the source punctuation should be replaced: %q", unchecked)
	}

	checked := renderMarkdown("- [x] 1.1 do the thing", 80)
	if !strings.Contains(checked, checkboxChecked) {
		t.Errorf("checked task lost its box: %q", checked)
	}
}

func TestOnlyTheCountedShapeBecomesABox(t *testing.T) {
	// The scanner counts `- [ ] ` at column zero and nothing else. Drawing a box
	// on a shape it does not count would make the display disagree with the
	// totals printed beside it.
	for _, line := range []string{
		"  - [ ] indented",
		"* [ ] asterisk",
		"- [X] uppercase",
		"text mentioning - [ ] mid-line",
	} {
		if got := renderMarkdown(line, 80); strings.Contains(got, checkboxUnchecked) ||
			strings.Contains(got, checkboxChecked) {
			t.Errorf("%q should not be drawn as a box, got %q", line, got)
		}
	}
}

// --- the line map ---

func TestLineMapCoversEveryRowOfAWrappedLine(t *testing.T) {
	long := "- [ ] 1.1 " + strings.Repeat("word ", 40)
	content := "## heading\n\n" + long + "\n- [ ] 1.2 short"

	rendered, lines := renderMarkdownLines(content, 40)
	rows := strings.Split(rendered, "\n")

	if len(lines) != 4 {
		t.Fatalf("mapped %d source lines, want 4", len(lines))
	}

	// Every row belongs to exactly one source line, and the ranges are
	// contiguous and in order.
	if lines[0].rowStart != 0 {
		t.Errorf("first line starts at row %d, want 0", lines[0].rowStart)
	}
	for i := 1; i < len(lines); i++ {
		if lines[i].rowStart != lines[i-1].rowEnd+1 {
			t.Errorf("line %d starts at %d, but line %d ended at %d",
				i, lines[i].rowStart, i-1, lines[i-1].rowEnd)
		}
	}
	if last := lines[len(lines)-1].rowEnd; last != len(rows)-1 {
		t.Errorf("last line ends at row %d, but there are %d rows", last, len(rows))
	}

	// The long task really did wrap, or this test proves nothing.
	if lines[2].rowEnd == lines[2].rowStart {
		t.Error("the long task did not wrap; the test needs a narrower pane")
	}
	if lines[2].text != long {
		t.Errorf("the map lost the source text: %q", lines[2].text)
	}
}

// --- toggling ---

func TestToggleTaskLine(t *testing.T) {
	got, ok := toggleTaskLine("- [ ] 1.1 a task")
	if !ok || got != "- [x] 1.1 a task" {
		t.Errorf("got %q ok=%v, want it checked", got, ok)
	}
	got, ok = toggleTaskLine("- [x] 1.1 a task")
	if !ok || got != "- [ ] 1.1 a task" {
		t.Errorf("got %q ok=%v, want it unchecked", got, ok)
	}
	if _, ok := toggleTaskLine("## not a task"); ok {
		t.Error("a non-task line should not toggle")
	}
}

func TestToggleInFileFlipsOneLine(t *testing.T) {
	path := tasksFile(t, "## 1. Group\n\n- [ ] 1.1 first\n- [ ] 1.2 second\n")

	if err := toggleTaskInFile(path, "- [ ] 1.2 second"); err != nil {
		t.Fatalf("toggle: %v", err)
	}
	want := "## 1. Group\n\n- [ ] 1.1 first\n- [x] 1.2 second\n"
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

// TestToggleKeepsAnEditMadeAfterTheScan is the reason the save re-reads.
func TestToggleKeepsAnEditMadeAfterTheScan(t *testing.T) {
	path := tasksFile(t, "- [ ] 1.1 first\n")

	// Someone saves from an editor after specgetty read the file: a task is
	// inserted above the one on screen, which also shifts every position.
	if err := os.WriteFile(path,
		[]byte("- [ ] 0.9 inserted above\n- [ ] 1.1 first\n- [ ] 1.2 also new\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := toggleTaskInFile(path, "- [ ] 1.1 first"); err != nil {
		t.Fatalf("toggle: %v", err)
	}

	got := readFile(t, path)
	want := "- [ ] 0.9 inserted above\n- [x] 1.1 first\n- [ ] 1.2 also new\n"
	if got != want {
		t.Errorf("the editor's work was not preserved:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestToggleRefusesWhenTheLineIsGone(t *testing.T) {
	path := tasksFile(t, "- [ ] 1.1 first\n")
	before := readFile(t, path)

	err := toggleTaskInFile(path, "- [ ] 9.9 never existed")
	if !errors.Is(err, errTaskLineGone) {
		t.Errorf("err = %v, want errTaskLineGone", err)
	}
	if after := readFile(t, path); after != before {
		t.Error("a refused toggle must leave the file byte-identical")
	}
}

func TestToggleRefusesWhenTheLineIsAmbiguous(t *testing.T) {
	path := tasksFile(t, "- [ ] same\n- [ ] other\n- [ ] same\n")
	before := readFile(t, path)

	err := toggleTaskInFile(path, "- [ ] same")
	if !errors.Is(err, errTaskLineAmbiguous) {
		t.Errorf("err = %v, want errTaskLineAmbiguous", err)
	}
	if after := readFile(t, path); after != before {
		t.Error("a refused toggle must leave the file byte-identical")
	}
}

func TestToggleDoesNotChangeFileLength(t *testing.T) {
	// A trailing newline is a final empty element after splitting; losing it
	// would trim a byte on every toggle.
	for _, content := range []string{
		"- [ ] 1.1 a\n",
		"- [ ] 1.1 a",       // no trailing newline
		"- [ ] 1.1 a\n\n\n", // several
	} {
		path := tasksFile(t, content)
		if err := toggleTaskInFile(path, "- [ ] 1.1 a"); err != nil {
			t.Fatalf("toggle %q: %v", content, err)
		}
		got := readFile(t, path)
		if len(got) != len(content) {
			t.Errorf("%q became %q: length %d -> %d", content, got, len(content), len(got))
		}
	}
}

// --- the atomic write ---

func TestAtomicWritePreservesMode(t *testing.T) {
	for _, mode := range []os.FileMode{0o644, 0o600, 0o664} {
		path := filepath.Join(t.TempDir(), "tasks.md")
		if err := os.WriteFile(path, []byte("- [ ] a\n"), mode); err != nil {
			t.Fatal(err)
		}
		// The umask may have cleared bits at creation, so the mode to preserve
		// is the one on disk, not the one requested.
		before, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		want := before.Mode().Perm()

		if err := toggleTaskInFile(path, "- [ ] a"); err != nil {
			t.Fatal(err)
		}
		after, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if after.Mode().Perm() != want {
			t.Errorf("mode %o became %o", want, after.Mode().Perm())
		}
	}
}

func TestAtomicWriteLeavesNoTemporaryFileBehind(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.md")
	if err := os.WriteFile(path, []byte("- [ ] a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := toggleTaskInFile(path, "- [ ] a"); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		names := []string{}
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("directory holds %v, want only tasks.md", names)
	}
}

func TestToggleReportsAFailedWrite(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root, which can write a read-only file")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.md")
	if err := os.WriteFile(path, []byte("- [ ] a\n"), 0o444); err != nil {
		t.Fatal(err)
	}
	// A read-only directory is what actually stops the rename.
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })

	if err := toggleTaskInFile(path, "- [ ] a"); err == nil {
		t.Error("expected a failure writing into a read-only directory")
	}
	if got := readFile(t, path); got != "- [ ] a\n" {
		t.Errorf("the file changed despite the failure: %q", got)
	}
}

// --- the cursor and the key handling, through Update ---

// taskModel opens a change on its tasks artifact, backed by a real file so a
// toggle can be checked against disk.
func taskModel(t *testing.T, tasks string) (model, string) {
	t.Helper()
	project := t.TempDir()
	dir := filepath.Join(project, "openspec", "changes", "my-change")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "tasks.md")
	if err := os.WriteFile(path, []byte(tasks), 0o644); err != nil {
		t.Fatal(err)
	}

	done, total := scanner.ParseTaskStats(tasks)
	m := makeListModel()
	m.width, m.height = 90, 24
	m.repoPaths = []string{project}
	m.displayNames = []string{"proj"}
	m.projects = scanner.ProjectMap{project: scanner.ProjectStatus{Info: scanner.ProjectInfo{
		Changes: []scanner.ChangeInfo{{
			Name: "my-change", DirName: "my-change",
			ArtifactFiles:    []string{"tasks.md"},
			ArtifactContents: map[string]string{"tasks.md": tasks},
			TasksTotal:       done, TasksDone: total,
		}},
	}}}
	m.level = levelChange
	m.changeArtifactTab = 0
	m.recalcLayout()
	m.syncDocument()
	return m, path
}

const threeTasks = "## 1. Group\n\n- [ ] 1.1 first\n- [ ] 1.2 second\n- [ ] 1.3 third\n"

func TestTasksPaneHasACursorAndOthersDoNot(t *testing.T) {
	m, _ := taskModel(t, threeTasks)
	if !m.docHasCursor() {
		t.Error("the tasks artifact should have a cursor")
	}

	// The config tab is a document too, and must not gain one.
	other := makeListModel()
	other.width, other.height = 90, 24
	other.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: scanner.ProjectInfo{
		ConfigFile: "project.md", ConfigContent: "# hello\n",
	}}}
	other.level = levelProject
	other.detailTab = tabProperties
	other.recalcLayout()
	other.syncDocument()
	if other.docHasCursor() {
		t.Error("the config tab must not have a cursor")
	}
}

func TestCursorMovesByTaskAndClamps(t *testing.T) {
	m, _ := taskModel(t, threeTasks)
	if m.docTasks.count() != 3 {
		t.Fatalf("%d tasks found, want 3", m.docTasks.count())
	}
	if m.docCursor != 0 {
		t.Fatalf("cursor starts at %d, want 0", m.docCursor)
	}

	m = press(m, tea.KeyPressMsg{Code: 'k', Text: "k"})
	if m.docCursor != 0 {
		t.Errorf("k at the top moved to %d", m.docCursor)
	}

	// One press per task, not one per source line. The heading and the blank
	// line above the first task are passed over rather than stopped on.
	m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	sel, ok := m.selectedTask()
	if !ok || !strings.Contains(sel.text, "1.2") {
		t.Errorf("one j selected %q, want the second task", sel.text)
	}

	for i := 0; i < 20; i++ {
		m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	}
	if m.docCursor != 2 {
		t.Errorf("cursor ran to %d, want it clamped at 2", m.docCursor)
	}
}

// TestTheCursorNeverStopsOnChrome pins the reason the unit changed: every
// position the cursor can hold is a task, so `j` `space` `j` `space` works
// through a list without a press landing on a heading or a blank line.
func TestTheCursorNeverStopsOnChrome(t *testing.T) {
	m, _ := taskModel(t, "## 1. Group\n\n- [ ] 1.1 first\n\n## 2. Other\n\n- [ ] 2.1 second\n")
	if m.docTasks.count() != 2 {
		t.Fatalf("%d tasks found, want 2", m.docTasks.count())
	}
	for i := 0; i < 6; i++ {
		sel, ok := m.selectedTask()
		if !ok {
			t.Fatal("no task selected")
		}
		if !isTaskLine(sel.text) {
			t.Fatalf("cursor stopped on %q, which is not a task", sel.text)
		}
		m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	}
}

// TestATasksFileWithNoTasksHasNoCursor pins the fallback: nothing to select
// means no cursor, and the pane scrolls by rows like any other document. It
// falls out of the count rather than being a case of its own.
func TestATasksFileWithNoTasksHasNoCursor(t *testing.T) {
	m, _ := taskModel(t, "## 1. Group\n\nProse, and not one checkbox.\n")
	if m.docTasks.count() != 0 {
		t.Fatalf("%d tasks found in a file with none", m.docTasks.count())
	}
	if m.docHasCursor() {
		t.Error("a tasks file with no tasks must not have a cursor")
	}
	if !m.docActive() {
		t.Error("it should still be scrollable as a document")
	}
}

func TestHighlightCoversEveryRowOfAWrappedTask(t *testing.T) {
	long := "- [ ] 1.1 " + strings.Repeat("word ", 40)
	m, _ := taskModel(t, long+"\n- [ ] 1.2 short\n")

	first, last, ok := m.selectedTaskRows()
	if !ok {
		t.Fatal("no task selected")
	}
	if last == first {
		t.Fatal("the task did not wrap; the test proves nothing")
	}

	rows := strings.Split(m.docViewport.View(), "\n")
	for i := first; i <= last && i < len(rows); i++ {
		if !strings.Contains(rows[i], "\x1b[") {
			t.Errorf("row %d of the selected task is not highlighted", i)
		}
	}
	if next := last + 1; next < len(rows) {
		plain := ansi.Strip(rows[next])
		if strings.TrimSpace(plain) != "" && rows[next] != plain && strings.Contains(rows[next], "30;42") {
			t.Error("the row after the selected task is highlighted too")
		}
	}
}

// TestHighlightCoversTheContinuationLines is the fault the bean reported: a
// task's indented continuation lines sat outside the band, so a task read as
// half selected.
func TestHighlightCoversTheContinuationLines(t *testing.T) {
	const tasks = "- [ ] 1.1 first line of the task\n" +
		"      second line, indented\n" +
		"      third line, indented\n" +
		"- [ ] 1.2 another task\n"
	m, _ := taskModel(t, tasks)

	if m.docTasks.count() != 2 {
		t.Fatalf("%d tasks found, want 2", m.docTasks.count())
	}

	first, last, ok := m.selectedTaskRows()
	if !ok {
		t.Fatal("no task selected")
	}
	if last-first < 2 {
		t.Fatalf("the item spans rows %d-%d, want the three lines of the task", first, last)
	}

	rows := strings.Split(m.docViewport.View(), "\n")
	for i := first; i <= last && i < len(rows); i++ {
		if !strings.Contains(rows[i], "30;42") {
			t.Errorf("row %d of the item is not highlighted: %q", i, ansi.Strip(rows[i]))
		}
	}
	// The next task is a different item and keeps out of the band.
	if next := last + 1; next < len(rows) && strings.Contains(rows[next], "30;42") {
		t.Error("the following task is highlighted too")
	}
}

// TestTheContinuationLinesAreNotSelectable pins the other half: the indented
// lines are drawn and highlighted with their task, and never selected on their
// own.
func TestTheContinuationLinesAreNotSelectable(t *testing.T) {
	const tasks = "- [ ] 1.1 first\n      continued\n      continued again\n- [ ] 1.2 second\n"
	m, _ := taskModel(t, tasks)

	m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	sel, ok := m.selectedTask()
	if !ok || !strings.Contains(sel.text, "1.2") {
		t.Errorf("one j from the first task selected %q, want the second task", sel.text)
	}
}

// TestAnIndentedCheckboxIsAContinuation pins the boundary against the rule the
// scanner already follows: it counts checkboxes at column zero, so an indented
// one is part of the task above rather than a task of its own.
func TestAnIndentedCheckboxIsAContinuation(t *testing.T) {
	m, _ := taskModel(t, "- [ ] 1.1 first\n      - [ ] not counted\n- [ ] 1.2 second\n")
	if m.docTasks.count() != 2 {
		t.Errorf("%d tasks found, want 2: an indented checkbox is a continuation",
			m.docTasks.count())
	}
}

func TestSpaceTogglesTheSelectedTaskOnDisk(t *testing.T) {
	m, path := taskModel(t, threeTasks)

	// Move to the second task. One press, the cursor selecting tasks only.
	for !strings.Contains(mustTask(t, m).text, "1.2") {
		if m.docCursor >= m.docTasks.count()-1 {
			break
		}
		m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	}

	m = press(m, tea.KeyPressMsg{Code: ' ', Text: " "})

	got := readFile(t, path)
	if !strings.Contains(got, "- [x] 1.2 second") {
		t.Errorf("1.2 was not checked:\n%s", got)
	}
	if strings.Count(got, "- [x]") != 1 {
		t.Errorf("exactly one task should be checked:\n%s", got)
	}
}

func mustTask(t *testing.T, m model) taskItem {
	t.Helper()
	sel, ok := m.selectedTask()
	if !ok {
		t.Fatal("no selected task")
	}
	return sel
}

func TestSpaceOnANonTaskLineWritesNothing(t *testing.T) {
	m, path := taskModel(t, threeTasks)
	before := readFile(t, path)

	// The cursor cannot rest on a heading any more: it selects tasks and
	// nothing else. What still has to hold is that `space` writes nothing where
	// there is no task under it, which is now a document with none in it.
	m, path = taskModel(t, "## 1. Group\n\nProse, and not one checkbox.\n")
	before = readFile(t, path)
	if _, ok := m.selectedTask(); ok {
		t.Fatal("expected no task to be selectable")
	}
	m = press(m, tea.KeyPressMsg{Code: ' ', Text: " "})

	if after := readFile(t, path); after != before {
		t.Error("space on a heading must not write")
	}
	if m.statusMsg != "" {
		t.Errorf("status = %q, want silence for a key that does not apply", m.statusMsg)
	}
}

func TestSpaceReportsWhenTheLineIsGone(t *testing.T) {
	m, path := taskModel(t, threeTasks)
	mustTask(t, m) // the cursor is on a task already; it selects nothing else

	// The file loses that task behind specgetty's back.
	if err := os.WriteFile(path, []byte("## 1. Group\n\n- [ ] 9.9 different\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := readFile(t, path)

	m = press(m, tea.KeyPressMsg{Code: ' ', Text: " "})

	if !strings.Contains(m.statusMsg, "changed on disk") {
		t.Errorf("status = %q, want it to report the file changed", m.statusMsg)
	}
	if after := readFile(t, path); after != before {
		t.Error("a refused toggle must leave the file alone")
	}
}

func TestSpaceIsInertWhileAConfirmationIsUp(t *testing.T) {
	// A confirmation modal owns the keyboard. This is pinned because an early
	// draft wired space into the modal handlers by accident, which would have
	// edited a file while a question was waiting for an answer.
	m, path := taskModel(t, threeTasks)
	mustTask(t, m) // the cursor is on a task already; it selects nothing else
	before := readFile(t, path)

	for _, setup := range []func(*model){
		func(mm *model) { mm.archiveState = archiveConfirming },
		func(mm *model) { mm.discardState = discardConfirming },
		func(mm *model) { mm.exportState = exportPrompting },
	} {
		blocked := m
		setup(&blocked)
		blocked = press(blocked, tea.KeyPressMsg{Code: ' ', Text: " "})
		if after := readFile(t, path); after != before {
			t.Error("space wrote to the file while a confirmation was up")
		}
	}
}

func TestSpaceIsInertWithoutACursor(t *testing.T) {
	m := makeListModel()
	m.width, m.height = 90, 24
	m.recalcLayout()
	m.syncDocument()

	if m.docHasCursor() {
		t.Fatal("the change list must not have a document cursor")
	}
	m = press(m, tea.KeyPressMsg{Code: ' ', Text: " "})
	if m.statusMsg != "" {
		t.Errorf("status = %q, want nothing", m.statusMsg)
	}
}

func TestCursorScrollsIntoView(t *testing.T) {
	var b strings.Builder
	for i := 1; i <= 60; i++ {
		fmt.Fprintf(&b, "- [ ] %d.1 a task\n", i)
	}
	m, _ := taskModel(t, b.String())

	height := m.docViewport.Height()
	if height < 3 {
		t.Fatalf("pane height %d is too small for this test", height)
	}

	for i := 0; i < 40; i++ {
		m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	}
	first, last, ok := m.selectedTaskRows()
	if !ok {
		t.Fatal("no task selected")
	}
	top := m.docViewport.YOffset()
	if first < top || last >= top+height {
		t.Errorf("selected rows %d-%d are outside the visible window %d-%d",
			first, last, top, top+height-1)
	}
}
