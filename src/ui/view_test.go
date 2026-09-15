package ui

import (
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/mipmip/specgetty/src/scanner"
)

// These tests assert structure, not pixels.
//
// A golden-file comparison of View() would break on every styling tweak and
// train whoever hits it to regenerate the file without reading it. Asserting
// that a modal appears when its state is set, that the size guard fires, and
// that overlays stack in the documented order survives cosmetic change and
// still catches a real regression.
//
// The cautionary tale is the picker misalignment of 2026-09-15: a test asserting
// every line had the same width passed both with and without the bug, because
// lipgloss pads the wrapped remainder. Assert the thing that actually breaks.

// makeViewModel returns a model big enough to render, with one project and two
// changes, sitting on the change list.
func makeViewModel() model {
	m := makeListModel()
	m.width, m.height = 100, 30
	m.recalcLayout()
	m.syncDocument()
	return m
}

func viewOf(m model) string {
	m.syncDocument()
	return m.View()
}

// --- the two early returns ---

func TestViewBeforeTheFirstWindowSize(t *testing.T) {
	// bubbletea renders once before it knows the terminal size.
	m := makeListModel()
	m.width, m.height = 0, 0
	if got := m.View(); got != "Initializing..." {
		t.Errorf("View() = %q, want the initialising placeholder", got)
	}

	m.width, m.height = 100, 0
	if got := m.View(); got != "Initializing..." {
		t.Errorf("with height 0, View() = %q, want the initialising placeholder", got)
	}
}

func TestViewRefusesATerminalThatIsTooSmall(t *testing.T) {
	for _, size := range []struct{ w, h int }{
		{59, 30}, {100, 19}, {59, 19},
	} {
		m := makeListModel()
		m.width, m.height = size.w, size.h
		got := m.View()
		if !strings.Contains(got, "Terminal too small") {
			t.Errorf("at %dx%d View() did not refuse: %q", size.w, size.h, got)
		}
	}

	// The documented minimum itself must be accepted.
	m := makeViewModel()
	m.width, m.height = 60, 20
	m.recalcLayout()
	if got := viewOf(m); strings.Contains(got, "Terminal too small") {
		t.Error("60x20 is the documented minimum and should render")
	}
}

// --- the ordinary view ---

func TestViewRendersTheProjectAndTheNavBar(t *testing.T) {
	m := makeViewModel()
	got := viewOf(m)

	for _, want := range []string{"alpha", "beta", "quit"} {
		if !strings.Contains(got, want) {
			t.Errorf("view is missing %q", want)
		}
	}
}

func TestViewFillsTheTerminalExactly(t *testing.T) {
	// padToHeight promises the frame is exactly the terminal height, and the
	// panels promise the width. A frame that is off by a line scrolls the
	// terminal on every redraw.
	for _, size := range []struct{ w, h int }{
		{60, 20}, {100, 30}, {200, 50},
	} {
		m := makeViewModel()
		m.width, m.height = size.w, size.h
		m.recalcLayout()

		lines := strings.Split(viewOf(m), "\n")
		if len(lines) != size.h {
			t.Errorf("at %dx%d the view is %d lines, want %d", size.w, size.h, len(lines), size.h)
		}
		for i, l := range lines {
			if w := lipgloss.Width(l); w > size.w {
				t.Errorf("at %dx%d line %d is %d columns, wider than the terminal", size.w, size.h, i, w)
				break
			}
		}
	}
}

func TestViewIncludesTheLogPanelOnlyWhenVisible(t *testing.T) {
	m := makeViewModel()
	m.logContent = "a log line that is quite distinctive"
	m.logViewport.SetContent(m.logContent)

	if strings.Contains(viewOf(m), "Log") {
		t.Error("the log panel should not be drawn while hidden")
	}

	m.logVisible = true
	m.recalcLayout()
	m.logViewport.SetContent(m.logContent)
	if !strings.Contains(viewOf(m), "Log") {
		t.Error("the log panel should be drawn once visible")
	}
}

// --- overlays ---

