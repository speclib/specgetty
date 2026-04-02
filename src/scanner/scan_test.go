package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	defaultConfig := `
scandirs:
  include:
    - /default/path
`

	t.Run("valid config file", func(t *testing.T) {
		f, err := os.CreateTemp("", "drs-test-*.yml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())

		content := `
scandirs:
  include:
    - /home/user/projects
    - /home/user/work
  exclude:
    - /home/user/projects/vendor
`
		if _, err := f.WriteString(content); err != nil {
			t.Fatal(err)
		}
		f.Close()

		config, err := ParseConfigFile(f.Name(), defaultConfig)
		if err != nil {
			t.Fatalf("ParseConfigFile: %v", err)
		}
		if len(config.ScanDirs.Include) != 2 {
			t.Errorf("got %d include dirs, want 2", len(config.ScanDirs.Include))
		}
		if config.ScanDirs.Include[0] != "/home/user/projects" {
			t.Errorf("include[0] = %q, want /home/user/projects", config.ScanDirs.Include[0])
		}
		if len(config.ScanDirs.Exclude) != 1 {
			t.Errorf("got %d exclude dirs, want 1", len(config.ScanDirs.Exclude))
		}
	})

	t.Run("file does not exist falls back to default", func(t *testing.T) {
		config, err := ParseConfigFile("/nonexistent/path/config.yml", defaultConfig)
		if err != nil {
			t.Fatalf("ParseConfigFile: %v", err)
		}
		if len(config.ScanDirs.Include) != 1 {
			t.Errorf("got %d include dirs, want 1", len(config.ScanDirs.Include))
		}
		if config.ScanDirs.Include[0] != "/default/path" {
			t.Errorf("include[0] = %q, want /default/path", config.ScanDirs.Include[0])
		}
	})

	t.Run("invalid YAML returns error", func(t *testing.T) {
		f, err := os.CreateTemp("", "drs-test-*.yml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())

		if _, err := f.WriteString("{{invalid yaml"); err != nil {
			t.Fatal(err)
		}
		f.Close()

		_, err = ParseConfigFile(f.Name(), defaultConfig)
		if err == nil {
			t.Error("expected error for invalid YAML, got nil")
		}
	})
}

func TestParseProjectInfo(t *testing.T) {
	t.Run("project with specs and changes", func(t *testing.T) {
		dir := t.TempDir()
		openspec := filepath.Join(dir, "openspec")
		os.MkdirAll(filepath.Join(openspec, "specs", "cap-a"), 0755)
		os.MkdirAll(filepath.Join(openspec, "specs", "cap-b"), 0755)
		os.MkdirAll(filepath.Join(openspec, "specs", "cap-c"), 0755)
		os.WriteFile(filepath.Join(openspec, "specs", "cap-a", "spec.md"), []byte("# Cap A"), 0644)
		os.WriteFile(filepath.Join(openspec, "specs", "cap-b", "spec.md"), []byte("# Cap B"), 0644)
		// cap-c has no spec.md
		os.MkdirAll(filepath.Join(openspec, "changes", "change-1"), 0755)
		os.MkdirAll(filepath.Join(openspec, "changes", "change-2"), 0755)
		os.MkdirAll(filepath.Join(openspec, "changes", "archive", "old-change"), 0755)

		info := ParseProjectInfo(dir)
		if info.SpecCount != 3 {
			t.Errorf("SpecCount = %d, want 3", info.SpecCount)
		}
		if len(info.SpecNames) != 3 {
			t.Errorf("SpecNames = %d, want 3", len(info.SpecNames))
		}
		if info.SpecNames[0] != "cap-a" {
			t.Errorf("SpecNames[0] = %q, want cap-a", info.SpecNames[0])
		}
		if info.SpecContents["cap-a"] != "# Cap A" {
			t.Errorf("SpecContents[cap-a] = %q, want '# Cap A'", info.SpecContents["cap-a"])
		}
		if _, ok := info.SpecContents["cap-c"]; ok {
			t.Error("cap-c should not have content (no spec.md)")
		}
		if len(info.ActiveChanges) != 2 {
			t.Errorf("ActiveChanges = %d, want 2", len(info.ActiveChanges))
		}
		if len(info.ArchivedChanges) != 1 {
			t.Errorf("ArchivedChanges = %d, want 1", len(info.ArchivedChanges))
		}
		if info.ArchivedChanges[0].Name != "old-change" {
			t.Errorf("ArchivedChanges[0].Name = %q, want old-change", info.ArchivedChanges[0].Name)
		}
	})

	t.Run("empty openspec directory", func(t *testing.T) {
		dir := t.TempDir()
		os.MkdirAll(filepath.Join(dir, "openspec"), 0755)

		info := ParseProjectInfo(dir)
		if info.SpecCount != 0 {
			t.Errorf("SpecCount = %d, want 0", info.SpecCount)
		}
		if len(info.ActiveChanges) != 0 {
			t.Errorf("ActiveChanges = %d, want 0", len(info.ActiveChanges))
		}
		if len(info.ArchivedChanges) != 0 {
			t.Errorf("ArchivedChanges = %d, want 0", len(info.ArchivedChanges))
		}
	})

	t.Run("no openspec directory", func(t *testing.T) {
		dir := t.TempDir()

		info := ParseProjectInfo(dir)
		if info.SpecCount != 0 {
			t.Errorf("SpecCount = %d, want 0", info.SpecCount)
		}
	})

	t.Run("project.md takes priority over config.yaml", func(t *testing.T) {
		dir := t.TempDir()
		openspec := filepath.Join(dir, "openspec")
		os.MkdirAll(openspec, 0755)
		os.WriteFile(filepath.Join(openspec, "project.md"), []byte("# My Project"), 0644)
		os.WriteFile(filepath.Join(openspec, "config.yaml"), []byte("schema: spec-driven"), 0644)

		info := ParseProjectInfo(dir)
		if info.ConfigFile != "project.md" {
			t.Errorf("ConfigFile = %q, want project.md", info.ConfigFile)
		}
		if info.ConfigContent != "# My Project" {
			t.Errorf("ConfigContent = %q, want '# My Project'", info.ConfigContent)
		}
	})

	t.Run("config.yaml when no project.md", func(t *testing.T) {
		dir := t.TempDir()
		openspec := filepath.Join(dir, "openspec")
		os.MkdirAll(openspec, 0755)
		os.WriteFile(filepath.Join(openspec, "config.yaml"), []byte("schema: spec-driven"), 0644)

		info := ParseProjectInfo(dir)
		if info.ConfigFile != "config.yaml" {
			t.Errorf("ConfigFile = %q, want config.yaml", info.ConfigFile)
		}
		if info.ConfigContent != "schema: spec-driven" {
			t.Errorf("ConfigContent = %q, want 'schema: spec-driven'", info.ConfigContent)
		}
	})

	t.Run("no config file", func(t *testing.T) {
		dir := t.TempDir()
		openspec := filepath.Join(dir, "openspec")
		os.MkdirAll(openspec, 0755)

		info := ParseProjectInfo(dir)
		if info.ConfigFile != "" {
			t.Errorf("ConfigFile = %q, want empty", info.ConfigFile)
		}
	})

	t.Run("changes with artifacts and tasks", func(t *testing.T) {
		dir := t.TempDir()
		openspec := filepath.Join(dir, "openspec")
		changeDir := filepath.Join(openspec, "changes", "my-change")
		os.MkdirAll(changeDir, 0755)
		os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte("# Proposal"), 0644)
		os.WriteFile(filepath.Join(changeDir, "design.md"), []byte("# Design"), 0644)
		os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte("- [x] 1.1 Done\n- [ ] 1.2 Todo\n- [x] 1.3 Done"), 0644)
		os.MkdirAll(filepath.Join(changeDir, "specs", "cap-a"), 0755)
		os.WriteFile(filepath.Join(changeDir, "specs", "cap-a", "spec.md"), []byte("# Spec A"), 0644)

		info := ParseProjectInfo(dir)
		if len(info.Changes) != 1 {
			t.Fatalf("Changes = %d, want 1", len(info.Changes))
		}
		ci := info.Changes[0]
		if ci.Name != "my-change" {
			t.Errorf("Name = %q, want my-change", ci.Name)
		}
		if len(ci.ArtifactFiles) != 3 {
			t.Errorf("ArtifactFiles = %d, want 3", len(ci.ArtifactFiles))
		}
		if ci.TasksTotal != 3 {
			t.Errorf("TasksTotal = %d, want 3", ci.TasksTotal)
		}
		if ci.TasksDone != 2 {
			t.Errorf("TasksDone = %d, want 2", ci.TasksDone)
		}
		if len(ci.SpecNames) != 1 {
			t.Errorf("SpecNames = %d, want 1", len(ci.SpecNames))
		}
		if ci.SpecContents["cap-a"] != "# Spec A" {
			t.Errorf("SpecContents[cap-a] = %q, want '# Spec A'", ci.SpecContents["cap-a"])
		}
	})

	t.Run("aggregate task stats across changes", func(t *testing.T) {
		dir := t.TempDir()
		openspec := filepath.Join(dir, "openspec")
		change1 := filepath.Join(openspec, "changes", "change-a")
		change2 := filepath.Join(openspec, "changes", "change-b")
		os.MkdirAll(change1, 0755)
		os.MkdirAll(change2, 0755)
		os.WriteFile(filepath.Join(change1, "tasks.md"), []byte("- [x] 1.1 Done\n- [ ] 1.2 Todo\n- [x] 1.3 Done"), 0644)
		os.WriteFile(filepath.Join(change2, "tasks.md"), []byte("- [ ] 1.1 Todo\n- [ ] 1.2 Todo\n- [x] 1.3 Done\n- [x] 1.4 Done"), 0644)

		info := ParseProjectInfo(dir)
		if info.TasksTotal != 7 {
			t.Errorf("TasksTotal = %d, want 7", info.TasksTotal)
		}
		if info.TasksDone != 4 {
			t.Errorf("TasksDone = %d, want 4", info.TasksDone)
		}
	})

	t.Run("aggregate task stats zero when no tasks", func(t *testing.T) {
		dir := t.TempDir()
		openspec := filepath.Join(dir, "openspec")
		os.MkdirAll(filepath.Join(openspec, "changes", "no-tasks"), 0755)
		os.WriteFile(filepath.Join(openspec, "changes", "no-tasks", "proposal.md"), []byte("# Proposal"), 0644)

		info := ParseProjectInfo(dir)
		if info.TasksTotal != 0 {
			t.Errorf("TasksTotal = %d, want 0", info.TasksTotal)
		}
		if info.TasksDone != 0 {
			t.Errorf("TasksDone = %d, want 0", info.TasksDone)
		}
	})

	t.Run("change without tasks.md", func(t *testing.T) {
		dir := t.TempDir()
		openspec := filepath.Join(dir, "openspec")
		changeDir := filepath.Join(openspec, "changes", "simple-change")
		os.MkdirAll(changeDir, 0755)
		os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte("# Proposal"), 0644)

		info := ParseProjectInfo(dir)
		if len(info.Changes) != 1 {
			t.Fatalf("Changes = %d, want 1", len(info.Changes))
		}
		if info.Changes[0].TasksTotal != 0 {
			t.Errorf("TasksTotal = %d, want 0", info.Changes[0].TasksTotal)
		}
	})
}

func TestParseTaskStats(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantTotal int
		wantDone  int
	}{
		{"mixed checkboxes", "- [x] Done\n- [ ] Todo\n- [x] Also done", 3, 2},
		{"all done", "- [x] A\n- [x] B", 2, 2},
		{"none done", "- [ ] A\n- [ ] B", 2, 0},
		{"no checkboxes", "# Just a header\nSome text", 0, 0},
		{"empty", "", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total, done := ParseTaskStats(tt.content)
			if total != tt.wantTotal {
				t.Errorf("total = %d, want %d", total, tt.wantTotal)
			}
			if done != tt.wantDone {
				t.Errorf("done = %d, want %d", done, tt.wantDone)
			}
		})
	}
}
