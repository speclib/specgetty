package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/mipmip/specgetty/src/scanner"
)

func TestMainPanelHeight(t *testing.T) {
	tests := []struct {
		name   string
		height int
	}{
		{"small terminal", 20},
		{"medium terminal", 40},
		{"large terminal", 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model{height: tt.height}
			got := m.mainPanelHeight()
			if got < 3 {
				t.Errorf("got %d, want >= 3", got)
			}
		})
	}
}

func TestLogPanelHeight(t *testing.T) {
	tests := []struct {
		name   string
		height int
	}{
		{"small terminal", 20},
		{"medium terminal", 40},
		{"large terminal", 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model{height: tt.height, logVisible: true}
			got := m.logPanelHeight()
			if got < 1 {
				t.Errorf("got %d, want >= 1", got)
			}
			if got > 10 {
				t.Errorf("got %d, want <= 10", got)
			}
		})
	}

	t.Run("hidden log returns 0", func(t *testing.T) {
		m := model{height: 40, logVisible: false}
		got := m.logPanelHeight()
		if got != 0 {
			t.Errorf("got %d, want 0 when log hidden", got)
		}
	})
}

func TestProjectDisplayNames(t *testing.T) {
	t.Run("unique basenames", func(t *testing.T) {
		paths := []string{"/home/user/project-a", "/home/user/project-b"}
		names := projectDisplayNames(paths)
		if names[0] != "project-a" {
			t.Errorf("got %q, want project-a", names[0])
		}
		if names[1] != "project-b" {
			t.Errorf("got %q, want project-b", names[1])
		}
	})

	t.Run("duplicate basenames get parent disambiguation", func(t *testing.T) {
		paths := []string{"/home/user/work/myapp", "/home/user/personal/myapp"}
		names := projectDisplayNames(paths)
		if !strings.Contains(names[0], "myapp") || !strings.Contains(names[0], "work") {
			t.Errorf("expected 'myapp (work)', got %q", names[0])
		}
		if !strings.Contains(names[1], "myapp") || !strings.Contains(names[1], "personal") {
			t.Errorf("expected 'myapp (personal)', got %q", names[1])
		}
	})

	t.Run("empty list", func(t *testing.T) {
		names := projectDisplayNames([]string{})
		if len(names) != 0 {
			t.Errorf("expected empty, got %v", names)
		}
	})
}

func TestTabOnlyMovesWhenTheLogPanelIsOpen(t *testing.T) {
	t.Run("nowhere to go with the log hidden", func(t *testing.T) {
		m := makeListModel()
		m.logVisible = false

		updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
		if um := updated.(model); um.focus == focusLog {
			t.Errorf("focus = %d, want it to stay out of the log", um.focus)
		}
	})

	t.Run("toggles to the log and back", func(t *testing.T) {
		m := makeListModel()
		m.logVisible = true

		updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
		m = updated.(model)
		if m.focus != focusLog {
			t.Fatalf("focus = %d, want the log panel", m.focus)
		}

		updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
		if um := updated.(model); um.focus == focusLog {
			t.Errorf("focus = %d, want it back out of the log", um.focus)
		}
	})
}

func TestRenderTabHeader(t *testing.T) {
	t.Run("specs tab active", func(t *testing.T) {
		m := model{detailTab: tabSpecs}
		got := m.renderTabHeader(80)
		if !strings.Contains(got, "specs") {
			t.Error("output missing 'specs'")
		}
		if !strings.Contains(got, "changes") {
			t.Error("output missing 'changes'")
		}
		if strings.Contains(got, "overview") {
			t.Error("output should not contain 'overview'")
		}
	})

	t.Run("all tab names present", func(t *testing.T) {
		m := model{detailTab: tabChanges}
		got := m.renderTabHeader(80)
		for _, name := range tabNames {
			if !strings.Contains(got, name) {
				t.Errorf("output missing tab %q", name)
			}
		}
	})
}

