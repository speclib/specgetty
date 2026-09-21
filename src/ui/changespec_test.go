package ui

import (
	"strings"
	"testing"
)

// --- 2.x comparing a delta against the spec it modifies ---

const liveSpec = `# a-capability Specification

## Purpose
What this capability is for, stated at enough length to count.

## Requirements

### Requirement: An old thing
The system SHALL do the old thing.

#### Scenario: Still this
- **WHEN** asked
- **THEN** it SHALL still do this

#### Scenario: Left alone
- **WHEN** ignored
- **THEN** nothing SHALL happen
`

// The delta keeps one scenario word for word, edits another, and adds a third.
const modifyingDelta = `## MODIFIED Requirements

### Requirement: An old thing
The system SHALL do the old thing, differently now.

#### Scenario: Still this
- **WHEN** asked
- **THEN** it SHALL still do this

#### Scenario: Left alone
- **WHEN** ignored
- **THEN** something SHALL happen after all

#### Scenario: Brand new
- **WHEN** asked twice
- **THEN** it SHALL also do this
`

func comparedTree(t *testing.T, delta, live string) specTree {
	t.Helper()
	tree, problems := buildChangeSpecTree("a-change",
		[]string{"a-capability"},
		map[string]string{"a-capability": delta},
		map[string]string{"a-capability": live})
	if len(problems) > 0 {
		t.Fatalf("should build: %v", problems)
	}
	return tree
}

func markOf(tree specTree, title string) int {
	for _, n := range tree.nodes {
		if n.kind == nodeScenario && n.title == title {
			return n.mark
		}
	}
	return -1
}

func TestAScenarioIsMarkedAgainstTheOneItRestates(t *testing.T) {
	tree := comparedTree(t, modifyingDelta, liveSpec)

	for _, c := range []struct {
		title string
		want  int
		why   string
	}{
		{"Still this", markUnchanged, "word for word the same"},
		{"Left alone", markEdited, "its THEN clause differs"},
		{"Brand new", markAdded, "it has no counterpart"},
	} {
		if got := markOf(tree, c.title); got != c.want {
			t.Errorf("%q: got mark %d, want %d, because it is %s", c.title, got, c.want, c.why)
		}
	}
}

func TestARewrapIsNotAnEdit(t *testing.T) {
	// A delta restates a requirement and the restatement is wrapped afresh.
	// Comparing the source lines would call every scenario edited.
	rewrapped := strings.Replace(modifyingDelta,
		"- **THEN** it SHALL still do this",
		"- **THEN** it SHALL still\n  do this", 1)
	tree := comparedTree(t, rewrapped, liveSpec)
	if got := markOf(tree, "Still this"); got != markUnchanged {
		t.Errorf("got mark %d, want unchanged: only the line breaks moved", got)
	}
	// Asserted in the same test so it cannot pass by marking everything
	// unchanged, which is what a comparison that does nothing would do.
	if got := markOf(tree, "Left alone"); got != markEdited {
		t.Errorf("got mark %d for a scenario whose words changed, want edited", got)
	}
}

func TestAModifiedRequirementCarriesItsOriginal(t *testing.T) {
	tree := comparedTree(t, modifyingDelta, liveSpec)
	for _, n := range tree.nodes {
		if n.kind == nodeRequirement && n.title == "An old thing" {
			if !n.comparable() {
				t.Fatal("a matched requirement should carry the text it modifies")
			}
			if !strings.Contains(n.old.body, "do the old thing.") {
				t.Errorf("the original should be the live spec's prose, got %q", n.old.body)
			}
			return
		}
	}
	t.Fatal("the modified requirement is missing")
}

func TestNoMatchMeansNoMarksRatherThanADefault(t *testing.T) {
	// A renamed heading reads as a different requirement, which is what
	// openspec does when it applies a delta.
	renamed := strings.Replace(modifyingDelta, "An old thing", "A renamed thing", 1)
	tree := comparedTree(t, renamed, liveSpec)

	for _, n := range tree.nodes {
		if n.kind == nodeRequirement && n.title == "A renamed thing" && n.comparable() {
			t.Error("a requirement with no counterpart should carry no original")
		}
		if n.kind == nodeScenario && n.mark != markNone {
			t.Errorf("scenario %q is marked %d with nothing to compare against",
				n.title, n.mark)
		}
	}
}

