package ui

// Every list in the application answers to the same keys as every document:
// `pgdown`, `pgup`, `ctrl+f`, `ctrl+b`, `ctrl+d`, `ctrl+u`, `gg` and `G`. Which
// surface they move is decided the way `j` and `k` already decide it, by where
// the keyboard is.
//
// The alternative, letting the paging keys ignore focus and always reach the
// document, shipped for one commit. It cannot survive a pageable list: on a
// split tab there is then a list and a document in play, and a key that ignores
// focus has nothing left to choose between them with.

// listPage is how far a page moves the list that holds the keyboard.
//
// It is the number of rows that list is showing, so a page lands about where
// the eye would and the key means the same thing at any terminal height. The
// expressions mirror what each renderer lays out; they are one fact written
// twice, and the tests assert the pairing rather than the arithmetic.
func (m model) listPage() int {
	if m.specLevel() {
		if !m.specStructured() {
			// A report has no outline to page. Its own scrolling goes through
			// the viewport, which docActive() routes to before this is asked.
			return 1
		}
		// A page of the outline is a pane of rows, and a node can occupy more
		// than one, so the page is however many nodes fit in that many rows.
		// With every node one row tall this reduces to the row count.
		rows := max(1, m.mainPanelHeight()-1-boxRows)
		outlineOuter, _ := specDetailSplit(m.panelContentWidth())
		return max(1, nodesInRows(m.specTree, outlineOuter-boxChrome, m.selectedSpecNode(), rows))
	}
	if m.level != levelProject {
		return 1
	}
	region := m.mainPanelHeight() - 5 // the project header and the tab bar

	switch m.detailTab {
	case tabChanges:
		// Inside the content box, less the table's own column header, less the
		// search prompt when one is on screen.
		rows := region - boxRows - 1
		if m.searchFocused || m.searchInput.Value() != "" {
			rows--
		}
		return max(1, rows)
	case tabSpecs, tabProperties:
		return max(1, region-boxRows)
	}
	return 1
}

// moveListCursor moves whichever list holds the keyboard by delta rows,
// stopping at either end.
func (m *model) moveListCursor(delta int) {
	if delta == 0 {
		return
	}
	if m.specLevel() {
		if m.focus != focusListPane {
			return
		}
		m.specNode = clampIndex(m.specNode+delta, len(m.specTree.nodes))
		m.rememberSpecNode()
		return
	}
	if m.level != levelProject {
		return
	}
	switch m.detailTab {
	case tabChanges:
		m.setChangeCursor(m.changeCursor + delta)
	case tabSpecs:
		if m.focus != focusListPane {
			return
		}
		m.specCursor = clampIndex(m.specCursor+delta, len(m.currentSpecNames()))
	case tabProperties:
		if m.focus != focusListPane {
			return
		}
		m.propSection = clampIndex(m.propSection+delta, len(m.currentSections()))
	}
}

// gotoListEnd sends the list that holds the keyboard to its first or last row.
func (m *model) gotoListEnd(last bool) {
	if m.specLevel() {
		if m.focus != focusListPane {
			return
		}
		to := 0
		if last {
			to = len(m.specTree.nodes) - 1
		}
		m.specNode = clampIndex(to, len(m.specTree.nodes))
		m.rememberSpecNode()
		return
	}
	if m.level != levelProject {
		return
	}
	end := func(n int) int {
		if last {
			return n - 1
		}
		return 0
	}
	switch m.detailTab {
	case tabChanges:
		m.setChangeCursor(end(len(m.currentRows())))
	case tabSpecs:
		if m.focus != focusListPane {
			return
		}
		m.specCursor = clampIndex(end(len(m.currentSpecNames())), len(m.currentSpecNames()))
	case tabProperties:
		if m.focus != focusListPane {
			return
		}
		m.propSection = clampIndex(end(len(m.currentSections())), len(m.currentSections()))
	}
}

// setChangeCursor moves the change list and remembers the selection, which is
// what keeps it across a rescan or a filter. A single-row move already does
// this, and a page move that did not would quietly lose the selection.
func (m *model) setChangeCursor(to int) {
	rows := len(m.currentRows())
	next := clampIndex(to, rows)
	if next == m.changeCursor {
		return
	}
	m.changeCursor = next
	m.changeArtifactTab = 0
	m.rememberSelection()
}

// pickerPage is how far a page moves the project picker, which sizes its own
// box rather than living in the detail panel.
func (m model) pickerPage() int {
	_, height := m.pickerBox()
	return max(1, height-4)
}

// movePickerCursor moves the picker by delta rows, stopping at either end.
func (m *model) movePickerCursor(delta int) {
	rows := len(m.pickerVisibleRows())
	if rows == 0 {
		return
	}
	m.pickerCursor = clampIndex(m.pickerCursor+delta, rows)
	m.pickerRemember()
}

// clampIndex keeps an index inside a list of n rows, or at zero when there are
// none. Paging past an end stops there rather than wrapping.
func clampIndex(i, n int) int {
	if n <= 0 {
		return 0
	}
	if i < 0 {
		return 0
	}
	if i >= n {
		return n - 1
	}
	return i
}

// nodesInRows counts how many nodes starting at `from` fit in `rows` drawn
// rows, which is what a page of the outline moves.
func nodesInRows(tree specTree, width, from, rows int) int {
	return ownersOfRows(outlineRows(tree, width)).fit(from, rows)
}
