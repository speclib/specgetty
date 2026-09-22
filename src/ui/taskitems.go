package ui

// The tasks pane is the fourth list in this application whose cursor indexes
// items while the pane counts drawn rows, and the only one that was written
// before lines.go existed. It selected a source line, which is not what a task
// is: a task wraps onto as many lines as it needs and indents the rest, so the
// cursor stepped through the remainder of a task, the highlight covered the
// first line of one, and a page key moved rows while the cursor stayed behind.
//
// Grouping the rendered lines into items puts it on the same footing as the
// change list, the spec outline and the properties list. itemLines does the
// arithmetic for all four.

// taskItem is one task the cursor can select.
//
// Only the checkbox line is carried. It is what a save matches on, and it is
// the only line of the item a toggle rewrites; the continuation lines are
// drawn, highlighted and scrolled with it but never edited.
type taskItem struct {
	text  string // the checkbox line verbatim
	index int    // its position in the source, counting from zero
}

// taskItems is the cursor's view of a rendered tasks document: the tasks it can
// select, and which drawn row belongs to which of them.
type taskItems struct {
	items []taskItem
	rows  itemLines
}

// count is how many tasks the cursor can select.
func (t taskItems) count() int { return len(t.items) }

// at returns the task at index i, if there is one.
func (t taskItems) at(i int) (taskItem, bool) {
	if i < 0 || i >= len(t.items) {
		return taskItem{}, false
	}
	return t.items[i], true
}

// taskItemsIn groups the source lines of a rendered tasks document into the
// items the cursor selects, and marks which drawn row belongs to which.
//
// An item is a checkbox line at column zero together with the indented lines
// that follow it, by continuesTaskItem's rule. Every other row is chrome: a
// heading, a blank line, or the `Tasks: n/m complete` prefix the artifact
// renderer puts above the document. Chrome is marked lineOwner, which is how
// the other three lists already mark theirs.
//
// totalRows is the number of rows the rendered document actually has, rather
// than the number the mapping accounts for. They differ by the prefix, and
// sizing the slice to the mapping would leave the last rows of a document
// unowned by anything, including by chrome.
func taskItemsIn(lines []sourceLine, totalRows int) taskItems {
	if totalRows < 0 {
		totalRows = 0
	}
	t := taskItems{rows: make(itemLines, totalRows)}
	for i := range t.rows {
		t.rows[i] = lineOwner
	}

	for i := 0; i < len(lines); i++ {
		if !isTaskLine(lines[i].text) {
			continue
		}
		owner := len(t.items)
		t.items = append(t.items, taskItem{
			text:  lines[i].text,
			index: lines[i].index,
		})

		last := i
		for last+1 < len(lines) && continuesTaskItem(lines[last+1].text) {
			last++
		}
		for r := lines[i].rowStart; r <= lines[last].rowEnd; r++ {
			if r >= 0 && r < len(t.rows) {
				t.rows[r] = owner
			}
		}
		i = last
	}
	return t
}