func TestViewDrawsThePicker(t *testing.T) {
	m := makeViewModel()
	if strings.Contains(viewOf(m), "Projects") {
		t.Error("the picker should not be drawn while closed")
	}

	m.pickerOpen = true
	m.pickerAll = testProjectRows()
	m.pickerLoaded = true
	m.pickerSync()
	if !strings.Contains(viewOf(m), "Projects") {
		t.Error("the picker should be drawn when open")
	}
}

func TestViewDrawsTheStartupPrompt(t *testing.T) {
	m := makeViewModel()
	m.askOpenPicker = true
	got := viewOf(m)
	if !strings.Contains(got, "No OpenSpec project here") {
		t.Errorf("the startup prompt is missing:\n%s", got)
	}
}

func TestViewDrawsTheScanningModal(t *testing.T) {
	m := makeViewModel()
	m.scanning = true
	if !strings.Contains(viewOf(m), "Scanning for OpenSpec sources") {
		t.Error("the scanning modal is missing")
	}
}

func TestViewDrawsAnError(t *testing.T) {
	m := makeViewModel()
	m.err = errors.New("the disk fell over")
	got := viewOf(m)
	if !strings.Contains(got, "the disk fell over") {
		t.Errorf("the error modal should carry the message:\n%s", got)
	}
}

// --- the three action state machines ---

// viewWithTasks puts a change with incomplete tasks under the cursor, so the
// confirmation modals take their warning branch.
func viewWithTasks(t *testing.T, done, total int) model {
	t.Helper()
	m := makeViewModel()
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: scanner.ProjectInfo{
		Changes: []scanner.ChangeInfo{
			{Name: "alpha", DirName: "alpha", TasksDone: done, TasksTotal: total},
		},
	}}}
	m.syncDocument()
	return m
}

func TestViewArchiveModals(t *testing.T) {
	t.Run("confirming, all tasks done", func(t *testing.T) {
		m := viewWithTasks(t, 5, 5)
		m.archiveState = archiveConfirming
		m.archiveChangeName = "alpha"
		got := viewOf(m)
		if !strings.Contains(got, `Archive "alpha"?`) {
			t.Errorf("missing the plain confirmation:\n%s", got)
		}
		if strings.Contains(got, "incomplete") {
			t.Error("a complete change should not warn about incomplete tasks")
		}
	})

	t.Run("confirming, tasks outstanding", func(t *testing.T) {
		m := viewWithTasks(t, 2, 5)
		m.archiveState = archiveConfirming
		m.archiveChangeName = "alpha"
		got := viewOf(m)
		if !strings.Contains(got, "3 incomplete task") {
			t.Errorf("the warning should count the outstanding tasks:\n%s", got)
		}
	})

	t.Run("running", func(t *testing.T) {
		m := makeViewModel()
		m.archiveState = archiveRunning
		if !strings.Contains(viewOf(m), "Archiving") {
			t.Error("the running modal is missing")
		}
	})

	t.Run("result, succeeded", func(t *testing.T) {
		m := makeViewModel()
		m.archiveState = archiveResult
		m.archiveResultOk = true
		m.archiveResultMsg = "archived alpha"
		got := viewOf(m)
		if !strings.Contains(got, "archived alpha") || !strings.Contains(got, "✓") {
			t.Errorf("a successful result should be marked as such:\n%s", got)
		}
	})

	t.Run("result, failed", func(t *testing.T) {
		m := makeViewModel()
		m.archiveState = archiveResult
		m.archiveResultOk = false
		m.archiveResultMsg = "openspec exploded"
		got := viewOf(m)
		if !strings.Contains(got, "openspec exploded") || !strings.Contains(got, "✗") {
			t.Errorf("a failed result should be marked as such:\n%s", got)
		}
	})
}

