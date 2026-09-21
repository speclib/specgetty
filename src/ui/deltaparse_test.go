package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- 1.x the delta grammar ---

const threeOps = `## ADDED Requirements

### Requirement: A new thing
The system SHALL do a new thing.

#### Scenario: It happens
- **WHEN** asked
- **THEN** it SHALL happen

## MODIFIED Requirements

### Requirement: An old thing
The system SHALL do the old thing, differently now.

#### Scenario: Still this
- **WHEN** asked
- **THEN** it SHALL still do this

#### Scenario: And now this
- **WHEN** asked twice
- **THEN** it SHALL also do this

## REMOVED Requirements

### Requirement: A gone thing
**Reason**: It was replaced by the new thing.
**Migration**: Use the new thing.
`

func TestADeltaParsesIntoItsThreeOperations(t *testing.T) {
	tree, problems := parseDelta("a-capability", threeOps)
	if len(problems) > 0 {
		t.Fatalf("should parse, got %v", problems)
	}

	var got []string
	for _, n := range tree.nodes {
		if n.kind == nodeRequirement {
			got = append(got, n.op+" "+n.title)
		}
	}
	want := []string{"ADDED A new thing", "MODIFIED An old thing", "REMOVED A gone thing"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("got %v, want %v", got, want)
	}

	if reqs, scens := tree.counts(); reqs != 3 || scens != 3 {
		t.Errorf("got %d requirements and %d scenarios, want 3 and 3", reqs, scens)
	}
}

func TestEveryNodeCarriesItsCapability(t *testing.T) {
	tree, _ := parseDelta("a-capability", threeOps)
	for _, n := range tree.nodes {
		if n.capability != "a-capability" {
			t.Errorf("node %q carries capability %q", n.title, n.capability)
		}
		if !strings.HasPrefix(n.path, "cap/a-capability/") {
			t.Errorf("node %q has path %q, which another delta could collide with", n.title, n.path)
		}
	}
}

func TestARemovalNeedsNoScenarioAndAnAdditionDoes(t *testing.T) {
	removal := "## REMOVED Requirements\n\n### Requirement: Gone\n**Reason**: Replaced.\n"
	if _, problems := parseDelta("c", removal); len(problems) > 0 {
		t.Errorf("a removal carries a reason instead of scenarios: %v", problems)
	}

	addition := "## ADDED Requirements\n\n### Requirement: New\nIt SHALL.\n"
	_, problems := parseDelta("c", addition)
	if len(problems) == 0 {
		t.Fatal("an added requirement with no scenario should be reported")
	}
	if !strings.Contains(problems[0].text, "scenario") {
		t.Errorf("the reason should name the missing scenario, got %q", problems[0].text)
	}
}

func TestARemovalKeepsItsReasonAndMigration(t *testing.T) {
	tree, _ := parseDelta("c", threeOps)
	var found *specNode
	for i := range tree.nodes {
		if tree.nodes[i].title == "A gone thing" {
			found = &tree.nodes[i]
		}
	}
	if found == nil {
		t.Fatal("the removed requirement is missing")
	}
	// It arrives as the requirement's own prose, which is where a removal's
	// reason sits: there is no scenario for it to be a part of.
	for _, want := range []string{"Reason", "replaced by the new thing", "Migration"} {
		if !strings.Contains(found.body, want) {
			t.Errorf("the removal's card loses %q:\n%s", want, found.body)
		}
	}
}

func TestPurposeIsOptionalInADelta(t *testing.T) {
	if _, problems := parseDelta("c", threeOps); len(problems) > 0 {
		t.Errorf("a delta for an existing capability writes no Purpose: %v", problems)
	}

	withPurpose := "## Purpose\nWhat this capability is for, at some length.\n\n" + threeOps
	tree, problems := parseDelta("c", withPurpose)
	if len(problems) > 0 {
		t.Fatalf("a delta introducing a capability writes one: %v", problems)
	}
	if tree.nodes[0].kind != nodePurpose {
		t.Error("the Purpose should lead the outline")
	}

	empty := "## Purpose\n\n" + threeOps
	if _, problems := parseDelta("c", empty); len(problems) == 0 {
		t.Error("an empty Purpose should be reported, archive copying it into the new spec")
	}
}

