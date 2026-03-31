package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSkip(t *testing.T) {
	tests := []struct {
		name     string
		needle   string
		haystack []string
		expected bool
	}{
		{
			name:     "full path match",
			needle:   "/home/user/node_modules",
			haystack: []string{"/home/user/node_modules"},
			expected: true,
		},
		{
			name:     "basename match",
			needle:   "/home/user/project/node_modules",
			haystack: []string{"node_modules"},
			expected: true,
		},
		{
			name:     "no match",
			needle:   "/home/user/project/src",
			haystack: []string{"vendor"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := skip(tt.needle, tt.haystack)
			if got != tt.expected {
				t.Errorf("skip(%q, %v) = %v, want %v", tt.needle, tt.haystack, got, tt.expected)
			}
		})
	}
}

func TestIsValidOpenSpecDir(t *testing.T) {
	t.Run("valid with config.yaml and specs/", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(""), 0644)
		os.Mkdir(filepath.Join(dir, "specs"), 0755)

		if !isValidOpenSpecDir(dir) {
			t.Error("expected valid, got invalid")
		}
	})

	t.Run("valid with project.md and archive/", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "project.md"), []byte(""), 0644)
		os.Mkdir(filepath.Join(dir, "archive"), 0755)

		if !isValidOpenSpecDir(dir) {
			t.Error("expected valid, got invalid")
		}
	})

	t.Run("invalid empty directory", func(t *testing.T) {
		dir := t.TempDir()

		if isValidOpenSpecDir(dir) {
			t.Error("expected invalid, got valid")
		}
	})

	t.Run("invalid with config.yaml but no specs/ or archive/", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(""), 0644)

		if isValidOpenSpecDir(dir) {
			t.Error("expected invalid, got valid")
		}
	})

	t.Run("invalid with specs/ but no config.yaml or project.md", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, "specs"), 0755)

		if isValidOpenSpecDir(dir) {
			t.Error("expected invalid, got valid")
		}
	})
}
