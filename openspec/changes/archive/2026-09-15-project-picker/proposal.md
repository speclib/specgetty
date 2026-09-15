## Why

Every launch pays for a full directory walk before showing anything. Measured on
a real configuration: **3.9s cold, 0.7s warm**, to discover 18 projects. The
project you are standing in, which is what you almost always want, costs about a
millisecond to resolve. The default is backwards.

The walk is also the wrong thing to repeat. It discovers a set of paths that
changes when you create a project, which is rarely. The per-project statistics,
which change constantly, are cheap: reading every artifact of all 18 projects
takes about 100ms.

Separately, "zoom" names a mode that stopped being a mode. Looking at one project
is not a special state to enter and leave; it is the ordinary way to use the
tool. The vocabulary survives from a design where the project list was the main
view.

## What Changes

- **BREAKING**: the single project view becomes the startup view. `spg` resolves
  the project from the working directory and opens it immediately.
- **BREAKING**: the project list is no longer a panel. It becomes a picker
  overlay opened with `p` from anywhere, closed with `esc`, which switches the
  current project with `enter`.
- **BREAKING**: the horizontal split layout is removed. The detail panel is the
  whole view.
- **BREAKING**: `--zoom` / `-z` is removed. `--view=single|all` replaces it and
  defaults to `single`. `--path` is retained and implies `--view=single`.
- The picker caches discovered project paths on disk, so only the first run pays
  for the walk. `r` inside the picker refreshes it.
- The picker filters with the same query grammar as the change list: fuzzy on
  project names, `'` for a literal name match, `:` for file paths and file
  contents inside projects.
- When no project can be resolved and none is selected, the view says so and
  points at `p` rather than showing an empty frame.
- Remove the vestigial file-list state (`filePaths`, `fileCursor`), which no
  longer has a renderer, and the requirement describing the listing it fed.

## Capabilities

### New Capabilities
- `project-picker`: the overlay, its table of projects with statistics, its
  filter, its refresh action, and what selecting a project does
- `project-cache`: the on-disk cache of discovered project paths, what
  invalidates it, and how it is refreshed
- `project-view`: the base view, how the startup project is resolved, the
  `--view` flag, the empty state, the persistent header and the minimum
  terminal size

### Modified Capabilities
- `openspec-scanning`: projects are presented in the picker rather than a left
  panel, and the `d`/`f` file listing requirement is removed because the view it
  describes no longer exists
- `fs-watch`: every scenario is phrased in terms of entering and leaving zoom
  mode and needs rewording for the project view

### Removed Capabilities
- `zoom-mode`: the concept is gone. Its surviving behaviour, a project filling
  the width and `s` rescanning only that project, belongs to `project-view`.
- `side-panel-layout`: there is no side panel. The persistent header and the
  minimum terminal size move to `project-view`, and basename display moves to
  `project-picker`.

## Impact

- `src/ui/ui.go`: delete `renderProjectList`, `leftPanelWidth`,
  `rightPanelWidth`, the split branch of `View()` and the `viewProjects` panel;
  `activeView` collapses to whether the log panel is focused; the level constants
  renumber to project and change; the picker adds a third key-intercepting mode
  alongside the confirm modals and the search prompt
- `src/ui/` new files for the picker overlay and the cache
- `src/scanner/`: a cache read and write around `Walk`
- `src/main.go`: `--view` replaces `--zoom`
- Non-goal: openspec stores. `--view=store` is deliberately not part of this
  change.
