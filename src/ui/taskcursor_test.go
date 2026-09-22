package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

// --- 1.x the boundary rule ---

func TestContinuesTaskItem(t *testing.T) {
	for _, c := range []struct {
		name string
		line string
		want bool
	}{
		{"an indented continuation", "      installing nix", true},
		{"a tab-indented continuation", "\tinstalling nix", true},
		{"an indented checkbox, which the scanner does not count", "  - [ ] sub", true},
		{"an unindented line", "prose at column zero", false},
		{"a blank line", "", false},
		{"a line of spaces", "    ", false},
		{"a heading", "## 2. The next group", false},
		{"the next unchecked task", "- [ ] 1.2 second", false},
		{"the next checked task", "- [x] 1.2 second", false},
	} {
		if got := continuesTaskItem(c.line); got != c.want {
			t.Errorf("%s: continuesTaskItem(%q) = %v, want %v", c.name, c.line, got, c.want)
		}
	}
}

// TestTheBoundaryAgainstThisRepositorysOwnTasks runs the rule over the shape
// OpenSpec's own generators write, which is the shape it has to survive.
func TestTheBoundaryAgainstThisRepositorysOwnTasks(t *testing.T) {
	const tasks = "## 1. The gate on every push\n" +
		"\n" +
		"- [ ] 1.1 Add `.github/workflows/check.yml` running on push and on pull request,\n" +
		"      installing nix and running `nix flake check` and nothing else\n" +
		"- [ ] 1.2 Verify the workflow fails when the gate fails, by pushing a branch with\n" +
		"      a deliberately failing test and reading the run, rather than assuming a\n" +
		"      non-zero exit is surfaced\n" +
		"\n" +
		"## 2. What CI measured is what is published\n" +
		"\n" +
		"- [x] 2.1 Take the total from the gate output\n"

	m, _ := taskModel(t, tasks)
	if got := m.docTasks.count(); got != 3 {
		t.Fatalf("%d tasks found, want 3", got)
	}
	for i, want := range []string{"1.1", "1.2", "2.1"} {
		task, ok := m.docTasks.at(i)
		if !ok || !strings.Contains(task.text, want) {
			t.Errorf("task %d is %q, want %s", i, task.text, want)
		}
	}
}

func TestTheBoundaryOnAwkwardFiles(t *testing.T) {
	for _, c := range []struct {
		name  string
		tasks string
		want  int
	}{
		{"a task on the last line with no trailing newline", "- [ ] 1.1 only", 1},
		{"a task at the end of the file with no blank after it",
			"## 1. G\n\n- [ ] 1.1 a\n      continued", 1},
		{"two checkboxes with nothing between them", "- [ ] a\n- [ ] b\n", 2},
		{"a file whose only checkbox is its last line", "prose\n\n- [ ] 1.1 last\n", 1},
		{"a heading ends the item", "- [ ] a\n## H\n- [ ] b\n", 2},
		{"a blank line ends the item", "- [ ] a\n\n      orphaned\n- [ ] b\n", 2},
		{"no checkbox at all", "## 1. G\n\nprose only\n", 0},
		{"an empty file", "", 0},
	} {
		m, _ := taskModel(t, c.tasks)
		if got := m.docTasks.count(); got != c.want {
			t.Errorf("%s: %d tasks, want %d", c.name, got, c.want)
		}
	}
}

// TestChromeRowsAreOwnedByNothing pins that headings, blank lines and the
// `Tasks: n/m complete` prefix are marked lineOwner, the way the other three
// lists mark theirs. Without it the prefix would be attributed to task zero and
// the first highlight would reach above the document.
func TestChromeRowsAreOwnedByNothing(t *testing.T) {
	m, _ := taskModel(t, threeTasks)

	first, _, ok := m.selectedTaskRows()
	if !ok {
		t.Fatal("no task selected")
	}
	if first < 1 {
		t.Fatalf("the first task starts at row %d, so there is no prefix to test", first)
	}
	for r := 0; r < first; r++ {
		if m.docTasks.rows[r] != lineOwner {
			t.Errorf("row %d above the first task is owned by item %d, want chrome",
				r, m.docTasks.rows[r])
		}
	}
}

