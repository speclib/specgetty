---
# specgetty-z97b
title: 'CI: run the gate on every push, publish badges'
status: in-progress
type: task
priority: normal
created_at: 2026-09-22T15:25:04Z
updated_at: 2026-09-22T15:25:13Z
parent: specgetty-ab9u
---

Nothing runs on push: .github/workflows holds only release.yml, on v* tags. No run says the tests pass, no number says what coverage is, and a pull request would be judged by reading it.

OpenSpec change `add-continuous-integration`. New capability `continuous-integration`: the gate on every push and pull request, the coverage number CI measured, the OpenSpec metrics, and the rule that decoration cannot fail the gate.

CI runs `nix flake check` and nothing else, the same command scripts/ship-change.sh runs, so that what CI enforces and what shipping enforces cannot drift apart.
