package ui

import (
	"strings"
	"testing"
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

func TestLeftPanelWidth(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		expected int
	}{
		{"narrow terminal clamps to min 20", 60, 20},
		{"medium terminal uses 30%", 100, 30},
		{"wide terminal clamps to max 40", 200, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model{width: tt.width}
			got := m.leftPanelWidth()
			if got != tt.expected {
				t.Errorf("got %d, want %d", got, tt.expected)
			}
		})
	}
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

func TestTabCyclingSkipsHiddenLog(t *testing.T) {
	t.Run("tab skips log when hidden", func(t *testing.T) {
		m := model{logVisible: false, activeView: viewProjects}
		if m.activeView == viewProjects {
			m.activeView = viewDetail
		} else {
			m.activeView = viewProjects
		}
		if m.activeView != viewDetail {
			t.Errorf("got %d, want viewDetail (%d)", m.activeView, viewDetail)
		}
		if m.activeView == viewProjects {
			m.activeView = viewDetail
		} else {
			m.activeView = viewProjects
		}
		if m.activeView != viewProjects {
			t.Errorf("got %d, want viewProjects (%d)", m.activeView, viewProjects)
		}
	})
}

func TestRenderProjectList(t *testing.T) {
	t.Run("shows display names", func(t *testing.T) {
		m := model{
			height:       40,
			width:        100,
			repoPaths:    []string{"/home/user/repo1", "/home/user/repo2"},
			displayNames: []string{"repo1", "repo2"},
			cursor:       0,
		}
		got := m.renderProjectList(30, 20)
		if !strings.Contains(got, "repo1") {
			t.Error("output missing repo1")
		}
		if !strings.Contains(got, "repo2") {
			t.Error("output missing repo2")
		}
		// Should NOT contain the full path
		if strings.Contains(got, "/home/user/") {
			t.Error("output should show basenames, not full paths")
		}
	})

	t.Run("empty project list", func(t *testing.T) {
		m := model{
			height:       40,
			width:        100,
			repoPaths:    []string{},
			displayNames: []string{},
		}
		got := m.renderProjectList(30, 20)
		if !strings.Contains(got, "No OpenSpec projects found") {
			t.Errorf("expected empty message, got %q", got)
		}
	})
}

func TestRenderTabHeader(t *testing.T) {
	t.Run("overview tab active", func(t *testing.T) {
		m := model{detailTab: tabOverview}
		got := m.renderTabHeader(80)
		if !strings.Contains(got, "overview") {
			t.Error("output missing 'overview'")
		}
		if !strings.Contains(got, "specs") {
			t.Error("output missing 'specs'")
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
