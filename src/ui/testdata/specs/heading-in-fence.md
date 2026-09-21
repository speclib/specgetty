<!-- A fenced code block holding lines that look like headings. openspec masks
     fences before every heading test. A spec that documents its own format
     reaches this; two files in the local corpus carry `# Page Title` inside a
     fence. -->
# documentation Specification

## Purpose
Describes what a generated page looks like, by showing one.

## Requirements

### Requirement: A page carries a title
Every generated page SHALL open with a level-one heading.

#### Scenario: The generated shape
- **WHEN** a page is generated
- **THEN** it SHALL look like this:

```markdown
# Page Title

## Requirements

### Requirement: not a requirement
#### Scenario: not a scenario
```

- **AND** nothing inside that block SHALL be read as structure
