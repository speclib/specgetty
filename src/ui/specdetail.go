package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// The spec detail view is an outline beside a card: the outline says what the
// spec contains, the card shows whichever node the cursor is on. It is the
// third navigation level, and it follows the same rule as the change view:
// `enter` descends, `esc` ascends.

// outlineRowStyle picks how one outline row is drawn.
//
// A main spec has only the requirement heading to distinguish, which is what it
// had before. A change's outline says more: what the change does to each
// requirement, and which of a modified requirement's scenarios it leaves alone.
func outlineRowStyle(n specNode) lipgloss.Style {
	switch n.kind {
	case nodeCapability:
		return capabilityStyle
	case nodeRequirement:
		switch n.op {
		case opAdded:
			return opAddedStyle
		case opModified:
			return opModifiedStyle
		case opRemoved:
			return opRemovedStyle
		case "":
			return sectionHeaderStyle
		}
		return sectionHeaderStyle
	case nodeScenario:
		switch n.mark {
		case markAdded:
			return opAddedStyle
		case markEdited:
			return opModifiedStyle
		case markUnchanged:
			return unchangedStyle
		}
	}
	return normalStyle
}

// specDetailSplit sizes the outline against the card.
//
// Four tenths, with a floor: an outline of scenario titles needs more room than
// a list of capability names, and a card of wrapped clauses needs more than a
// spec's raw markdown. Both the renderer and docRegion ask this, so the card is
// wrapped to the width it is drawn at.
func specDetailSplit(width int) (outlineOuter, cardOuter int) {
	outlineOuter = width * 4 / 10
	if outlineOuter < 22 {
		outlineOuter = 22
	}
	if max := width - 24; outlineOuter > max {
		outlineOuter = max
	}
	// Below about 46 columns neither minimum can be met. The outline is what
	// gives way: a card of no columns shows no spec at all, which is worse than
	// an outline of clipped labels.
	if outlineOuter < 1 {
		outlineOuter = 1
	}
	cardOuter = width - outlineOuter - 1
	if cardOuter < 1 {
		cardOuter = 1
	}
	return outlineOuter, cardOuter
}

// outlineRow is one drawn row of the outline, and the node it belongs to.
//
// A label wraps rather than being cut, so one node can occupy several rows. The
// cursor still moves one node per keystroke, and every row of the selected node
// is highlighted.
type outlineRow struct {
	text string
	node int
}

// markedTree reports whether a tree came from a change's deltas.
//
// Only such a tree has capability roots, and only such a tree needs the gutter
// the marks sit in. A main spec is drawn exactly as it was before there were
// marks at all.
func markedTree(tree specTree) bool {
	for _, n := range tree.nodes {
		if n.kind == nodeCapability {
			return true
		}
	}
	return false
}

// outlineRows lays the tree out at a given width.
func outlineRows(tree specTree, width int) []outlineRow {
	gutter := markedTree(tree)

	var rows []outlineRow
	for i, n := range tree.nodes {
		indent, hang := "", "  "
		switch {
		case n.kind == nodeScenario && gutter:
			indent, hang = "  ", "      "
		case n.kind == nodeScenario:
			indent, hang = "  ", "    "
		}
		// The mark sits in a gutter every row of a marked tree carries, so a
		// requirement with a mark and one without still line up. The gutter is
		// two columns whether or not this node fills them.
		if gutter && n.kind != nodeCapability {
			sigil, _ := outlineMark(n)
			if sigil == "" {
				sigil = " "
			}
			indent += sigil + " "
			hang += "  "
		}

		// Wrapped against the hanging indent, which is the wider of the two: a
		// continuation row that only got its indent after being wrapped could
		// come out wider than the pane, and lipgloss would then wrap it again
		// into a row the row count did not know about.
		wrapped := strings.Split(ansi.Wrap(n.title, max(1, width-len(hang)), " "), "\n")
		for j, line := range wrapped {
			prefix := indent
			if j > 0 {
				// A continuation never repeats the mark: the gutter it would
				// sit in is part of the hanging indent.
				prefix = hang
			}
			rows = append(rows, outlineRow{text: prefix + strings.TrimSpace(line), node: i})
		}
	}
	return rows
}

