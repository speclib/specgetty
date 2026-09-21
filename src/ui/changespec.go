package ui

import (
	"fmt"
	"strings"
)

// A change's spec deltas read as one outline: a node per capability the change
// touches, its requirements under it marked with what the change does to them,
// and each requirement's scenarios under that.
//
// The marks are the point. More than half of what a MODIFIED requirement
// restates is text it does not change, and a reader has no way to tell which
// half is which without comparing by eye.

// buildChangeSpecTree reads every delta of a change into one outline.
//
// live is the project's main specs, which for a change that has not been
// archived are the text the deltas modify. For an archived change it is the
// result of applying them, so nothing is passed and nothing is compared.
func buildChangeSpecTree(name string, specNames []string, contents map[string]string,
	live map[string]string) (specTree, []specProblem) {

	t := specTree{name: name}
	var problems []specProblem

	for _, cap := range specNames {
		sub, probs := parseDelta(cap, contents[cap])
		if len(probs) > 0 {
			// Named, so a change with several deltas says which file is at
			// fault rather than reporting line numbers against nothing.
			for _, p := range probs {
				problems = append(problems,
					specProblem{line: p.line, text: cap + ": " + p.text})
			}
			continue
		}

		t.nodes = append(t.nodes, specNode{
			kind:       nodeCapability,
			title:      cap,
			capability: cap,
			body:       capabilitySummary(sub),
			path:       pathOf(cap, "capability"),
		})
		if live != nil {
			compareDelta(&sub, live[cap])
		}
		t.nodes = append(t.nodes, sub.nodes...)
	}

	if len(problems) > 0 {
		return specTree{}, problems
	}
	return t, nil
}

// capabilitySummary says what the change does to one capability.
func capabilitySummary(sub specTree) string {
	counts := map[string]int{}
	var order []string
	for _, n := range sub.nodes {
		if n.kind != nodeRequirement {
			continue
		}
		if counts[n.op] == 0 {
			order = append(order, n.op)
		}
		counts[n.op]++
	}
	if len(order) == 0 {
		return "This change alters no requirement of this capability."
	}
	// One line per operation rather than a sentence listing them. A card has
	// the room, and three counts read faster stacked than joined by commas
	// that repeat the noun.
	var lines []string
	for _, op := range order {
		lines = append(lines, fmt.Sprintf("%s %s", opVerb(op),
			plural(counts[op], "requirement", "requirements")))
	}
	return "This change " + strings.Join(lines, "\n\n")
}

// opVerb says an operation as something the change does, rather than as the
// heading it was written under.
func opVerb(op string) string {
	switch op {
	case opAdded:
		return "adds"
	case opModified:
		return "modifies"
	case opRemoved:
		return "removes"
	case opRenamed:
		return "renames"
	}
	return "lists " + strings.ToLower(op) + ":"
}

// compareDelta attaches each MODIFIED requirement's original, and marks its
// scenarios against it.
//
// A requirement is matched by its heading, exactly, which is what OpenSpec
// itself does when it applies a delta. A heading that was reworded therefore
// reads as a different requirement here too, which is correct even where it is
// unhelpful: the archive would not apply it either.
func compareDelta(sub *specTree, liveContent string) {
	if liveContent == "" {
		return
	}
	original, problems := parseSpec(sub.name, liveContent)
	if len(problems) > 0 {
		// A main spec that does not parse is not an original this view can
		// compare against. No marks is right; wrong marks would not be.
		return
	}

	byTitle := map[string]int{}
	for i, n := range original.nodes {
		if n.kind == nodeRequirement {
			byTitle[n.title] = i
		}
	}

	for i := range sub.nodes {
		n := &sub.nodes[i]
		if n.kind != nodeRequirement || n.op != opModified {
			continue
		}
		at, ok := byTitle[n.title]
		if !ok {
			continue
		}
		old := original.nodes[at]
		n.old = &old

		// The original's scenarios, which are the nodes following it up to the
		// next requirement.
		oldScen := map[string]specNode{}
		for j := at + 1; j < len(original.nodes); j++ {
			if original.nodes[j].kind != nodeScenario {
				break
			}
			oldScen[original.nodes[j].title] = original.nodes[j]
		}

		for j := i + 1; j < len(sub.nodes); j++ {
			s := &sub.nodes[j]
			if s.kind != nodeScenario {
				break
			}
			prev, found := oldScen[s.title]
			if !found {
				s.mark = markAdded
				continue
			}
			before := prev
			s.old = &before
			if samePartsText(prev.parts, s.parts) {
				s.mark = markUnchanged
			} else {
				s.mark = markEdited
			}
		}
	}
}

