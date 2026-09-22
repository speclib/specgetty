---
# specgetty-y662
title: 'README: a landing page, not a manual'
status: todo
type: task
priority: normal
created_at: 2026-09-22T15:25:04Z
updated_at: 2026-09-22T15:25:04Z
parent: specgetty-ab9u
---

The README is 535 lines and 89% of it is reference manual. The content is good; it is in the wrong place.

OpenSpec change `restructure-the-readme`, skip_specs. Cuts the README to about 120 lines, moves 478 lines into five pages under docs/, and adds the five badges.

Three failures it fixes: OpenSpec is named eleven times and never explained or linked; the install section documents only `go install @master` while releases, archives and a nix flake exist; 2.5 MB of GIFs load on the landing page.

Ships after `add-continuous-integration`, which publishes three of the five badges.