// ownersOfRows marks which node each drawn row belongs to.
//
// The outline draws no chrome, so every row has an owner. It reads the same
// shared arithmetic the change list and the properties list do.
func ownersOfRows(rows []outlineRow) itemLines {
	owners := make(itemLines, len(rows))
	for i, r := range rows {
		owners[i] = r.node
	}
	return owners
}

// rowRangeOfNode returns the first and last drawn row a node occupies.
func rowRangeOfNode(rows []outlineRow, node int) (first, last int) {
	return ownersOfRows(rows).span(node)
}

// renderSpecOutline draws the outline, scrolled so the whole selected node is
// visible.
func renderSpecOutline(tree specTree, selected, width, height int, lit bool) string {
	rows := outlineRows(tree, width)
	// A node taller than the pane shows its top rather than its bottom, which
	// is what the shared offset does for every list that has one.
	offset := ownersOfRows(rows).offsetFor(selected, height)
	end := offset + height
	if end > len(rows) {
		end = len(rows)
	}

	var b strings.Builder
	for i := offset; i < end; i++ {
		if i > offset {
			b.WriteString("\n")
		}
		r := rows[i]
		switch {
		case r.node == selected && lit:
			b.WriteString(selectedStyle.Width(width).Render(r.text))
		case r.node == selected:
			b.WriteString(dimSelectedStyle.Width(width).Render(r.text))
		default:
			b.WriteString(outlineRowStyle(tree.nodes[r.node]).Render(r.text))
		}
	}
	return b.String()
}

// renderSpecCard renders the node the cursor is on.
//
// A requirement's card carries its own prose and not its scenarios: the
// scenarios are nodes of their own, and repeating them would make the card the
// document the view exists to replace.
func renderSpecCard(n specNode, width int) string {
	pad := cardPadding(width)
	inner := max(1, width-2*pad)
	left := strings.Repeat(" ", pad)

	var b strings.Builder
	write := func(style lipgloss.Style, text string) {
		for _, line := range strings.Split(ansi.Wrap(text, inner, " "), "\n") {
			b.WriteString(left + style.Render(line) + "\n")
		}
	}

	b.WriteString("\n")
	switch n.kind {
	case nodePurpose:
		write(sectionHeaderStyle, "Purpose")
		b.WriteString("\n")
		writeProse(&b, n.body, left, inner)
	case nodeRequirement:
		write(sectionHeaderStyle, "Requirement: "+n.title)
		b.WriteString("\n")
		writeProse(&b, n.body, left, inner)
	default:
		write(sectionHeaderStyle, "Scenario: "+n.title)
		for _, part := range n.parts {
			b.WriteString("\n")
			if part.kind == partProse {
				// Content in a shape the clause layout does not recognise is
				// drawn as the prose it is, in its place among the clauses.
				// Leaving it out would be the bug this view shipped with.
				writeProse(&b, part.text, left, inner)
				continue
			}
			// By role rather than one style for all four: the condition, the
			// assertion and the continuation of a clause read apart, which is
			// what lets the shape of a scenario be read without reading it.
			b.WriteString(left + clauseStyleFor(part.keyword).Render(part.keyword) + "\n")
			// A hanging indent, so every row after the first still reads as
			// belonging to the keyword above it.
			clauseIndent := left + "   "
			wrapped := ansi.Wrap(renderInlineMarkdown(part.text), max(1, inner-3), " ")
			for _, line := range strings.Split(wrapped, "\n") {
				b.WriteString(clauseIndent + line + "\n")
			}
		}
	}
	return b.String()
}

// writeProse lays a node's body out, one paragraph at a time.
func writeProse(b *strings.Builder, body, left string, inner int) {
	first := true
	for _, para := range strings.Split(body, "\n\n") {
		joined := strings.Join(strings.Fields(strings.ReplaceAll(para, "\n", " ")), " ")
		if joined == "" {
			continue
		}
		// Between paragraphs, not after the last one. The caller separates
		// whatever comes next, and a trailing blank here would double it.
		if !first {
			b.WriteString("\n")
		}
		first = false
		// Styled before it is wrapped. The other way round, a backticked span
		// that happens to straddle the wrap is two halves with one backtick
		// each, and neither half is a span any more: the marks stay on screen
		// and the colour never arrives.
		for _, line := range strings.Split(ansi.Wrap(renderInlineMarkdown(joined), inner, " "), "\n") {
			b.WriteString(left + line + "\n")
		}
	}
}

