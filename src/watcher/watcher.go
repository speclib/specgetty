package watcher

import (
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

// New creates a Watcher that recursively watches dir and all its subdirectories.
// It returns the watcher and any error from setup.
func New(dir string) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		fsw:    fsw,
		events: make(chan struct{}, 1),
		done:   make(chan struct{}),
	}

	// Add dir and all subdirectories
	if err := w.addRecursive(dir); err != nil {
		fsw.Close()
		return nil, err
	}

	go w.loop()
	return w, nil
}

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
