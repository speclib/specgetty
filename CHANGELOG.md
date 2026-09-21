# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- A spec is read by OpenSpec's own rules, transcribed from its parser: the
  `## Requirements` section, `### Requirement:` inside it, any `#### ` heading
  with content under it, headings inside code fences ignored, and a delta header
  treated as the error it is in a main spec. A file that does not follow them now
  opens as a report naming every reason and its line, with `E` to open it in your
  editor, instead of opening as half an outline. The whole file stays readable as
  markdown on the specs tab.

### Fixed

- Scenarios no longer show up empty. A scenario's content is now kept whatever
  shape it is written in, rather than only the bulleted `- **WHEN**` form the
  OpenSpec template shows. Across the 605 specs on the machine this was written
  on, 254 scenario cards were blank and every one of them was content that had
  been dropped.
- `left` and `right` no longer do anything in the spec detail view. They moved
  the project tab bar underneath it, which is not on screen there, and took the
  keyboard off the card mid-read.

## [0.7.1] - 2026-09-21

### Added

- `E` opens the file a pane is showing in your own editor, from a change's
  artifacts, a spec, or the project configuration. The editor is `$VISUAL`, then
  `$EDITOR`, and a value with arguments such as `code -w` works. With neither
  set the key says so and does nothing rather than guessing. The interface
  yields the terminal while the editor runs and reads the project again when it
  exits.

### Fixed

- In a project that keeps its specs in a store, toggling a task checkbox failed
  with `could not save: no such file or directory`, and `Y` copied a path that
  looked right and did not exist. Both built their path from the repository you
  started in rather than from the store the content was read from.

## [0.7.0] - 2026-09-21

### Changed

- **BREAKING** The change list shows active and archived changes together,
  grouped, with the active group first and each group counting its rows. The
  three filter modes are gone with it: the `f` key, the `change_mode` config
  setting and the `--change-mode` option no longer exist, because which changes
  you are looking at is now said by the group a row sits in rather than by a
  label in the nav bar. A configuration still carrying `change_mode` is reported
  at startup.
- Archived changes are ordered newest first. They were ordered oldest first, by
  accident of the directory naming, so the change you had just archived sat at
  the bottom of the list.
- The archive date is a default column, blank on an active change. The column
  saying whether a change is active or archived is no longer a default, since
  the group header says it, and remains available in `change_fields`.

### Added

- `enter` on the specs tab opens a spec at its own navigation level: an outline
  of its requirements and scenarios beside a card showing whichever one the
  cursor is on, with its clauses laid out under their keywords rather than as
  raw markdown. `tab` moves the keyboard between the two halves, `esc` returns
  to the list, and the cursor keeps its place across a rescan.
- `e` asks where to put the zip. The prompt opens on a new `export_dir` config
  setting, or your home directory when it is unset, and `tab` completes a path.
  You type a directory and the filename stays generated. A directory that does
  not exist is refused rather than created, and an existing file is confirmed
  before it is replaced, where it used to be overwritten silently.
- The search prompt names its three matchers while it is focused and empty, so
  `'exact` and `:inside` stop being things you have to read the README to learn.
  It goes on the first keystroke and is dropped where the line is too narrow.
- A name search that finds nothing suggests the same term as a contents search,
  which is the moment you are already looking for another way. Both the change
  list and the project picker.

### Changed

- One blank line separates the active and archived groups in the change list, so
  the boundary reads as a break rather than as another row. It is drawn whether
  or not either group has rows in it.

### Removed

- `edit_command` from the config file. It was parsed and never read by anything,
  so it never had an effect to lose. A config that still sets it is reported.
- The log panel, its `l` key and its place in the focus ring. One scan wrote 35
  lines into it, 30 of them per-project timings, and the three real errors it
  could carry were buried under them in a panel that had to be asked for. Log
  output is discarded while the interface is running; `--debug` prints exactly
  what it always did.

### Fixed

- No modal can be wider than the terminal. The export prompt and the startup
  question both had fixed widths that overflowed a 60-column terminal.
- A `scandirs.include` entry ending in `*` is no longer walked as a path in its
  own right, only as what it expands to. The pattern is not a directory, so
  walking it failed on every scan, and it failed the whole scan outright with
  `--ignore_dir_errors=false`, which was therefore broken for any configuration
  using a glob.

- `pgup`, `pgdown`, `ctrl+f`, `ctrl+b`, `ctrl+d`, `ctrl+u`, `gg` and `G` now
  work in every list, not only in documents. They act on whatever holds the
  keyboard, by the same rule `j` and `k` follow, so the change list, the spec
  list, the properties rows and the project picker can all be paged and jumped
  through. The change list had thirty-eight rows and no way down it but `j`.
- In the combined view, nothing indicated which rows were archived unless you
  had configured the `archived` column. Grouping makes it positional.
- `spg --config <unparseable.yml> <dir>` crashed. Directory arguments replace
  the configured scan directories, so a configuration that failed to parse is
  now started from empty rather than dereferenced.

### Added

- The config tab is now `properties`, and reports what a project is rather than
  only what is in its config file. Its rows sit in a list beside their content,
  the same split the specs tab uses: the one configuration that applies, a row
  for each workflow schema the project's changes record, and the store the
  content comes from.
- Workflow schemas are reported: which ones the project's changes use and how
  many changes are on each, plus each schema's source, artifact chain and apply
  rule. A schema that overrides a built-in says so. Definitions are located with
  `openspec schema which`, so the rows need the `openspec` CLI for their detail
  but not for the counts.
