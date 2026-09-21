package scanner

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	taskDoneRe   = regexp.MustCompile(`(?m)^- \[x\] `)
	taskUndoneRe = regexp.MustCompile(`(?m)^- \[ \] `)
)

// ParseTaskStats counts done and total checkbox tasks in markdown content.
func ParseTaskStats(content string) (total, done int) {
	done = len(taskDoneRe.FindAllString(content, -1))
	undone := len(taskUndoneRe.FindAllString(content, -1))
	total = done + undone
	return
}

type FileEntry struct {
	Path  string
	IsDir bool
}

type ChangeInfo struct {
	Name string
	// DirName is the directory this change was read from. For an archived
	// change that still carries its `YYYY-MM-DD-` prefix, while Name has it
	// stripped for display. Anything that needs to find the change on disk
	// again must use DirName.
	DirName          string
	ArtifactFiles    []string          // sorted .md filenames
	ArtifactContents map[string]string // filename → content
	TasksTotal       int
	TasksDone        int
	SpecNames        []string
	SpecContents     map[string]string
	ArchiveDate      time.Time // zero for active changes, mtime for archived
}

type ProjectInfo struct {
	SpecCount       int
	SpecNames       []string
	SpecContents    map[string]string
	ActiveChanges   []string
	Changes         []ChangeInfo
	ArchivedChanges []ChangeInfo
	ConfigFile      string // "project.md", "config.yaml", "config.yml", or ""
	ConfigContent   string
	TasksTotal      int // aggregate across all active changes
	TasksDone       int

	// Where the content was read from, and where the reading started. They
	// differ only when a store declaration was followed, and every filesystem
	// operation belongs to Root while the header belongs to Origin.
	Root   string
	Origin string

	// Set when Root is a store. StoreID is what the store calls itself, which
	// is not necessarily what its directory is called.
	StoreID string
	Store   *StoreInfo

	// The configuration of the repo the resolution started from, which a store
	// keeps none of. Empty unless Origin and Root differ.
	OriginConfigFile    string
	OriginConfigContent string

	// Set when a store declaration could not be followed. The project is not
	// empty in that case; it is unreadable, and saying so is the difference.
	StoreProblem *StoreProblem
}

// FromStore reports whether this project reads its content from a store.
func (i ProjectInfo) FromStore() bool {
	return i.StoreID != "" || i.StoreProblem != nil
}

// ResolvedElsewhere reports whether the content came from a directory other
// than the one the user is standing in.
func (i ProjectInfo) ResolvedElsewhere() bool {
	return i.Origin != "" && i.Root != "" && i.Origin != i.Root
}

type ProjectStatus struct {
	Files    []FileEntry
	Info     ProjectInfo
	ScanTime time.Duration
}

type ProjectMap map[string]ProjectStatus

