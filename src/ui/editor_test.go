package ui

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

func TestTheEditorComesFromTheEnvironment(t *testing.T) {
	for _, c := range []struct {
		name           string
		visual, editor string
		wantName       string
		wantReason     bool
	}{
		{name: "both set, VISUAL wins", visual: "nvim", editor: "ed", wantName: "nvim"},
		{name: "only EDITOR", editor: "ed", wantName: "ed"},
		{name: "only VISUAL", visual: "nvim", wantName: "nvim"},
		{name: "neither", wantReason: true},
		{name: "VISUAL empty falls through", visual: "", editor: "ed", wantName: "ed"},
		{name: "VISUAL blank falls through", visual: "   ", editor: "ed", wantName: "ed"},
		{name: "both blank", visual: " ", editor: "\t", wantReason: true},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("VISUAL", c.visual)
			t.Setenv("EDITOR", c.editor)

			cmd, reason := resolveEditor("/p/openspec/project.md")

			if c.wantReason {
				if reason == "" {
					t.Fatalf("no reason given; command was %+v", cmd)
				}
				if cmd.name != "" || len(cmd.args) != 0 {
					t.Errorf("a process was described anyway: %+v", cmd)
				}
				for _, want := range []string{"VISUAL", "EDITOR"} {
					if !strings.Contains(reason, want) {
						t.Errorf("the reason should name $%s: %q", want, reason)
					}
				}
				// Nothing is guessed.
				for _, guess := range []string{"vi", "nano", "vim"} {
					if cmd.name == guess {
						t.Errorf("fell back to %q", guess)
					}
				}
				return
			}

			if reason != "" {
				t.Fatalf("reason %q, want a command", reason)
			}
			if cmd.name != c.wantName {
				t.Errorf("command %q, want %q", cmd.name, c.wantName)
			}
		})
	}
}

func TestAnEditorValueCarryingArgumentsWorks(t *testing.T) {
	for _, c := range []struct {
		value    string
		wantName string
		wantArgs []string
	}{
		{"nvim", "nvim", []string{"/p/f.md"}},
		{"code -w", "code", []string{"-w", "/p/f.md"}},
		{"emacsclient -nw -a ''", "emacsclient", []string{"-nw", "-a", "''", "/p/f.md"}},
		{"  nvim   -R  ", "nvim", []string{"-R", "/p/f.md"}},
	} {
		t.Setenv("VISUAL", "")
		t.Setenv("EDITOR", c.value)

		cmd, reason := resolveEditor("/p/f.md")
		if reason != "" {
			t.Fatalf("%q: reason %q", c.value, reason)
		}
		if cmd.name != c.wantName {
			t.Errorf("%q: command %q, want %q", c.value, cmd.name, c.wantName)
		}
		if strings.Join(cmd.args, "\x00") != strings.Join(c.wantArgs, "\x00") {
			t.Errorf("%q: args %q, want %q", c.value, cmd.args, c.wantArgs)
		}
		// The file is the last argument, and no position within it is passed.
		if len(cmd.args) == 0 || cmd.args[len(cmd.args)-1] != "/p/f.md" {
			t.Errorf("%q: the file must be the last argument, got %q", c.value, cmd.args)
		}
		for _, a := range cmd.args[:len(cmd.args)-1] {
			if strings.HasPrefix(a, "+") || strings.Contains(a, "--goto") {
				t.Errorf("%q: a position was passed: %q", c.value, a)
			}
		}
	}
}

// --- the key ---

// fakeEditor replaces the launcher for the duration of a test.
//
// Nothing in this suite may start a real editor: it would take the terminal the
// test runner is using, and it would depend on an editor being installed.
func fakeEditor(t *testing.T, exit error) *[]*exec.Cmd {
	t.Helper()
	var started []*exec.Cmd
	previous := runEditor
	runEditor = func(c *exec.Cmd, done func(error) tea.Msg) tea.Cmd {
		started = append(started, c)
		return func() tea.Msg { return done(exit) }
	}
	t.Cleanup(func() { runEditor = previous })
	return &started
}