func TestViewDiscardModals(t *testing.T) {
	t.Run("confirming, tasks outstanding", func(t *testing.T) {
		m := viewWithTasks(t, 1, 4)
		m.discardState = discardConfirming
		m.discardChangeName = "alpha"
		got := viewOf(m)
		if !strings.Contains(got, "3 incomplete task") || !strings.Contains(got, "Discard anyway") {
			t.Errorf("missing the discard warning:\n%s", got)
		}
	})

	t.Run("confirming, all tasks done", func(t *testing.T) {
		m := viewWithTasks(t, 4, 4)
		m.discardState = discardConfirming
		m.discardChangeName = "alpha"
		if !strings.Contains(viewOf(m), `Discard "alpha"?`) {
			t.Error("missing the plain discard confirmation")
		}
	})

	t.Run("running", func(t *testing.T) {
		m := makeViewModel()
		m.discardState = discardRunning
		if !strings.Contains(viewOf(m), "Discarding") {
			t.Error("the running modal is missing")
		}
	})

	t.Run("result", func(t *testing.T) {
		m := makeViewModel()
		m.discardState = discardResult
		m.discardResultOk = true
		m.discardResultMsg = "discarded alpha"
		if !strings.Contains(viewOf(m), "discarded alpha") {
			t.Error("the result modal should carry the message")
		}

		m.discardResultOk = false
		if !strings.Contains(viewOf(m), "✗") {
			t.Error("a failed discard should be marked as such")
		}
	})
}

func TestViewExportModals(t *testing.T) {
	t.Setenv("HOME", "/home/someone")

	t.Run("confirming shows the destination", func(t *testing.T) {
		m := makeViewModel()
		m.exportState = exportConfirming
		m.exportChangeName = "alpha"
		got := viewOf(m)
		if !strings.Contains(got, `Export "alpha"?`) {
			t.Errorf("missing the export confirmation:\n%s", got)
		}
		// The spec requires the modal to show where the zip will land.
		if !strings.Contains(got, "/home/someone") {
			t.Errorf("the modal should show the destination path:\n%s", got)
		}
	})

	t.Run("running", func(t *testing.T) {
		m := makeViewModel()
		m.exportState = exportRunning
		if !strings.Contains(viewOf(m), "Exporting") {
			t.Error("the running modal is missing")
		}
	})

	t.Run("result", func(t *testing.T) {
		m := makeViewModel()
		m.exportState = exportResult
		m.exportResultOk = true
		m.exportResultMsg = "Exported to /home/someone/alpha.zip"
		if !strings.Contains(viewOf(m), "alpha.zip") {
			t.Error("the result modal should name the zip")
		}

		m.exportResultOk = false
		m.exportResultMsg = "could not write"
		got := viewOf(m)
		if !strings.Contains(got, "could not write") || !strings.Contains(got, "✗") {
			t.Errorf("a failed export should be marked as such:\n%s", got)
		}
	})
}

// --- stacking order ---

func TestConfirmationModalDrawsOverThePicker(t *testing.T) {
	// The comment in View says confirmation modals are drawn after the picker
	// so they sit on top. A question awaiting an answer must not be buried.
	// A change with all its tasks done, so the modal takes its plain branch and
	// the assertion is about precedence rather than about which variant shows.
	m := viewWithTasks(t, 4, 4)
	m.pickerOpen = true
	m.pickerAll = testProjectRows()
	m.pickerLoaded = true
	m.pickerSync()
	m.archiveState = archiveConfirming
	m.archiveChangeName = "alpha"

	got := viewOf(m)
	if !strings.Contains(got, `Archive "alpha"?`) {
		t.Errorf("the confirmation should be visible over the picker:\n%s", got)
	}

	// Note what is NOT asserted here: that the picker is still visible behind
	// the modal. It is not. placeOverlay takes a `background` argument and
	// never uses it, so each overlay replaces the whole screen rather than
	// compositing onto it. Tracked separately; this test pins only the
	// precedence, which holds either way.
}

func TestViewStaysWithinTheTerminalWithEveryOverlayUp(t *testing.T) {
	// Overlays are composed by writing over the frame, so a modal wider or
	// taller than the terminal would push the frame out of shape.
	m := makeViewModel()
	m.pickerOpen = true
	m.pickerAll = testProjectRows()
	m.pickerLoaded = true
	m.pickerSync()
	m.scanning = true
	m.err = errors.New("something went wrong")
	m.archiveState = archiveConfirming
	m.archiveChangeName = "alpha"

	lines := strings.Split(viewOf(m), "\n")
	if len(lines) != m.height {
		t.Errorf("the view is %d lines with overlays up, want %d", len(lines), m.height)
	}
	for i, l := range lines {
		if w := lipgloss.Width(l); w > m.width {
			t.Errorf("line %d is %d columns, wider than the terminal", i, w)
			break
		}
	}
}
