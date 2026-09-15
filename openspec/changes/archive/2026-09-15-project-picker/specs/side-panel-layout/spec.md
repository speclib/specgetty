## REMOVED Requirements

### Requirement: Horizontal split layout
**Reason**: The left panel is removed. The project view fills the width, as
specified by `project-view`, which also carries the minimum terminal size.

### Requirement: Project names as basenames
**Reason**: Moved to `project-picker`, which is now the only place project names
are listed.

### Requirement: Detail panel shows project info and file listing
**Reason**: The persistent header moves to `project-view`. The file listing half
of this requirement describes a view that no longer exists.

### Requirement: Navigation works across both panels
**Reason**: There is one panel. Switching focus between a project list and a
detail panel has nothing left to switch between.
