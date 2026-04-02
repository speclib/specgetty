## 1. Scanner: config file detection and reading

- [x] 1.1 Add `ConfigFile` (filename) and `ConfigContent` (raw string) fields to `ProjectInfo`
- [x] 1.2 In `ParseProjectInfo`, detect and read `project.md` (priority) or `config.yaml`, store content and filename
- [x] 1.3 Add test for `ParseProjectInfo` config file detection with both, only yaml, only md, and neither

## 2. UI: config tab rendering

- [x] 2.1 Add `renderConfigTab(width, height int) string` that dispatches to markdown or YAML renderer based on `ConfigFile`
- [x] 2.2 Add `renderMarkdown(content string, width int) string` — line-by-line: style `#` headers, `- `/`* ` lists, `**bold**`, `_italic_`
- [x] 2.3 Add `renderYAML(content string, width int) string` — line-by-line: style keys before `:`, dim `#` comments, color values
- [x] 2.4 Wire `renderConfigTab` into `renderDetailPanel` for `tabConfig` case
- [x] 2.5 Show dimmed file source indicator at top of config tab content

## 3. Tests and verification

- [x] 3.1 Add test for `renderMarkdown` with headers, lists, bold
- [x] 3.2 Add test for `renderYAML` with keys, comments, values
- [x] 3.3 Verify build and all tests pass