- An open change names the schema it records.
- A repo that declares a store is told which of its own configuration keys have
  no effect. OpenSpec reads such a file for `store:` alone and takes `schema`,
  `context`, `rules` and `operations` from the store, and warns about none of
  them.

- OpenSpec stores are supported. A repo whose `openspec/` holds only a
  `config.yaml` naming a `store:` now opens on that store's specs and changes
  instead of reporting an empty project, resolved through the registry at
  `$XDG_DATA_HOME/openspec/stores/registry.yaml` without calling the `openspec`
  binary or touching the network.
- The project header marks a project whose content comes from a store, and the
  config tab explains the mark: it gains sub-tabs for the repo's own
  configuration, the store's shared one, and a store report with the id, the
  root path, the registered remote and branch, and the store working copy's
  local git state (uncommitted changes, and how far ahead or behind its last
  known upstream ref it is, read without fetching).
- The project picker lists the repos you work in. A repo that declares a store
  and keeps no specs of its own is a row like any other, carrying that store's
  spec, change and task counts, and a `store` column names the store it reads
  from. A registered store is not listed on its own: it is where content lives
  rather than where work happens. Run `spg` inside a store, or pass `--path`, to
  open one directly.
- A store declaration that cannot be followed now says so where the project
  would otherwise look empty, naming the declared id, the file that declared it
  and the reason, rather than showing zero specs and zero changes.

### Changed

- The archive, discard and export actions act on the store a project resolves
  to. Discarding a change in a store-backed project used to create an unused
  `openspec/changes/discarded/` directory in the pointing repo and move nothing.
- Filesystem watching covers both the store's tree and the pointing repo's, so
  editing a `store:` key is picked up and the view follows it.
- `openspec/config.yml` is now read wherever `openspec/config.yaml` is, matching
  what OpenSpec itself accepts.
- A directory named `openspec` holding neither a configuration nor any content
  is no longer treated as a project. Without that rule a store kept at
  `~/openspec` would make the home directory capture every project beneath it.
- A directory counts as a store only when the store registry resolves its
  declared id to that same directory. A clone, or a folder the registry has
  moved on from, keeps its copy of `.openspec-store/store.yaml` and used to be
  shown under the registered store's name.

## [0.6.0] - 2026-09-18

### Changed

- Every view now has a column of air inside the panel border, on both sides.
  Only the project header had it before, so the tab bars, the change list, the
  specs split, the config pane and every markdown document sat flush against
  the frame, and wrapped lines ended on the border itself. The change list's
  selected row now has the same gutter the project picker already had.
- The nav bar indents its text by a column at each end, and the gap between the
  key hints and the version is painted rather than left bare. That gap used to
  be a hole in the middle of the strip.
- Table columns are separated by two blank columns instead of one, in the
  change list and the project picker. A value that filled its column used to
  sit a single space from the value beside it. At narrow terminals a column is
  now dropped two columns sooner, since the wider gaps come out of the same
  budget.
- The content of each tab is drawn in a border of its own, directly under the
  tab bar, so the chips read as tabs belonging to what is below them. An open
  change gets the same treatment around its artifact. The border costs two rows,
  so a 20-row terminal shows nine changes where it showed eleven.
- The config tab's filename now sits above that border rather than inside it,
  with the tab bar. It names what is in the box, the same job the tab chips do.
- The specs tab draws its list and its document in separate borders, and the
  one holding the keyboard is lit. `tab` moves between them, and until now the
  only sign it had was the selected spec dimming. A border is lit whenever the
  keyboard is inside it, so the panel border keeps meaning exactly what it did.

## [0.5.0] - 2026-09-18

### Added

- Tick tasks off without leaving specgetty. Open a change's tasks, move the
  cursor with `j`/`k`, and press space. Checkboxes draw as `▢` and `▣`, the
  selected task is highlighted across every row it wraps onto, and the file is
  saved immediately.
- A toggle is applied to the file as it is on disk at that moment, so an edit
  saved from your editor since the last scan is never discarded. If the task has
  moved or vanished, specgetty says so instead of guessing.

### Changed

- The change list calls its two states **active** and archived everywhere. It
  previously said "open" in the nav bar and "active" in the empty state for the
  same thing, and "open" also means the change you have opened.
- `change_mode` in the config, or `--change-mode`, sets which changes the list
  starts on and returns to when you switch project: `active`, `archived` or
  `active+archived`. `f` still cycles between them.
- Building from source now needs Go 1.25 or newer, which the move to bubbletea
  v2 requires. Nothing about how specgetty behaves has changed.

### Fixed

- A scan directory that cannot be read no longer ends the program. An include
  ending in `*` whose parent is missing, an unmounted drive or a path that
  moved, used to exit specgetty from inside the full-screen view. It is now
  logged and skipped, which is what `--ignore_dir_errors` always promised.
- An empty entry under `scandirs.include` or `scandirs.exclude` no longer
  crashes the scan. A YAML dash with nothing after it now matches nothing.

## [0.4.0] - 2026-09-15

### Added

- Copy a reference to the selected change: `y` puts its name on the clipboard,
  `Y` puts the absolute path of its directory. A one-line message reports what
  was copied, and says so plainly if no clipboard tool is available rather than
  looking like it worked.

### Fixed

- Exporting an archived change works. It looked for the change directory under
  its display name, without the date prefix the directory actually carries, so
  it always reported the source as not found. Broken since the feature shipped.

## [0.3.0] - 2026-09-15

### Added

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
