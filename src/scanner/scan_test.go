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
		os.MkdirAll(filepath.Join(openspec, "changes", "change-1"), 0755)
		os.MkdirAll(filepath.Join(openspec, "changes", "change-2"), 0755)
		os.MkdirAll(filepath.Join(openspec, "archive", "old-change"), 0755)

		info := ParseProjectInfo(dir)
		if info.SpecCount != 3 {
			t.Errorf("SpecCount = %d, want 3", info.SpecCount)
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
}
