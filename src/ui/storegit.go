package ui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/mipmip/specgetty/src/scanner"
)

// A store's git state is read from its `.git/` directory, which lies outside
// every watched `openspec/` tree. A commit, fetch or checkout made in another
// terminal therefore changes it without any watched file changing, and until
// this it went stale until the project was scanned again.
//
// Watching `.git/` was the obvious symmetry with watching the registry and is a
// trap: a fetch rewrites refs in bulk and index operations touch files many
// times a second, so it would turn a quiet project into a rescan loop, and a
// rescan re-reads every spec and change.
//
// So it is read on entering the properties tab instead. That is the rule
// `schema-inspection` already states for schema definitions: nothing is read
// until the tab is asked for, and a session that never opens it never pays.

// storeGitMsg carries a re-read of the open store's git state.
type storeGitMsg struct {
	project string
	git     scanner.StoreGit
}

// refreshStoreGit re-reads the open project's store git state, if it has one.
//
// Nothing is read for a project holding its own content: there is no store to
// report on, and the properties tab says so without asking git anything.
func (m *model) refreshStoreGit() []tea.Cmd {
	key := m.currentKey()
	if key == "" {
		return nil
	}
	info := m.projects[key].Info
	if info.Store == nil {
		return nil
	}
	root := info.Store.Root
	return []tea.Cmd{func() tea.Msg {
		// Local refs only, with no fetch and no network access, which is what
		// ReadStoreGit does and what config-tab-display requires of it.
		return storeGitMsg{project: key, git: scanner.ReadStoreGit(root)}
	}}
}

// applyStoreGit records a re-read against the project it was read for.
//
// An answer for a project that is no longer open is dropped rather than shown
// against the wrong one, which is the rule the schema answers already follow.
func (m *model) applyStoreGit(msg storeGitMsg) {
	st, ok := m.projects[msg.project]
	if !ok || st.Info.Store == nil {
		return
	}
	store := *st.Info.Store
	g := msg.git
	store.Git = &g
	st.Info.Store = &store
	m.projects[msg.project] = st
}
