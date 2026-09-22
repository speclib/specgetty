## MODIFIED Requirements

### Requirement: Startup view selectable by flag
The application SHALL accept a `--view` flag selecting which view opens first.
Flags that selected a change list filter mode are gone with the modes.

The command line SHALL take no positional arguments. A directory to open is
named with `--path`, and the directories the picker searches come from
`scandirs.include` in the configuration file. An argument that is neither SHALL
be refused rather than ignored, because a silently dropped argument looks like
one that was honoured.

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

#### Scenario: Directories as arguments are gone
- **WHEN** the user runs `spg ~/work ~/clients`
- **THEN** the application SHALL report that directories are not accepted, and
  SHALL name `--path` and `scandirs.include` as what to use instead

#### Scenario: Nothing is scanned on the strength of an argument
- **WHEN** an argument is given
- **THEN** no project SHALL be opened and no directory SHALL be walked, the
  application having refused before it starts

#### Scenario: A configuration that cannot be read is still an error
- **GIVEN** a configuration file that does not parse
- **WHEN** the application is run with any arguments at all
- **THEN** it SHALL report the configuration error, rather than continuing on
  defaults