// cardPadding is the air inside the card, which gives way as the pane narrows.
//
// One formula rather than a table, so there is no width at which the card is
// suddenly cramped or suddenly roomy.
func cardPadding(width int) int {
	pad := width / 12
	if pad > 5 {
		pad = 5
	}
	if pad < 1 {
		pad = 1
	}
	return pad
}

// openSelectedSpec descends into the spec under the specs-tab cursor.
//
// A spec that does not fit the structure does not open: it reports why and
// leaves the level and the cursor alone, because an outline of nothing is
// worse than staying where you are.
func (m model) openSelectedSpec() model {
	key := m.currentKey()
	if key == "" {
		return m
	}
	info := m.projects[key].Info
	if m.specCursor >= len(info.SpecNames) {
		return m
	}
	name := info.SpecNames[m.specCursor]

	// enter descends in every case. Whether the file can be structured is
	// answered by the view it opens: a file usually fails in more than one way,
	// and a transient line on the nav bar can hold neither the reasons nor the
	// lines they sit on nor the key that opens an editor on them.
	tree, problems := parseSpec(name, info.SpecContents[name])
	m.specTree = tree
	m.specProblems = problems
	m.specName = name
	m.specNode = 0
	m.specNodePath = ""
	if len(tree.nodes) > 0 {
		m.specNodePath = tree.nodes[0].path
	}
	m.level = levelSpec
	m.focus = focusListPane
	m.statusMsg = ""
	return m
}

// renderSpecReport draws the reasons a file is not a spec.
//
// Every reason, each with the line it sits on, rather than the first: a file
// that does not fit usually does not fit in more than one way, and naming one
// sends the reader back for the next after each repair. It is fed through the
// document viewport so a file with many faults can be scrolled rather than
// clipped.
func renderSpecReport(name string, problems []specProblem, width int) string {
	pad := cardPadding(width)
	inner := max(1, width-2*pad)
	left := strings.Repeat(" ", pad)

	var b strings.Builder
	b.WriteString("\n")
	for _, line := range strings.Split(
		ansi.Wrap("specgetty cannot read "+name+" as a spec", inner, " "), "\n") {
		b.WriteString(left + sectionHeaderStyle.Render(line) + "\n")
	}
	b.WriteString("\n")
	for _, line := range strings.Split(ansi.Wrap(
		"These are the rules openspec itself reads a spec by, so validate, "+
			"list and archive cannot see this file either.", inner, " "), "\n") {
		b.WriteString(left + dimStyle.Render(line) + "\n")
	}

	for _, p := range problems {
		b.WriteString("\n")
		if p.line > 0 {
			b.WriteString(left + sectionHeaderStyle.Render(fmt.Sprintf("line %d", p.line)) + "\n")
		}
		for _, line := range strings.Split(ansi.Wrap(renderInlineMarkdown(p.text), inner, " "), "\n") {
			b.WriteString(left + line + "\n")
		}
	}

	b.WriteString("\n")
	for _, line := range strings.Split(ansi.Wrap(
		"E opens this file in your editor. esc returns to the specs tab, where "+
			"the whole file is readable as markdown.", inner, " "), "\n") {
		b.WriteString(left + dimStyle.Render(line) + "\n")
	}
	return b.String()
}

// specStructured reports whether the open spec could be read as one.
func (m model) specStructured() bool {
	return len(m.specTree.nodes) > 0
}

