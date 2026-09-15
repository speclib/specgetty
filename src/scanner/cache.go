package scanner

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// cacheVersion guards the on-disk format. A file written by a different version
// is treated as absent rather than migrated.
const cacheVersion = 1

// Cache remembers which project paths a walk discovered.
//
// It holds paths and nothing else. Discovering them is the expensive part of a
// scan, measured at 0.7s warm and about 3.9s cold for 18 projects, and it only
// changes when a project is created or deleted. The per-project statistics are
// cheap by comparison, roughly 100ms to read every artifact of those same 18
// projects, and they change every time a file is saved. So statistics are
// always read from disk and never served from here.
type Cache struct {
	Version   int       `yaml:"version"`
	ScannedAt time.Time `yaml:"scanned_at"`

	// ScanDirs records the question these paths answer. If the configured
	// directories change, the cached list describes a different search and must
	// not be reused.
	ScanDirs struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"scandirs"`

	Paths []string `yaml:"paths"`
}

// CachePath returns the cache file location, following the XDG Base Directory
// Specification through os.UserCacheDir.
func CachePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "specgetty", "projects.yaml"), nil
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// LoadCache reads the cached project paths.
//
// It reports ok=false whenever the cache cannot be trusted: missing, unreadable,
// malformed, written by another version, or built from different scan
// directories. Every one of those is an ordinary outcome, not an error worth
// failing the application over, so the caller simply walks instead.
//
// Paths that no longer exist on disk are dropped. That costs one stat per entry
// and saves a whole walk to notice a deleted project.
func LoadCache(path string, config *Config) ([]string, bool) {
	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, false
	}

	var c Cache
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, false
	}
	if c.Version != cacheVersion {
		return nil, false
	}
	if !sameStrings(c.ScanDirs.Include, config.ScanDirs.Include) ||
		!sameStrings(c.ScanDirs.Exclude, config.ScanDirs.Exclude) {
		return nil, false
	}

	paths := make([]string, 0, len(c.Paths))
	for _, p := range c.Paths {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			paths = append(paths, p)
		}
	}
	return paths, true
}

// SaveCache writes the discovered paths, recording the scan directories they
// came from. A failure to write is not worth failing the application over: the
// next run simply walks again.
func SaveCache(path string, config *Config, paths []string, now time.Time) error {
	c := Cache{
		Version:   cacheVersion,
		ScannedAt: now,
		Paths:     paths,
	}
	c.ScanDirs.Include = config.ScanDirs.Include
	c.ScanDirs.Exclude = config.ScanDirs.Exclude

	b, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
