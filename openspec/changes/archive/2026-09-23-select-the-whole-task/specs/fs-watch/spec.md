## ADDED Requirements

### Requirement: A refresh of content already on screen interrupts nothing
Reading a project for the first time and re-reading one that is already on
screen are two different events and SHALL be presented as two different events.

The first happens with an empty screen, may search the whole filesystem, and is
the only thing the user is waiting for, so it raises the scan indicator and
answers no key but `q` and `ctrl+c`. The second happens with the content on
screen and the user's hands on the keyboard. Raising a modal over what is being
read, and discarding the keys pressed while it is up, makes the application
appear to stutter and lose presses at the exact moment the user is working
quickly.

#### Scenario: A rescan while a project is open
- **WHEN** the open project is read again, whether because a watched file
  changed or because the application asked for it after writing one
- **THEN** no modal SHALL be raised, and what is on screen SHALL stay on screen
  until the new content replaces it

#### Scenario: Keys during a rescan
- **WHEN** a key is pressed while the open project is being read again
- **THEN** it SHALL be handled, rather than discarded

#### Scenario: The first scan still says so
- **WHEN** the application is searching for projects and has none to show
- **THEN** the scan indicator SHALL be raised as it is today, and only `q` and
  `ctrl+c` SHALL be answered

#### Scenario: Repeated writes
- **WHEN** the user changes several things in quick succession, each write
  prompting a read
- **THEN** every keystroke SHALL take effect, and the display SHALL come to show
  the final state of the file