// editorPane is one pane and the file it should hand over.
type editorPane struct {
	name string
	set  func(m *model)
	file func(root string) string
}

func filePanes() []editorPane {
	return []editorPane{
		{"a change's proposal", func(m *model) {
			m.level = levelChange
			m.changeArtifactTab = 0
		}, func(root string) string {
			return filepath.Join(root, "openspec", "changes", "my-change", "proposal.md")
		}},
		{"a change's tasks", func(m *model) {
			m.level = levelChange
			m.changeArtifactTab = 1
		}, func(root string) string {
			return filepath.Join(root, "openspec", "changes", "my-change", "tasks.md")
		}},
		{"a spec", func(m *model) { m.detailTab = tabSpecs }, func(root string) string {
			return filepath.Join(root, "openspec", "specs", "some-capability", "spec.md")
		}},
		{"the project configuration", func(m *model) { m.detailTab = tabProperties },
			func(root string) string {
				return filepath.Join(root, "openspec", "project.md")
			}},
	}
}

func TestEOpensTheFileThePaneIsShowing(t *testing.T) {
	for _, p := range filePanes() {
		t.Run(p.name, func(t *testing.T) {
			t.Setenv("VISUAL", "")
			t.Setenv("EDITOR", "nvim")
			started := fakeEditor(t, nil)

			m, _, root, _ := onDiskStoreModel(t, threeTasks)
			p.set(&m)
			m.recalcLayout()
			m.syncDocument()

			um := press(m, tea.KeyPressMsg{Code: 'E', Text: "E"})

			if len(*started) != 1 {
				t.Fatalf("%d processes started, want 1 (status %q)", len(*started), um.statusMsg)
			}
			c := (*started)[0]
			want := p.file(root)
			if got := c.Args[len(c.Args)-1]; got != want {
				t.Errorf("opened %q, want %q", got, want)
			}
			if _, err := os.Stat(want); err != nil {
				t.Errorf("the path handed over does not resolve: %v", err)
			}
			if !strings.HasSuffix(c.Path, "nvim") && c.Path != "nvim" {
				t.Errorf("ran %q, want the editor from the environment", c.Path)
			}
		})
	}
}

func TestEIsInertWhereThereIsNoSingleFile(t *testing.T) {
	t.Run("a change's specs sub-tab", func(t *testing.T) {
		t.Setenv("EDITOR", "nvim")
		started := fakeEditor(t, nil)

		m, _, _, _ := onDiskStoreModel(t, threeTasks)
		info := m.projects[m.repoPaths[0]].Info
		info.Changes[0].SpecNames = []string{"one", "two"}
		info.Changes[0].SpecContents = map[string]string{"one": "# a\n", "two": "# b\n"}
		m.projects[m.repoPaths[0]] = scanner.ProjectStatus{Info: info}
		m.level = levelChange
		r, ok := m.selectedRow()
		if !ok {
			t.Fatal("no change selected")
		}
		m.changeArtifactTab = len(r.artifactTabNames()) - 1
		m.recalcLayout()
		m.syncDocument()

		um := press(m, tea.KeyPressMsg{Code: 'E', Text: "E"})
		if len(*started) != 0 {
			t.Errorf("a process was started for a pane showing every delta at once")
		}
		if um.statusMsg != "" {
			t.Errorf("status = %q, want nothing said", um.statusMsg)
		}
	})

	t.Run("the schema and store rows", func(t *testing.T) {
		t.Setenv("EDITOR", "nvim")
		started := fakeEditor(t, nil)

		m, _, _, _ := onDiskStoreModel(t, threeTasks)
		info := m.projects[m.repoPaths[0]].Info
		info.SchemaUsage = []scanner.SchemaUsage{{Name: "spec-driven", IsDefault: true, Changes: 1}}
		m.projects[m.repoPaths[0]] = scanner.ProjectStatus{Info: info}
		m.detailTab = tabProperties
		m.recalcLayout()

		for i, sec := range m.currentSections() {
			if sec.kind == sectionConfig {
				continue
			}
			m.propSection = i
			m.syncDocument()
			press(m, tea.KeyPressMsg{Code: 'E', Text: "E"})
			if len(*started) != 0 {
				t.Errorf("row %q started a process for an assembled report", sec.label)
			}
		}
	})
}

