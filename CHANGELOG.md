# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Building from source now needs Go 1.25 or newer, which the move to bubbletea
  v2 requires. Nothing about how specgetty behaves has changed.

## [0.4.0] - 2026-09-15

## [0.3.0] - 2026-09-15

### Added

- Copy a reference to the selected change: `y` puts its name on the clipboard,
  `Y` puts the absolute path of its directory. A one-line message reports what
  was copied, and says so plainly if no clipboard tool is available rather than
  looking like it worked.
- Export a change as a zip with `e`, from the change list or on an archived
  change. Writes `~/<change-name>-<date>.zip` with the change directory intact,
  so the full proposal, design, tasks and specs can be attached to a pull
  request in another repository.
- Project picker on `p`, available from anywhere. Lists every project with its
  spec, change and task counts, filters with `/`, and switches project on
  `enter`.
- The picker remembers what it found in `~/.cache/specgetty/projects.yaml`, so
  only the first run pays for scanning your disk. `r` looks again. The counts it
  shows are always read fresh, so only the list of projects can be stale.
- `--view=single|all` chooses which view opens first.
- Search the change list with `/`: fuzzy on change names, `'` for a literal name
  match, and `:` to search the text inside proposal, design, tasks and spec
  files. A row matched on file contents names the files that matched.
- Cycle the change list between open, archived and both with `f`.
- Choose which columns the change list shows through `change_fields` in the
  config file or `--change-fields` on the command line.

### Fixed

- Markdown documents can be scrolled. Opening a change whose proposal was longer
  than the panel used to put the rest of it out of reach; `j`/`k`, the page keys,
  `ctrl-d`/`ctrl-u` and `gg`/`G` now move through it, and the panel title reports
  how far down you are.
- Long lines are wrapped to the pane instead of being cut. Content was
  previously lost inside the visible height as well as below it, in every
  markdown pane including the specs tab.
- The config tab scrolls the same way, keeping the file name fixed above it.
- Exporting an archived change works. It looked for the change directory under
  its display name, without the date prefix the directory actually carries, so
  it always reported the source as not found. Broken since the feature shipped.
- The specs tab scrolls too. `tab` moves the keyboard between the spec list and
  the spec content: with the list focused `j`/`k` change spec as before, with the
  content focused they scroll it. The selected spec dims while the content has
  the keys, and the title reports the reading position.

### Changed

- **Breaking**: `spg` now opens the project you are standing in, immediately,
  instead of scanning your disk first. If there is no project there it offers
  the picker.
- **Breaking**: the project list panel is gone, replaced by the picker overlay.
  The window is one panel now.
- **Breaking**: `--zoom` / `-z` is removed. Use `--view=single`, which is the
  default, or `--path`. The word "zoom" is gone from the UI: looking at one
  project is the ordinary way to use the tool, not a mode.
- The change list now fills the panel as a table, and a change opens one level
  deeper with `enter`. Its artifact sub-tabs no longer share left/right with the
  project tab bar, so pressing right on the last sub-tab stays put instead of
  jumping to another tab.
- The archive tab is gone. Archived changes live in the one change list behind
  the `f` filter.
- Tabs are now changes, specs, config, reachable with `1`, `2` and `3`.
- `esc` is the only way back up a level, and it does nothing at the project
  view, which is the floor.

## [0.2.0] - 2026-04-02

- feat: auto-rescan in zoom mode — filesystem watcher (fsnotify) monitors the openspec directory and triggers rescan on changes
- feat: config tab renders `project.md` as markdown or `config.yaml` with syntax highlighting
- feat: specs tab with split view — spec list on left, rendered markdown on right
- feat: changes tab with artifact sub-navigation — browse proposal, design, tasks, and specs per change
- feat: task progress display in changes list (e.g. `18/20`)
- feat: archive tab reuses changes layout for browsing archived changes with date display
- feat: zoom mode — press enter to give detail panel full terminal width, escape to return
- feat: `--zoom` / `-z` and `--path` / `-p` CLI flags for starting in zoom mode
- feat: persistent project header (path + stats) visible across all tabs, replacing overview tab
- feat: single-project rescan with `s` key in zoom mode
- feat: show aggregate task progress (done/total) in persistent project header
- feat: archive changes from TUI — press `a` on changes tab, with confirmation and incomplete task warnings
- refactor: replace `ArchivedChange` struct with unified `ChangeInfo` for active and archived changes
- refactor: shared `renderChangeList` function used by both changes and archive tabs

## [0.1.5] - 2026-03-31

- fix: nix flake

## [0.1.4] - 2026-03-31
- 
- feat: horizontal split layout — project list on left, detail panel on right
- feat: project list shows basenames with parent-dir disambiguation for duplicates
- feat: detail panel with tab system — overview, specs, changes, config, search (only overview functional)
- feat: overview tab shows project stats, active changes, and recently archived changes
- feat: tab switching with left/right arrows and 1-5 number keys
- fix: validate openspec directories require (config.yaml or project.md) and (specs/ or archive/)
- fix: scanning modal text updated to "Scanning for OpenSpec sources..."
- chore: remove enter/open keybinding and tmux popup functionality
- chore: improve contrast on inactive tabs and stats text

## [0.1.3] - 2026-03-31

## [0.1.2] - 2026-03-30

## [0.1.1] - 2026-03-30

- **BREAKING**: Replace git repo scanning with OpenSpec project detection
- **BREAKING**: Remove git status, diff panel, and all go-git dependency
- feat: scan for `openspec/` directories instead of `.git/`
- feat: show OpenSpec directory contents with d/f indicators
- feat: rename UI panels to "Projects" and "Contents"
- chore: remove gitignore config section (no longer relevant)

## [0.1.0] - 2026-03-30

- Fork from mipmip/dirty-repo-scanner
- Rename project to specgetty (binary: spg)
- Add MIT license
