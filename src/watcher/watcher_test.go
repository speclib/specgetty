package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatcherEmitsOnFileChange(t *testing.T) {
	dir := t.TempDir()

	w, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()

	// Write a file
	if err := os.WriteFile(filepath.Join(dir, "test.md"), []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	select {
	case <-w.Events():
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestWatcherDebounce(t *testing.T) {
	dir := t.TempDir()

	w, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()

	// Write multiple files rapidly
	for i := 0; i < 5; i++ {
		name := filepath.Join(dir, "file"+string(rune('a'+i))+".md")
		if err := os.WriteFile(name, []byte("data"), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	// Should get exactly one event (debounced)
	select {
	case <-w.Events():
		// got first event
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for debounced event")
	}

	// Should NOT get a second event quickly
	select {
	case <-w.Events():
		t.Fatal("got unexpected second event — debounce failed")
	case <-time.After(500 * time.Millisecond):
		// good, no second event
	}
}

func TestWatcherNewSubdirectory(t *testing.T) {
	dir := t.TempDir()

	w, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()

	// Create a subdirectory
	subdir := filepath.Join(dir, "changes", "new-change")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	// Wait for the directory creation event
	select {
	case <-w.Events():
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for dir creation event")
	}

	// Now write a file in the new subdirectory — should also be detected
	if err := os.WriteFile(filepath.Join(subdir, "spec.md"), []byte("spec"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	select {
	case <-w.Events():
		// success — new subdirectory is being watched
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for event in new subdirectory")
	}
}

func TestWatcherClose(t *testing.T) {
	dir := t.TempDir()

	w, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// Events channel should be closed
	_, ok := <-w.Events()
	if ok {
		t.Fatal("expected events channel to be closed")
	}
}