func TestAMainSpecHeadingInADeltaIsReported(t *testing.T) {
	_, problems := parseDelta("c", "## Requirements\n\n"+threeOps)
	if len(problems) == 0 {
		t.Fatal("`## Requirements` is not read in a change and should be reported")
	}
	if !strings.Contains(problems[0].text, "applied") {
		t.Errorf("the reason should say the requirements will not be applied, got %q",
			problems[0].text)
	}
}

func TestARequirementUnderNoDeltaHeaderIsReported(t *testing.T) {
	orphan := "### Requirement: Loose\nIt SHALL.\n\n#### Scenario: S\n- **WHEN** a\n- **THEN** b\n"
	_, problems := parseDelta("c", orphan)
	if len(problems) == 0 {
		t.Fatal("a requirement under no delta header is never applied and should be reported")
	}
}

func TestAnUnknownOperationIsCarriedThroughRatherThanHidden(t *testing.T) {
	// One archived change in this repository carries `## NEW Requirements`,
	// which openspec does not apply. Hiding what is under it would make that
	// change unreadable for no benefit.
	odd := "## NEW Requirements\n\n### Requirement: Odd\nIt SHALL.\n\n" +
		"#### Scenario: S\n- **WHEN** a\n- **THEN** b\n"
	tree, problems := parseDelta("c", odd)
	if len(problems) > 0 {
		t.Fatalf("should still be read: %v", problems)
	}
	if len(tree.nodes) == 0 || tree.nodes[0].op != "NEW" {
		t.Fatalf("the operation should be carried as written, got %+v", tree.nodes)
	}
	if sigil, _ := outlineMark(tree.nodes[0]); sigil != "n" {
		t.Errorf("an unknown operation should still mark its row, got %q", sigil)
	}
}

func TestAFencedHeadingInADeltaIsNotAHeading(t *testing.T) {
	fenced := "## ADDED Requirements\n\n### Requirement: R\nIt SHALL.\n\n" +
		"#### Scenario: S\n- **WHEN** a\n- **THEN** b\n\n" +
		"```\n#### Scenario: Not one\n```\n"
	tree, problems := parseDelta("c", fenced)
	if len(problems) > 0 {
		t.Fatalf("should parse: %v", problems)
	}
	if _, scens := tree.counts(); scens != 1 {
		t.Errorf("got %d scenarios, want 1: a heading inside a fence is content", scens)
	}
}

// TestEveryDeltaInThisRepositoryIsRead is task 1.6. It walks the changes of
// this project, which are the only deltas guaranteed to be here, and asserts
// every one of them parses.
func TestEveryDeltaInThisRepositoryIsRead(t *testing.T) {
	root := filepath.Join("..", "..", "openspec", "changes")
	ops := map[string]int{}
	files, reqs, scens := 0, 0, 0

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		if !strings.Contains(filepath.ToSlash(path), "/specs/") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		files++
		cap := filepath.Base(filepath.Dir(path))
		tree, problems := parseDelta(cap, string(b))
		if len(problems) > 0 {
			t.Errorf("%s does not parse as a delta: %v", path, problems)
			return nil
		}
		for _, n := range tree.nodes {
			switch n.kind {
			case nodeRequirement:
				reqs++
				ops[n.op]++
			case nodeScenario:
				scens++
				if len(n.parts) == 0 {
					t.Errorf("%s: scenario %q is empty", path, n.title)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if files == 0 {
		t.Fatal("no delta files found; the corpus this asserts against is missing")
	}
	t.Logf("%d delta files, %d requirements, %d scenarios, per operation: %v",
		files, reqs, scens, ops)
	if ops[opAdded] == 0 || ops[opModified] == 0 || ops[opRemoved] == 0 {
		t.Errorf("the corpus should exercise all three common operations, got %v", ops)
	}
}
