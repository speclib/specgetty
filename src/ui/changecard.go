package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// A node of a change's delta that modifies something has three cards, not one:
// what the change proposes, what it modifies, and the difference between them.
// The difference is what the reader came for, so it is what opens.

// Which of a comparable node's three views the card is showing.
const (
	viewDiff = iota
	viewOld
	viewNew
)

var cardViewNames = []string{"diff", "old", "new"}

// renderCardViewRow draws the three-way chooser above the card.
//
// The change view one level up draws its artifact sub-tabs this way and moves
// them with the same two keys, so the row is the idiom already learned rather
// than a new one.
func renderCardViewRow(view int) string {
	var b strings.Builder
	for i, name := range cardViewNames {
		if i > 0 {
			b.WriteString(" ")
		}
		if i == view {
			b.WriteString(activeTabStyle.Render(name))
		} else {
			b.WriteString(inactiveTabStyle.Render(name))
		}
	}
	return b.String()
}

// renderChangeCard draws the node the cursor is on, in the view asked for.
//
// A node with nothing to compare against ignores the view and draws itself,
// which is how the absence of the row and the absence of a choice stay the same
// fact.
func renderChangeCard(n specNode, view, width int) string {
	if n.kind == nodeCapability {
		return renderCapabilityCard(n, width)
	}
	if !n.comparable() {
		return renderSpecCard(n, width)
	}
	switch view {
	case viewOld:
		return renderSpecCard(*n.old, width)
	case viewNew:
		return renderSpecCard(n, width)
	default:
		return renderDiffCard(n, width)
	}
}

// renderCapabilityCard names the capability and what the change does to it.
func renderCapabilityCard(n specNode, width int) string {
	pad := cardPadding(width)
	inner := max(1, width-2*pad)
	left := strings.Repeat(" ", pad)

	var b strings.Builder
	b.WriteString("\n")
	for _, line := range strings.Split(ansi.Wrap(n.title, inner, " "), "\n") {
		b.WriteString(left + capabilityStyle.Render(line) + "\n")
	}
	b.WriteString("\n")
	writeProse(&b, n.body, left, inner)
	return b.String()
}

// renderDiffCard shows what the change alters, word by word.
//
// Word level rather than line level: a delta restates a requirement and the
// restatement is wrapped afresh, so a line diff would mark every line of a
// paragraph changed to report one altered clause.
func renderDiffCard(n specNode, width int) string {
	pad := cardPadding(width)
	inner := max(1, width-2*pad)
	left := strings.Repeat(" ", pad)

	var b strings.Builder
	b.WriteString("\n")

	title := "Requirement: " + n.title
	if n.kind == nodeScenario {
		title = "Scenario: " + n.title
	}
	for _, line := range strings.Split(ansi.Wrap(title, inner, " "), "\n") {
		b.WriteString(left + sectionHeaderStyle.Render(line) + "\n")
	}

	if n.kind == nodeScenario {
		writeClauseDiff(&b, n.old.parts, n.parts, left, inner)
		return b.String()
	}
	b.WriteString("\n")
	writeDiff(&b, n.old.body, n.body, left, inner)
	return b.String()
}

// writeClauseDiff lays two scenarios out as one, part by part.
//
// Part by part rather than as two flattened texts, so the keyword still sits on
// its own row with the clause indented under it. A reader moving between the
// difference, the original and the new text sees the same shape in all three,
// and only the marks move.
func writeClauseDiff(b *strings.Builder, old, new []specPart, left string, inner int) {
	n := len(old)
	if len(new) > n {
		n = len(new)
	}
	clauseIndent := left + "   "
	for i := 0; i < n; i++ {
		b.WriteString("\n")
		switch {
		case i >= len(new):
			// A part the change drops entirely.
			writeClausePart(b, old[i], removedWordStyle, left, clauseIndent, inner,
				func(_, o string) string { return markRun(removedWordStyle, splitWords(o)) })
		case i >= len(old):
			writeClausePart(b, new[i], addedWordStyle, left, clauseIndent, inner,
				func(_, o string) string { return markRun(addedWordStyle, splitWords(o)) })
		default:
			style := specKeywordStyle
			if old[i].keyword != new[i].keyword {
				style = addedWordStyle
			}
			writeClausePart(b, new[i], style, left, clauseIndent, inner,
				func(_, _ string) string { return diffWords(old[i].text, new[i].text) })
		}
	}
}

// writeClausePart draws one part: its keyword on a row of its own where it has
// one, and the text beneath it at a hanging indent.
func writeClausePart(b *strings.Builder, p specPart, kwStyle lipgloss.Style,
	left, clauseIndent string, inner int, text func(keyword, body string) string) {

	body := text(p.keyword, p.text)
	indent, width := left, inner
	if p.kind == partClause {
		b.WriteString(left + kwStyle.Render(p.keyword) + "\n")
		indent, width = clauseIndent, max(1, inner-3)
	}
	for _, line := range strings.Split(ansi.Wrap(body, width, " "), "\n") {
		b.WriteString(indent + line + "\n")
	}
}

