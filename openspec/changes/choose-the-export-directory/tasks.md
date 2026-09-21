## 1. Configuration

- [ ] 1.1 Replace `EditCommand` with an export directory in the parsed configuration, and document `export_dir` in `src/config.yml` where `edit_command` was documented
- [ ] 1.2 Expand `~` and environment variables in the export directory, wherever it came from. Proven for both forms, from the configuration and from the prompt
- [ ] 1.3 Add `edit_command` to the retired keys, and generalise the report so it names every retired key it finds rather than one. Proven with a configuration carrying both, asserting both are named in a stable order
- [ ] 1.4 Confirm a configuration carrying neither reports nothing

## 2. The prompt

- [ ] 2.1 Replace the export confirmation with a prompt holding an editable directory, opening on the configured one, or on the home directory when none is configured
- [ ] 2.2 Show the generated filename beside the field as fixed text, so what is editable and what is not is visible
- [ ] 2.3 Complete the directory against the filesystem on `tab`, feeding the text input's own suggestions rather than building completion
- [ ] 2.4 Export on `enter`, cancel on `esc` leaving nothing written
- [ ] 2.5 Name the submit, cancel and complete keys in the modal, since a field cannot imply them
- [ ] 2.6 Give every printable key to the field while the prompt is up, so `a`, `d`, `e`, `q` and the rest do not reach the actions they are bound to. Proven for each of those keys
- [ ] 2.7 Keep the frame exactly its terminal's rows and columns with the prompt up, at three widths including the 60-column minimum
- [ ] 2.8 Raise no other modal while the prompt is up

## 3. Refusing what cannot be used

- [ ] 3.1 Refuse a directory that does not exist, naming it, and leave what was typed in the field to be corrected. Proven by asserting the prompt is still up and still holds the text
- [ ] 3.2 Refuse a path that exists and is not a directory, the same way
- [ ] 3.3 Create nothing on either path. Proven by asserting the filesystem is unchanged
- [ ] 3.4 Ask before replacing an existing file, naming the file
- [ ] 3.5 Write nothing and return to the prompt when the user declines, with the directory still in the field
- [ ] 3.6 Replace the file and report where it wrote when the user agrees

## 4. Writing it

- [ ] 4.1 Write into the chosen directory, with the filename this capability already specifies. Proven for an active change and an archived one, asserting the archived name loses its date prefix
- [ ] 4.2 Confirm a redirected export leaves the configuration alone, and that the next export opens on the configured directory rather than on what was typed
- [ ] 4.3 Report the full path it wrote on success, and the reason on failure

## 5. Verification

- [ ] 5.1 Export a change from this project into a directory that is not home, and confirm the zip is there and opens
- [ ] 5.2 Export it again and confirm the replace question appears, and that declining leaves the first zip untouched
- [ ] 5.3 Type a directory that does not exist and confirm the message names it and nothing is created
- [ ] 5.4 Start `spg` with a configuration still carrying `edit_command` and confirm it is reported
- [ ] 5.5 Revert task 3.1 and confirm the refusal test fails. Check the build succeeds first
- [ ] 5.6 `nix flake check` passes, coverage floors included

## 6. Notes

- [ ] 6.1 Document `export_dir` and the prompt in the README, where the export key is described
- [ ] 6.2 Create a follow-up bean for the retired-setting requirement's home. It lives in `change-list-view` because `change_mode` was introduced there, and now serves `edit_command`, which has nothing to do with the change list. Moving it means removing and re-adding it under a capability that does not exist yet
