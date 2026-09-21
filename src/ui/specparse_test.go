package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const handWritten = `# thing Specification

## Purpose
What this capability is for.
A second line of it.

### Requirement: The first thing
Its prose, which the card shows.

#### Scenario: One clause
- **WHEN** something happens
- **THEN** something else does

#### Scenario: A clause over three lines
- **GIVEN** a clause whose text runs on
  across a second source line
  and a third
- **AND** another one
`

func TestParseSpecBuildsTheTree(t *testing.T) {
	tree, err := parseSpec("thing", handWritten)
	if err != nil {
		t.Fatal(err)
	}
	if tree.name != "thing" {
		t.Errorf("name = %q", tree.name)
	}

	want := []struct {
		kind  int
		title string
	}{
		{nodePurpose, "Purpose"},
		{nodeRequirement, "The first thing"},
		{nodeScenario, "One clause"},
		{nodeScenario, "A clause over three lines"},
	}
	if len(tree.nodes) != len(want) {
		t.Fatalf("got %d nodes, want %d: %+v", len(tree.nodes), len(want), tree.nodes)
	}
	for i, w := range want {
		if tree.nodes[i].kind != w.kind || tree.nodes[i].title != w.title {
			t.Errorf("node %d = %d/%q, want %d/%q", i,
				tree.nodes[i].kind, tree.nodes[i].title, w.kind, w.title)
		}
	}

	if !strings.Contains(tree.nodes[0].body, "What this capability is for") {
		t.Errorf("purpose body = %q", tree.nodes[0].body)
	}
	if !strings.Contains(tree.nodes[1].body, "Its prose") {
		t.Errorf("requirement body = %q", tree.nodes[1].body)
	}
	if strings.Contains(tree.nodes[1].body, "Scenario") {
		t.Error("a requirement's body must exclude its scenarios")
	}

	req, scen := tree.counts()
	if req != 1 || scen != 2 {
		t.Errorf("counts = %d/%d, want 1/2", req, scen)
	}
}

func TestParseSpecReadsClauses(t *testing.T) {
	tree, err := parseSpec("thing", handWritten)
	if err != nil {
		t.Fatal(err)
	}

	one := tree.nodes[2]
	if len(one.clauses) != 2 {
		t.Fatalf("got %d clauses, want 2: %+v", len(one.clauses), one.clauses)
	}
	if one.clauses[0].keyword != "WHEN" || one.clauses[0].text != "something happens" {
		t.Errorf("clause 0 = %+v", one.clauses[0])
	}
	if one.clauses[1].keyword != "THEN" {
		t.Errorf("clause 1 = %+v", one.clauses[1])
	}

	// A clause may span several source lines, joined into one.
	multi := tree.nodes[3]
	if len(multi.clauses) != 2 {
		t.Fatalf("got %d clauses, want 2: %+v", len(multi.clauses), multi.clauses)
	}
	if multi.clauses[0].keyword != "GIVEN" {
		t.Errorf("clause 0 = %+v", multi.clauses[0])
	}
	joined := multi.clauses[0].text
	for _, want := range []string{"runs on", "second source line", "and a third"} {
		if !strings.Contains(joined, want) {
			t.Errorf("the continuation lines must be joined, missing %q in %q", want, joined)
		}
	}
	if strings.Contains(joined, "\n") {
		t.Errorf("a clause is one line of text, got %q", joined)
	}
	if multi.clauses[1].keyword != "AND" {
		t.Errorf("clause 1 = %+v", multi.clauses[1])
	}
}

func TestParseSpecCoversEveryClauseKeyword(t *testing.T) {
	src := `## Purpose
p

### Requirement: r
text

#### Scenario: s
- **GIVEN** a
- **WHEN** b
- **THEN** c
- **AND** d
`
	tree, err := parseSpec("k", src)
	if err != nil {
		t.Fatal(err)
	}
	got := []string{}
	for _, c := range tree.nodes[2].clauses {
		got = append(got, c.keyword)
	}
	want := []string{"GIVEN", "WHEN", "THEN", "AND"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseSpecRefusesWhatItCannotNavigate(t *testing.T) {
	cases := []struct{ name, src string }{
		{"empty", ""},
		{"only whitespace", "\n\n   \n"},
		{"no requirements", "## Purpose\njust prose\n"},
		{"a requirement with no scenarios", "## Purpose\np\n\n### Requirement: r\ntext only\n"},
		{"the second requirement has none", `## Purpose
p

### Requirement: one
text

#### Scenario: s
- **WHEN** a
- **THEN** b

### Requirement: two
text and nothing else
`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := parseSpec("x", tc.src); err == nil {
				t.Error("want a refusal rather than an outline of nothing")
			}
		})
	}
}

func TestNodePathsSurviveAReParse(t *testing.T) {
	tree, err := parseSpec("thing", handWritten)
	if err != nil {
		t.Fatal(err)
	}
	target := tree.nodes[3].path

	// An edit above the node: the prose grows, the headings do not move.
	edited := strings.Replace(handWritten, "Its prose, which the card shows.",
		"Its prose, rewritten and\nnow two lines long.", 1)
	again, err := parseSpec("thing", edited)
	if err != nil {
		t.Fatal(err)
	}
	if got := again.indexOfPath(target); got != 3 {
		t.Errorf("the node is at %d after a re-parse, want 3", got)
	}

	// The node itself renamed: the path is gone, and saying so is the point.
	renamed := strings.Replace(handWritten, "#### Scenario: A clause over three lines",
		"#### Scenario: Renamed", 1)
	third, err := parseSpec("thing", renamed)
	if err != nil {
		t.Fatal(err)
	}
	if got := third.indexOfPath(target); got != -1 {
		t.Errorf("got %d, want -1 for a node that is gone", got)
	}
}

func TestNodePathsAreUnique(t *testing.T) {
	// Two requirements may carry a scenario of the same name, so a scenario's
	// path has to include the requirement above it.
	src := `## Purpose
p

### Requirement: one
text

#### Scenario: same name
- **WHEN** a
- **THEN** b

### Requirement: two
text

#### Scenario: same name
- **WHEN** c
- **THEN** d
`
	tree, err := parseSpec("x", src)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, n := range tree.nodes {
		if seen[n.path] {
			t.Errorf("path %q is not unique", n.path)
		}
		seen[n.path] = true
	}
}

func TestTheParserAcceptsEveryLiveSpec(t *testing.T) {
	paths, err := filepath.Glob("../../openspec/specs/*/spec.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Skip("not this project")
	}

	totalReq, totalScen := 0, 0
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		name := filepath.Base(filepath.Dir(p))
		tree, err := parseSpec(name, string(b))
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		req, scen := tree.counts()
		totalReq += req
		totalScen += scen
	}
	t.Logf("%d specs, %d requirements, %d scenarios", len(paths), totalReq, totalScen)
	if totalReq == 0 || totalScen == 0 {
		t.Error("the corpus must parse into something")
	}
	if totalScen < totalReq {
		t.Errorf("got %d scenarios for %d requirements; every requirement has at least one",
			totalScen, totalReq)
	}
}
