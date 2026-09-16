package ui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/mipmip/specgetty/src/scanner"
)

// fakeClipboard replaces the real writer for the duration of a test.
//
// Nothing in this suite may touch the developer's actual clipboard: clobbering
// what someone had copied is a bad neighbour, and depending on wl-copy being
// installed would fail on machines that are perfectly fine.
func fakeClipboard(t *testing.T, err error) *string {
	t.Helper()
	var got string
	previous := writeClipboard
	writeClipboard = func(s string) error {
		if err != nil {
			return err
		}
		got = s
		return nil
	}
	t.Cleanup(func() { writeClipboard = previous })
	return &got
}

// copyModel puts a project with one active and one archived change on the
// change list, backed by real directories so a copied path can be resolved.
func copyModel(t *testing.T, mode int) (model, string) {
	t.Helper()
	project := t.TempDir()

	active := filepath.Join(project, "openspec", "changes", "my-feature")
	archived := filepath.Join(project, "openspec", "changes", "archive", "2026-04-02-old-thing")
	for _, d := range []string{active, archived} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	m := makeListModel()
	m.width, m.height = 100, 30
	m.repoPaths = []string{project}
	m.displayNames = []string{"proj"}
	m.projects = scanner.ProjectMap{project: scanner.ProjectStatus{Info: scanner.ProjectInfo{
		Changes: []scanner.ChangeInfo{
			{Name: "my-feature", DirName: "my-feature"},
		},
		ArchivedChanges: []scanner.ChangeInfo{
			{Name: "old-thing", DirName: "2026-04-02-old-thing"},
		},
	}}}
	m.listMode = mode
	m.recalcLayout()
	m.syncDocument()
	return m, project
}

func TestCopyNameCopiesTheDisplayName(t *testing.T) {
	got := fakeClipboard(t, nil)
	m, _ := copyModel(t, modeOpen)

	um := press(m, tea.KeyPressMsg{Code: 'y', Text: "y"})

	if *got != "my-feature" {
		t.Errorf("copied %q, want the display name", *got)
	}
	if !strings.Contains(um.statusMsg, "my-feature") {
		t.Errorf("status = %q, want it to name what was copied", um.statusMsg)
	}
}

func TestCopyPathOfAnActiveChangeResolves(t *testing.T) {
	got := fakeClipboard(t, nil)
	m, project := copyModel(t, modeOpen)

	press(m, tea.KeyPressMsg{Code: 'Y', Text: "Y"})

	want := filepath.Join(project, "openspec", "changes", "my-feature")
	if *got != want {
		t.Errorf("copied %q, want %q", *got, want)
	}
	if _, err := os.Stat(*got); err != nil {
		t.Errorf("the copied path does not resolve: %v", err)
	}
}

// TestCopyPathOfAnArchivedChangeResolves is the one that matters.
//
// The displayed name has the date prefix stripped; the directory keeps it.
// Building the path from Name produces something that does not exist, which is
// the defect doExportChange shipped with. This asserts against the filesystem
// rather than against a string, so it fails if the path is wrong for any reason.
func TestCopyPathOfAnArchivedChangeResolves(t *testing.T) {
	got := fakeClipboard(t, nil)
	m, project := copyModel(t, modeArchived)

	press(m, tea.KeyPressMsg{Code: 'Y', Text: "Y"})

	if _, err := os.Stat(*got); err != nil {
		t.Fatalf("the copied path does not resolve: %q: %v", *got, err)
	}
	if !strings.Contains(*got, "2026-04-02-old-thing") {
		t.Errorf("copied %q, want the directory name including its date prefix", *got)
	}

	want := filepath.Join(project, "openspec", "changes", "archive", "2026-04-02-old-thing")
	if *got != want {
		t.Errorf("copied %q, want %q", *got, want)
	}
}

func TestCopyReportsAFailure(t *testing.T) {
	fakeClipboard(t, errors.New("no clipboard tool found"))
	m, _ := copyModel(t, modeOpen)

	um := press(m, tea.KeyPressMsg{Code: 'y', Text: "y"})

	if !strings.Contains(um.statusMsg, "could not copy") {
		t.Errorf("status = %q, want it to report the failure", um.statusMsg)
	}
	if !strings.Contains(um.statusMsg, "no clipboard tool found") {
		t.Errorf("status = %q, want it to carry the reason", um.statusMsg)
	}
}

func TestCopyDoesNothingWithoutASelection(t *testing.T) {
	got := fakeClipboard(t, nil)
	m, _ := copyModel(t, modeOpen)
	m.searchInput.SetValue("nothing-matches-this")
	m.syncCursor()

	for _, k := range []rune{'y', 'Y'} {
		um := press(m, tea.KeyPressMsg{Code: k, Text: string(k)})
		if *got != "" {
			t.Errorf("%c copied %q with nothing selected", k, *got)
		}
		if um.statusMsg != "" {
			t.Errorf("%c reported %q with nothing selected", k, um.statusMsg)
		}
	}
}

