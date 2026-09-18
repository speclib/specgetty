package watcher

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

const debounceDuration = 200 * time.Millisecond

// Watcher monitors an openspec directory tree for changes and emits
// debounced notifications on a channel.
type Watcher struct {
	fsw    *fsnotify.Watcher
	events chan struct{}
	done   chan struct{}
}

// New creates a Watcher that recursively watches every given directory and all
// their subdirectories. It returns the watcher and any error from setup.
//
// More than one tree is watched when a project reads its content from a store:
// the store's tree, where specs and changes move, and the originating repo's,
// where the declaration pointing at the store lives. The second holds one file,
// but editing that file changes everything on screen without touching the
// store at all.
//
// A directory that does not exist is skipped rather than failing the whole
// watcher, so a project with one readable tree still gets notifications.
func New(dirs ...string) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		fsw:    fsw,
		events: make(chan struct{}, 1),
		done:   make(chan struct{}),
	}

	watched := 0
	seen := make(map[string]bool, len(dirs))
	for _, dir := range dirs {
		if dir == "" || seen[dir] {
			continue
		}
		seen[dir] = true
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			continue
		}
		if err := w.addRecursive(dir); err != nil {
			fsw.Close()
			return nil, err
		}
		watched++
	}
	if watched == 0 {
		fsw.Close()
		return nil, errNothingToWatch
	}

	go w.loop()
	return w, nil
}

// errNothingToWatch is returned when none of the given directories exist. The
// caller logs it and carries on without a watcher rather than failing to open
// the project.
var errNothingToWatch = errors.New("watcher: no directory to watch")

// Events returns a channel that receives a value each time a debounced
// filesystem change is detected. The channel is closed when the watcher stops.
func (w *Watcher) Events() <-chan struct{} {
	return w.events
}

// Close stops the watcher and closes the events channel.
func (w *Watcher) Close() error {
	err := w.fsw.Close()
	<-w.done // wait for loop to exit
	return err
}

func (w *Watcher) addRecursive(dir string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip inaccessible entries
		}
		if d.IsDir() {
			return w.fsw.Add(path)
		}
		return nil
	})
}

func (w *Watcher) loop() {
	defer close(w.events)
	defer close(w.done)

	var timer *time.Timer
	var timerC <-chan time.Time

	for {
		select {
		case event, ok := <-w.fsw.Events:
			if !ok {
				return
			}

			// If a new directory is created, watch it too
			if event.Has(fsnotify.Create) {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					_ = w.addRecursive(event.Name)
				}
			}

			// Reset debounce timer
			if timer == nil {
				timer = time.NewTimer(debounceDuration)
				timerC = timer.C
			} else {
				timer.Reset(debounceDuration)
			}

		case <-timerC:
			timer = nil
			timerC = nil
			// Non-blocking send — if a previous event hasn't been consumed yet, skip
			select {
			case w.events <- struct{}{}:
			default:
			}

		case _, ok := <-w.fsw.Errors:
			if !ok {
				return
			}
			// Errors are non-fatal for our use case; continue watching
		}
	}
}
