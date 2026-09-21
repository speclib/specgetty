package ui

import (
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// The editor is the user's, not ours. `E` hands it the file the pane is showing
// and gets out of the way: specgetty reads files, and the one thing it writes
// is a task checkbox.

// editorCommand is the process to run, already split into a command and its
// arguments with the file last.
type editorCommand struct {
	name string
	args []string
}

// resolveEditor builds the command for a file from the environment.
//
// `$VISUAL` first, then `$EDITOR`. The convention is older than this program:
// `EDITOR` names a line editor that works without a full terminal and `VISUAL`
// the full-screen one, so a user with both set has already said which they want
// for an interactive edit. Reading only `EDITOR`, the shorter thing to write,
// would open `ed` for someone whose environment is correct.
//
// Nothing is guessed when neither is set. A fallback to `vi` is wrong on a
// machine that has no `vi`, and when it is wrong the failure arrives as a
// process that would not start rather than as an answer.
func resolveEditor(file string) (editorCommand, string) {
	value := firstSet("VISUAL", "EDITOR")
	if value == "" {
		return editorCommand{}, "no editor configured: set $VISUAL or $EDITOR"
	}

	// Split on whitespace, with no shell involved. `EDITOR="emacsclient -nw"`
	// and `EDITOR="code -w"` are both ordinary, and passing the whole string to
	// exec would look for a binary with a space in its name.
	//
	// A shell would handle quoting too, at the cost of running the file path
	// through shell parsing, where a path holding a space or a quote becomes an
	// injection question. This way the path is one argument that nothing
	// parses; the case it does not handle, a quoted path inside the variable,
	// is stated in the spec rather than left to be discovered.
	// firstSet has already rejected a value that is only whitespace, so there
	// is at least one field here.
	fields := strings.Fields(value)

	return editorCommand{
		name: fields[0],
		args: append(fields[1:], file),
	}, ""
}

// firstSet returns the first environment variable that holds something other
// than whitespace. A variable set to the empty string is unset: it is what a
// shell leaves behind when a value is cleared.
func firstSet(names ...string) string {
	for _, name := range names {
		if v := strings.TrimSpace(os.Getenv(name)); v != "" {
			return v
		}
	}
	return ""
}

// editorDoneMsg is delivered when the editor has exited and the application has
// the terminal back.
type editorDoneMsg struct{ err error }

// runEditor is the seam a test replaces, so that pressing `E` can be asserted
// on without a real process and without a real terminal to hand over.
//
// tea.ExecProcess is documented for exactly this case: it pauses the program,
// runs the command against the terminal, and resumes when it exits. Nothing
// here draws while that is true, which is what keeps the editor from being
// drawn over.
var runEditor = func(c *exec.Cmd, done func(error) tea.Msg) tea.Cmd {
	return tea.ExecProcess(c, done)
}

// openInEditor hands the file the pane is showing to the user's editor.
//
// A pane with no single file does nothing, which is the whole guard: the key
// asks the document what file it came from rather than asking which level and
// tab are active, so the panes that carry it follow from one fact instead of
// from a list each new pane has to be added to.
func (m model) openInEditor() (model, tea.Cmd) {
	if m.docPath == "" {
		return m, nil
	}

	cmd, reason := resolveEditor(m.docPath)
	if reason != "" {
		// One line on the nav bar, the way the copy keys report. A modal is the
		// wrong weight for something instant and harmless.
		m.statusMsg = reason
		return m, nil
	}

	m.statusMsg = ""
	return m, runEditor(exec.Command(cmd.name, cmd.args...), func(err error) tea.Msg {
		return editorDoneMsg{err: err}
	})
}

// editorFinished resumes after the editor has exited.
//
// The open project is read again whether or not the file changed. The
// filesystem watcher would usually notice the edit on its own, but it is
// allowed to fail to start, and then nothing on screen would move. The terminal
// may also have been resized while the editor had it; bubbletea reports the new
// size on resume, so the existing WindowSizeMsg path covers that.
func (m model) editorFinished(msg editorDoneMsg) (model, tea.Cmd) {
	if msg.err != nil {
		m.statusMsg = "editor: " + msg.err.Error()
	}
	if key := m.currentKey(); key != "" {
		return m, m.doScanSingle(key)
	}
	return m, nil
}
