## MODIFIED Requirements

### Requirement: Project view fills the terminal
The project view SHALL occupy the full terminal width and height, with no
project list beside it and no panel below it.

#### Scenario: Layout
- **WHEN** the project view is displayed
- **THEN** it SHALL span the full width, with the nav bar below it and nothing
  between them

#### Scenario: Terminal too small
- **WHEN** the terminal is narrower than 60 columns or shorter than 20 rows
- **THEN** the application SHALL display a message indicating the terminal is
  too small
