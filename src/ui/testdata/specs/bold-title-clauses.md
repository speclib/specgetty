<!-- Bold, title case, with and without a bullet. Taken from
     /home/pim/gh.mipmip/quiqr-desktop/openspec/specs/docusaurus-setup/spec.md
     and scaffold-model/spec.md. That file passes `openspec validate --strict`
     and specgetty rendered all 14 of its scenarios as bare titles. 900 lines
     in the local corpus are written one of these two ways. -->
# docusaurus-setup Specification

## Purpose
Builds the documentation site as a workspace package of the monorepo.

## Requirements

### Requirement: Docusaurus Workspace Package
The project MUST include a Docusaurus documentation site at `packages/docs`.

#### Scenario: Initialize Docusaurus Package

**Given** the monorepo root
**When** the workspace is initialized
**Then** a Docusaurus package exists at `packages/docs`

#### Scenario: Scaffold single from YAML file
- **Given** a site is mounted
- **When** the user triggers "Scaffold Single from File"
- **Then** the system infers field types and creates a model definition
