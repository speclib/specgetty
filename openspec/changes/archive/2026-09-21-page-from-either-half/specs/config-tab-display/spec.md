## ADDED Requirements

### Requirement: The properties tab pages from either half
The properties tab's content SHALL respond to the paging and jump keys whether
the row list or the content holds the keyboard.

Its list is a short fixed set of rows and the document beside it is what the tab
is for, so requiring the keyboard to be moved before a page can be turned was
the regression this restores.

#### Scenario: Paging straight after opening the tab
- **GIVEN** the properties tab has just been opened, so its row list holds the
  keyboard
- **WHEN** the user presses `pgdown`
- **THEN** the row's content SHALL scroll

#### Scenario: Choosing a row still needs the list
- **GIVEN** the properties tab's row list holds the keyboard
- **WHEN** the user presses `j`
- **THEN** the selected row SHALL move and the content SHALL NOT scroll
