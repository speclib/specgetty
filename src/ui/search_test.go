package ui

import (
	"reflect"
	"testing"

	"github.com/mipmip/specgetty/src/scanner"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		raw       string
		wantKind  int
		wantTerm  string
		wantCased bool
	}{
		{"", matchFuzzyName, "", false},
		{"export", matchFuzzyName, "export", false},
		{"Export", matchFuzzyName, "Export", true},
		{"'export", matchLiteralName, "export", false},
		{"'Export", matchLiteralName, "Export", true},
		{":export", matchBody, "export", false},
		{":Export", matchBody, "Export", true},
		{"'", matchLiteralName, "", false},
		{":", matchBody, "", false},
	}
	for _, tt := range tests {
		got := parseQuery(tt.raw)
		if got.kind != tt.wantKind || got.term != tt.wantTerm || got.caseSensitive != tt.wantCased {
			t.Errorf("parseQuery(%q) = {kind:%d term:%q cased:%v}, want {kind:%d term:%q cased:%v}",
				tt.raw, got.kind, got.term, got.caseSensitive, tt.wantKind, tt.wantTerm, tt.wantCased)
		}
	}
}

func TestQueryEmpty(t *testing.T) {
	if !parseQuery("").empty() {
		t.Error("empty raw query should be empty")
	}
	// A bare sigil carries no term, so it filters nothing.
	if !parseQuery(":").empty() {
		t.Error("bare : should be empty")
	}
	if parseQuery("x").empty() {
		t.Error("non-empty term should not be empty")
	}
}

func rowsFor(names ...string) []changeRow {
	rows := make([]changeRow, len(names))
	for i, n := range names {
		rows[i] = changeRow{ci: scanner.ChangeInfo{Name: n}}
	}
	return rows
}

func names(rows []filtered[changeRow]) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.row.ci.Name
	}
	return out
}

// plainNames reads names off unwrapped rows.
func plainNames(rows []changeRow) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.ci.Name
	}
	return out
}

func TestFilterRowsFuzzyOnName(t *testing.T) {
	rows := rowsFor("export-change-as-zip", "change-list-full-view", "discard-change-action")

	got := names(filterRows(rows, parseQuery("expzip")))
	if !reflect.DeepEqual(got, []string{"export-change-as-zip"}) {
		t.Errorf("fuzzy subsequence match = %v, want [export-change-as-zip]", got)
	}

	// Every name contains the subsequence "change", so all three survive.
	if got := filterRows(rows, parseQuery("change")); len(got) != 3 {
		t.Errorf("got %d rows, want 3: %v", len(got), names(got))
	}
}

func TestFilterRowsRanksStrongestFirst(t *testing.T) {
	// "export" leads one name and is buried in the other, so the leading match
	// scores higher and must sort first.
	rows := rowsFor("some-other-export-thing", "export-change-as-zip")
	got := names(filterRows(rows, parseQuery("export")))
	if len(got) != 2 {
		t.Fatalf("got %d rows, want 2: %v", len(got), got)
	}
	if got[0] != "export-change-as-zip" {
		t.Errorf("strongest match = %q, want export-change-as-zip (order: %v)", got[0], got)
	}
}

func TestFilterRowsSmartCase(t *testing.T) {
	rows := rowsFor("Export-Change", "export-change")

	// All-lowercase query folds case, so both match.
	if got := filterRows(rows, parseQuery("export")); len(got) != 2 {
		t.Errorf("lowercase query matched %d rows, want 2: %v", len(got), names(got))
	}

	// An uppercase character makes the query case-sensitive.
	got := names(filterRows(rows, parseQuery("Export")))
	if !reflect.DeepEqual(got, []string{"Export-Change"}) {
		t.Errorf("cased query = %v, want [Export-Change]", got)
	}
}

func TestFilterRowsLiteralName(t *testing.T) {
	rows := rowsFor("export-change-as-zip", "expzip-not-literal")

	// The fuzzy form matches both; the literal form matches only the substring.
	got := names(filterRows(rows, parseQuery("'export")))
	if !reflect.DeepEqual(got, []string{"export-change-as-zip"}) {
		t.Errorf("literal name query = %v, want [export-change-as-zip]", got)
	}

	if got := filterRows(rows, parseQuery("'xpc")); len(got) != 0 {
		t.Errorf("literal query matched %v, want nothing (it is not a substring)", names(got))
	}
}

func TestFilterRowsLiteralKeepsInputOrder(t *testing.T) {
	rows := rowsFor("b-change", "a-change")
	got := names(filterRows(rows, parseQuery("'change")))
	if !reflect.DeepEqual(got, []string{"b-change", "a-change"}) {
		t.Errorf("literal results = %v, want input order [b-change a-change]", got)
	}
}

func TestFilterRowsBodyNamesMatchedFiles(t *testing.T) {
	rows := []changeRow{
		{ci: scanner.ChangeInfo{
			Name:          "alpha",
			ArtifactFiles: []string{"design.md", "proposal.md"},
			ArtifactContents: map[string]string{
				"proposal.md": "we should zip the change",
				"design.md":   "nothing relevant here",
			},
			SpecNames:    []string{"export-change"},
			SpecContents: map[string]string{"export-change": "creates a zip file"},
		}},
		{ci: scanner.ChangeInfo{
			Name:             "beta",
			ArtifactFiles:    []string{"proposal.md"},
			ArtifactContents: map[string]string{"proposal.md": "unrelated text"},
		}},
	}

	got := filterRows(rows, parseQuery(":zip"))
	if len(got) != 1 {
		t.Fatalf("got %d rows, want 1: %v", len(got), names(got))
	}
	if got[0].row.ci.Name != "alpha" {
		t.Errorf("matched %q, want alpha", got[0].row.ci.Name)
	}
	want := []string{"export-change", "proposal"}
	if !reflect.DeepEqual(got[0].matched, want) {
		t.Errorf("matched labels = %v, want %v", got[0].matched, want)
	}
}

func TestFilterRowsBodyMatchesNameWithoutFileHint(t *testing.T) {
	rows := []changeRow{
		{ci: scanner.ChangeInfo{
			Name:             "zip-things",
			ArtifactContents: map[string]string{"proposal.md": "nothing"},
		}},
	}
	got := filterRows(rows, parseQuery(":zip"))
	if len(got) != 1 {
		t.Fatalf("got %d rows, want 1", len(got))
	}
	if len(got[0].matched) != 0 {
		t.Errorf("matched labels = %v, want none for a name-only match", got[0].matched)
	}
}

func TestFilterRowsEmptyQueryReturnsEverything(t *testing.T) {
	rows := rowsFor("a", "b", "c")
	if got := filterRows(rows, parseQuery("")); len(got) != 3 {
		t.Errorf("empty query returned %d rows, want 3", len(got))
	}
}

func TestFilterRowsClearsStaleMatchedFiles(t *testing.T) {
	// Match reasons are returned alongside the row rather than written onto it,
	// so a name match cannot inherit hints from an earlier body search.
	rows := []changeRow{{ci: scanner.ChangeInfo{Name: "alpha"}}}
	got := filterRows(rows, parseQuery("alpha"))
	if len(got) != 1 {
		t.Fatalf("got %d rows, want 1", len(got))
	}
	if got[0].matched != nil {
		t.Errorf("matched labels = %v, want nil after a name match", got[0].matched)
	}
}
