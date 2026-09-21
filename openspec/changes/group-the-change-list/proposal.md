## Why

Bean `specgetty-m1qx`. One list holds active and archived changes, and a hidden
key cycles which of them it contains. Three things follow from that, and all
three are visible in this project right now.

**You cannot tell what you are looking at.** In `active+archived` mode with the
default columns the list reads:

```
  name                                       tasks    specs
  properties-tab                             0/47     5        <- active
  fix-openspec-detection-false-positives     6/6      1        <- archived
  replace-git-with-openspec-scanning         24/24    1        <- archived
```

Nothing distinguishes them. `change-list-view` says otherwise, in a scenario
that is not met as shipped: "Both mode ... each row SHALL indicate whether it is
active or archived". The `archived` column exists and is not in the defaults.

**The archive is upside down.** Archived changes are appended in `os.ReadDir`
order, which sorts by the date-prefixed directory name, and nothing re-sorts
them. The change archived most recently is last:

```
  row  1   2026-03-31-fix-openspec-detection-false-positives
  row 35   2026-09-21-list-the-repos-not-the-stores     <- shipped today
```

**The mode is invisible.** `f` cycles three states whose only indication is a
label in the nav bar. Being in `archived` and not realising it looks exactly
like a project whose active work has vanished.

## What Changes

- Active and archived changes are shown together, always, in one table, grouped
  under headers that name the group and count it. Active first
- The filter modes go: the `f` key, the `change_mode` config key, the
  `--change-mode` option and the three mode names. The group a row sits in is
  what the modes were trying to say, and it says it by position rather than by
  a label somewhere else
- A configuration still carrying `change_mode` is reported as no longer having
  an effect. YAML ignores unknown keys silently, and the requirement being
  removed here is the one that says an unusable setting SHALL be reported rather
  than swallowed
- Archived changes are ordered newest first. Active changes stay ordered by
  name. Neither is user-sortable: what was wanted from sorting was to see recent
  work first, and choosing the other fixed order gives that without a sort
  mechanism
- The archived date joins the default columns, blank on active rows, which is
  what the field already renders for them. The `archived` state column leaves
  the defaults as redundant against the group header, and stays configurable
- A group is shown with its count even when the count is zero, so that "nothing
  in flight" is an answer rather than an absence

Not in scope: sorting. No column becomes sortable, and the two orders above are
fixed. `specgetty-vru8` is where sortable columns would belong if they are ever
wanted, together with the columns that need new scanner capability.

Not in scope: hiding the archive. Dropping the modes means a project with a
large archive always renders it. The group header says where you are and `/`
filters; that is accepted rather than overlooked.

## Capabilities

### Modified Capabilities
- `change-list-view`: one grouped list replaces the three filter modes, the
  default columns gain the date, and the two group orders are fixed
- `changes-tab`: the tab lists both groups, and its empty state no longer refers
  to a filter mode
- `project-view`: `--change-mode` is gone, the way `--zoom` is gone

## Impact

- `src/ui/changerow.go`: the three modes, their names and `emptyListMessage` go;
  `buildRows` stops taking a mode and starts producing groups
- `src/ui/fields.go`: `ResolveListMode` goes; `defaultFields` gains `date`
- `src/ui/ui.go`: `listMode` and `defaultMode` leave the model, the `f` handler
  and its nav bar entry go, and the change list gets a renderer that draws group
  headers between rows
- `src/scanner/scan.go`: `Config.ChangeMode` goes
- `src/main.go`: the `--change-mode` flag goes
- `src/config.yml`: the documented `change_mode` block goes
- `renderTable` and the project picker are untouched. The grouped renderer
  reuses `layoutFields` and `fitCell` rather than teaching the shared function
  about groups

## Rollback

Its own commit. Reverting restores the modes, the key, the flag and the config
key together, since they are one feature. A configuration written in the
meantime cannot be invalidated by the revert, because nothing new is written to
one.
