package scanner

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// makeProjectAt builds a minimal OpenSpec project.
//
// isValidOpenSpecDir requires a config.yaml or project.md marker beside specs/
// or changes/, so a bare specs tree is not detected at all.
func makeProjectAt(t *testing.T, root, name string) string {
	t.Helper()
	dir := filepath.Join(root, name)
	spec := filepath.Join(dir, "openspec", "specs", "thing")
	if err := os.MkdirAll(spec, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "openspec", "config.yaml"),
		[]byte("schema: spec-driven\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(spec, "spec.md"), []byte("# thing\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func configFor(include, exclude []string, followSymlinks bool) *Config {
	c := &Config{FollowSymlinks: followSymlinks}
	c.ScanDirs.Include = include
	c.ScanDirs.Exclude = exclude
	return c
}

// collectWalk runs Walk and returns the discovered project paths, sorted.
func collectWalk(t *testing.T, config *Config) []string {
	t.Helper()
	results := make(chan string, 64)
	var found []string
	done := make(chan struct{})
	go func() {
		for p := range results {
			found = append(found, p)
		}
		close(done)
	}()
	if err := Walk(context.Background(), config, results, true); err != nil {
		t.Fatalf("Walk: %v", err)
	}
	<-done
	sort.Strings(found)
	return found
}

func TestWalkFindsProjectsAndIgnoresOtherDirectories(t *testing.T) {
	root := t.TempDir()
	alpha := makeProjectAt(t, root, "alpha")
	beta := makeProjectAt(t, root, "beta")

	// Not a project: no openspec/ at all.
	if err := os.MkdirAll(filepath.Join(root, "plain", "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	// Not a project: an openspec/ without the required marker file.
	if err := os.MkdirAll(filepath.Join(root, "bare", "openspec", "specs"), 0o755); err != nil {
		t.Fatal(err)
	}

	got := collectWalk(t, configFor([]string{root}, nil, false))
	want := []string{alpha, beta}
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("found %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("found %v, want %v", got, want)
			break
		}
	}
}

func TestWalkExcludeMatchesFullPathOrBasename(t *testing.T) {
	root := t.TempDir()
	keep := makeProjectAt(t, root, "keep")
	dropByPath := makeProjectAt(t, root, "drop-by-path")
	makeProjectAt(t, root, "vendor")

	t.Run("leading slash compares the whole path", func(t *testing.T) {
		got := collectWalk(t, configFor([]string{root}, []string{dropByPath}, false))
		for _, p := range got {
			if p == dropByPath {
				t.Errorf("%s should have been excluded by full path", p)
			}
		}
		if len(got) != 2 {
			t.Errorf("found %v, want 2 projects", got)
		}
	})

	t.Run("no leading slash compares the basename only", func(t *testing.T) {
		got := collectWalk(t, configFor([]string{root}, []string{"vendor"}, false))
		for _, p := range got {
			if filepath.Base(p) == "vendor" {
				t.Errorf("%s should have been excluded by basename", p)
			}
		}
		if len(got) != 2 {
			t.Errorf("found %v, want 2 projects", got)
		}
	})

	t.Run("a basename entry does not match a different path", func(t *testing.T) {
		got := collectWalk(t, configFor([]string{root}, []string{"keep"}, false))
		for _, p := range got {
			if p == keep {
				t.Errorf("%s should have been excluded", p)
			}
		}
	})
}

func TestWalkFollowSymlinks(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	target := makeProjectAt(t, outside, "linked")

	link := filepath.Join(root, "link-to-project")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	t.Run("off", func(t *testing.T) {
		if got := collectWalk(t, configFor([]string{root}, nil, false)); len(got) != 0 {
			t.Errorf("found %v, want nothing with followsymlinks off", got)
		}
	})

	t.Run("on", func(t *testing.T) {
		if got := collectWalk(t, configFor([]string{root}, nil, true)); len(got) != 1 {
			t.Errorf("found %v, want the linked project with followsymlinks on", got)
		}
	})
}

func TestWalkGlobInclude(t *testing.T) {
	root := t.TempDir()
	makeProjectAt(t, root, "capp")
	makeProjectAt(t, root, "cbpp")
	makeProjectAt(t, root, "other")

	got := collectWalk(t, configFor([]string{filepath.Join(root, "c*")}, nil, false))
	if len(got) != 2 {
		t.Errorf("found %v, want the two projects whose names begin with c", got)
	}
	for _, p := range got {
		if filepath.Base(p) == "other" {
			t.Error("the glob should not have matched 'other'")
		}
	}
}

func TestSkipComparisonModes(t *testing.T) {
	tests := []struct {
		needle   string
		haystack []string
		want     bool
	}{
		{"/a/b/vendor", []string{"/a/b/vendor"}, true}, // full path
		{"/a/b/vendor", []string{"/a/b/other"}, false}, // full path, no match
		{"/a/b/vendor", []string{"vendor"}, true},      // basename
		{"/a/b/vendor", []string{"b"}, false},          // basename must be the last element
		{"/a/b/vendor", nil, false},
	}
	for _, tt := range tests {
		if got := skip(tt.needle, tt.haystack); got != tt.want {
			t.Errorf("skip(%q, %v) = %v, want %v", tt.needle, tt.haystack, got, tt.want)
		}
	}
}