// samePartsText reports whether two scenarios say the same thing.
//
// Compared on the parsed parts rather than on the source lines, so a rewrap
// that changed where the lines break is not reported as an edit. That is the
// ordinary shape of a delta: the requirement is restated, and the restatement
// is wrapped afresh.
func samePartsText(a, b []specPart) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].kind != b[i].kind || a[i].keyword != b[i].keyword ||
			normaliseSpace(a[i].text) != normaliseSpace(b[i].text) {
			return false
		}
	}
	return true
}

func normaliseSpace(s string) string { return strings.Join(strings.Fields(s), " ") }

// outlineMark is the sigil a node carries in the outline, and the style it is
// drawn in. An empty sigil is a node with nothing to say about itself.
func outlineMark(n specNode) (sigil string, style int) {
	switch n.kind {
	case nodeRequirement:
		switch n.op {
		case opAdded:
			return "+", markAdded
		case opModified:
			return "~", markEdited
		case opRemoved:
			return "-", markUnchanged
		case "":
			return "", markNone
		default:
			// An operation nothing recognises still says what it said.
			return strings.ToLower(n.op[:1]), markNone
		}
	case nodeScenario:
		switch n.mark {
		case markAdded:
			return "+", markAdded
		case markEdited:
			return "~", markEdited
		}
	}
	return "", markNone
}

// comparable reports whether a node has an original to show beside itself.
func (n specNode) comparable() bool { return n.old != nil }

// The specs sub-tab of an open change is the one that lists several files. It
// is the last sub-tab when a change carries spec deltas, which is what
// artifactTabNames appends it as.

// changeSpecsTabIndex is the sub-tab that lists the change's deltas, or -1.
func (m model) changeSpecsTabIndex() int {
	r, ok := m.selectedRow()
	if !ok || len(r.ci.SpecNames) == 0 {
		return -1
	}
	return len(r.ci.ArtifactFiles)
}

// onChangeSpecsTab reports whether the open change is showing its deltas.
func (m model) onChangeSpecsTab() bool {
	at := m.changeSpecsTabIndex()
	return at >= 0 && m.changeArtifactTab == at
}

// openChangeSpecs descends into the open change's spec deltas.
//
// enter descends in every case, as it does from the specs tab. Whether the
// deltas can be structured is answered by the view it opens, a report naming
// each reason being worth more than a transient line on the nav bar.
func (m model) openChangeSpecs() model {
	r, ok := m.selectedRow()
	if !ok || len(r.ci.SpecNames) == 0 {
		return m
	}

	// An archived change's deltas have already been applied to the project's
	// specs, so those specs are the result rather than the original. Nothing is
	// passed, and nothing is compared.
	var live map[string]string
	if !r.archived {
		live = m.projects[m.currentKey()].Info.SpecContents
	}

	tree, problems := buildChangeSpecTree(r.ci.Name, r.ci.SpecNames, r.ci.SpecContents, live)
	m.specTree = tree
	m.specProblems = problems
	m.specName = r.ci.Name
	m.specNode = 0
	m.specNodePath = ""
	if len(tree.nodes) > 0 {
		m.specNodePath = tree.nodes[0].path
	}
	m.level = levelChangeSpec
	m.focus = focusListPane
	m.cardView = viewDiff
	m.statusMsg = ""
	return m
}

// selectedNode returns the node the outline cursor is on.
func (m model) selectedNode() (specNode, bool) {
	if len(m.specTree.nodes) == 0 {
		return specNode{}, false
	}
	return m.specTree.nodes[m.selectedSpecNode()], true
}

// moveCardView steps the three-way chooser, stopping at either end.
//
// A node with nothing to compare has no chooser, so the keys do nothing there:
// the absence of the row and the absence of a choice are the same fact.
func (m *model) moveCardView(delta int) {
	n, ok := m.selectedNode()
	if !ok || !n.comparable() {
		return
	}
	to := m.cardView + delta
	if to < viewDiff {
		to = viewDiff
	}
	if to > viewNew {
		to = viewNew
	}
	m.cardView = to
}

// reparseOpenChangeSpecs rebuilds a change's outline after a rescan, keeping
// the cursor on the node it was on.
//
// A delta edited into a shape this view cannot read keeps the outline it opened
// with, rather than emptying the view under a reader mid-sentence. That is the
// rule reparseOpenSpec already follows for a main spec.
func (m *model) reparseOpenChangeSpecs() {
	r, ok := m.selectedRow()
	if !ok || len(r.ci.SpecNames) == 0 {
		return
	}
	var live map[string]string
	if !r.archived {
		live = m.projects[m.currentKey()].Info.SpecContents
	}
	tree, problems := buildChangeSpecTree(r.ci.Name, r.ci.SpecNames, r.ci.SpecContents, live)
	if len(problems) > 0 {
		return
	}
	m.specTree = tree
	m.syncSpecNode()
}