// TestTheMappingCoversEveryDrawnRow guards the off-by-the-prefix the shift in
// renderChangeArtifact invites: a mapping sized to the source rather than to
// the rendered document leaves its last rows outside the slice entirely.
func TestTheMappingCoversEveryDrawnRow(t *testing.T) {
	m, _ := taskModel(t, threeTasks)

	// The document's own rows, not the viewport's: the viewport pads to its
	// height, and padding belongs to no document.
	d, ok := m.currentDoc()
	if !ok {
		t.Fatal("no document on screen")
	}
	rows := strings.Count(d.content, "\n") + 1
	if len(m.docTasks.rows) != rows {
		t.Errorf("the mapping has %d rows, the document has %d",
			len(m.docTasks.rows), rows)
	}
	_, last, _ := m.selectedTaskRows()
	if last >= len(m.docTasks.rows) {
		t.Errorf("the selected task ends at row %d, past the %d-row mapping",
			last, len(m.docTasks.rows))
	}
}

// --- 2.5 the toggle still writes the checkbox line only ---

func TestATogglesLeavesTheContinuationLinesAlone(t *testing.T) {
	const tasks = "- [ ] 1.1 first\n" +
		"      a continuation that must not move\n" +
		"      and another\n" +
		"- [ ] 1.2 second\n"
	m, path := taskModel(t, tasks)

	m = press(m, tea.KeyPressMsg{Code: ' ', Text: " "})

	got := readFile(t, path)
	want := "- [x] 1.1 first\n" +
		"      a continuation that must not move\n" +
		"      and another\n" +
		"- [ ] 1.2 second\n"
	if got != want {
		t.Errorf("the file is\n%q\nwant\n%q", got, want)
	}
}

// --- 3.x the paging keys reach the cursor ---

// manyTasks builds a document with more tasks than any pane can show.
func manyTasks(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, "- [ ] %d.1 a task\n", i)
	}
	return b.String()
}