type Config struct {
	ScanDirs struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"scandirs"`
	FollowSymlinks bool   `yaml:"followsymlinks"`
	EditCommand    string `yaml:"edit_command"`
	ChangeFields   string `yaml:"change_fields"`
	ChangeMode     string `yaml:"change_mode"`
}

func DumpConfig(config *Config) error {
	b, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func ParseConfigFile(filename, defaultConfig string) (*Config, error) {
	b, err := ioutil.ReadFile(filepath.Clean(filename))
	switch {
	case err == nil:
	case os.IsNotExist(err):
		b = ([]byte)(defaultConfig)
	default:
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// ListOpenSpecContents recursively reads the openspec/ directory under dir
// and returns relative paths with a dir/file indicator.
func ListOpenSpecContents(dir string) ([]FileEntry, error) {
	openspecDir := filepath.Join(dir, "openspec")
	var entries []FileEntry

	err := filepath.WalkDir(openspecDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, path)
		}
		rel, err := filepath.Rel(openspecDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		entries = append(entries, FileEntry{
			Path:  rel,
			IsDir: d.IsDir(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// parseChangeDir reads a change directory and returns a ChangeInfo.
func parseChangeDir(dir string, name string) ChangeInfo {
	ci := ChangeInfo{
		Name:             name,
		DirName:          filepath.Base(dir),
		ArtifactContents: make(map[string]string),
		SpecContents:     make(map[string]string),
	}

	// Discover .md files
	if files, err := os.ReadDir(dir); err == nil {
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
				continue
			}
			ci.ArtifactFiles = append(ci.ArtifactFiles, f.Name())
			if b, err := os.ReadFile(filepath.Join(dir, f.Name())); err == nil {
				ci.ArtifactContents[f.Name()] = string(b)
			}
		}
	}
	sort.Strings(ci.ArtifactFiles)

	// Parse task stats from tasks.md if present
	if tasksContent, ok := ci.ArtifactContents["tasks.md"]; ok {
		ci.TasksTotal, ci.TasksDone = ParseTaskStats(tasksContent)
	}

	// Read specs within the change
	specsDir := filepath.Join(dir, "specs")
	if specEntries, err := os.ReadDir(specsDir); err == nil {
		for _, se := range specEntries {
			if se.IsDir() {
				ci.SpecNames = append(ci.SpecNames, se.Name())
				specFile := filepath.Join(specsDir, se.Name(), "spec.md")
				if b, err := os.ReadFile(specFile); err == nil {
					ci.SpecContents[se.Name()] = string(b)
				}
			}
		}
	}
	sort.Strings(ci.SpecNames)

	return ci
}

// ParseProjectInfo reads the openspec/ directory structure to extract stats.
func ParseProjectInfo(dir string) ProjectInfo {
	info := ProjectInfo{}
	openspecDir := filepath.Join(dir, "openspec")

	// Specs: subdirectories in openspec/specs/, read spec.md contents
	specsDir := filepath.Join(openspecDir, "specs")
	info.SpecContents = make(map[string]string)
	if entries, err := os.ReadDir(specsDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				info.SpecCount++
				info.SpecNames = append(info.SpecNames, e.Name())
				specFile := filepath.Join(specsDir, e.Name(), "spec.md")
				if b, err := os.ReadFile(specFile); err == nil {
					info.SpecContents[e.Name()] = string(b)
				}
			}
		}
	}
	sort.Strings(info.SpecNames)

	// Active changes: subdirectories in openspec/changes/
	changesDir := filepath.Join(openspecDir, "changes")
	if entries, err := os.ReadDir(changesDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() || e.Name() == "archive" || e.Name() == "discarded" {
				continue
			}
			info.ActiveChanges = append(info.ActiveChanges, e.Name())
			ci := parseChangeDir(filepath.Join(changesDir, e.Name()), e.Name())
			info.Changes = append(info.Changes, ci)
		}
	}
	sort.Slice(info.Changes, func(i, j int) bool {
		return info.Changes[i].Name < info.Changes[j].Name
	})

	// Aggregate task stats across active changes
	for _, ci := range info.Changes {
		info.TasksTotal += ci.TasksTotal
		info.TasksDone += ci.TasksDone
	}

	// Archived changes: subdirectories in openspec/archive/
	archiveDir := filepath.Join(openspecDir, "changes", "archive")
	if entries, err := os.ReadDir(archiveDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				dirName := e.Name()
				displayName := dirName
				var archiveDate time.Time
				// Parse YYYY-MM-DD- prefix from directory name
				if len(dirName) >= 11 && dirName[4] == '-' && dirName[7] == '-' && dirName[10] == '-' {
					if t, err := time.Parse("2006-01-02", dirName[:10]); err == nil {
						archiveDate = t
						displayName = dirName[11:]
					}
				}
				ci := parseChangeDir(filepath.Join(archiveDir, dirName), displayName)
				ci.ArchiveDate = archiveDate
				info.ArchivedChanges = append(info.ArchivedChanges, ci)
			}
		}
	}

	// Config file: project.md takes priority, then the two spellings OpenSpec
	// accepts for the YAML configuration, .yaml before .yml.
	info.ConfigFile, info.ConfigContent = readProjectConfig(dir)

	// Identity. A directory is a store when it carries the metadata file, and
	// that is checked with one stat before the registry is opened, so an
	// ordinary project never pays for a registry parse.
	info.Root = dir
	info.Origin = dir
	if si := storeInfoAt(dir); si != nil {
		info.StoreID = si.ID
		info.Store = si
	}

	return info
}

// readProjectConfig returns the display name and content of a directory's
// OpenSpec configuration, or two empty strings when it has none.
func readProjectConfig(dir string) (name, content string) {
	openspecDir := filepath.Join(dir, "openspec")
	for _, candidate := range []string{"project.md", "config.yaml", "config.yml"} {
		if b, err := os.ReadFile(filepath.Join(openspecDir, candidate)); err == nil {
			return candidate, string(b)
		}
	}
	return "", ""
}

// ScanResolved reads the project reached from startDir, following a store
// declaration when there is one.
//
// It returns the key the result is filed under, which is where the reading
// STARTED, not where it ended. A project is identified by the directory a
// person stands in, so two repos sharing a store are two projects; the root
// their content came from travels inside the info as Root, and every read and
// every write goes there.
func ScanResolved(startDir string) (string, ProjectStatus, error) {
	return scanResolved(startDir, newRootCache(), true)
}

// rootCache memoises the content of a root within one scan.
//
// Two repos declaring the same store is the ordinary case now that a repo is
// what gets listed, and the store's `changes/archive/` can hold ninety
// directories. Reading it once per row rather than once per scan would make the
// scan cost scale with how many repos share a store, which is the wrong thing
// for it to scale with.
type rootCache map[string]cachedRoot

type cachedRoot struct {
	files []FileEntry
	info  ProjectInfo
}

func newRootCache() rootCache { return make(rootCache) }

// scanResolved reads the project reached from startDir.
//
// withGit is false for discovery. The git state costs several subprocesses per
// store and is only wanted for the project actually being opened.
func scanResolved(startDir string, cache rootCache, withGit bool) (string, ProjectStatus, error) {
	res, ok := ResolveRoot(startDir)
	if !ok {
		return "", ProjectStatus{}, nil
	}

	cached, hit := cache[res.Root]
	if !hit {
		files, err := ListOpenSpecContents(res.Root)
		if err != nil {
			// A store whose registered path went missing between resolution
			// and reading leaves nothing to list. The problem is the report,
			// not the listing error.
			if res.Problem == nil {
				return "", ProjectStatus{}, err
			}
			files = nil
		}
		cached = cachedRoot{files: files, info: ParseProjectInfo(res.Root)}
		cache[res.Root] = cached
	}

	// The content is shared between every repo reading this root; everything
	// below belongs to this one repo and is written onto a copy.
	info := cached.info
	info.Root = res.Root
	info.Origin = res.Origin
	info.StoreID = res.StoreID
	info.StoreProblem = res.Problem
	info.Store = nil

	if res.Store != nil {
		// Copied, not aliased: Origin and Git differ per repo, and two repos
		// sharing a store must not write over each other's.
		store := *res.Store
		store.Origin = res.Origin
		if withGit {
			g := ReadStoreGit(store.Root)
			store.Git = &g
		}
		info.Store = &store
	}

	if info.ResolvedElsewhere() {
		info.OriginConfigFile, info.OriginConfigContent = readProjectConfig(res.Origin)
	}

	return res.Origin, ProjectStatus{Files: cached.files, Info: info}, nil
}

// ScanPaths parses a known list of project paths without walking the
// filesystem to discover them. It is what the cache enables: discovery is the
// expensive half of a scan, parsing is the cheap half, so a cached list of
// paths still gets fully current statistics.
//
// Each path is resolved, so a repo that declares a store arrives carrying that
// store's specs, changes and file listing rather than the nothing it holds
// itself. A row reporting zero for a project whose specs are one directory away
// would be worse than no row.
func ScanPaths(paths []string) ProjectMap {
	results := make(ProjectMap, len(paths))
	cache := newRootCache()
	for _, d := range paths {
		_, st, err := scanResolved(d, cache, false)
		if err != nil || st.Info.Root == "" {
			continue
		}
		// Keyed by the path that was asked for, not by the root: two repos
		// sharing a store are two rows, and keying by root would collapse them
		// into one.
		results[d] = st
	}
	return results
}

// Scan finds all OpenSpec projects in directories specified by config
func Scan(config *Config, ignore_dir_errors bool) (ProjectMap, error) {
	ctx := context.Background()
	projects := make(chan string, 1000)

	type walkResult struct {
		err      error
		duration time.Duration
	}
	ch := make(chan walkResult)
	go func() {
		start := time.Now()
		err := Walk(ctx, config, projects, ignore_dir_errors)
		ch <- walkResult{
			err:      err,
			duration: time.Since(start),
		}
	}()

	results := make(ProjectMap)
	cache := newRootCache()
	totalScanDuration := time.Duration(0)
	for d := range projects {
		start := time.Now()

		// Resolved, so a repo that declares a store arrives carrying that
		// store's statistics and file listing. The cache means a store shared
		// by several repos is read once for the scan.
		_, st, err := scanResolved(d, cache, false)
		if err != nil {
			return nil, err
		}
		if st.Info.Root == "" {
			continue
		}

		duration := time.Since(start)
		log.Println(d, duration)

		totalScanDuration += duration
		st.ScanTime = duration
		results[d] = st
	}

	w := <-ch
	log.Println("walkDuration:", w.duration)
	log.Println("scanDuration:", totalScanDuration)
	return results, w.err
}
