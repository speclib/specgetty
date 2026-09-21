package scanner

import (
	"context"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
	"golang.org/x/sync/errgroup"
)

func skip(needle string, haystack []string) bool {
	for _, f := range haystack {
		// A YAML dash with nothing after it parses to an empty string. It
		// excludes nothing, and it must not be indexed.
		if f == "" {
			continue
		}

		//FULL PATH COMPARISON
		if strings.HasPrefix(f, "/") {
			if f == needle {
				return true
			}

			//PARTIAL PATH COMPARISON
		} else {
			if f == path.Base(needle) {
				return true
			}
		}
	}
	return false
}

// walkone descends a single directory tree looking for OpenSpec projects
func walkone(ctx context.Context, dir string, config *Config, stores map[string]StoreBackend, results chan string) error {
	err := godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted:            true,
		ScratchBuffer:       make([]byte, godirwalk.MinimumScratchBufferSize),
		FollowSymbolicLinks: config.FollowSymlinks,
		ErrorCallback: func(path string, err error) godirwalk.ErrorAction {
			patherr, ok := err.(*os.PathError)
			if ok {
				switch patherr.Unwrap().Error() {
				case "no such file or directory":
					return godirwalk.SkipNode

				case "too many levels of symbolic links":
					return godirwalk.SkipNode
				}
			}
			log.Printf("ERROR: %s: %v", path, err)
			return godirwalk.Halt
		},
		Callback: func(path string, ent *godirwalk.Dirent) error {

			select {
			case <-ctx.Done():
				return filepath.SkipDir
			default:
			}

			if skip(path, config.ScanDirs.Exclude) {
				return godirwalk.SkipThis
			}
			if ent.IsSymlink() && !config.FollowSymlinks {
				return godirwalk.SkipThis
			}

			if ent.Name() != "openspec" {
				return nil
			}
			isDir, _ := ent.IsDirOrSymlinkToDir()
			if !isDir {
				return nil
			}

			if !isValidOpenSpecDir(path, stores) {
				return nil
			}

			results <- filepath.Dir(path)
			return godirwalk.SkipThis // don't descend further
		},
	})
	return err
}

// isValidOpenSpecDir checks whether an openspec/ directory belongs to a
// directory worth listing.
//
// A configuration is required and content is not. A repo declaring a `store:`
// keeps neither `specs/` nor `changes/`, and it is still a project: it is where
// a person works, and it is the only place that repo's own context and rules
// can be read from. The store it reads from is excluded instead, because a
// store is where content lives rather than where work happens, and every repo
// reading from one already stands for it in the list.
//
// stores is the registry, read once for the whole walk. A nil or empty map
// excludes nothing, which is what a machine with no stores should see.
func isValidOpenSpecDir(path string, stores map[string]StoreBackend) bool {
	hasConfig := false
	for _, name := range []string{"config.yaml", "config.yml", "project.md"} {
		if _, err := os.Stat(filepath.Join(path, name)); err == nil {
			hasConfig = true
			break
		}
	}
	if !hasConfig {
		return false
	}

	return StoreInfoWith(stores, filepath.Dir(path)) == nil
}

// Walk finds all OpenSpec projects in the directories specified in config
func Walk(ctx context.Context, config *Config, results chan string, ignore_dir_errors bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Every return path has to close the channel. A caller ranging over it
	// blocks forever otherwise, and the error paths below return early.
	defer close(results)

	// The registry decides which discovered directories are stores, and every
	// candidate asks the same question, so it is read once here rather than
	// once per candidate. A registry that cannot be read excludes nothing:
	// failing to list projects because of it would be worse than listing a
	// store alongside them.
	stores, _, err := LoadRegistry(RegistryPath())
	if err != nil {
		log.Printf("ERROR: store registry: %v", err)
		stores = nil
	}

	// Copied rather than aliased: the expansion below appends to this list, and
	// appending to config.ScanDirs.Include would reach into the caller's config.
	completeIncludeList := make([]string, 0, len(config.ScanDirs.Include))

	var errors errgroup.Group

	for i := range config.ScanDirs.Include {
		j := i // copy loop variable
		globPath := config.ScanDirs.Include[j]

		// An empty include matches nothing, and slicing it would panic.
		if globPath == "" {
			continue
		}
		completeIncludeList = append(completeIncludeList, globPath)

		if strings.HasSuffix(globPath, "*") {
			parent := filepath.Dir(globPath)
			baseGlob := path.Base(globPath[0 : len(globPath)-1])

			entries, err := os.ReadDir(parent)
			if err != nil {
				// The same treatment the walk below gives a directory it
				// cannot read. Expanding a glob used to call log.Fatal here,
				// which took the whole process down and ignored the flag that
				// exists to prevent exactly that.
				if !ignore_dir_errors {
					return err
				}
				log.Printf("ERROR: %s: %v", globPath, err)
				continue
			}

			for _, e := range entries {
				if strings.HasPrefix(e.Name(), baseGlob) {
					completeIncludeList = append(completeIncludeList, parent+"/"+e.Name())
				}
			}

		}
	}

	for i := range completeIncludeList {
		j := i // copy loop variable
		globPath := completeIncludeList[j]

		errors.Go(func() error {
			err := walkone(ctx, globPath, config, stores, results)
			if err == filepath.SkipDir {
				cancel()
			} else if err != nil {
				if ignore_dir_errors {
					log.Printf("ERROR: %s: %v", globPath, err)
					return nil
				} else {
					return err
				}
			}
			return nil
		})
	}

	return errors.Wait()
}