func TestEBelongsToAnOverlay(t *testing.T) {
	t.Run("the picker", func(t *testing.T) {
		t.Setenv("EDITOR", "nvim")
		started := fakeEditor(t, nil)
		m, _, _, _ := onDiskStoreModel(t, threeTasks)
		m.detailTab = tabSpecs
		m.recalcLayout()
		m.syncDocument()
		m.pickerOpen = true
		m.pickerAll = testProjectRows()
		m.pickerLoaded = true

		press(m, tea.KeyPressMsg{Code: 'E', Text: "E"})
		if len(*started) != 0 {
			t.Error("a process was started with the picker open")
		}
	})

	t.Run("a confirmation", func(t *testing.T) {
		t.Setenv("EDITOR", "nvim")
		started := fakeEditor(t, nil)
		m, _, _, _ := onDiskStoreModel(t, threeTasks)
		m.level = levelChange
		m.recalcLayout()
		m.syncDocument()
		m.archiveState = archiveConfirming

		press(m, tea.KeyPressMsg{Code: 'E', Text: "E"})
		if len(*started) != 0 {
			t.Error("a process was started while a confirmation awaited an answer")
		}
	})

	t.Run("the search prompt", func(t *testing.T) {
		t.Setenv("EDITOR", "nvim")
		started := fakeEditor(t, nil)
		m, _, _, _ := onDiskStoreModel(t, threeTasks)
		m.level = levelChange
		m.recalcLayout()
		m.syncDocument()
		m.level = levelProject
		m.detailTab = tabChanges
		m.searchFocused = true
		m.searchInput.Focus()

		um := press(m, tea.KeyPressMsg{Code: 'E', Text: "E"})
		if len(*started) != 0 {
			t.Error("a process was started while a filter was being typed")
		}
		if um.searchInput.Value() != "E" {
			t.Errorf("query = %q, want the E typed into the filter", um.searchInput.Value())
		}
	})
}

func TestENoEditorConfiguredReportsAndClears(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "")
	started := fakeEditor(t, nil)

	m, _, _, _ := onDiskStoreModel(t, threeTasks)
	m.detailTab = tabSpecs
	m.recalcLayout()
	m.syncDocument()

	um := press(m, tea.KeyPressMsg{Code: 'E', Text: "E"})
	if len(*started) != 0 {
		t.Error("a process was started with no editor configured")
	}
	if !strings.Contains(um.statusMsg, "no editor configured") {
		t.Errorf("status = %q, want it to say no editor is configured", um.statusMsg)
	}
	if !strings.Contains(um.View().Content, "no editor configured") {
		t.Error("the message should be visible in the view")
	}

	after := press(um, tea.KeyPressMsg{Code: 'j', Text: "j"})
	if after.statusMsg != "" {
		t.Errorf("status = %q, want it cleared by the next keystroke", after.statusMsg)
	}
}

