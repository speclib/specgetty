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
		os.MkdirAll(filepath.Join(openspec, "specs", "nested", "inner"), 0755)
		os.WriteFile(filepath.Join(openspec, "specs", "cap-a", "spec.md"), []byte("# Cap A"), 0644)
		os.WriteFile(filepath.Join(openspec, "specs", "cap-b", "spec.md"), []byte("# Cap B"), 0644)
		os.WriteFile(filepath.Join(openspec, "specs", "nested", "inner", "spec.md"), []byte("# Inner"), 0644)
		// cap-c holds no spec.md and groups nothing, so it is not a capability
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
		if info.SpecNames[2] != "nested/inner" {
			t.Errorf("SpecNames[2] = %q, want nested/inner", info.SpecNames[2])
		}
		if info.SpecContents["cap-a"] != "# Cap A" {
			t.Errorf("SpecContents[cap-a] = %q, want '# Cap A'", info.SpecContents["cap-a"])
		}
		if info.SpecContents["nested/inner"] != "# Inner" {
			t.Errorf("SpecContents[nested/inner] = %q, want '# Inner'",
				info.SpecContents["nested/inner"])
		}
		if _, ok := info.SpecContents["cap-c"]; ok {
			t.Error("cap-c should not be a capability (no spec.md)")
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

func TestCapabilities(t *testing.T) {
	write := func(t *testing.T, dir string, files ...string) {
		t.Helper()
		for _, f := range files {
			path := filepath.Join(dir, filepath.FromSlash(f))
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte("# Spec"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	specs := func(subs ...string) string {
		dir := t.TempDir()
		write(t, dir, subs...)
		return dir
	}
	equal := func(t *testing.T, got, want []string) {
		t.Helper()
		if len(got) != len(want) {
			t.Fatalf("capabilities = %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("capabilities = %v, want %v", got, want)
			}
		}
	}

	t.Run("two nested capabilities", func(t *testing.T) {
		dir := specs("content/editing/spec.md", "content/pages/spec.md")
		equal(t, capabilities(dir), []string{"content/editing", "content/pages"})
	})

	t.Run("a container holding no spec of its own", func(t *testing.T) {
		dir := specs("tenancy/organisations/spec.md", "tenancy/projects/spec.md")
		equal(t, capabilities(dir), []string{
			"tenancy/organisations", "tenancy/projects",
		})
	})

	t.Run("a container that is also a capability", func(t *testing.T) {
		dir := specs("platform/spec.md", "platform/release/spec.md")
		equal(t, capabilities(dir), []string{"platform", "platform/release"})
	})

	t.Run("a flat tree is unchanged", func(t *testing.T) {
		dir := specs("cap-a/spec.md", "cap-b/spec.md", "cap-c")
		equal(t, capabilities(dir), []string{"cap-a", "cap-b"})
	})

	t.Run("no specs directory", func(t *testing.T) {
		equal(t, capabilities(filepath.Join(t.TempDir(), "specs")), nil)
	})

	t.Run("an unreadable subtree is skipped", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("root reads a directory whatever its mode says")
		}
		dir := specs("content/pages/spec.md", "locked/hidden/spec.md")
		locked := filepath.Join(dir, "locked")
		if err := os.Chmod(locked, 0000); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.Chmod(locked, 0755) })
		equal(t, capabilities(dir), []string{"content/pages"})
	})
}

func TestNestedDeltaNamesMatchLiveSpecNames(t *testing.T) {
	dir := t.TempDir()
	openspec := filepath.Join(dir, "openspec")
	liveDir := filepath.Join(openspec, "specs", "tenancy", "projects")
	deltaDir := filepath.Join(openspec, "changes", "nested-change", "specs", "tenancy", "projects")
	for _, d := range []string{liveDir, deltaDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}
	os.WriteFile(filepath.Join(liveDir, "spec.md"), []byte("# Projects"), 0644)
	os.WriteFile(filepath.Join(deltaDir, "spec.md"), []byte("## MODIFIED Requirements"), 0644)

	info := ParseProjectInfo(dir)
	if len(info.Changes) != 1 {
		t.Fatalf("Changes = %d, want 1", len(info.Changes))
	}
	if len(info.Changes[0].SpecNames) != 1 {
		t.Fatalf("delta names = %v, want one", info.Changes[0].SpecNames)
	}
	name := info.Changes[0].SpecNames[0]
	if name != "tenancy/projects" {
		t.Errorf("delta name = %q, want tenancy/projects", name)
	}
	if _, ok := info.SpecContents[name]; !ok {
		t.Errorf("no live spec named %q, so the delta has no original; specs are %v",
			name, info.SpecNames)
	}
}
