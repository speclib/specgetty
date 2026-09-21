## MODIFIED Requirements

### Requirement: A configuration setting that no longer applies is reported
A configuration file may still carry a setting a released version accepted and
this one does not. YAML ignores keys a program does not know, so such a key
would otherwise sit in a configuration doing nothing. The application SHALL
report every one it finds, naming the key and why it no longer has an effect.

#### Scenario: A configuration still naming a change mode
- **WHEN** the configuration file carries a `change_mode` key
- **THEN** the application SHALL report that the key no longer has an effect,
  and SHALL start normally

#### Scenario: A configuration without it
- **WHEN** the configuration file carries no retired key
- **THEN** nothing SHALL be reported

#### Scenario: A configuration naming an edit command
- **WHEN** the configuration file carries an `edit_command` key
- **THEN** the application SHALL report it the same way, it being a setting that
  shipped, was parsed, and was never read by anything

#### Scenario: More than one at once
- **WHEN** the configuration file carries several retired keys
- **THEN** each SHALL be reported, in a stable order
