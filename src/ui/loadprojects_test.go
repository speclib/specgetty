package ui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mipmip/specgetty/src/scanner"
)

// makeProject creates a minimal OpenSpec project on disk.
func makeProject(t *testing.T, root, name string, specs ...string) string {
	t.Helper()
	dir := filepath.Join(root, name)

	// isValidOpenSpecDir requires a config.yaml or project.md marker alongside
	// specs/ or changes/, so a bare specs tree is not detected as a project.
	if err := os.MkdirAll(filepath.Join(dir, "openspec"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "openspec", "config.yaml"),
		[]byte("schema: spec-driven\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, s := range specs {
		specDir := filepath.Join(dir, "openspec", "specs", s)
		if err := os.MkdirAll(specDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(specDir, "spec.md"),
			[]byte("# "+s+"\n\ncontent for "+s+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestLoadProjectsWalksThenServesFromCache(t *testing.T) {
	root := t.TempDir()
	cache := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cache)

	makeProject(t, root, "alpha", "one", "two")
	makeProject(t, root, "beta", "three")

	config := &scanner.Config{}
	config.ScanDirs.Include = []string{root}

	// First load has no cache, so it walks and writes one.
	msg := loadProjects(config, true, false)()
	loaded, ok := msg.(pickerLoadedMsg)
	if !ok {
		t.Fatalf("got %T, want pickerLoadedMsg", msg)
	}
	if loaded.err != nil {
		t.Fatalf("load failed: %v", loaded.err)
	}
	if len(loaded.rows) != 2 {
		t.Fatalf("found %d projects, want 2: %v", len(loaded.rows), projectNames(wrap(loaded.rows)))
	}

	cacheFile := filepath.Join(cache, "specgetty", "projects.yaml")
	if _, err := os.Stat(cacheFile); err != nil {
		t.Fatalf("the walk should have written a cache: %v", err)
	}

	// Statistics come from disk, not the cache.
	byName := map[string]projectRow{}
	for _, r := range loaded.rows {
		byName[r.display] = r
	}
	if got := byName["alpha"].info.SpecCount; got != 2 {
		t.Errorf("alpha SpecCount = %d, want 2", got)
	}

	// A project created after the walk stays hidden until a refresh, which is
	// the deliberate trade for not walking on every open.
	makeProject(t, root, "gamma", "four")

	cached := loadProjects(config, true, false)().(pickerLoadedMsg)
	if len(cached.rows) != 2 {
		t.Errorf("cached load found %d projects, want the 2 from the cache", len(cached.rows))
	}

	refreshed := loadProjects(config, true, true)().(pickerLoadedMsg)
	if len(refreshed.rows) != 3 {
		t.Errorf("refresh found %d projects, want 3 including the new one", len(refreshed.rows))
	}
}

func TestLoadProjectsPicksUpFreshStatisticsFromACachedPath(t *testing.T) {
	root := t.TempDir()
	cache := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cache)

	dir := makeProject(t, root, "alpha", "one")

	config := &scanner.Config{}
	config.ScanDirs.Include = []string{root}

	first := loadProjects(config, true, false)().(pickerLoadedMsg)
	if first.rows[0].info.SpecCount != 1 {
		t.Fatalf("SpecCount = %d, want 1", first.rows[0].info.SpecCount)
	}

	// Add a spec without rewalking. The cache holds paths only, so the count
	// must still be current on the next load.
	specDir := filepath.Join(dir, "openspec", "specs", "two")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte("# two\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	second := loadProjects(config, true, false)().(pickerLoadedMsg)
	if got := second.rows[0].info.SpecCount; got != 2 {
		t.Errorf("SpecCount = %d, want 2: statistics must never come from the cache", got)
	}
}
