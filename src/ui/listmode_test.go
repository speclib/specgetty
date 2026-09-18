package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/mipmip/specgetty/src/scanner"
)

func TestResolveListModePrecedence(t *testing.T) {
	t.Run("default when neither is set", func(t *testing.T) {
		got, err := ResolveListMode("", "")
		if err != nil || got != modeActive {
			t.Errorf("got %d err=%v, want modeActive", got, err)
		}
	})

	t.Run("config used when no flag", func(t *testing.T) {
		got, err := ResolveListMode("", "archived")
		if err != nil || got != modeArchived {
			t.Errorf("got %d err=%v, want modeArchived", got, err)
		}
	})

	t.Run("flag beats config", func(t *testing.T) {
		got, err := ResolveListMode("active+archived", "archived")
		if err != nil || got != modeBoth {
			t.Errorf("got %d err=%v, want modeBoth", got, err)
		}
	})

	t.Run("whitespace tolerated", func(t *testing.T) {
		got, err := ResolveListMode("  archived  ", "")
		if err != nil || got != modeArchived {
			t.Errorf("got %d err=%v, want modeArchived", got, err)
		}
	})
}

func TestResolveListModeAcceptsEveryNameTheNavBarShows(t *testing.T) {
	// The config vocabulary and the on-screen vocabulary are the same list, so
	// this cannot drift into accepting a word the UI never shows.
	for want, name := range listModeNames {
		got, err := ResolveListMode(name, "")
		if err != nil {
			t.Errorf("%q: %v", name, err)
			continue
		}
		if got != want {
			t.Errorf("%q resolved to %d, want %d", name, got, want)
		}
	}
}

func TestResolveListModeRejectsUnknown(t *testing.T) {
	_, err := ResolveListMode("open", "")
	if err == nil {
		t.Fatal("expected an error: 'open' is the old word and must not be accepted")
	}
	for _, want := range []string{"open", "--change-mode", "active", "archived"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should mention %q", err, want)
		}
	}
}

func TestResolveListModeNamesTheConfigAsItsSource(t *testing.T) {
	_, err := ResolveListMode("", "nonsense")
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "config file") {
		t.Errorf("error %q should point at the config file, not the flag", err)
	}
}

func TestChangeListOpensInTheConfiguredMode(t *testing.T) {
	m := makeListModel()
	m.listMode, m.defaultMode = modeBoth, modeBoth
	if got := len(m.currentRows()); got != 3 {
		t.Errorf("got %d rows, want all 3 with the mode set to both", got)
	}
}

func TestProjectSwitchReturnsToTheConfiguredDefault(t *testing.T) {
	// The default here is deliberately not modeActive. With modeActive the test
	// could not tell "returned to the configured default" apart from "reset to
	// the old hardcoded value", and would pass either way.
	m := makePickerModel()
	m.defaultMode = modeBoth
	m.listMode = modeBoth
	m.pickerOpen = true
	m.pickerCursor = 1

	// The user narrows the view by hand before switching.
	m.listMode = modeArchived

	// Switching project goes through the picker, which is the only way now.
	m = press(m, tea.KeyPressMsg{Code: tea.KeyEnter})

	if m.listMode != modeBoth {
		t.Errorf("listMode = %d after switching project, want the configured default %d",
			m.listMode, modeBoth)
	}
}

func TestFStillCyclesFromAnyDefault(t *testing.T) {
	for start := range listModeNames {
		m := makeListModel()
		m.listMode, m.defaultMode = start, start
		seen := map[int]bool{start: true}
		for i := 0; i < 3; i++ {
			m = press(m, tea.KeyPressMsg{Code: 'f', Text: "f"})
			seen[m.listMode] = true
		}
		if len(seen) != 3 {
			t.Errorf("starting at %d, f reached %d modes, want all 3", start, len(seen))
		}
		if m.listMode != start {
			t.Errorf("starting at %d, three presses of f ended at %d", start, m.listMode)
		}
	}
}

func TestNoModeIsSpelledOpenAnywhereTheUserLooks(t *testing.T) {
	// "open" already means the change you drilled into. It must not also mean
	// "not archived", which is why this asserts on what reaches the screen.
	m := makeListModel()
	m.width = 220
	m.listMode, m.defaultMode = modeActive, modeActive

	if got := m.renderNavBar(); strings.Contains(got, "mode:open") {
		t.Errorf("the nav bar still says mode:open:\n%s", got)
	}
	if got := knownFields["archived"].value(changeRow{}); got == "open" {
		t.Error(`the state column still renders "open"`)
	}
	for _, name := range listModeNames {
		if name == "open" || name == "open+archived" {
			t.Errorf("listModeNames still carries %q", name)
		}
	}

	r := changeRow{ci: scanner.ChangeInfo{Name: "x"}}
	m.level = levelChange
	if got := m.renderChangeDetail(r, 0, m.panelContentWidth(), m.mainPanelHeight()); strings.Contains(got, "(open)") {
		t.Errorf("the open-change header still says (open):\n%s", got)
	}
}
