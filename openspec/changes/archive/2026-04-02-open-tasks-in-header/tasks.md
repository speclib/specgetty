## 1. Scanner: aggregate task fields

- [x] 1.1 Add `TasksTotal int` and `TasksDone int` fields to `ProjectInfo`
- [x] 1.2 In `ParseProjectInfo`, after parsing all active changes, sum `TasksTotal` and `TasksDone` from each `ChangeInfo` into the new `ProjectInfo` fields
- [x] 1.3 Add test: project with multiple changes containing tasks verifies correct aggregation
- [x] 1.4 Add test: project with no tasks verifies zero values

## 2. UI: display in persistent header

- [x] 2.1 Update stats line in `renderDetailPanel` to conditionally append `Tasks: <done>/<total>` when `info.TasksTotal > 0`
- [x] 2.2 Add test: verify stats line includes task count when tasks exist
- [x] 2.3 Add test: verify stats line omits task count when no tasks

## 3. Verification

- [x] 3.1 Verify build and all tests pass