// renderSpecDetail draws the outline and the card, each in its own border, by
// the same split the specs and properties tabs use.
func (m model) renderSpecDetail(width, height int) string {
	var b strings.Builder
	// The spec's name stays put above the two halves, so the view is never
	// anonymous. It names the content, which is what keeps it outside the
	// borders.
	b.WriteString(headerStyle.Render(m.specName))
	b.WriteString("\n")

	boxHeight := height - 1
	if boxHeight < boxRows+1 {
		boxHeight = boxRows + 1
	}
	rows := boxHeight - boxRows
	if rows < 1 {
		rows = 1
	}

	if !m.specStructured() {
		// One panel, not two: there is no outline to put beside anything.
		body := truncateContent(m.docViewport.View(), rows)
		var fit []string
		for _, line := range strings.Split(body, "\n") {
			fit = append(fit, ansi.Truncate(line, max(1, width-boxChrome), ""))
		}
		b.WriteString(contentBox(width, boxHeight, true, strings.Join(fit, "\n")))
		return b.String()
	}

	outlineOuter, cardOuter := specDetailSplit(width)

	// The chooser sits above the card, inside its half, the way a change's
	// artifact sub-tabs sit above the document they choose between. It costs
	// the card a row, which docRegion has already taken off.
	cardRows, viewRow := rows, ""
	if m.cardViewRowShown() {
		viewRow = renderCardViewRow(m.cardView)
		cardRows--
		if cardRows < 1 {
			cardRows = 1
		}
	}

	outline := renderSpecOutline(m.specTree, m.selectedSpecNode(),
		outlineOuter-boxChrome, rows, m.focus == focusListPane)
	// Truncated to the pane as well as to the row count: a resize renders once
	// with a viewport still sized for the width before it, and a card wider
	// than its box would push the frame past the terminal.
	card := truncateContent(m.docViewport.View(), cardRows)
	var fitted []string
	for _, line := range strings.Split(card, "\n") {
		fitted = append(fitted, ansi.Truncate(line, max(1, cardOuter-boxChrome), ""))
	}
	card = strings.Join(fitted, "\n")
	if viewRow != "" {
		card = ansi.Truncate(viewRow, max(1, cardOuter-boxChrome), "") + "\n" + card
	}

	left := contentBox(outlineOuter, boxHeight, m.focus == focusListPane, outline)
	right := contentBox(cardOuter, boxHeight, m.focus == focusContentPane, card)
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right))
	return b.String()
}

// selectedSpecNode clamps the cursor to the tree, so a re-parse that shortened
// it cannot land past the end.
func (m model) selectedSpecNode() int {
	return clampIndex(m.specNode, len(m.specTree.nodes))
}

// syncSpecNode keeps the cursor on the node it was on across a re-parse.
//
// A rescan rewrites the tree whenever the file changes, and an edit above the
// cursor would otherwise move it. The path is what survives; when the node it
// named is gone, the index is clamped instead.
func (m *model) syncSpecNode() {
	if !m.specLevel() || len(m.specTree.nodes) == 0 {
		return
	}
	if m.specNodePath != "" {
		if i := m.specTree.indexOfPath(m.specNodePath); i >= 0 {
			m.specNode = i
			return
		}
	}
	m.specNode = clampIndex(m.specNode, len(m.specTree.nodes))
	m.specNodePath = m.specTree.nodes[m.specNode].path
}

// rememberSpecNode records which node the cursor is on, by path.
func (m *model) rememberSpecNode() {
	if !m.specLevel() || len(m.specTree.nodes) == 0 {
		return
	}
	was := m.specNodePath
	m.specNodePath = m.specTree.nodes[m.selectedSpecNode()].path
	if m.specNodePath != was {
		// Each node is read from its difference first. Carrying the choice
		// forward would show the reader an original they did not ask for.
		m.cardView = viewDiff
	}
}

// reparseOpenSpec rebuilds the tree after a rescan, keeping the cursor.
func (m *model) reparseOpenSpec() {
	if m.level == levelChangeSpec {
		m.reparseOpenChangeSpecs()
		return
	}
	if m.level != levelSpec {
		return
	}
	key := m.currentKey()
	if key == "" {
		return
	}
	info := m.projects[key].Info
	content, ok := info.SpecContents[m.specName]
	if !ok {
		return
	}

	tree, problems := parseSpec(m.specName, content)
	if len(problems) == 0 {
		// Including from the report state: a file repaired in an editor
		// structures itself under the reader without them leaving the view.
		m.specTree = tree
		m.specProblems = nil
		m.syncSpecNode()
		return
	}
	if len(m.specTree.nodes) > 0 {
		// A spec edited into a shape this view cannot structure keeps the tree
		// it opened with rather than emptying the view mid-read.
		return
	}
	m.specProblems = problems
}
