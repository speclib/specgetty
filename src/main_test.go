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
