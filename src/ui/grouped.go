package ui

import (
	"fmt"
	"strings"
)

// The change list is one table with headers between its groups. The columns are
// shared, because the archive date renders blank for an active change, so what
// differs between the groups is which rows they hold and nothing else.
//
// This is the change list's own renderer rather than a capability of
// renderTable, which the project picker shares and has no use for headers. It
// reuses layoutFields and fitCell, so the columns cannot drift between the two.

// groupLine is one drawn line: either a header or a change.
//
// The cursor never indexes these. It indexes changes, because the selection
// key, the search filter, `g`, `G` and all three actions already work that way
// and none of them should have to know a header exists. Only the arithmetic
// below maps between the two.
// What a drawn line is. A change, a group header, or the blank line that
// separates one group from the one above it.
const (
	lineChange = iota
	lineHeader
	lineSpacer
)

type groupLine struct {
	kind   int
	header string
	row    filtered[changeRow]
	// rowIndex is the change's position in the filtered list, and is meaningful
	// only for lineChange. This is the whole mapping, in one field.
	rowIndex int
}

// groupLines lays the filtered rows out under their headers.
//
// A group keeps its header and a count even when it holds nothing: "nothing in
// flight" is an answer, and an absent header would be indistinguishable from a
// filter having hidden it.
func groupLines(groups []changeGroup, rows []filtered[changeRow]) []groupLine {
	// Which group each surviving row belongs to, by the key it already carries.
	// The filter may have dropped rows, so membership is looked up rather than
	// assumed from position.
	groupOf := make(map[string]int, len(rows))
	for gi, g := range groups {
		for _, r := range g.rows {
			groupOf[r.key()] = gi
		}
	}

	buckets := make([][]int, len(groups))
	for i, r := range rows {
		gi, ok := groupOf[r.row.key()]
		if !ok {
			continue
		}
		buckets[gi] = append(buckets[gi], i)
	}

	var lines []groupLine
	for gi, g := range groups {
		// Every group but the first is preceded by a blank line, whether or not
		// either group has rows in it. A rule with no exception is one less
		// thing on screen to explain.
		if gi > 0 {
			lines = append(lines, groupLine{kind: lineSpacer, rowIndex: -1})
		}
		lines = append(lines, groupLine{
			kind:     lineHeader,
			header:   fmt.Sprintf("%s (%d)", g.label, len(buckets[gi])),
			rowIndex: -1,
		})
		for _, i := range buckets[gi] {
			lines = append(lines, groupLine{kind: lineChange, row: rows[i], rowIndex: i})
		}
	}
	return lines
}

// lineOfRow returns the drawn line a change sits on, or -1 when the filter has
// removed it.
func lineOfRow(lines []groupLine, rowIndex int) int {
	for i, l := range lines {
		if l.kind == lineChange && l.rowIndex == rowIndex {
			return i
		}
	}
	return -1
}

// renderGroupedTable draws the column header, then the groups with their
// headers and rows, scrolled so the selected change stays visible.
func renderGroupedTable(groups []changeGroup, rows []filtered[changeRow],
	defs []fieldDef[changeRow], cursor, width, height int) string {

	hintWidth := 0
	for _, r := range rows {
		if len(r.matched) > 0 {
			hintWidth = min(28, width/3)
			break
		}
	}
	// The hint is a column like any other, so it is separated like one. It used
	// to be appended flush, and looked right only because the last default
	// column held values shorter than its width; the archive date fills its
	// column exactly.
	if hintWidth > 0 {
		hintWidth += columnGap
	}
	tableWidth := width - hintWidth
	if tableWidth < 1 {
		tableWidth = width
		hintWidth = 0
	}

	kept, widths := layoutFields(defs, tableWidth)

	var b strings.Builder
	headerCells := make([]string, len(kept))
	for i, d := range kept {
		headerCells[i] = fitCell(d.header, widths[i])
	}
	header := fitCell(strings.Join(headerCells, columnSep), tableWidth)
	if hintWidth > 0 {
		header += columnSep + fitCell("matched", hintWidth-columnGap)
	}
	b.WriteString(dimStyle.Render(header))

	bodyHeight := height - 1
	if bodyHeight < 1 {
		return b.String()
	}

	lines := groupLines(groups, rows)

	// The offset is computed in drawn lines, not in change indices, because a
	// header occupies a line the cursor cannot land on. This is where every
	// off-by-one in grouping would live.
	offset := 0
	if at := lineOfRow(lines, cursor); at >= bodyHeight {
		offset = at - bodyHeight + 1
	}
	// The pane never opens on a blank line. The offset puts the selected change
	// on the last visible row, so moving the window down by one keeps it in
	// range while dropping the spacer that would otherwise head the pane.
	if offset < len(lines) && lines[offset].kind == lineSpacer {
		offset++
	}
	end := offset + bodyHeight
	if end > len(lines) {
		end = len(lines)
	}

	for i := offset; i < end; i++ {
		b.WriteString("\n")
		l := lines[i]

		switch l.kind {
		case lineSpacer:
			// Nothing at all: no style, so no background is painted across the
			// width of a line that is meant to read as a gap.
			continue
		case lineHeader:
			b.WriteString(sectionHeaderStyle.Render(fitCell(l.header, width)))
			continue
		}

		cells := make([]string, len(kept))
		for j, d := range kept {
			cells[j] = fitCell(d.value(l.row.row), widths[j])
		}
		base := fitCell(strings.Join(cells, columnSep), tableWidth)

		hint := ""
		if hintWidth > 0 {
			hint = columnSep + fitCell(strings.Join(l.row.matched, ", "), hintWidth-columnGap)
		}

		if l.rowIndex == cursor {
			b.WriteString(selectedStyle.Width(width).Render(fitCell(base+hint, width)))
		} else {
			b.WriteString(normalStyle.Render(base + hint))
		}
	}

	return b.String()
}
