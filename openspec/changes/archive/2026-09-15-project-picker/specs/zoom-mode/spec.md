## REMOVED Requirements

### Requirement: Enter zoom mode
**Reason**: Looking at one project is no longer a mode to enter. It is the
startup view, specified by `project-view`.

### Requirement: Exit zoom mode
**Reason**: There is no mode to leave. `esc` at the project view does nothing,
because it is the floor of the navigation stack.

### Requirement: Full functionality in zoom mode
**Reason**: There is no reduced mode to contrast with, so there is nothing to
qualify as full.

### Requirement: Zoom mode shows project context
**Reason**: Replaced by the persistent project header in `project-view`, which
is shown whatever tab is active.

### Requirement: Scanning scoped to zoomed project
**Reason**: Moved to `project-view`, where `s` rescans the open project. The
whole-disk walk now belongs to the picker's refresh action.

### Requirement: Start zoomed via CLI flag
**Reason**: `--zoom` is removed. Startup is specified by `project-view` through
`--view` and `--path`.