func TestStatsLineIncludesTasks(t *testing.T) {
	t.Run("shows tasks when present", func(t *testing.T) {
		m := model{
			width:        100,
			height:       40,
			repoPaths:    []string{"/test/project"},
			displayNames: []string{"project"},
			cursor:       0,
			projects: scanner.ProjectMap{
				"/test/project": scanner.ProjectStatus{
					Info: scanner.ProjectInfo{
						SpecCount:     2,
						ActiveChanges: []string{"ch1"},
						TasksTotal:    10,
						TasksDone:     4,
					},
				},
			},
		}
		got := m.renderDetailPanel(80, 30)
		if !strings.Contains(got, "Tasks: 4/10") {
			t.Errorf("expected 'Tasks: 4/10' in output, got:\n%s", got)
		}
	})

	t.Run("omits tasks when none", func(t *testing.T) {
		m := model{
			width:        100,
			height:       40,
			repoPaths:    []string{"/test/project"},
			displayNames: []string{"project"},
			cursor:       0,
			projects: scanner.ProjectMap{
				"/test/project": scanner.ProjectStatus{
					Info: scanner.ProjectInfo{
						SpecCount:     2,
						ActiveChanges: []string{"ch1"},
						TasksTotal:    0,
						TasksDone:     0,
					},
				},
			},
		}
		got := m.renderDetailPanel(80, 30)
		if strings.Contains(got, "Tasks:") {
			t.Errorf("expected no 'Tasks:' in output, got:\n%s", got)
		}
	})
}

func TestRenderMarkdown(t *testing.T) {
	t.Run("headers are present", func(t *testing.T) {
		content := "# Title\n\nSome text\n\n## Subtitle"
		got := renderMarkdown(content, 80)
		if !strings.Contains(got, "Title") {
			t.Error("output missing 'Title'")
		}
		if !strings.Contains(got, "Subtitle") {
			t.Error("output missing 'Subtitle'")
		}
	})

	t.Run("list items are present", func(t *testing.T) {
		content := "- item one\n- item two"
		got := renderMarkdown(content, 80)
		if !strings.Contains(got, "item one") {
			t.Error("output missing 'item one'")
		}
		if !strings.Contains(got, "item two") {
			t.Error("output missing 'item two'")
		}
	})

	t.Run("bold text is present", func(t *testing.T) {
		content := "This is **bold** text"
		got := renderMarkdown(content, 80)
		if !strings.Contains(got, "bold") {
			t.Error("output missing 'bold'")
		}
		// Should not contain the ** markers in plain form
		if strings.Contains(got, "**bold**") {
			t.Error("output still contains raw ** markers")
		}
	})
}

func makeArchiveTestModel() model {
	return model{
		width:     100,
		height:    40,
		repoPaths: []string{"/test/project"},
		cursor:    0,
		detailTab: tabChanges,
		projects: scanner.ProjectMap{
			"/test/project": scanner.ProjectStatus{
				Info: scanner.ProjectInfo{
					ActiveChanges: []string{"my-change"},
					Changes: []scanner.ChangeInfo{
						{Name: "my-change", TasksTotal: 5, TasksDone: 3},
					},
				},
			},
		},
	}
}

func TestArchiveKeyStartsConfirming(t *testing.T) {
	m := makeArchiveTestModel()
	m.focus = focusDetail

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	um := updated.(model)
	if um.archiveState != archiveConfirming {
		t.Errorf("archiveState = %d, want %d (archiveConfirming)", um.archiveState, archiveConfirming)
	}
	if um.archiveChangeName != "my-change" {
		t.Errorf("archiveChangeName = %q, want my-change", um.archiveChangeName)
	}
}

func TestArchiveConfirmY(t *testing.T) {
	m := makeArchiveTestModel()
	m.archiveState = archiveConfirming
	m.archiveChangeName = "my-change"

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'y', Text: "y"})
	um := updated.(model)
	if um.archiveState != archiveRunning {
		t.Errorf("archiveState = %d, want %d (archiveRunning)", um.archiveState, archiveRunning)
	}
}

func TestArchiveConfirmN(t *testing.T) {
	m := makeArchiveTestModel()
	m.archiveState = archiveConfirming
	m.archiveChangeName = "my-change"

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	um := updated.(model)
	if um.archiveState != archiveIdle {
		t.Errorf("archiveState = %d, want %d (archiveIdle)", um.archiveState, archiveIdle)
	}
}

func TestRenderYAML(t *testing.T) {
	t.Run("keys and values present", func(t *testing.T) {
		content := "schema: spec-driven\nname: my-project"
		got := renderYAML(content, 80)
		if !strings.Contains(got, "schema") {
			t.Error("output missing 'schema'")
		}
		if !strings.Contains(got, "spec-driven") {
			t.Error("output missing 'spec-driven'")
		}
	})

	t.Run("comments present", func(t *testing.T) {
		content := "# This is a comment\nkey: value"
		got := renderYAML(content, 80)
		if !strings.Contains(got, "This is a comment") {
			t.Error("output missing comment text")
		}
		if !strings.Contains(got, "key") {
			t.Error("output missing 'key'")
		}
	})
}
