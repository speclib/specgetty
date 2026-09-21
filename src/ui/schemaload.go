package ui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/mipmip/specgetty/src/scanner"
)

// Locating a schema's definition costs a subprocess of roughly a second, which
// is why it happens at most once per schema per project opened, only once the
// properties tab has been asked for, and concurrently when there is more than
// one schema to read.
//
// Nothing re-reads a schema while a project is open. The watcher covers the
// root's openspec/ tree and so already fires when a schema file changes, but a
// schema is a workflow definition rather than content: it changes once in a
// project's life where changes and specs move all day. Opening the project
// again is how an edit is seen.

// schemaState is what is known about one schema's definition. Both fields nil
// means it is being read.
type schemaState struct {
	details *scanner.SchemaDetails
	problem *scanner.SchemaProblem
}

// loaded reports whether the read has finished, either way.
func (s schemaState) loaded() bool { return s.details != nil || s.problem != nil }

// schemaMsg carries one finished read back to the model.
//
// The project is carried so that an answer arriving after the user has switched
// project is dropped rather than shown against the wrong one.
type schemaMsg struct {
	project string
	name    string
	details *scanner.SchemaDetails
	problem *scanner.SchemaProblem
}

// loadSchema reads one schema's definition.
func loadSchema(project, root, name string) tea.Cmd {
	return func() tea.Msg {
		d, p := scanner.ResolveSchema(root, name)
		return schemaMsg{project: project, name: name, details: d, problem: p}
	}
}

// ensureSchemasLoaded starts a read for every schema the open project uses that
// has not been read yet, and returns the commands to run.
//
// Calling it again while reads are in flight starts nothing new, so opening the
// properties tab repeatedly costs one subprocess per schema and no more.
func (m *model) ensureSchemasLoaded() []tea.Cmd {
	key := m.currentKey()
	if key == "" {
		return nil
	}
	info := m.projects[key].Info

	// A different project owns different schemas, so what was read for the last
	// one is dropped rather than shown against this one.
	if m.schemaFor != key {
		m.schemaFor = key
		m.schemas = nil
	}
	if m.schemas == nil {
		m.schemas = make(map[string]schemaState)
	}

	var cmds []tea.Cmd
	for _, u := range info.SchemaUsage {
		if _, started := m.schemas[u.Name]; started {
			continue
		}
		// An entry with neither field set marks the read as in flight, which is
		// both what the pane shows and what stops a second one starting.
		m.schemas[u.Name] = schemaState{}
		cmds = append(cmds, loadSchema(key, info.Root, u.Name))
	}
	return cmds
}

// schemaStateOf returns what is known about one schema for the open project.
func (m model) schemaStateOf(name string) schemaState {
	if m.schemaFor != m.currentKey() {
		return schemaState{}
	}
	return m.schemas[name]
}
