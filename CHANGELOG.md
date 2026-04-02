# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
