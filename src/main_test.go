package main

import (
	"os"
	"path/filepath"
	"testing"
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
