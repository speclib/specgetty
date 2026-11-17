# Spec Delta

## MODIFIED Requirements

### Requirement: Specs tab displays a list of specs
The specs tab SHALL display a vertical list of spec names on the left side of
the detail panel. A name with more than one segment SHALL be drawn as a header
for its first segment, carrying how many capabilities that segment holds, and
a row for the rest of the path.

#### Scenario: Project with specs
- **WHEN** the user switches to the specs tab for a project that has specs
- **THEN** the left side SHALL list all spec names sorted alphabetically

#### Scenario: Project with no specs
- **WHEN** the user switches to the specs tab for a project with no specs
- **THEN** the tab SHALL display "No specs found"

#### Scenario: A capability with a parent segment
- **WHEN** the project holds capabilities `content/editing` and
  `content/pages`
- **THEN** the list SHALL show one header for `content`, saying it holds two,
  with `editing` and `pages` as the rows beneath it

#### Scenario: A capability with no parent segment
- **WHEN** the project holds a capability named `archive-from-ui`
- **THEN** it SHALL be drawn as a row, with no header of its own

#### Scenario: Where a group appears
- **WHEN** the list draws groups beside capabilities that have no parent
- **THEN** each group's header SHALL sit where its first capability sorts, and
  one group SHALL be separated from the next by a blank line

#### Scenario: A row wider than the pane
- **WHEN** a header or a row is wider than the list pane
- **THEN** it SHALL be truncated to the pane rather than widening it

### Requirement: Specs tab displays selected spec content
The specs tab SHALL display the content of the selected spec's `spec.md` file
on the right side, rendered as styled markdown.

#### Scenario: Spec selected
- **WHEN** a spec is highlighted in the list
- **THEN** the right side SHALL show the spec.md content rendered with markdown
  styling (headers, lists, bold, italic)

#### Scenario: Spec without spec.md
- **WHEN** a capability is listed but its spec.md cannot be read
- **THEN** the right side SHALL display "No spec.md found"

### Requirement: The spec list can be paged and jumped through
The spec list SHALL move its cursor by a page, by half a page, and to either
end, while it holds the keyboard.

#### Scenario: Paging the spec list
- **WHEN** the spec list holds the keyboard and a page or half-page key is
  pressed
- **THEN** the selected spec SHALL move, and the spec content SHALL NOT scroll

#### Scenario: A page counts the headers it draws
- **WHEN** the list is grouped and a page key is pressed
- **THEN** the cursor SHALL move by the capabilities that fit in the rows the
  pane shows, the headers and blank lines they draw included

#### Scenario: Jumping the spec list
- **WHEN** the spec list holds the keyboard and `gg` or `G` is pressed
- **THEN** the selection SHALL move to the first or last spec

#### Scenario: The cursor lands on capabilities only
- **WHEN** the list is grouped and the cursor is moved by any key
- **THEN** it SHALL select a capability, never a header or a blank line

#### Scenario: The content still pages when it holds the keyboard
- **WHEN** the spec content holds the keyboard and a page key is pressed
- **THEN** the content SHALL scroll, as it always has
