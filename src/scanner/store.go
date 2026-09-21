package scanner

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// OpenSpec 1.10 lets a repo keep no specs and no changes of its own and point
// at a store instead. The names below are OpenSpec's, and resolution here
// mirrors its own rules so that specgetty and the CLI never disagree about
// which directory a project's content lives in.
const (
	storeMetadataDir  = ".openspec-store"
	storeMetadataFile = "store.yaml"
	storesDirName     = "stores"
	storeRegistryFile = "registry.yaml"
)

// Problem codes. They are what the config tab reports, and what the tests
// assert on, so they are named rather than being free text.
const (
	ProblemMalformedPointer = "malformed_pointer"
	ProblemNoRegistry       = "no_registry"
	ProblemUnknownStore     = "unknown_store"
	ProblemMissingPath      = "missing_path"
	ProblemStoreIdentity    = "store_identity"
)

// StoreProblem is a store declaration that could not be followed.
//
// Every field is here so the report can name both halves of the failure: the
// id that was asked for and the file that asked for it. Without the file, a
// user with several repos pointing at one store cannot tell which one is wrong.
type StoreProblem struct {
	Code   string
	ID     string
	File   string
	Detail string
}

func (p *StoreProblem) Error() string {
	if p == nil {
		return ""
	}
	return p.Detail
}

// StoreBackend is a registry entry's backend block.
type StoreBackend struct {
	Type      string `yaml:"type"`
	LocalPath string `yaml:"local_path"`
	Remote    string `yaml:"remote"`
	Branch    string `yaml:"branch"`
}

type registryEntry struct {
	Backend StoreBackend `yaml:"backend"`
}

type registryFile struct {
	Version int                      `yaml:"version"`
	Stores  map[string]registryEntry `yaml:"stores"`
}

// StoreMetadata is a store's identity file, which is what makes a directory a
// store rather than an ordinary OpenSpec root.
type StoreMetadata struct {
	Version int    `yaml:"version"`
	ID      string `yaml:"id"`
	Remote  string `yaml:"remote"`
}

// StoreGit is what can be learned about a store's working copy without
// touching the network.
//
// Ahead and Behind compare HEAD against the upstream ref as the local
// repository currently holds it. Nothing here fetches, so the comparison is
// against whatever the last fetch left behind, and the display says so.
type StoreGit struct {
	IsRepo        bool
	OriginURL     string
	Dirty         bool
	DirtyKnown    bool
	Ahead         int
	Behind        int
	TrackingKnown bool
}

// StoreInfo is everything the config tab reports about a store.
type StoreInfo struct {
	ID        string
	Root      string
	Origin    string // the repo the resolution started from, empty when there was none
	Remote    string // recorded by the registry
	Branch    string // recorded by the registry
	Canonical string // recorded by the store's own metadata
	Git       *StoreGit
}

// Resolution is the answer to "where does this project's content live".
//
// Root and Origin differ only when a store declaration was followed. Keeping
// both is the whole point: Root is where every read and every write goes,
// Origin is where the user stands and what the header names.
type Resolution struct {
	Root    string
	Origin  string
	StoreID string
	Store   *StoreInfo
	Problem *StoreProblem
}

// FromStore reports whether the content came from somewhere other than the
// directory the user is standing in.
func (r Resolution) FromStore() bool {
	return r.StoreID != "" || r.Problem != nil
}

// DataDir returns the directory OpenSpec keeps machine-level data in.
//
// XDG_DATA_HOME wins wherever it is set, which is OpenSpec's rule on every
// platform rather than a Linux-only convention.
func DataDir() string {
	if x := os.Getenv("XDG_DATA_HOME"); x != "" {
		return filepath.Join(x, "openspec")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "share", "openspec")
}

// RegistryPath returns the store registry's location.
func RegistryPath() string {
	d := DataDir()
	if d == "" {
		return ""
	}
	return filepath.Join(d, storesDirName, storeRegistryFile)
}

