package scanner

import (
	"context"
	"log"
	"os"
	"strings"
	"path/filepath"
	"path"

	"github.com/karrick/godirwalk"
	"golang.org/x/sync/errgroup"
)

func skip(needle string, haystack []string) bool {
	for _, f := range haystack {
		//FULL PATH COMPARISON
		if(f[0:1]=="/"){
			if f == needle {
				return true
			}

		//PARTIAL PATH COMPARISON
		} else{
			if f == path.Base(needle) {
				return true
			}
		}
	}
	return false
}

// walkone descends a single directory tree looking for OpenSpec projects
func walkone(ctx context.Context, dir string, config *Config, results chan string) error {
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

			if !isValidOpenSpecDir(path) {
				return nil
			}

			results <- filepath.Dir(path)
			return godirwalk.SkipThis // don't descend further
		},
	})
	return err
}

// isValidOpenSpecDir checks whether an openspec/ directory contains the required
// markers: (config.yaml OR project.md) AND (specs/ OR archive/).
func isValidOpenSpecDir(path string) bool {
	hasConfig := false
	if _, err := os.Stat(filepath.Join(path, "config.yaml")); err == nil {
		hasConfig = true
	}
	if !hasConfig {
		if _, err := os.Stat(filepath.Join(path, "project.md")); err == nil {
			hasConfig = true
		}
	}
	if !hasConfig {
		return false
	}

	hasStructure := false
	if info, err := os.Stat(filepath.Join(path, "specs")); err == nil && info.IsDir() {
		hasStructure = true
	}
	if !hasStructure {
		if info, err := os.Stat(filepath.Join(path, "archive")); err == nil && info.IsDir() {
			hasStructure = true
		}
	}
	return hasStructure
}

// Walk finds all OpenSpec projects in the directories specified in config
func Walk(ctx context.Context, config *Config, results chan string, ignore_dir_errors bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	completeIncludeList := config.ScanDirs.Include

	var errors errgroup.Group

	for i := range config.ScanDirs.Include {
		j := i // copy loop variable
		globPath := config.ScanDirs.Include[j]

		if(string(globPath[len(globPath)-1:]) == "*"){
			parent := filepath.Dir(globPath)
			baseGlob := path.Base(globPath[0:len(globPath)-1])

			entries, err := os.ReadDir(parent)
			if err != nil {
				log.Fatal(err)
			}

			for _, e := range entries {
				if strings.HasPrefix(e.Name(), baseGlob) {
					completeIncludeList = append(completeIncludeList, parent + "/" + e.Name())
				}
			}

		}
	}

	for i := range completeIncludeList {
		j := i // copy loop variable
		globPath := completeIncludeList[j]

		errors.Go(func() error {
			err := walkone(ctx, globPath, config, results)
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

	err := errors.Wait()
	close(results)
	return err
}
