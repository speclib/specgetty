## ADDED Requirements

### Requirement: Overview displays project stats
The overview tab SHALL display summary statistics for the selected OpenSpec project.

#### Scenario: Project with specs and changes
- **WHEN** a project is selected that has 5 specs, 2 active changes, and 3 archived changes
- **THEN** the overview SHALL display: spec count (5), active changes (2), archived changes (3)

#### Scenario: Project with empty openspec directory
- **WHEN** a project is selected whose openspec/ has no specs, changes, or archive
- **THEN** the overview SHALL display zero counts without errors

### Requirement: Overview lists active changes
The overview tab SHALL list all active changes by name under an "Active Changes" heading.

#### Scenario: Active changes exist
- **WHEN** a project has directories under `openspec/changes/`
- **THEN** each change directory name SHALL be listed under "Active Changes"

#### Scenario: No active changes
- **WHEN** a project has no directories under `openspec/changes/`
- **THEN** the "Active Changes" section SHALL display "None"

### Requirement: Overview lists recently archived changes
The overview tab SHALL list archived changes with their dates under a "Recently Archived" heading.

#### Scenario: Archived changes exist
- **WHEN** a project has directories under `openspec/archive/`
- **THEN** each archived change SHALL be listed with its directory modification date

#### Scenario: No archived changes
- **WHEN** a project has no directories under `openspec/archive/`
- **THEN** the "Recently Archived" section SHALL display "None"
