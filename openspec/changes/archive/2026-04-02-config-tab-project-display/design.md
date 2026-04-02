## Context

The detail panel has 5 tabs. Only overview is functional. The config tab should display the project's primary configuration file with appropriate formatting.

OpenSpec projects use either `config.yaml` (structured config) or `project.md` (markdown description), or both. When both exist, `project.md` takes priority as it's the more human-readable format.

## Goals / Non-Goals

**Goals:**
- Read config file contents at scan time and store in ProjectInfo
- Render markdown with basic styling: headers (bold/colored), lists (indented), bold/italic text
- Render YAML with syntax highlighting: keys in one color, values in another, comments dimmed
- Scrollable content via the detail viewport when content exceeds panel height

**Non-Goals:**
- Full markdown rendering (images, tables, code blocks with language detection)
- YAML validation or schema checking
- Editing configuration from the TUI

## Decisions

### 1. Read file content at scan time
Store the raw file content and file type (md/yaml/none) in `ProjectInfo`. This avoids re-reading files on every tab switch.

### 2. Simple line-by-line rendering, no markdown parser dependency
For markdown: detect `#` headers, `- ` list items, `**bold**`, `_italic_` patterns per line and apply lipgloss styles. This is intentionally simple — good enough for config files without adding a dependency.

For YAML: detect `key:` patterns (style the key), `#` comments (dim), and string values. Line-by-line processing, no YAML parser needed for display.

### 3. Priority: project.md over config.yaml
If both exist, show project.md. Show a dim note at the top indicating which file is being displayed.

## Risks / Trade-offs

- **Trade-off**: Simple line-by-line rendering won't handle multi-line markdown constructs (e.g. nested lists, multi-paragraph items). Acceptable for config files which are typically simple.
