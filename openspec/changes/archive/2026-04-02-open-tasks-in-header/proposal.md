## Why

The persistent header shows project stats (specs, changes, archived) but doesn't surface open task counts. When managing multiple projects, knowing at a glance how much work remains across active changes is valuable without having to drill into each change.

## What Changes

- Aggregate `TasksDone` and `TasksTotal` across all active changes per project
- Display the aggregate task count in the persistent header stats line
- Only show the task stat when there are tasks (avoid clutter on projects with no tasks)

## Capabilities

### New Capabilities

- `aggregate-task-stats`: Compute and expose total open tasks across all active changes in a project

### Modified Capabilities

## Impact

- `src/scanner/scan.go`: Add aggregate fields to `ProjectInfo`
- `src/ui/ui.go`: Update stats line format in `renderDetailPanel`