func TestCopyIsInertOnOtherTabs(t *testing.T) {
	got := fakeClipboard(t, nil)
	m, _ := copyModel(t, modeOpen)

	for _, tab := range []int{tabSpecs, tabConfig} {
		m.detailTab = tab
		m.syncDocument()
		for _, k := range []rune{'y', 'Y'} {
			press(m, tea.KeyPressMsg{Code: k, Text: string(k)})
			if *got != "" {
				t.Errorf("%c copied %q on tab %d", k, *got, tab)
			}
		}
	}
}

func TestCopyIsInertWhileAnOverlayHoldsTheKeyboard(t *testing.T) {
	t.Run("picker open", func(t *testing.T) {
		got := fakeClipboard(t, nil)
		m, _ := copyModel(t, modeOpen)
		m.pickerOpen = true
		m.pickerAll = testProjectRows()
		m.pickerLoaded = true

		press(m, tea.KeyPressMsg{Code: 'y', Text: "y"})
		if *got != "" {
			t.Errorf("y copied %q with the picker open", *got)
		}
	})

	t.Run("search prompt focused", func(t *testing.T) {
		got := fakeClipboard(t, nil)
		m, _ := copyModel(t, modeOpen)
		m.searchFocused = true
		m.searchInput.Focus()

		um := press(m, tea.KeyPressMsg{Code: 'y', Text: "y"})
		if *got != "" {
			t.Errorf("y copied %q while typing a filter", *got)
		}
		if um.searchInput.Value() != "y" {
			t.Errorf("query = %q, want the y to be typed into the filter", um.searchInput.Value())
		}
	})

	t.Run("confirmation modal up", func(t *testing.T) {
		got := fakeClipboard(t, nil)
		m, _ := copyModel(t, modeOpen)
		m.archiveState = archiveConfirming

		press(m, tea.KeyPressMsg{Code: 'y', Text: "y"})
		if *got != "" {
			t.Errorf("y copied %q while a confirmation was awaiting an answer", *got)
		}
	})
}

// --- the status line ---

func TestStatusLineAppearsInTheViewAndThenClears(t *testing.T) {
	fakeClipboard(t, nil)
	m, _ := copyModel(t, modeOpen)

	um := press(m, tea.KeyPressMsg{Code: 'y', Text: "y"})
	if !strings.Contains(um.View().Content, "copied name") {
		t.Errorf("the status line should be visible in the view:\n%s", um.View().Content)
	}
	// It replaces the nav bar rather than being added to it, so the frame keeps
	// its height.
	if lines := strings.Split(um.View().Content, "\n"); len(lines) != um.height {
		t.Errorf("the view is %d lines with a status message, want %d", len(lines), um.height)
	}

	after := press(um, tea.KeyPressMsg{Code: tea.KeyDown})
	if after.statusMsg != "" {
		t.Errorf("status = %q, want it cleared by the next keystroke", after.statusMsg)
	}
	if !strings.Contains(after.View().Content, "quit") {
		t.Error("the nav bar should be back once the message is gone")
	}
}

func TestNavBarAdvertisesTheCopyKeys(t *testing.T) {
	m, _ := copyModel(t, modeOpen)
	m.width = 220 // wide enough that no hint is dropped
	if !strings.Contains(m.renderNavBar(), "copy") {
		t.Errorf("the nav bar should advertise the copy keys:\n%s", m.renderNavBar())
	}

	// With nothing selected there is nothing to copy, so no hint.
	m.searchInput.SetValue("nothing-matches-this")
	m.syncCursor()
	if strings.Contains(m.renderNavBar(), "copy") {
		t.Errorf("the nav bar should not offer copy with no selection:\n%s", m.renderNavBar())
	}
}

func TestChangeDirPathUsesDirNameNotName(t *testing.T) {
	// A direct test of the rule, independent of the key handling, because this
	// is the thing that has been got wrong before.
	archived := changeRow{
		ci:       scanner.ChangeInfo{Name: "old-thing", DirName: "2026-04-02-old-thing"},
		archived: true,
	}
	got := changeDirPath("/proj", archived)
	if got != "/proj/openspec/changes/archive/2026-04-02-old-thing" {
		t.Errorf("got %q, want the path built from DirName", got)
	}

	active := changeRow{ci: scanner.ChangeInfo{Name: "feat", DirName: "feat"}}
	if got := changeDirPath("/proj", active); got != "/proj/openspec/changes/feat" {
		t.Errorf("got %q, want the active change path", got)
	}
}
