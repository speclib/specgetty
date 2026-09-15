package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func navBarFor(m model) string {
	return m.renderNavBar()
}

func TestNavBarAtProjectsLevel(t *testing.T) {
	m := makeListModel()
	m.level = levelProjects
	m.activeView = viewProjects

	got := navBarFor(m)
	for _, want := range []string{"quit", "open project", "switch"} {
		if !strings.Contains(got, want) {
			t.Errorf("projects nav bar missing %q, got:\n%s", want, got)
		}
	}
	// Change actions belong to the change list, not here.
	if strings.Contains(got, "a archive") {
		t.Errorf("projects nav bar should not offer archive, got:\n%s", got)
	}
}

func TestNavBarAtChangeListOffersActions(t *testing.T) {
	m := makeListModel()
	m.width = 200 // wide enough that nothing is dropped

	got := navBarFor(m)
	for _, want := range []string{"search", "mode:open", "archive", "discard", "export", "view"} {
		if !strings.Contains(got, want) {
			t.Errorf("change list nav bar missing %q, got:\n%s", want, got)
		}
	}
}

func TestNavBarShowsCurrentListMode(t *testing.T) {
	m := makeListModel()
	m.width = 200
	m.listMode = modeBoth

	got := navBarFor(m)
	if !strings.Contains(got, "mode:"+listModeNames[modeBoth]) {
		t.Errorf("nav bar should show the active mode, got:\n%s", got)
	}
}

func TestNavBarHidesArchiveOnArchivedRow(t *testing.T) {
	m := makeListModel()
	m.width = 200
	m.listMode = modeArchived // gamma only, which is archived

	got := navBarFor(m)
	// Match the action, not the mode label: "mode:archived" contains "archive".
	if strings.Contains(got, "a archive") {
		t.Errorf("archive is a no-op on an archived row and should not be offered, got:\n%s", got)
	}
	// Export still works on archived changes.
	if !strings.Contains(got, "export") {
		t.Errorf("export should still be offered for an archived row, got:\n%s", got)
	}
}

func TestNavBarWhileTypingShowsPromptKeys(t *testing.T) {
	m := makeListModel()
	m.width = 200
	m.searchFocused = true

	got := navBarFor(m)
	if !strings.Contains(got, "clear filter") {
		t.Errorf("prompt nav bar should offer clear filter, got:\n%s", got)
	}
	// The ordinary actions are unavailable while typing, so advertising them
	// would be a lie.
	if strings.Contains(got, "a archive") || strings.Contains(got, "quit") {
		t.Errorf("prompt nav bar should not advertise unavailable actions, got:\n%s", got)
	}
}

func TestNavBarInOpenChange(t *testing.T) {
	m := makeListModel()
	m.width = 200
	m.level = levelChange

	got := navBarFor(m)
	if !strings.Contains(got, "artifact") {
		t.Errorf("open change nav bar should mention artifact sub-tabs, got:\n%s", got)
	}
	if !strings.Contains(got, "back to list") {
		t.Errorf("open change nav bar should offer going back, got:\n%s", got)
	}
	if strings.Contains(got, "tabs") {
		t.Errorf("project tab hints do not apply inside a change, got:\n%s", got)
	}
}

func TestNavBarDropsHintsRatherThanOverlapping(t *testing.T) {
	m := makeListModel()
	m.width = 60

	got := navBarFor(m)
	if w := lipgloss.Width(got); w != 60 {
		t.Errorf("nav bar width = %d, want exactly 60", w)
	}
	// The version must survive, since it is what the keys are trimmed for.
	if !strings.Contains(got, "specgetty") {
		t.Errorf("version text should never be overrun, got:\n%s", got)
	}
	// A narrow bar keeps the essentials and drops the utilities.
	if !strings.Contains(got, "quit") {
		t.Errorf("quit should survive truncation, got:\n%s", got)
	}
	if strings.Contains(got, "jump") {
		t.Errorf("gg/G jump should be dropped first at width 60, got:\n%s", got)
	}
}

func TestNavBarWiderShowsMore(t *testing.T) {
	narrow := makeListModel()
	narrow.width = 60
	wide := makeListModel()
	wide.width = 200

	if lipgloss.Width(navBarFor(wide)) <= lipgloss.Width(navBarFor(narrow)) {
		t.Error("a wider terminal should render a wider nav bar")
	}
	if strings.Count(navBarFor(wide), " ") <= strings.Count(navBarFor(narrow), " ") {
		t.Error("a wider terminal should show more hints")
	}
}
