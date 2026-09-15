package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func testConfig(include, exclude []string) *Config {
	c := &Config{}
	c.ScanDirs.Include = include
	c.ScanDirs.Exclude = exclude
	return c
}

func TestCacheRoundTrip(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "sub", "projects.yaml")
	cfg := testConfig([]string{"/a", "/b"}, []string{"junk"})

	// Real directories, because LoadCache drops paths that do not exist.
	p1 := filepath.Join(dir, "p1")
	p2 := filepath.Join(dir, "p2")
	for _, p := range []string{p1, p2} {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	if err := SaveCache(file, cfg, []string{p1, p2}, time.Now()); err != nil {
		t.Fatalf("SaveCache: %v", err)
	}

	got, ok := LoadCache(file, cfg)
	if !ok {
		t.Fatal("LoadCache reported the cache unusable")
	}
	if len(got) != 2 || got[0] != p1 || got[1] != p2 {
		t.Errorf("got %v, want [%s %s]", got, p1, p2)
	}
}

func TestCacheMissingFile(t *testing.T) {
	cfg := testConfig(nil, nil)
	if _, ok := LoadCache(filepath.Join(t.TempDir(), "absent.yaml"), cfg); ok {
		t.Error("a missing cache should report ok=false, not succeed")
	}
}

func TestCacheMalformedFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "projects.yaml")
	if err := os.WriteFile(file, []byte("this: is: not: valid: yaml:\n\t- broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, ok := LoadCache(file, testConfig(nil, nil)); ok {
		t.Error("a malformed cache should report ok=false rather than failing the app")
	}
}

func TestCacheVersionMismatch(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "projects.yaml")
	if err := os.WriteFile(file, []byte("version: 99\npaths: [/tmp]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, ok := LoadCache(file, testConfig(nil, nil)); ok {
		t.Error("a cache from another format version should be treated as absent")
	}
}

func TestCacheInvalidatedByChangedScanDirs(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "projects.yaml")
	p := filepath.Join(dir, "p")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}

	built := testConfig([]string{"/a"}, []string{"junk"})
	if err := SaveCache(file, built, []string{p}, time.Now()); err != nil {
		t.Fatal(err)
	}

	if _, ok := LoadCache(file, built); !ok {
		t.Fatal("the cache should load against the config it was built from")
	}

	t.Run("include changed", func(t *testing.T) {
		if _, ok := LoadCache(file, testConfig([]string{"/a", "/c"}, []string{"junk"})); ok {
			t.Error("the cached paths answer a different question and should be discarded")
		}
	})

	t.Run("exclude changed", func(t *testing.T) {
		if _, ok := LoadCache(file, testConfig([]string{"/a"}, nil)); ok {
			t.Error("the cached paths answer a different question and should be discarded")
		}
	})
}

func TestCacheDropsVanishedProjects(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "projects.yaml")
	cfg := testConfig([]string{"/a"}, nil)

	alive := filepath.Join(dir, "alive")
	if err := os.MkdirAll(alive, 0o755); err != nil {
		t.Fatal(err)
	}
	gone := filepath.Join(dir, "gone")

	if err := SaveCache(file, cfg, []string{alive, gone}, time.Now()); err != nil {
		t.Fatal(err)
	}

	got, ok := LoadCache(file, cfg)
	if !ok {
		t.Fatal("LoadCache reported the cache unusable")
	}
	if len(got) != 1 || got[0] != alive {
		t.Errorf("got %v, want only the surviving project %s", got, alive)
	}
}

func TestCacheIgnoresAFileWhereAProjectShouldBe(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "projects.yaml")
	cfg := testConfig(nil, nil)

	notADir := filepath.Join(dir, "regular-file")
	if err := os.WriteFile(notADir, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := SaveCache(file, cfg, []string{notADir}, time.Now()); err != nil {
		t.Fatal(err)
	}

	got, _ := LoadCache(file, cfg)
	if len(got) != 0 {
		t.Errorf("got %v, want nothing: a file is not a project directory", got)
	}
}

func TestCachePathFollowsXDG(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "/tmp/specgetty-cache-test")
	got, err := CachePath()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join("/tmp/specgetty-cache-test", "specgetty", "projects.yaml")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
