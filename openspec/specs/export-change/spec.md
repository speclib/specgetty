# export-change Specification

## Purpose
Allows users to export a change (active or archived) as a zip file for attaching to foreign PRs.

## Requirements

### Requirement: Export keybinding
The system SHALL provide an `e` keybinding that initiates export when the changes tab or archive tab is active and a change is selected.

#### Scenario: Press e on changes tab with a change selected
- **WHEN** the user presses `e` while the changes tab is focused and a change is highlighted
- **THEN** the system SHALL show a confirmation modal for exporting that change

#### Scenario: Press e on archive tab with a change selected
- **WHEN** the user presses `e` while the archive tab is focused and an archived change is highlighted
- **THEN** the system SHALL show a confirmation modal for exporting that archived change

#### Scenario: Press e when no changes exist
- **WHEN** the user presses `e` on the changes or archive tab with no entries
- **THEN** the system SHALL do nothing

#### Scenario: Press e on a different tab
- **WHEN** the user presses `e` while a tab other than changes or archive is active
- **THEN** the system SHALL do nothing

### Requirement: Zip creation
The system SHALL create a zip file containing the entire change directory with
preserved internal structure, in the directory chosen for the export.

#### Scenario: Active change export
- **WHEN** exporting an active change named `my-feature`
- **THEN** the zip SHALL contain a root folder `my-feature/` with all files and subdirectories from `openspec/changes/my-feature/`

#### Scenario: Archived change export
- **WHEN** exporting an archived change with directory name `2026-04-02-my-feature`
- **THEN** the zip SHALL contain a root folder `my-feature/` (date prefix stripped) with all files and subdirectories from `openspec/changes/archive/2026-04-02-my-feature/`

#### Scenario: Zip filename format
- **WHEN** a change is exported on date 2026-04-22
- **THEN** the zip file SHALL be named `<semantic-name>-2026-04-22.zip` where semantic-name has any date prefix stripped

#### Scenario: Zip destination
- **WHEN** a change is exported
- **THEN** the zip file SHALL be written to the directory chosen for that
  export, which is the configured one unless it was edited at the prompt

#### Scenario: Existing file at destination
- **WHEN** a zip file with the same name already exists in that directory
- **THEN** the system SHALL ask before replacing it rather than replacing it
  silently

### Requirement: Result feedback
The system SHALL display the export result in a modal that dismisses on any keypress.

#### Scenario: Successful export
- **WHEN** the zip file is created successfully
- **THEN** the system SHALL display a success modal showing the destination path

#### Scenario: Failed export
- **WHEN** the zip creation fails (e.g., permission error, disk full)
- **THEN** the system SHALL display a failure modal with the error message

#### Scenario: Dismiss result modal
- **WHEN** the user presses any key while the result modal is shown
- **THEN** the modal SHALL be dismissed and the UI SHALL return to normal state

### Requirement: Nav bar hint
The nav bar SHALL show an `e export` hint when the changes or archive tab is active and entries exist.

#### Scenario: Changes tab active with changes
- **WHEN** the changes tab is active and the project has active changes
- **THEN** the nav bar SHALL include the `e export` keybinding hint

#### Scenario: Archive tab active with archived changes
- **WHEN** the archive tab is active and the project has archived changes
- **THEN** the nav bar SHALL include the `e export` keybinding hint

#### Scenario: No entries on current tab
- **WHEN** the changes or archive tab is active but has no entries
- **THEN** the nav bar SHALL NOT include the `e export` hint

### Requirement: The export directory is configured and editable
The directory a change is exported to SHALL come from configuration, and SHALL
be editable at the moment of export. The configured value is what the prompt
opens on.

Someone who always exports to the same place sets it once and types nothing
again; someone who does not can still redirect a single export without changing
their configuration.

#### Scenario: A configured directory
- **WHEN** the configuration names an export directory
- **THEN** the prompt SHALL open on it

#### Scenario: No configured directory
- **WHEN** the configuration names none
- **THEN** the prompt SHALL open on the user's home directory, which is where
  every export went before it could be chosen

#### Scenario: A directory written with a shorthand
- **WHEN** a configured or typed directory begins with `~`, or contains an
  environment variable
- **THEN** it SHALL be expanded before it is used, because both are how people
  write a path they mean

#### Scenario: Editing the directory for one export
- **WHEN** the directory is edited at the prompt
- **THEN** that export SHALL use what was typed and the configuration SHALL be
  unchanged

#### Scenario: The next export
- **WHEN** a further export is started after one was redirected
- **THEN** its prompt SHALL open on the configured directory again, not on what
  was typed last

### Requirement: Only a directory is typed
The prompt SHALL accept a directory and SHALL NOT accept a filename. The name of
the zip is generated, and remains what this capability already specifies.

The generated name strips an archived change's date prefix and stamps the export
date. Both do work that retyping would lose.

#### Scenario: What the prompt shows
- **WHEN** the prompt is displayed
- **THEN** it SHALL show the directory as editable text and the generated
  filename beside it as fixed text

#### Scenario: Completing a path
- **WHEN** the user presses `tab` in the prompt
- **THEN** the directory SHALL be completed against the filesystem

#### Scenario: Cancelling
- **WHEN** the user presses `esc` at the prompt
- **THEN** no file SHALL be written and the change list SHALL return unchanged

### Requirement: A destination that cannot be used is refused
An export SHALL NOT create a directory, and SHALL NOT replace an existing file
without being told to.

Creating a directory silently turns a typo into a directory. Replacing a file
silently was defensible while the destination was a generated name in the user's
own home; it is a different act once the path is typed.

#### Scenario: A directory that does not exist
- **WHEN** the typed directory does not exist
- **THEN** the export SHALL be refused with a message naming the directory, and
  what was typed SHALL remain in the prompt to be corrected

#### Scenario: A path that is not a directory
- **WHEN** the typed path exists and is not a directory
- **THEN** the export SHALL be refused the same way

#### Scenario: A file already at the destination
- **WHEN** a file already exists at the destination
- **THEN** the user SHALL be asked to confirm replacing it, naming the file

#### Scenario: Declining to replace it
- **WHEN** the user declines
- **THEN** nothing SHALL be written and the prompt SHALL return with the
  directory still in it, so another can be chosen

#### Scenario: Agreeing to replace it
- **WHEN** the user agrees
- **THEN** the file SHALL be replaced and the export SHALL report where it wrote