// cardText flattens a scenario into the text a diff compares, keeping each
// clause's keyword with the clause so a moved keyword reads as a change.
func cardText(n specNode) string {
	var parts []string
	for _, p := range n.parts {
		if p.kind == partClause {
			parts = append(parts, p.keyword+" "+p.text)
			continue
		}
		parts = append(parts, p.text)
	}
	return strings.Join(parts, "\n\n")
}

// writeDiff lays two texts out as one, marking what differs.
//
// Removed words are struck through in the removed colour and added words are
// drawn in the added colour. Everything untouched is drawn as ordinary prose,
// which is most of it, and is the point.
func writeDiff(b *strings.Builder, old, new string, left string, inner int) {
	for _, para := range diffParagraphs(old, new) {
		if para == "" {
			continue
		}
		b.WriteString("\n")
		for _, line := range strings.Split(ansi.Wrap(para, inner, " "), "\n") {
			b.WriteString(left + line + "\n")
		}
	}
}

// diffParagraphs renders each paragraph of the new text with what changed in it
// marked, and any paragraph dropped entirely shown as removed.
func diffParagraphs(old, new string) []string {
	oldParas := splitParagraphs(old)
	newParas := splitParagraphs(new)

	// Paired by position, which is what a restated requirement does: it keeps
	// the shape and edits inside it. A paragraph with no counterpart is wholly
	// added or wholly removed, and says so.
	var out []string
	n := len(oldParas)
	if len(newParas) > n {
		n = len(newParas)
	}
	for i := 0; i < n; i++ {
		switch {
		case i >= len(newParas):
			out = append(out, markRun(removedWordStyle, splitWords(oldParas[i])))
		case i >= len(oldParas):
			out = append(out, markRun(addedWordStyle, splitWords(newParas[i])))
		default:
			out = append(out, diffWords(oldParas[i], newParas[i]))
		}
	}
	return out
}

func splitParagraphs(s string) []string {
	var out []string
	for _, p := range strings.Split(s, "\n\n") {
		if j := strings.Join(strings.Fields(p), " "); j != "" {
			out = append(out, j)
		}
	}
	return out
}

func splitWords(s string) []string { return strings.Fields(s) }

func markRun(style lipgloss.Style, words []string) string {
	if len(words) == 0 {
		return ""
	}
	return style.Render(strings.Join(words, " "))
}

// diffWords marks the words of new that old does not have, and shows the words
// of old that new dropped.
//
// A longest-common-subsequence walk over words. The texts being compared are a
// paragraph each, so the quadratic table is a few thousand cells at worst.
func diffWords(old, new string) string {
	a, c := splitWords(old), splitWords(new)

	// lcs[i][j] is the length of the longest common subsequence of a[i:] and
	// c[j:], which is what the walk below follows.
	lcs := make([][]int, len(a)+1)
	for i := range lcs {
		lcs[i] = make([]int, len(c)+1)
	}
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(c) - 1; j >= 0; j-- {
			if a[i] == c[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
				continue
			}
			lcs[i][j] = max(lcs[i+1][j], lcs[i][j+1])
		}
	}

	var b strings.Builder
	var same, removed, added []string
	flush := func() {
		if s := markRun(removedWordStyle, removed); s != "" {
			b.WriteString(s + " ")
		}
		if s := markRun(addedWordStyle, added); s != "" {
			b.WriteString(s + " ")
		}
		removed, added = nil, nil
	}
	flushSame := func() {
		if len(same) > 0 {
			b.WriteString(renderInlineMarkdown(strings.Join(same, " ")) + " ")
			same = nil
		}
	}

	i, j := 0, 0
	for i < len(a) || j < len(c) {
		switch {
		case i < len(a) && j < len(c) && a[i] == c[j]:
			flush()
			same = append(same, a[i])
			i, j = i+1, j+1
		case j < len(c) && (i == len(a) || lcs[i][j+1] >= lcs[i+1][j]):
			flushSame()
			added = append(added, c[j])
			j++
		default:
			flushSame()
			removed = append(removed, a[i])
			i++
		}
	}
	flushSame()
	flush()
	return strings.TrimSpace(b.String())
}

// cardViewRowShown reports whether the chooser is on screen, which costs the
// card a row.
//
// Only a node with an original has one, so the row appears and disappears as
// the cursor moves. That is how the view says there is nothing to compare,
// without a line saying so on every other node.
func (m model) cardViewRowShown() bool {
	if m.level != levelChangeSpec || !m.specStructured() {
		return false
	}
	n, ok := m.selectedNode()
	return ok && n.comparable()
}