func TestAnArchivedChangeIsNotCompared(t *testing.T) {
	// The live spec of an archived change is the result of applying it. A
	// delta equal to it must not be reported as altering nothing.
	tree, problems := buildChangeSpecTree("a-change",
		[]string{"a-capability"},
		map[string]string{"a-capability": modifyingDelta},
		nil)
	if len(problems) > 0 {
		t.Fatalf("should build: %v", problems)
	}
	for _, n := range tree.nodes {
		if n.comparable() {
			t.Errorf("node %q carries an original for an archived change", n.title)
		}
		if n.kind == nodeScenario && n.mark != markNone {
			t.Errorf("scenario %q is marked for an archived change", n.title)
		}
	}
}

func TestAnUnparsableOriginalIsNotComparedAgainst(t *testing.T) {
	// Wrong marks would be worse than none.
	tree := comparedTree(t, modifyingDelta, "not a spec at all\n")
	for _, n := range tree.nodes {
		if n.comparable() {
			t.Errorf("node %q compared against a file that is not a spec", n.title)
		}
	}
}

func TestOnlyModifiedRequirementsAreCompared(t *testing.T) {
	added := strings.Replace(modifyingDelta, "## MODIFIED", "## ADDED", 1)
	tree := comparedTree(t, added, liveSpec)
	for _, n := range tree.nodes {
		if n.kind == nodeRequirement && n.comparable() {
			t.Error("an added requirement modifies nothing and has no original")
		}
	}
}

// --- 3.1 the outline covers every delta, in a stable order ---

func TestEveryDeltaOfTheChangeIsListedUnderItsCapability(t *testing.T) {
	names := []string{"cap-a", "cap-b"}
	contents := map[string]string{"cap-a": threeOps, "cap-b": threeOps}

	var first []string
	for round := 0; round < 2; round++ {
		tree, problems := buildChangeSpecTree("a-change", names, contents, nil)
		if len(problems) > 0 {
			t.Fatalf("should build: %v", problems)
		}
		var order []string
		caps := 0
		for _, n := range tree.nodes {
			if n.kind == nodeCapability {
				caps++
				order = append(order, n.title)
			}
		}
		if caps != 2 {
			t.Fatalf("got %d capability nodes, want 2", caps)
		}
		if round == 0 {
			first = order
			continue
		}
		if strings.Join(order, ",") != strings.Join(first, ",") {
			t.Errorf("the order moved between parses: %v then %v", first, order)
		}
	}
}

func TestACapabilityCardSaysWhatTheChangeDoesToIt(t *testing.T) {
	tree, _ := buildChangeSpecTree("a-change", []string{"cap-a"},
		map[string]string{"cap-a": threeOps}, nil)
	for _, n := range tree.nodes {
		if n.kind != nodeCapability {
			continue
		}
		for _, want := range []string{"adds 1 requirement", "modifies 1 requirement",
			"removes 1 requirement"} {
			if !strings.Contains(n.body, want) {
				t.Errorf("the capability card should report %q, got %q", want, n.body)
			}
		}
		return
	}
	t.Fatal("no capability node")
}

func TestADeltaThatCannotBeReadNamesItsFile(t *testing.T) {
	_, problems := buildChangeSpecTree("a-change",
		[]string{"cap-a", "cap-b"},
		map[string]string{"cap-a": threeOps, "cap-b": "### Requirement: Loose\nIt SHALL.\n"},
		nil)
	if len(problems) == 0 {
		t.Fatal("a delta under no operation header should be reported")
	}
	if !strings.HasPrefix(problems[0].text, "cap-b: ") {
		t.Errorf("a change with several deltas must say which file: %q", problems[0].text)
	}
}

// TestTheComparisonReadsOnlyWhatIsInMemory is task 2.5. Both sides of the
// comparison are already scanned: the delta from the change and the original
// from the project's specs. Nothing here opens a file or starts a process,
// which is what makes the comparison affordable on every keystroke.
func TestTheComparisonReadsOnlyWhatIsInMemory(t *testing.T) {
	// Capability names that exist on no disk. If anything were read rather
	// than taken from the maps, there would be nothing to read.
	tree := comparedTree(t, modifyingDelta, liveSpec)
	if got := markOf(tree, "Still this"); got != markUnchanged {
		t.Fatalf("got mark %d for content that exists only in this test", got)
	}
}
