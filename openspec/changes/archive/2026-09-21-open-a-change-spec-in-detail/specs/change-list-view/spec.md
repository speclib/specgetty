## ADDED Requirements

### Requirement: Enter on the specs sub-tab opens the change's deltas
The specs sub-tab of an open change SHALL treat `enter` the way the specs tab
does: it descends into what the pane is listing. The view it opens is described
by `change-spec-detail-view`.

#### Scenario: Enter on the specs sub-tab
- **GIVEN** a change is open with its specs sub-tab active
- **WHEN** the user presses `enter`
- **THEN** the change's spec deltas SHALL open in the change spec detail view

#### Scenario: Enter on an artifact sub-tab
- **WHEN** a change is open with a proposal, design or tasks sub-tab active and
  the user presses `enter`
- **THEN** nothing SHALL happen, those panes being one document with nothing
  below them

#### Scenario: Returning from the view
- **WHEN** the user leaves the change spec detail view with `esc`
- **THEN** the change SHALL be shown again with its specs sub-tab active
