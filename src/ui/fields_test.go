package ui

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mipmip/specgetty/src/scanner"
)

func TestResolveFieldsPrecedence(t *testing.T) {
	t.Run("default when neither is set", func(t *testing.T) {
		got, err := ResolveFields("", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, defaultFields) {
			t.Errorf("got %v, want %v", got, defaultFields)
		}
	})

	t.Run("config used when no flag", func(t *testing.T) {
		got, err := ResolveFields("", "name,tasks")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, []string{"name", "tasks"}) {
			t.Errorf("got %v, want [name tasks]", got)
		}
	})

	t.Run("flag beats config", func(t *testing.T) {
		got, err := ResolveFields("name,specs", "name,tasks,archived")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, []string{"name", "specs"}) {
			t.Errorf("got %v, want [name specs]", got)
		}
	})

	t.Run("whitespace tolerated", func(t *testing.T) {
		got, err := ResolveFields(" name , tasks ", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, []string{"name", "tasks"}) {
			t.Errorf("got %v, want [name tasks]", got)
		}
	})
}

func TestResolveFieldsUnknownName(t *testing.T) {
	_, err := ResolveFields("name,bogus", "")
	if err == nil {
		t.Fatal("expected an error for an unknown field")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Errorf("error %q should name the offending field", err)
	}
	// The message must list what is valid, or a typo is a guessing game.
	if !strings.Contains(err.Error(), "tasks") {
		t.Errorf("error %q should list the valid fields", err)
	}
	if !strings.Contains(err.Error(), "--change-fields") {
		t.Errorf("error %q should name where the bad value came from", err)
	}
}

func TestResolveFieldsReportsConfigAsSource(t *testing.T) {
	_, err := ResolveFields("", "name,bogus")
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "config file") {
		t.Errorf("error %q should point at the config file, not the flag", err)
	}
}

func TestLayoutFieldsDropsColumnsFromTheRight(t *testing.T) {
	fields := []string{"name", "tasks", "specs", "archived"}

	wide, widths := layoutFields(fields, 100)
	if !reflect.DeepEqual(wide, fields) {
		t.Errorf("at width 100 kept %v, want all of %v", wide, fields)
	}
	if widths[0] < minFlexWidth {
		t.Errorf("flexible name column got %d, want at least %d", widths[0], minFlexWidth)
	}

	// At 30 the fixed columns leave the name under minFlexWidth, so the
	// rightmost column is dropped rather than squeezing the name further.
	narrow, _ := layoutFields(fields, 30)
	if len(narrow) >= len(fields) {
		t.Errorf("at width 30 kept %v, want fewer than %d columns", narrow, len(fields))
	}
	// Whatever survives must be a prefix: columns drop from the right.
	for i, f := range narrow {
		if f != fields[i] {
			t.Errorf("kept %v, want a prefix of %v", narrow, fields)
			break
		}
	}
}

func TestLayoutFieldsKeepsNameWhenVeryNarrow(t *testing.T) {
	kept, widths := layoutFields([]string{"name", "tasks", "specs"}, 10)
	if len(kept) != 1 || kept[0] != "name" {
		t.Errorf("kept %v, want just [name] at width 10", kept)
	}
	if widths[0] < 1 {
		t.Errorf("name width = %d, want at least 1", widths[0])
	}
}

func TestFitCellPadsAndTruncates(t *testing.T) {
	if got := fitCell("ab", 5); got != "ab   " {
		t.Errorf("fitCell(ab,5) = %q, want padded to 5", got)
	}
	if got := fitCell("abcdef", 4); got != "abc…" {
		t.Errorf("fitCell(abcdef,4) = %q, want abc…", got)
	}
	if got := fitCell("abc", 0); got != "" {
		t.Errorf("fitCell(abc,0) = %q, want empty", got)
	}
}

func TestFieldValues(t *testing.T) {
	open := changeRow{ci: scanner.ChangeInfo{
		Name: "my-change", TasksTotal: 5, TasksDone: 3,
		SpecNames: []string{"a", "b"},
	}}
	archived := changeRow{ci: scanner.ChangeInfo{Name: "old"}, archived: true}

	if got := knownFields["tasks"].value(open); got != "3/5" {
		t.Errorf("tasks = %q, want 3/5", got)
	}
	if got := knownFields["specs"].value(open); got != "2" {
		t.Errorf("specs = %q, want 2", got)
	}
	if got := knownFields["archived"].value(open); got != "open" {
		t.Errorf("archived = %q, want open", got)
	}
	if got := knownFields["archived"].value(archived); got != "archived" {
		t.Errorf("archived = %q, want archived", got)
	}
	// A change with no tasks.md shows an empty cell rather than 0/0.
	if got := knownFields["tasks"].value(archived); got != "" {
		t.Errorf("tasks for a change without tasks.md = %q, want empty", got)
	}
	if got := knownFields["specs"].value(archived); got != "" {
		t.Errorf("specs for a change without specs = %q, want empty", got)
	}
}