func TestEReportsWhenTheEditorFails(t *testing.T) {
	for _, c := range []struct {
		name string
		err  error
		want string
	}{
		{"could not be started", errors.New(`exec: "nope": executable file not found in $PATH`), "not found"},
		{"exited non-zero", errors.New("exit status 1"), "exit status 1"},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("EDITOR", "nope")
			fakeEditor(t, c.err)

			m, _, _, _ := onDiskStoreModel(t, threeTasks)
			m.detailTab = tabSpecs
			m.recalcLayout()
			m.syncDocument()

			um, cmd := m.openInEditor()
			if cmd == nil {
				t.Fatal("no command was issued")
			}
			msg := cmd()
			done, ok := msg.(editorDoneMsg)
			if !ok {
				t.Fatalf("the callback returned %T, want editorDoneMsg", msg)
			}

			after, _ := um.editorFinished(done)
			if !strings.Contains(after.statusMsg, c.want) {
				t.Errorf("status = %q, want it to carry %q", after.statusMsg, c.want)
			}
		})
	}
}

func TestTheProjectIsReadAgainWhenTheEditorExits(t *testing.T) {
	// A project that resolves for real, so the command can be run and its
	// result inspected rather than only its existence asserted.
	m, _ := taskModel(t, threeTasks)
	project := m.repoPaths[0]
	m.detailTab = tabChanges
	m.recalcLayout()
	m.syncDocument()

	after, cmd := m.editorFinished(editorDoneMsg{})
	if cmd == nil {
		t.Fatal("no rescan was issued")
	}
	if after.statusMsg != "" {
		t.Errorf("status = %q, want nothing said after a clean exit", after.statusMsg)
	}
	// The command reads the project, which is what puts an edit on screen
	// without the user asking for it.
	msg, ok := cmd().(scanMsg)
	if !ok {
		t.Fatalf("the command returned %T, want a scan", cmd())
	}
	if msg.err != nil {
		t.Fatalf("the scan failed: %v", msg.err)
	}
	if _, found := msg.projects[project]; !found {
		t.Errorf("the scan did not cover the open project; it covered %d others",
			len(msg.projects))
	}
}

// --- 5.1 the nav bar ---

func TestTheNavBarListsEExactlyWhereTheKeyApplies(t *testing.T) {
	withFile := filePanes()
	for _, p := range withFile {
		m, _, _, _ := onDiskStoreModel(t, threeTasks)
		m.width = 260 // wide enough that no hint is dropped
		p.set(&m)
		m.recalcLayout()
		m.syncDocument()

		if !strings.Contains(ansi.Strip(m.renderNavBar()), "E edit") {
			t.Errorf("%s shows a file, so the nav bar should list E:\n%s",
				p.name, ansi.Strip(m.renderNavBar()))
		}
	}

	without := []editorPane{
		{name: "the change list", set: func(m *model) { m.detailTab = tabChanges }},
		{name: "a change's specs sub-tab", set: func(m *model) {
			info := m.projects[m.repoPaths[0]].Info
			info.Changes[0].SpecNames = []string{"one", "two"}
			info.Changes[0].SpecContents = map[string]string{"one": "# a\n", "two": "# b\n"}
			m.projects[m.repoPaths[0]] = scanner.ProjectStatus{Info: info}
			m.level = levelChange
			m.changeArtifactTab = 2
		}},
		{name: "the store row", set: func(m *model) {
			m.detailTab = tabProperties
			m.propSection = len(m.currentSections()) - 1
		}},
	}
	for _, p := range without {
		m, _, _, _ := onDiskStoreModel(t, threeTasks)
		m.width = 260
		p.set(&m)
		m.recalcLayout()
		m.syncDocument()

		if strings.Contains(ansi.Strip(m.renderNavBar()), "E edit") {
			t.Errorf("%s shows no single file, so the nav bar should not list E:\n%s",
				p.name, ansi.Strip(m.renderNavBar()))
		}
	}
}

func TestNothingIsReadAgainWithNoProjectOpen(t *testing.T) {
	m := newModel(&scanner.Config{}, true, "0.0.0")

	after, cmd := m.editorFinished(editorDoneMsg{})
	if cmd != nil {
		t.Error("there is no project to read again")
	}
	if after.statusMsg != "" {
		t.Errorf("status = %q, want nothing said", after.statusMsg)
	}
}