// LoadRegistry reads the store registry.
//
// present is false when there is no registry file at all, which is a different
// report from a registry that exists and does not hold the wanted id.
func LoadRegistry(path string) (stores map[string]StoreBackend, present bool, err error) {
	if path == "" {
		return nil, false, nil
	}
	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	var rf registryFile
	if err := yaml.Unmarshal(b, &rf); err != nil {
		return nil, true, err
	}
	stores = make(map[string]StoreBackend, len(rf.Stores))
	for id, e := range rf.Stores {
		stores[id] = e.Backend
	}
	return stores, true, nil
}

// PointerState says what a configuration file has to say about a store.
type PointerState int

const (
	PointerAbsent     PointerState = iota // no store key
	PointerDeclared                       // store: <id>
	PointerNotAString                     // store: is a list, a mapping or empty
	PointerUnreadable                     // the file is not YAML
)

// StorePointer is a configuration file's store declaration.
type StorePointer struct {
	State PointerState
	ID    string
	File  string // the file it was read from, empty when there is none
}

// ConfigFilePath returns the YAML configuration of an OpenSpec directory.
//
// OpenSpec accepts both spellings and prefers .yaml. specgetty knew only the
// first, which mattered the moment the store declaration started living in
// that file.
func ConfigFilePath(dir string) string {
	yamlPath := filepath.Join(dir, "openspec", "config.yaml")
	if _, err := os.Stat(yamlPath); err == nil {
		return yamlPath
	}
	ymlPath := filepath.Join(dir, "openspec", "config.yml")
	if _, err := os.Stat(ymlPath); err == nil {
		return ymlPath
	}
	return ""
}

// ReadStorePointer reads a directory's store declaration.
func ReadStorePointer(dir string) StorePointer {
	path := ConfigFilePath(dir)
	if path == "" {
		return StorePointer{State: PointerAbsent}
	}
	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return StorePointer{State: PointerUnreadable, File: path}
	}

	// A document that is not a mapping carries no declaration. It is an
	// imperfect config, not a malformed pointer, which is the distinction
	// OpenSpec draws too.
	var doc interface{}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return StorePointer{State: PointerUnreadable, File: path}
	}
	if _, ok := doc.(map[interface{}]interface{}); !ok {
		return StorePointer{State: PointerAbsent, File: path}
	}

	var probe struct {
		Store interface{} `yaml:"store"`
	}
	if err := yaml.Unmarshal(b, &probe); err != nil {
		return StorePointer{State: PointerUnreadable, File: path}
	}
	switch v := probe.Store.(type) {
	case nil:
		return StorePointer{State: PointerAbsent, File: path}
	case string:
		if strings.TrimSpace(v) == "" {
			return StorePointer{State: PointerNotAString, File: path}
		}
		return StorePointer{State: PointerDeclared, ID: strings.TrimSpace(v), File: path}
	default:
		return StorePointer{State: PointerNotAString, File: path}
	}
}

// ReadStoreMetadata reads a directory's store identity file.
//
// ok is false both when the directory is not a store and when its metadata
// cannot be read or names no id. A root whose metadata is unreadable is still
// a root, so the caller shows it as an ordinary project rather than hiding it.
func ReadStoreMetadata(root string) (StoreMetadata, bool) {
	path := filepath.Join(root, storeMetadataDir, storeMetadataFile)
	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return StoreMetadata{}, false
	}
	var md StoreMetadata
	if err := yaml.Unmarshal(b, &md); err != nil {
		return StoreMetadata{}, false
	}
	if strings.TrimSpace(md.ID) == "" {
		return StoreMetadata{}, false
	}
	md.ID = strings.TrimSpace(md.ID)
	return md, true
}