func pressKey(m model, name string) model {
	switch name {
	case "pgdown":
		return press(m, tea.KeyPressMsg{Code: tea.KeyPgDown})
	case "pgup":
		return press(m, tea.KeyPressMsg{Code: tea.KeyPgUp})
	case "ctrl+d":
		return press(m, tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	case "ctrl+u":
		return press(m, tea.KeyPressMsg{Code: 'u', Mod: tea.ModCtrl})
	case "ctrl+f":
		return press(m, tea.KeyPressMsg{Code: 'f', Mod: tea.ModCtrl})
	case "ctrl+b":
		return press(m, tea.KeyPressMsg{Code: 'b', Mod: tea.ModCtrl})
	case "G":
		return press(m, tea.KeyPressMsg{Code: 'G', Text: "G", Mod: tea.ModShift})
	case "gg":
		m = press(m, tea.KeyPressMsg{Code: 'g', Text: "g"})
		return press(m, tea.KeyPressMsg{Code: 'g', Text: "g"})
	}
	panic("unknown key " + name)
}

// TestThePageKeysMoveTheCursor is the fault the bean reported: the page keys
// moved rows and left the cursor behind, so `space` afterwards acted on a task
// that was no longer on screen.
func TestThePageKeysMoveTheCursor(t *testing.T) {
	for _, key := range []string{"pgdown", "ctrl+f", "ctrl+d"} {
		m, _ := taskModel(t, manyTasks(60))
		before := m.docCursor

		m = pressKey(m, key)

		if m.docCursor == before {
			t.Errorf("%s did not move the cursor", key)
		}
		if m.docCursor <= before {
			t.Errorf("%s moved the cursor backwards, to %d", key, m.docCursor)
		}
	}
	for _, key := range []string{"pgup", "ctrl+b", "ctrl+u"} {
		m, _ := taskModel(t, manyTasks(60))
		m.docCursor = 40
		m.syncDocument()
		before := m.docCursor

		m = pressKey(m, key)

		if m.docCursor >= before {
			t.Errorf("%s moved the cursor to %d, want it back from %d", key, m.docCursor, before)
		}
	}
}

// TestTheCursorIsNeverLeftOffScreen is the assertion the whole of section 3
// exists for.
func TestTheCursorIsNeverLeftOffScreen(t *testing.T) {
	for _, key := range []string{"pgdown", "ctrl+f", "ctrl+d", "pgup", "ctrl+b", "ctrl+u", "G", "gg"} {
		m, _ := taskModel(t, manyTasks(60))
		// Start in the middle so both directions have somewhere to go.
		m.docCursor = 30
		m.syncDocument()
		m.scrollCursorIntoView()

		m = pressKey(m, key)

		first, last, ok := m.selectedTaskRows()
		if !ok {
			t.Fatalf("%s: no task selected", key)
		}
		top := m.docViewport.YOffset()
		height := m.docViewport.Height()
		if first < top || last >= top+height {
			t.Errorf("%s: cursor rows %d-%d are outside the window %d-%d",
				key, first, last, top, top+height-1)
		}
	}
}

func TestGAndGGReachTheFirstAndLastTask(t *testing.T) {
	m, _ := taskModel(t, manyTasks(60))

	m = pressKey(m, "G")
	if want := m.docTasks.count() - 1; m.docCursor != want {
		t.Errorf("G left the cursor at %d, want the last task at %d", m.docCursor, want)
	}

	m = pressKey(m, "gg")
	if m.docCursor != 0 {
		t.Errorf("gg left the cursor at %d, want the first task", m.docCursor)
	}
}

func TestPagingPastAnEndStopsThere(t *testing.T) {
	m, _ := taskModel(t, manyTasks(60))
	for i := 0; i < 30; i++ {
		m = pressKey(m, "pgdown")
	}
	if want := m.docTasks.count() - 1; m.docCursor != want {
		t.Errorf("paging to the end left the cursor at %d, want %d", m.docCursor, want)
	}
	for i := 0; i < 30; i++ {
		m = pressKey(m, "pgup")
	}
	if m.docCursor != 0 {
		t.Errorf("paging back left the cursor at %d, want 0", m.docCursor)
	}
}

// TestATallerPaneMovesFurther pins the rule the document-viewer spec already
// states for the three lists: a page is what the surface can show.
func TestATallerPaneMovesFurther(t *testing.T) {
	short, _ := taskModel(t, manyTasks(200))
	short.height = 20
	short.recalcLayout()
	short.syncDocument()

	tall, _ := taskModel(t, manyTasks(200))
	tall.height = 50
	tall.recalcLayout()
	tall.syncDocument()

	if tall.docViewport.Height() <= short.docViewport.Height() {
		t.Fatalf("the two panes are %d and %d rows; the test proves nothing",
			short.docViewport.Height(), tall.docViewport.Height())
	}
	if tall.docPage() <= short.docPage() {
		t.Errorf("a page is %d tasks in a %d-row pane and %d in a %d-row pane",
			tall.docPage(), tall.docViewport.Height(),
			short.docPage(), short.docViewport.Height())
	}
}

// TestAPageOverTallTasksMovesFewerOfThem is the other half of that rule: a task
// occupying several rows counts once, not once per row.
func TestAPageOverTallTasksMovesFewerOfThem(t *testing.T) {
	var tall strings.Builder
	for i := 1; i <= 60; i++ {
		fmt.Fprintf(&tall, "- [ ] %d.1 a task\n      with a continuation line\n      and another one\n", i)
	}
	short, _ := taskModel(t, manyTasks(60))
	deep, _ := taskModel(t, tall.String())

	if deep.docPage() >= short.docPage() {
		t.Errorf("a page is %d tall tasks and %d one-line tasks; tall tasks should page fewer",
			deep.docPage(), short.docPage())
	}
	if deep.docPage() < 1 {
		t.Error("a page must always move at least one task")
	}
}

// TestHalfAPageIsHalfAPage pins the pairing rather than the arithmetic.
func TestHalfAPageIsHalfAPage(t *testing.T) {
	full, _ := taskModel(t, manyTasks(200))
	half, _ := taskModel(t, manyTasks(200))

	full = pressKey(full, "pgdown")
	half = pressKey(half, "ctrl+d")

	if half.docCursor >= full.docCursor {
		t.Errorf("ctrl+d moved to %d and pgdown to %d; half should be less",
			half.docCursor, full.docCursor)
	}
	if half.docCursor < 1 {
		t.Error("half a page must still move at least one task")
	}
}

// TestADocumentWithoutACursorStillScrollsByRows pins that nothing outside the
// tasks pane changed: the proposal and the design artifacts page their rows as
// they always did.
func TestADocumentWithoutACursorStillScrollsByRows(t *testing.T) {
	m, _ := taskModel(t, manyTasks(60))
	// Make the document a cursorless one by emptying the mapping the way an
	// artifact other than tasks.md does.
	m.docTasks = taskItems{}
	if m.docHasCursor() {
		t.Fatal("the pane should have no cursor for this test")
	}
	if !m.docActive() {
		t.Fatal("the pane should still be an active document")
	}

	before := m.docViewport.YOffset()
	m = pressKey(m, "pgdown")
	if m.docViewport.YOffset() <= before {
		t.Error("a document without a cursor should still page its rows")
	}
}

// --- 4.x the refresh is silent ---

func TestARefreshRaisesNoModal(t *testing.T) {
	m := makeViewModel()
	m.refreshing = true
	if strings.Contains(viewOf(m), "Scanning for OpenSpec sources") {
		t.Error("a refresh must not raise the scan indicator")
	}
}

func TestAFirstScanStillRaisesTheModal(t *testing.T) {
	m := makeViewModel()
	m.scanning = true
	if !strings.Contains(viewOf(m), "Scanning for OpenSpec sources") {
		t.Error("a first scan still says what it is doing")
	}
}

func TestKeysAreHandledDuringARefresh(t *testing.T) {
	m, _ := taskModel(t, manyTasks(60))
	m.refreshing = true
	before := m.docCursor

	m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})

	if m.docCursor == before {
		t.Error("a key pressed during a refresh was discarded")
	}
}

