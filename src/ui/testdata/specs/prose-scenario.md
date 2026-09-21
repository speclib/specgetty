<!-- A scenario whose whole content is a paragraph carrying no keyword. Taken
     from quiqr-desktop/openspec/specs/unified-config/spec.md, scenario
     "Site settings (Future)". Valid per openspec; rendered blank before. -->
# unified-config Specification

## Purpose
Describes where configuration is read from and which file wins.

## Requirements

### Requirement: Configuration sources
The application SHALL read configuration from one file per scope.

#### Scenario: Site settings (Future)

Site-specific settings files (`site_settings_*.json`) are supported in the
codebase but not currently used in the simplified architecture. This feature
is deferred to a future change.

#### Scenario: The ordinary case
- **WHEN** one configuration file is present
- **THEN** it SHALL be the one used
