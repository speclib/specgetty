## MODIFIED Requirements

### Requirement: Startup view selectable by flag
The application SHALL accept a `--view` flag selecting which view opens first.
Flags that selected a change list filter mode are gone with the modes.

#### Scenario: Default
- **WHEN** no `--view` flag is given
- **THEN** the application SHALL behave as `--view=single`

#### Scenario: Opening at the picker
- **WHEN** the user runs `spg --view=all`
- **THEN** the project picker SHALL be open when the application starts

#### Scenario: Explicit path
- **WHEN** the user runs `spg --path /some/project`
- **THEN** that project SHALL be opened, as though `--view=single` had resolved
  to it

#### Scenario: The zoom flag is gone
- **WHEN** the user runs `spg --zoom`
- **THEN** the application SHALL report that the flag is not recognised

#### Scenario: The change-mode flag is gone
- **WHEN** the user runs `spg --change-mode=active`
- **THEN** the application SHALL report that the flag is not recognised