func TestKeysAreStillRefusedDuringAFirstScan(t *testing.T) {
	m, _ := taskModel(t, manyTasks(60))
	m.scanning = true
	before := m.docCursor

	m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})

	if m.docCursor != before {
		t.Error("a first scan has nothing to act on and should refuse keys")
	}
}

// TestABurstOfTogglesLosesNoKeystroke is the fault behind "flicker when change
// checkboxes": each toggle wrote a file, the watcher brought the flag back up,
// and the presses made in that window were dropped.
func TestABurstOfTogglesLosesNoKeystroke(t *testing.T) {
	m, path := taskModel(t, "- [ ] 1.1 a\n- [ ] 1.2 b\n- [ ] 1.3 c\n")

	for i := 0; i < 3; i++ {
		m = press(m, tea.KeyPressMsg{Code: ' ', Text: " "})
		// The watcher notices each write while the previous read is in flight.
		m, _ = issuedScan(m, fsChangeMsg{})
		m = press(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	}

	got := readFile(t, path)
	if n := strings.Count(got, "- [x]"); n != 3 {
		t.Errorf("%d tasks were ticked, want 3:\n%s", n, got)
	}
}

// TestAToggleTakesPartInTheCoalescing pins task 4.4: the toggle's own read sets
// the refresh flag, so the watcher's notice of that same write is remembered
// rather than starting a read of its own.
func TestAToggleTakesPartInTheCoalescing(t *testing.T) {
	m, _ := taskModel(t, threeTasks)

	m = press(m, tea.KeyPressMsg{Code: ' ', Text: " "})
	if !m.refreshing {
		t.Fatal("a toggle's own read should go through the refresh flag")
	}
	if m.scanning {
		t.Error("a toggle must not raise the first-scan flag")
	}

	m, started := issuedScan(m, fsChangeMsg{})
	if started && !m.scanPending {
		t.Error("the watcher's notice should coalesce rather than start a read")
	}
	if !m.scanPending {
		t.Error("the watcher's notice should be remembered")
	}
}

// TestAStaleModelCannotCauseAWrongWrite is task 4.6. Letting keys through
// during a refresh means a key can be handled against a model one read behind.
// The toggle re-reads the file and refuses a line it cannot identify, which is
// the safety net rather than the key guard.
func TestAStaleModelCannotCauseAWrongWrite(t *testing.T) {
	m, path := taskModel(t, threeTasks)
	m.refreshing = true

	// The file loses the selected task behind specgetty's back, so the model
	// is describing a line that is no longer there.
	writeFile(t, path, "## 1. Group\n\n- [ ] 9.9 different\n")
	before := readFile(t, path)

	m = press(m, tea.KeyPressMsg{Code: ' ', Text: " "})

	if after := readFile(t, path); after != before {
		t.Errorf("a stale model wrote to the file:\n%s", after)
	}
	if !strings.Contains(m.statusMsg, "changed on disk") {
		t.Errorf("status = %q, want it to report the file changed", m.statusMsg)
	}
}

// --- 5.x the nav bar ---

func TestTheNavBarAdvertisesSpaceOnTheTasksPane(t *testing.T) {
	m, _ := taskModel(t, threeTasks)
	if !m.docHasCursor() {
		t.Fatal("the tasks pane should have a cursor")
	}
	bar := ansi.Strip(m.renderNavBar())
	if !strings.Contains(bar, "space") {
		t.Errorf("the nav bar does not offer the toggle:\n%s", bar)
	}
	if !strings.Contains(bar, "navigate") {
		t.Errorf("the nav bar should call jk navigate on a pane with a cursor:\n%s", bar)
	}
	if strings.Contains(bar, "scroll") {
		t.Errorf("the nav bar still calls jk scroll on a pane with a cursor:\n%s", bar)
	}
}

func TestTheNavBarOffersNoToggleWithoutACursor(t *testing.T) {
	m, _ := taskModel(t, threeTasks)
	m.docTasks = taskItems{}
	if m.docHasCursor() {
		t.Fatal("the pane should have no cursor for this test")
	}
	bar := ansi.Strip(m.renderNavBar())
	if strings.Contains(bar, "space") {
		t.Errorf("the nav bar offers a toggle where there is nothing to toggle:\n%s", bar)
	}
	if !strings.Contains(bar, "scroll") {
		t.Errorf("a document without a cursor scrolls, and should say so:\n%s", bar)
	}
}

// TestTheToggleHintSurvivesANarrowTerminal pins task 5.3: hints are dropped
// from the end, so the action of the pane has to sit ahead of the paging keys.
func TestTheToggleHintSurvivesANarrowTerminal(t *testing.T) {
	// The rule is relative, not absolute. A bar narrow enough drops every hint
	// in the end; what must never happen is a paging hint surviving while the
	// action of the pane does not.
	for width := 40; width <= 140; width += 4 {
		m, _ := taskModel(t, threeTasks)
		m.width = width
		m.recalcLayout()

		bar := ansi.Strip(m.renderNavBar())
		paging := strings.Contains(bar, "^f^b") ||
			strings.Contains(bar, "^d^u") ||
			strings.Contains(bar, "gg/G")
		if paging && !strings.Contains(bar, "space") {
			t.Errorf("at %d columns a paging hint survived and the toggle did not:\n%s",
				width, bar)
		}
		if w := ansi.StringWidth(bar); w > width {
			t.Errorf("the nav bar is %d columns wide in a %d-column terminal", w, width)
		}
	}
}

// TestTheToggleHintIsShownAtAnOrdinaryWidth is the positive half: the relative
// rule above passes on a bar that offers nothing at all, so one width where
// the hints do fit is asserted directly.
func TestTheToggleHintIsShownAtAnOrdinaryWidth(t *testing.T) {
	m, _ := taskModel(t, threeTasks)
	m.width = 120
	m.recalcLayout()

	bar := ansi.Strip(m.renderNavBar())
	if !strings.Contains(bar, "space") {
		t.Errorf("a 120-column bar does not offer the toggle:\n%s", bar)
	}
}
