package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

// exportModel is a project with one active change, backed by real directories
// so an export can be run end to end.
func exportModel(t *testing.T, cfg *scanner.Config) (model, string) {
	t.Helper()
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, "openspec", "changes", "my-change"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, "openspec", "changes", "my-change", "tasks.md"),
		[]byte("- [x] done\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	m := newModel(cfg, true, "0.0.0")
	m.width, m.height = 100, 24
	m.repoPaths = []string{project}
	m.displayNames = []string{"p"}
	m.detailTab = tabChanges
	m.focus = focusDetail
	m.fields = append([]string(nil), defaultFields...)
	m.projects = scanner.ProjectMap{project: scanner.ProjectStatus{Info: scanner.ProjectInfo{
		Root: project, Origin: project,
		ActiveChanges: []string{"my-change"},
		Changes:       []scanner.ChangeInfo{{Name: "my-change", DirName: "my-change"}},
	}}}
	m.recalcLayout()
	return m, project
}

func openPrompt(t *testing.T, cfg *scanner.Config) (model, string) {
	t.Helper()
	m, project := exportModel(t, cfg)
	updated, _ := m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	um := updated.(model)
	if um.exportState != exportPrompting {
		t.Fatalf("exportState = %d, want the prompt", um.exportState)
	}
	return um, project
}

// --- 2.1 the prompt opens on the configured directory ---

func TestThePromptOpensOnTheConfiguredDirectory(t *testing.T) {
	dir := t.TempDir()
	m, _ := openPrompt(t, &scanner.Config{ExportDir: dir})
	if got := m.exportDirInput.Value(); got != dir {
		t.Errorf("got %q, want the configured directory %q", got, dir)
	}
}

func TestThePromptOpensOnHomeWhenNothingIsConfigured(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	m, _ := openPrompt(t, &scanner.Config{})
	if got := m.exportDirInput.Value(); got != home {
		t.Errorf("got %q, want the home directory %q", got, home)
	}
}

// --- 2.6 every printable key belongs to the field ---

func TestThePromptSwallowsTheActionKeys(t *testing.T) {
	// A path called `~/archive` must not archive a change.
	m, _ := openPrompt(t, &scanner.Config{ExportDir: t.TempDir()})
	m.exportDirInput.SetValue("")

	for _, k := range []rune{'a', 'd', 'e', 'q', 'y', 'n', 'l', 'p', 's', '/'} {
		updated, _ := m.Update(tea.KeyPressMsg{Code: k, Text: string(k)})
		um := updated.(model)
		if um.exportState != exportPrompting {
			t.Fatalf("%q left the prompt, state %d", k, um.exportState)
		}
		if um.archiveState != archiveIdle || um.discardState != discardIdle {
			t.Errorf("%q reached an action", k)
		}
		if um.pickerOpen || um.searchFocused {
			t.Errorf("%q raised something else", k)
		}
		if !strings.Contains(um.exportDirInput.Value(), string(k)) {
			t.Errorf("%q did not reach the field, value is %q", k, um.exportDirInput.Value())
		}
		m = um
		m.exportDirInput.SetValue("")
	}
}

func TestEscCancelsTheExport(t *testing.T) {
	m, project := openPrompt(t, &scanner.Config{ExportDir: t.TempDir()})
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	um := updated.(model)

	if um.exportState != exportIdle {
		t.Errorf("state = %d, want idle", um.exportState)
	}
	entries, _ := os.ReadDir(project)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".zip") {
			t.Error("nothing may be written on cancel")
		}
	}
}

// --- 3.1 to 3.3 refusals ---

func TestADirectoryThatDoesNotExistIsRefused(t *testing.T) {
	m, _ := openPrompt(t, &scanner.Config{ExportDir: t.TempDir()})
	missing := filepath.Join(t.TempDir(), "not-here")
	m.exportDirInput.SetValue(missing)

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	um := updated.(model)

	if um.exportState != exportPrompting {
		t.Fatalf("state = %d, want the prompt still up", um.exportState)
	}
	if um.exportDirInput.Value() != missing {
		t.Errorf("what was typed must stay to be corrected, got %q", um.exportDirInput.Value())
	}
	if !strings.Contains(um.exportProblem, missing) {
		t.Errorf("the message must name it, got %q", um.exportProblem)
	}
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Error("nothing may be created")
	}
	if !strings.Contains(ansi.Strip(viewOf(um)), "does not exist") {
		t.Error("and it must be on screen")
	}
}

