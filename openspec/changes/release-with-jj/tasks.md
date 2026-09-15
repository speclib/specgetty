## 1. Implementation

- [ ] 1.1 Detect the VCS once at the top of `scripts/release.sh`: jj when `jj root` succeeds in the working directory, git otherwise, and report which one was chosen
- [ ] 1.2 Replace the clean-tree check with one that covers both tools, and that fails on untracked files as well as modified ones (`git status --porcelain` rather than `git diff --quiet`; `jj status` for jj)
- [ ] 1.3 Under jj, make the release commit with `jj commit -m "chore: release vX.Y.Z"`
- [ ] 1.4 Under jj, move the bookmark to the release commit with `jj bookmark set <branch> -r @-`, and push it with `jj git push --bookmark <branch>`
- [ ] 1.5 Under jj, create the tag with `jj tag set vX.Y.Z -r @-`; push it with `git push origin vX.Y.Z`, which works because a colocated workspace exports refs to git automatically (`jj git push` carries bookmarks only, never tags)
- [ ] 1.6 Keep the existing git path exactly as it is, so a plain git clone releases the way it does today
- [ ] 1.7 Make the CHANGELOG edit portable: the `\n` in a `sed` replacement is GNU-only and inserts a literal `n` on BSD, which is the userland on the macOS builds this project ships
- [ ] 1.8 Give every `sed -i` a form that works on both userlands, or replace the in-place edits with a small write-and-move helper
- [ ] 1.9 Reuse the existing branch detection for the bookmark name, so a jj release targets the same branch a git release would

## 2. Verification

- [ ] 2.1 `bash -n scripts/release.sh` parses
- [ ] 2.2 Dry run the whole thing in a throwaway jj colocated repo built in a temp directory: bump, commit, bookmark, tag, and push to a local bare remote, then assert with `git ls-remote` that both the branch and the `v*` tag arrived
- [ ] 2.3 Do the same in a throwaway plain git repo and confirm the result is identical to today's behaviour
- [ ] 2.4 Check the dirty-tree guard in both: one with a modified tracked file, one with only an untracked file, and confirm neither writes to `src/VERSION` or `CHANGELOG.md`
- [ ] 2.5 Confirm the changelog heading lands on its own line, by running the edit and reading the file rather than trusting the `sed`
- [ ] 2.6 `nix flake check` still passes afterwards, which is what catches the two vendor hashes drifting apart

## 3. Notes

- [ ] 3.1 The `tinychange` schema comes from https://github.com/speclib/openspec-tinychange-schema and is vendored at `openspec/schemas/tinychange/`; collaborators who do not have it can install it from that repo
