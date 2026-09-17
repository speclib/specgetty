package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mipmip/specgetty/src/scanner"
)

func TestFindOpenSpecProject(t *testing.T) {
	t.Run("finds openspec in current dir", func(t *testing.T) {
		dir := t.TempDir()
		os.MkdirAll(filepath.Join(dir, "openspec"), 0755)

		got := findOpenSpecProject(dir)
		if got != dir {
			t.Errorf("got %q, want %q", got, dir)
		}
	})

	t.Run("finds openspec in parent dir", func(t *testing.T) {
		parent := t.TempDir()
		os.MkdirAll(filepath.Join(parent, "openspec"), 0755)
		child := filepath.Join(parent, "subdir")
		os.MkdirAll(child, 0755)

		got := findOpenSpecProject(child)
		if got != parent {
			t.Errorf("got %q, want %q", got, parent)
		}
	})

	t.Run("returns empty when no openspec found", func(t *testing.T) {
		dir := t.TempDir()

		got := findOpenSpecProject(dir)
		if got != "" {
			t.Errorf("got %q, want empty", got)
		}
	})
}

func TestFindOpenSpecProjectIgnoresAFile(t *testing.T) {
	// A file named openspec is not a project; only a directory counts.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "openspec"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	if got := findOpenSpecProject(dir); got != "" {
		t.Errorf("got %q, want empty for a file named openspec", got)
	}
}

func TestFindOpenSpecProjectNearestWins(t *testing.T) {
	// Walking up must stop at the first project, not continue to an outer one.
	outer := t.TempDir()
	if err := os.MkdirAll(filepath.Join(outer, "openspec"), 0755); err != nil {
		t.Fatal(err)
	}
	inner := filepath.Join(outer, "nested")
	if err := os.MkdirAll(filepath.Join(inner, "openspec"), 0755); err != nil {
		t.Fatal(err)
	}
	if got := findOpenSpecProject(inner); got != inner {
		t.Errorf("got %q, want the nearest project %q", got, inner)
	}
}

func TestFindOpenSpecProjectRejectsBadPath(t *testing.T) {
	if got := findOpenSpecProject(""); got != "" && !filepath.IsAbs(got) {
		t.Errorf("got %q, want either empty or an absolute path", got)
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	got := getDefaultConfigPath()
	if !filepath.IsAbs(got) {
		t.Errorf("got %q, want an absolute path", got)
	}
	if filepath.Base(got) != "config.yml" {
		t.Errorf("got %q, want it to end in config.yml", got)
	}
	if filepath.Base(filepath.Dir(got)) != "specgetty" {
		t.Errorf("got %q, want it inside a specgetty directory", got)
	}
}

func TestGetDefaultConfigPathFallsBackToHome(t *testing.T) {
	// With XDG_CONFIG_HOME unset, os.UserConfigDir falls back to ~/.config.
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "/tmp/specgetty-test-home")

	got := getDefaultConfigPath()
	want := filepath.Join("/tmp/specgetty-test-home", ".config", "specgetty", "config.yml")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolveStartView(t *testing.T) {
	tests := []struct {
		view, path string
		want       string
		wantErr    bool
	}{
		{"single", "", "single", false},
		{"all", "", "all", false},
		{"bogus", "", "", true},
		{"", "", "", true},
		// An explicit path names one project, so it settles the view.
		{"all", "/some/project", "single", false},
		{"single", "/some/project", "single", false},
	}
	for _, tt := range tests {
		got, err := resolveStartView(tt.view, tt.path)
		if tt.wantErr {
			if err == nil {
				t.Errorf("resolveStartView(%q,%q) returned no error, want one", tt.view, tt.path)
			}
			continue
		}
		if err != nil {
			t.Errorf("resolveStartView(%q,%q): %v", tt.view, tt.path, err)
			continue
		}
		if got != tt.want {
			t.Errorf("resolveStartView(%q,%q) = %q, want %q", tt.view, tt.path, got, tt.want)
		}
	}
}

func TestResolveStartViewErrorNamesValidValues(t *testing.T) {
	_, err := resolveStartView("bogus", "")
	if err == nil {
		t.Fatal("expected an error")
	}
	for _, want := range []string{"bogus", "single", "all"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should mention %q", err, want)
		}
	}
}

