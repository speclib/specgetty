package ui

// Three lists in this application share one shape: the cursor indexes items,
// the pane counts drawn lines, and an item may occupy more than one line.
//
//	the change list    cursor: changes   drawn: rows, group headers, spacers
//	the spec outline   cursor: nodes     drawn: wrapped label rows
//	the properties     cursor: sections  drawn: rows and group headers
//
// Two instances was a coincidence. Three is where the arithmetic should be one
// thing, and this is it: a slice with one entry per drawn line holding the item
// that line belongs to, and -1 where the line is chrome.
//
// The slice is deliberately all it carries. What a chrome line actually is, a
// header or a spacer, belongs to the renderer that drew it, and the change list
// needs to know that difference where nothing else does.

// lineOwner marks a drawn line that belongs to no item.
const lineOwner = -1

// itemLines is one entry per drawn line: the item it belongs to, or lineOwner.
type itemLines []int

// span returns the first and last drawn line an item occupies, or -1, -1 when
// the item is not drawn at all, which a filter can do.
func (l itemLines) span(item int) (first, last int) {
	first, last = -1, -1
	for i, owner := range l {
		if owner != item {
			continue
		}
		if first < 0 {
			first = i
		}
		last = i
	}
	return first, last
}

// fit reports how many items, counting from `from`, fill a pane of `rows` lines.
//
// At least one: a page that could move nothing because the next item is taller
// than the pane would leave the cursor stuck rather than moving it to the item
// it cannot fully show.
//
// With every item one line tall this is the row count, which is the property
// the change list's page key relies on.
func (l itemLines) fit(from, rows int) int {
	used, moved := 0, 0
	for item := from; ; item++ {
		first, last := l.span(item)
		if first < 0 {
			break
		}
		h := last - first + 1
		if used+h > rows && moved > 0 {
			break
		}
		used += h
		moved++
	}
	if moved < 1 {
		moved = 1
	}
	return moved
}

// offsetFor returns the scroll offset that brings the whole of an item into
// view in a pane of `height` lines.
//
// An item taller than the pane shows its top rather than its bottom: the first
// line is where reading starts, and a clause scrolled to its last row reads as
// a fragment with no beginning.
func (l itemLines) offsetFor(item, height int) int {
	first, last := l.span(item)
	if first < 0 || height < 1 {
		return 0
	}
	offset := 0
	if last >= height {
		offset = last - height + 1
	}
	if first < offset {
		offset = first
	}
	return offset
}

// ownedLines builds the slice for a list whose items each occupy one line, in
// order. It is what a renderer that draws no chrome would hand over.
func ownedLines(items int) itemLines {
	l := make(itemLines, items)
	for i := range l {
		l[i] = i
	}
	return l
}