func TestAPathThatIsNotADirectoryIsRefused(t *testing.T) {
	m, _ := openPrompt(t, &scanner.Config{ExportDir: t.TempDir()})
	file := filepath.Join(t.TempDir(), "a-file")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	m.exportDirInput.SetValue(file)

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	um := updated.(model)
	if um.exportState != exportPrompting || !strings.Contains(um.exportProblem, "not a directory") {
		t.Errorf("state %d, problem %q", um.exportState, um.exportProblem)
	}
}

func TestTypingClearsTheProblem(t *testing.T) {
	m, _ := openPrompt(t, &scanner.Config{ExportDir: t.TempDir()})
	m.exportDirInput.SetValue(filepath.Join(t.TempDir(), "not-here"))
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(model)
	if m.exportProblem == "" {
		t.Fatal("expected a problem")
	}
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	if updated.(model).exportProblem != "" {
		t.Error("correcting the path clears the message")
	}
}

// --- 3.4 to 3.6 replacing ---

func TestAnExistingFileIsConfirmedBeforeItIsReplaced(t *testing.T) {
	dest := t.TempDir()
	m, _ := openPrompt(t, &scanner.Config{ExportDir: dest})
	existing := filepath.Join(dest, exportFileName("my-change"))
	if err := os.WriteFile(existing, []byte("not a zip"), 0o600); err != nil {
		t.Fatal(err)
	}

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	um := updated.(model)
	if um.exportState != exportReplacing {
		t.Fatalf("state = %d, want the replace question", um.exportState)
	}
	if !strings.Contains(ansi.Strip(viewOf(um)), exportFileName("my-change")) {
		t.Error("the question must name the file")
	}

	// Declining returns to the prompt with the directory still in it.
	declined, _ := um.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	dm := declined.(model)
	if dm.exportState != exportPrompting {
		t.Errorf("state = %d, want back at the prompt", dm.exportState)
	}
	if dm.exportDirInput.Value() != dest {
		t.Errorf("the directory must still be there, got %q", dm.exportDirInput.Value())
	}
	if b, _ := os.ReadFile(existing); string(b) != "not a zip" {
		t.Error("declining must write nothing")
	}
}

func TestAgreeingReplacesTheFile(t *testing.T) {
	dest := t.TempDir()
	m, _ := openPrompt(t, &scanner.Config{ExportDir: dest})
	existing := filepath.Join(dest, exportFileName("my-change"))
	if err := os.WriteFile(existing, []byte("not a zip"), 0o600); err != nil {
		t.Fatal(err)
	}

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	updated, cmd := updated.(model).Update(tea.KeyPressMsg{Code: 'y', Text: "y"})
	if updated.(model).exportState != exportRunning {
		t.Fatalf("state = %d, want running", updated.(model).exportState)
	}

	var result *exportMsg
	for _, msg := range runCmd(cmd) {
		if e, ok := msg.(exportMsg); ok {
			result = &e
		}
	}
	if result == nil || !result.ok {
		t.Fatalf("export failed: %+v", result)
	}
	if b, _ := os.ReadFile(existing); string(b) == "not a zip" {
		t.Error("the file must have been replaced")
	}
	if !strings.Contains(result.output, dest) {
		t.Errorf("the result reports where it wrote, got %q", result.output)
	}
}

// --- 4.1 and 4.2 writing, and the next export ---

func TestExportWritesIntoTheChosenDirectory(t *testing.T) {
	configured := t.TempDir()
	chosen := t.TempDir()
	m, _ := openPrompt(t, &scanner.Config{ExportDir: configured})
	m.exportDirInput.SetValue(chosen)

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	for _, msg := range runCmd(cmd) {
		if e, ok := msg.(exportMsg); ok && !e.ok {
			t.Fatalf("export failed: %s", e.output)
		}
	}
	if _, err := os.Stat(filepath.Join(chosen, exportFileName("my-change"))); err != nil {
		t.Errorf("the zip is not in the chosen directory: %v", err)
	}
	if entries, _ := os.ReadDir(configured); len(entries) != 0 {
		t.Error("and nothing is in the configured one")
	}

	// The configuration is untouched, so the next export opens on it again.
	um := updated.(model)
	if um.config.ExportDir != configured {
		t.Errorf("the configuration changed to %q", um.config.ExportDir)
	}
	um.exportState = exportIdle
	again, _ := um.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	if got := again.(model).exportDirInput.Value(); got != configured {
		t.Errorf("the next prompt opens on %q, want the configured %q", got, configured)
	}
}

func TestExportExpandsAShorthandDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.Mkdir(filepath.Join(home, "Downloads"), 0o755); err != nil {
		t.Fatal(err)
	}
	m, _ := openPrompt(t, &scanner.Config{})
	m.exportDirInput.SetValue("~/Downloads")

	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	for _, msg := range runCmd(cmd) {
		if e, ok := msg.(exportMsg); ok && !e.ok {
			t.Fatalf("export failed: %s", e.output)
		}
	}
	if _, err := os.Stat(filepath.Join(home, "Downloads", exportFileName("my-change"))); err != nil {
		t.Errorf("the tilde was not expanded: %v", err)
	}
}

// --- 2.3 completion ---

func TestTabCompletesTheDirectory(t *testing.T) {
	base := t.TempDir()
	if err := os.Mkdir(filepath.Join(base, "exports"), 0o755); err != nil {
		t.Fatal(err)
	}
	m, _ := openPrompt(t, &scanner.Config{ExportDir: base})
	m.exportDirInput.SetValue(filepath.Join(base, "exp"))

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	um := updated.(model)
	if len(um.exportDirInput.AvailableSuggestions()) == 0 {
		t.Fatal("tab must offer what is on disk")
	}
	if got := um.exportDirInput.AvailableSuggestions()[0]; got != filepath.Join(base, "exports") {
		t.Errorf("got %q, want the directory beside it", got)
	}
}

// --- 2.8 one modal at a time ---

func TestThePromptRaisesNothingElse(t *testing.T) {
	m, _ := openPrompt(t, &scanner.Config{ExportDir: t.TempDir()})
	for _, k := range []rune{'p', 'a', 'd', 'e'} {
		updated, _ := m.Update(tea.KeyPressMsg{Code: k, Text: string(k)})
		um := updated.(model)
		if um.pickerOpen || um.archiveState != archiveIdle || um.discardState != discardIdle {
			t.Errorf("%q raised another modal", k)
		}
	}
}

// --- 2.7 the frame ---

func TestThePromptFitsTheFrame(t *testing.T) {
	for _, size := range []struct{ w, h int }{{60, 20}, {92, 30}, {120, 50}} {
		m, _ := openPrompt(t, &scanner.Config{ExportDir: "/some/where"})
		m.width, m.height = size.w, size.h
		m.recalcLayout()

		lines := strings.Split(m.renderFrame(), "\n")
		if len(lines) != size.h {
			t.Errorf("%dx%d: got %d rows, want %d", size.w, size.h, len(lines), size.h)
		}
		for _, l := range lines {
			if w := ansi.StringWidth(l); w > size.w {
				t.Errorf("%dx%d: a row is %d columns wide", size.w, size.h, w)
				break
			}
		}
	}
}

var _ = textinput.New

func TestNoModalPushesTheFrameOutOfShape(t *testing.T) {
	// `modal-presentation` has required this since it shipped, and two fixed
	// widths broke it at the 60-column minimum: the export prompt asked for 66
	// columns and the startup question for 62.
	states := []struct {
		name  string
		apply func(*model)
	}{
		{"the startup question", func(m *model) { m.askOpenPicker = true }},
		{"scanning", func(m *model) { m.scanning = true }},
		{"an error", func(m *model) { m.err = os.ErrNotExist }},
		{"archiving", func(m *model) { m.archiveState = archiveConfirming }},
		{"discarding", func(m *model) { m.discardState = discardConfirming }},
		{"the export prompt", func(m *model) {
			m.exportState = exportPrompting
			m.exportChangeName = "a-change-with-a-fairly-long-name"
			m.exportDirInput.SetValue("/some/deep/directory/somewhere")
		}},
		{"replacing", func(m *model) {
			m.exportState = exportReplacing
			m.exportChangeName = "a-change-with-a-fairly-long-name"
			m.exportDirInput.SetValue("/some/deep/directory/somewhere")
		}},
		{"a result", func(m *model) {
			m.exportState = exportResult
			m.exportResultMsg = strings.Repeat("long ", 40)
		}},
	}
	for _, st := range states {
		for _, size := range []struct{ w, h int }{{60, 20}, {92, 30}} {
			m, _ := exportModel(t, &scanner.Config{ExportDir: "/some/where"})
			m.width, m.height = size.w, size.h
			m.recalcLayout()
			st.apply(&m)

			lines := strings.Split(m.renderFrame(), "\n")
			if len(lines) != size.h {
				t.Errorf("%s at %dx%d: got %d rows, want %d", st.name, size.w, size.h, len(lines), size.h)
			}
			for _, l := range lines {
				if w := ansi.StringWidth(l); w > size.w {
					t.Errorf("%s at %dx%d: a row is %d columns wide", st.name, size.w, size.h, w)
					break
				}
			}
		}
	}
}
