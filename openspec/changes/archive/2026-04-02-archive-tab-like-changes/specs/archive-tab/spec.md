## ADDED Requirements

### Requirement: Archive tab lists archived changes with artifacts
The archive tab SHALL display archived changes with the same layout as the changes tab.

#### Scenario: Project with archived changes
- **WHEN** the user switches to the archive tab
- **THEN** the left side SHALL list archived change names with their archive date

#### Scenario: No archived changes
- **WHEN** the project has no archived changes
- **THEN** the archive tab SHALL display "No archived changes"

### Requirement: Archive tab shows artifact sub-tabs
When an archived change is selected, the right side SHALL show artifact sub-tabs identical to the changes tab.

#### Scenario: Browsing archived change artifacts
- **WHEN** an archived change is selected
- **THEN** the user SHALL be able to browse proposal, design, tasks, specs sub-tabs with rendered markdown

### Requirement: Archive tab navigation
The user SHALL navigate the archive tab with j/k and left/right the same way as the changes tab.

#### Scenario: Navigate archived changes
- **WHEN** the archive tab is active
- **THEN** j/k SHALL move the archive cursor and left/right SHALL move the artifact sub-tab
