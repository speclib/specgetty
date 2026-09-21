package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// The spec detail view is an outline beside a card: the outline says what the
// spec contains, the card shows whichever node the cursor is on. It is the
// third navigation level, and it follows the same rule as the change view:
// `enter` descends, `esc` ascends.

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

// outlineRows lays the tree out at a given width.
func outlineRows(tree specTree, width int) []outlineRow {
	var rows []outlineRow
	for i, n := range tree.nodes {
		indent, hang := "", "  "
		if n.kind == nodeScenario {
			indent, hang = "  ", "    "
		}
		// Wrapped against the hanging indent, which is the wider of the two: a
		// continuation row that only got its indent after being wrapped could
		// come out wider than the pane, and lipgloss would then wrap it again
		// into a row the row count did not know about.
		wrapped := strings.Split(ansi.Wrap(n.title, max(1, width-len(hang)), " "), "\n")
		for j, line := range wrapped {
			prefix := indent
			if j > 0 {
				prefix = hang
			}
			rows = append(rows, outlineRow{text: prefix + strings.TrimSpace(line), node: i})
		}
	}
	return rows
}

// rowRangeOfNode returns the first and last drawn row a node occupies.
func rowRangeOfNode(rows []outlineRow, node int) (first, last int) {
	first, last = -1, -1
	for i, r := range rows {
		if r.node != node {
			continue
		}
		if first < 0 {
			first = i
		}
		last = i
	}
	return
}

// renderSpecOutline draws the outline, scrolled so the whole selected node is
// visible.
func renderSpecOutline(tree specTree, selected, width, height int, lit bool) string {
	rows := outlineRows(tree, width)
	first, last := rowRangeOfNode(rows, selected)

	offset := 0
	if last >= height {
		offset = last - height + 1
	}
	// A node taller than the pane shows its top rather than its bottom.
	if first >= 0 && first < offset {
		offset = first
	}
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
		case tree.nodes[r.node].kind == nodeRequirement:
			b.WriteString(sectionHeaderStyle.Render(r.text))
		default:
			b.WriteString(normalStyle.Render(r.text))
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
		for _, c := range n.clauses {
			b.WriteString("\n")
			b.WriteString(left + specKeywordStyle.Render(c.keyword) + "\n")
			// A hanging indent, so every row after the first still reads as
			// belonging to the keyword above it.
			clauseIndent := left + "   "
			for _, line := range strings.Split(ansi.Wrap(c.text, max(1, inner-3), " "), "\n") {
				b.WriteString(clauseIndent + renderInlineMarkdown(line) + "\n")
			}
		}
	}
	return b.String()
}

// writeProse lays a node's body out, one paragraph at a time.
func writeProse(b *strings.Builder, body, left string, inner int) {
	for _, para := range strings.Split(body, "\n\n") {
		joined := strings.Join(strings.Fields(strings.ReplaceAll(para, "\n", " ")), " ")
		if joined == "" {
			continue
		}
		for _, line := range strings.Split(ansi.Wrap(joined, inner, " "), "\n") {
			b.WriteString(left + renderInlineMarkdown(line) + "\n")
		}
		b.WriteString("\n")
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

	tree, err := parseSpec(name, info.SpecContents[name])
	if err != nil {
		m.statusMsg = "Cannot open " + name + ": " + err.Error()
		return m
	}

	m.specTree = tree
	m.specNode = 0
	m.specNodePath = tree.nodes[0].path
	m.level = levelSpec
	m.focus = focusListPane
	m.statusMsg = ""
	return m
}

// renderSpecDetail draws the outline and the card, each in its own border, by
// the same split the specs and properties tabs use.
func (m model) renderSpecDetail(width, height int) string {
	var b strings.Builder
	// The spec's name stays put above the two halves, so the view is never
	// anonymous. It names the content, which is what keeps it outside the
	// borders.
	b.WriteString(headerStyle.Render(m.specTree.name))
	b.WriteString("\n")

	boxHeight := height - 1
	if boxHeight < boxRows+1 {
		boxHeight = boxRows + 1
	}
	rows := boxHeight - boxRows
	if rows < 1 {
		rows = 1
	}

	outlineOuter, cardOuter := specDetailSplit(width)
	outline := renderSpecOutline(m.specTree, m.selectedSpecNode(),
		outlineOuter-boxChrome, rows, m.focus == focusListPane)
	// Truncated to the pane as well as to the row count: a resize renders once
	// with a viewport still sized for the width before it, and a card wider
	// than its box would push the frame past the terminal.
	card := truncateContent(m.docViewport.View(), rows)
	var fitted []string
	for _, line := range strings.Split(card, "\n") {
		fitted = append(fitted, ansi.Truncate(line, max(1, cardOuter-boxChrome), ""))
	}
	card = strings.Join(fitted, "\n")

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
	if m.level != levelSpec || len(m.specTree.nodes) == 0 {
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
	if m.level == levelSpec && len(m.specTree.nodes) > 0 {
		m.specNodePath = m.specTree.nodes[m.selectedSpecNode()].path
	}
}

// reparseOpenSpec rebuilds the tree after a rescan, keeping the cursor.
func (m *model) reparseOpenSpec() {
	if m.level != levelSpec {
		return
	}
	key := m.currentKey()
	if key == "" {
		return
	}
	info := m.projects[key].Info
	content, ok := info.SpecContents[m.specTree.name]
	if !ok {
		return
	}
	tree, err := parseSpec(m.specTree.name, content)
	if err != nil {
		// A spec edited into a shape this view cannot navigate keeps the tree
		// it opened with rather than emptying the view mid-read.
		return
	}
	m.specTree = tree
	m.syncSpecNode()
}