func TestZoomFlagIsGone(t *testing.T) {
	// --zoom named a mode that no longer exists. Its absence is part of the
	// contract, not an oversight, so it is asserted rather than remembered.
	for _, f := range appFlags() {
		for _, name := range f.Names() {
			if name == "zoom" || name == "z" {
				t.Errorf("flag %q is still defined; --zoom was removed with zoom mode", name)
			}
		}
	}
}

func TestViewFlagExists(t *testing.T) {
	found := false
	for _, f := range appFlags() {
		for _, name := range f.Names() {
			if name == "view" {
				found = true
			}
		}
	}
	if !found {
		t.Error("--view is missing; it replaces --zoom")
	}
}

func TestExpandScanDirs(t *testing.T) {
	t.Setenv("SPECGETTY_TEST_ROOT", "/somewhere")

	config := &scanner.Config{}
	config.ScanDirs.Include = []string{"$SPECGETTY_TEST_ROOT/code", "/literal"}
	config.ScanDirs.Exclude = []string{"$SPECGETTY_TEST_ROOT/junk", "vendor"}

	expandScanDirs(config)

	if config.ScanDirs.Include[0] != "/somewhere/code" {
		t.Errorf("include[0] = %q, want it expanded", config.ScanDirs.Include[0])
	}
	if config.ScanDirs.Include[1] != "/literal" {
		t.Errorf("include[1] = %q, want it untouched", config.ScanDirs.Include[1])
	}
	if config.ScanDirs.Exclude[0] != "/somewhere/junk" {
		t.Errorf("exclude[0] = %q, want it expanded", config.ScanDirs.Exclude[0])
	}
	if config.ScanDirs.Exclude[1] != "vendor" {
		t.Errorf("exclude[1] = %q, want it untouched", config.ScanDirs.Exclude[1])
	}
}

func TestExpandScanDirsLeavesAnUnsetVariableEmpty(t *testing.T) {
	// os.ExpandEnv turns an unset variable into an empty string, which is worth
	// knowing: a typo in the config silently produces a path of "/code".
	t.Setenv("SPECGETTY_TEST_UNSET", "")
	config := &scanner.Config{}
	config.ScanDirs.Include = []string{"$SPECGETTY_TEST_UNSET/code"}

	expandScanDirs(config)

	if config.ScanDirs.Include[0] != "/code" {
		t.Errorf("got %q, want the unset variable to vanish", config.ScanDirs.Include[0])
	}
}

func TestResolveStartupPath(t *testing.T) {
	t.Run("an explicit path is made absolute", func(t *testing.T) {
		got := resolveStartupPath("single", "some/where")
		if !filepath.IsAbs(got) {
			t.Errorf("got %q, want an absolute path", got)
		}
		if filepath.Base(got) != "where" {
			t.Errorf("got %q, want it to end at the given path", got)
		}
	})

	t.Run("view=all resolves nothing", func(t *testing.T) {
		if got := resolveStartupPath("all", "/some/where"); got != "" {
			t.Errorf("got %q, want nothing: the picker opens instead", got)
		}
	})

	t.Run("found by walking up from the working directory", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "openspec"), 0o755); err != nil {
			t.Fatal(err)
		}
		nested := filepath.Join(root, "a", "b")
		if err := os.MkdirAll(nested, 0o755); err != nil {
			t.Fatal(err)
		}
		t.Chdir(nested)

		got := resolveStartupPath("single", "")
		// Compare resolved paths: a temp dir can sit behind a symlink.
		wantResolved, _ := filepath.EvalSymlinks(root)
		gotResolved, _ := filepath.EvalSymlinks(got)
		if gotResolved != wantResolved {
			t.Errorf("got %q, want the project root %q", gotResolved, wantResolved)
		}
	})

	t.Run("nothing found outside a project", func(t *testing.T) {
		t.Chdir(t.TempDir())
		if got := resolveStartupPath("single", ""); got != "" {
			t.Errorf("got %q, want nothing outside a project", got)
		}
	})
}

func TestChangeModeFlagExists(t *testing.T) {
	found := false
	for _, f := range appFlags() {
		for _, name := range f.Names() {
			if name == "change-mode" {
				found = true
			}
		}
	}
	if !found {
		t.Error("--change-mode is missing")
	}
}