// hasPlanningShape reports whether an OpenSpec directory holds content of its
// own. Content present locally outranks any pointer beside it.
func hasPlanningShape(dir string) bool {
	for _, sub := range []string{"specs", "changes"} {
		if info, err := os.Stat(filepath.Join(dir, "openspec", sub)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

// hasConfigFile reports whether an OpenSpec directory holds a configuration of
// any of the three spellings specgetty recognises.
func hasConfigFile(dir string) bool {
	if ConfigFilePath(dir) != "" {
		return true
	}
	_, err := os.Stat(filepath.Join(dir, "openspec", "project.md"))
	return err == nil
}

// qualifies reports whether a directory's openspec/ is one the walk may stop
// at. A directory named openspec with nothing in it is not: without this, a
// store kept at ~/openspec would make the home directory capture every project
// beneath it.
func qualifies(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, "openspec"))
	if err != nil || !info.IsDir() {
		return false
	}
	return hasPlanningShape(dir) || hasConfigFile(dir)
}

// FindRoot walks up from startDir to the nearest qualifying OpenSpec
// directory, returning "" when there is none.
func FindRoot(startDir string) string {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return ""
	}
	for {
		if qualifies(dir) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// storeInfoAt describes a directory when it is a store, and returns nil when
// it is not.
//
// The registry is only read once the identity file has been found, so an
// ordinary project costs one failed stat rather than a registry parse.
func storeInfoAt(root string) *StoreInfo {
	if _, ok := ReadStoreMetadata(root); !ok {
		return nil
	}
	stores, _, err := LoadRegistry(RegistryPath())
	if err != nil {
		return nil
	}
	return StoreInfoWith(stores, root)
}

// StoreInfoWith describes a directory when the given registry makes it a store.
//
// An identity file is not enough on its own. A clone, or a directory the
// registry has since moved on from, keeps its copy of that file, and trusting
// it would let a leftover wear the registered store's name. The declaration
// path already demanded that the registry key and the metadata id agree; this
// applies the same test to a directory met head on, so the two cannot disagree
// about what a store is.
//
// The registry is passed in rather than read, because discovery asks this of
// every candidate directory and reads the registry once for the whole walk.
func StoreInfoWith(stores map[string]StoreBackend, root string) *StoreInfo {
	md, ok := ReadStoreMetadata(root)
	if !ok {
		return nil
	}
	backend, found := stores[md.ID]
	if !found || !samePath(backend.LocalPath, root) {
		return nil
	}
	return &StoreInfo{
		ID: md.ID, Root: root, Canonical: md.Remote,
		Remote: backend.Remote, Branch: backend.Branch,
	}
}

// samePath compares two paths as the filesystem sees them.
//
// The registry records a canonicalised path while a walk reports whatever it
// descended through, so a store reached behind a symlink would otherwise fail
// to match the entry that names it.
func samePath(a, b string) bool {
	if a == b {
		return true
	}
	if a == "" || b == "" {
		return false
	}
	ra, err := filepath.EvalSymlinks(a)
	if err != nil {
		ra = filepath.Clean(a)
	}
	rb, err := filepath.EvalSymlinks(b)
	if err != nil {
		rb = filepath.Clean(b)
	}
	return ra == rb
}

// ResolveRoot answers where the content of the project at startDir lives.
//
// The rules are OpenSpec's: walk up to the nearest qualifying openspec/
// directory; content present there outranks a pointer beside it; otherwise a
// declaration is followed through the registry. ok is false when there is no
// OpenSpec directory to be found at all.
func ResolveRoot(startDir string) (Resolution, bool) {
	dir := FindRoot(startDir)
	if dir == "" {
		return Resolution{}, false
	}

	// Content of its own settles it. A store declaration sitting beside real
	// specs is dead, which is what OpenSpec warns about rather than follows.
	if hasPlanningShape(dir) {
		return localResolution(dir), true
	}

	pointer := ReadStorePointer(dir)
	switch pointer.State {
	case PointerUnreadable:
		return Resolution{
			Root: dir, Origin: dir,
			Problem: &StoreProblem{
				Code: ProblemMalformedPointer, File: pointer.File,
				Detail: "the configuration file could not be read as YAML",
			},
		}, true
	case PointerNotAString:
		return Resolution{
			Root: dir, Origin: dir,
			Problem: &StoreProblem{
				Code: ProblemMalformedPointer, File: pointer.File,
				Detail: "the store key must be a single store id",
			},
		}, true
	case PointerAbsent:
		// A freshly initialised project: no content yet and nothing to follow.
		return localResolution(dir), true
	}

	return followPointer(dir, pointer), true
}

// localResolution describes a directory that is its own root.
func localResolution(dir string) Resolution {
	r := Resolution{Root: dir, Origin: dir}
	if si := storeInfoAt(dir); si != nil {
		r.StoreID = si.ID
		r.Store = si
	}
	return r
}

// followPointer resolves a declared store id through the registry.
func followPointer(origin string, pointer StorePointer) Resolution {
	fail := func(code, detail string) Resolution {
		return Resolution{
			Root: origin, Origin: origin,
			Problem: &StoreProblem{Code: code, ID: pointer.ID, File: pointer.File, Detail: detail},
		}
	}

	stores, present, err := LoadRegistry(RegistryPath())
	if err != nil {
		return fail(ProblemNoRegistry, "the store registry could not be read: "+err.Error())
	}
	if !present || len(stores) == 0 {
		return fail(ProblemNoRegistry, "no stores are registered on this machine")
	}
	backend, found := stores[pointer.ID]
	if !found {
		return fail(ProblemUnknownStore, "it is not among the registered stores: "+strings.Join(sortedKeys(stores), ", "))
	}
	root := backend.LocalPath
	if root == "" {
		return fail(ProblemMissingPath, "its registry entry records no local path")
	}
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		return fail(ProblemMissingPath, "its registered path no longer exists: "+root)
	}

	md, ok := ReadStoreMetadata(root)
	if !ok {
		return fail(ProblemStoreIdentity, "its registered path holds no readable store identity: "+root)
	}
	if md.ID != pointer.ID {
		return fail(ProblemStoreIdentity, "its metadata declares the id "+md.ID+", which is not the id it is registered under")
	}

	si := &StoreInfo{ID: md.ID, Root: root, Origin: origin, Canonical: md.Remote,
		Remote: backend.Remote, Branch: backend.Branch}
	return Resolution{Root: root, Origin: origin, StoreID: md.ID, Store: si}
}

func sortedKeys(m map[string]StoreBackend) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// ReadStoreGit reports a store working copy's local state.
//
// Every command here reads local files and refs. None of them fetches, so the
// ahead and behind counts are against the upstream ref as it currently stands
// on disk, which is what the display says.
func ReadStoreGit(root string) StoreGit {
	var g StoreGit
	if _, err := exec.LookPath("git"); err != nil {
		return g
	}
	if out, err := gitAt(root, "rev-parse", "--is-inside-work-tree"); err != nil || out != "true" {
		return g
	}
	g.IsRepo = true

	if out, err := gitAt(root, "remote", "get-url", "origin"); err == nil {
		g.OriginURL = out
	}

	// Scoped to the store's own subtree, because several stores can share one
	// working copy and a sibling's edits are not this store's state.
	if out, err := gitAt(root, "status", "--porcelain", "--", "."); err == nil {
		g.DirtyKnown = true
		g.Dirty = out != ""
	}

	if out, err := gitAt(root, "rev-list", "--left-right", "--count", "HEAD...@{upstream}"); err == nil {
		fields := strings.Fields(out)
		if len(fields) == 2 {
			ahead, errA := strconv.Atoi(fields[0])
			behind, errB := strconv.Atoi(fields[1])
			if errA == nil && errB == nil {
				g.TrackingKnown = true
				g.Ahead = ahead
				g.Behind = behind
			}
		}
	}
	return g
}

func gitAt(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	// A repository asking for credentials would block the UI. Nothing here
	// needs the network, so refuse to be prompted.
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_OPTIONAL_LOCKS=0")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
