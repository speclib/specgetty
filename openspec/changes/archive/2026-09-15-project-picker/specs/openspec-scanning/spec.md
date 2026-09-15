## MODIFIED Requirements

### Requirement: Display projects in TUI
The TUI SHALL present discovered OpenSpec projects in the project picker, sorted
alphabetically by path.

#### Scenario: Projects found
- **WHEN** projects have been discovered and the picker is opened
- **THEN** the picker SHALL list them, and selecting one SHALL open it in the
  project view

#### Scenario: No projects found
- **WHEN** no OpenSpec projects have been discovered
- **THEN** the picker SHALL display a message indicating none were found

## REMOVED Requirements

### Requirement: Display OpenSpec contents without git status
**Reason**: The flat `d`/`f` listing of a project's `openspec/` directory was
replaced by the specs, changes and config tabs, and has had no renderer for some
time. The state feeding it is removed with this change. The requirement that no
git or diff functionality exists is carried by `Remove go-git dependency`, which
is unaffected.
