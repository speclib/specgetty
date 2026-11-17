# Spec Delta

## ADDED Requirements

### Requirement: Capabilities are discovered by the spec.md they hold
A capability is a directory under a project's `specs/` that directly holds a
`spec.md`, named by its path relative to `specs/` with `/` between segments. The
scanner SHALL discover capabilities at any depth, SHALL NOT treat a directory
holding no `spec.md` as one, and SHALL reach a change's deltas by the same rule.
A project's spec count SHALL be the number of capabilities found.

#### Scenario: Nested main specs
- **WHEN** a project holds `specs/content/editing/spec.md` and
  `specs/content/pages/spec.md`
- **THEN** both SHALL be reported, named `content/editing` and `content/pages`

#### Scenario: A container is not a capability
- **WHEN** `specs/content` holds no `spec.md` of its own
- **THEN** `content` SHALL NOT be reported as a capability, whatever it holds
  beneath it

#### Scenario: Nested deltas
- **WHEN** a change holds `specs/tenancy/projects/spec.md` and
  `specs/publish/act/spec.md`
- **THEN** it SHALL report two deltas, named `tenancy/projects` and
  `publish/act`, and neither SHALL replace the other

#### Scenario: A capability and a group of the same name
- **WHEN** a project holds both `specs/platform/spec.md` and
  `specs/platform/release/spec.md`
- **THEN** both SHALL be reported, named `platform` and `platform/release`

#### Scenario: The count is the capabilities, not the directories
- **WHEN** a project holds 52 directories with a `spec.md` under 14 containers
- **THEN** its spec count SHALL be 52

#### Scenario: A project with no nesting
- **WHEN** a project holds only `specs/archive-from-ui/spec.md`
- **THEN** it SHALL be reported as `archive-from-ui`, exactly as before

#### Scenario: A subtree that cannot be read
- **WHEN** a directory beneath `specs/` cannot be read
- **THEN** the scan SHALL report the capabilities it could read rather than
  none, and SHALL NOT end the process
