## Why

Bean `specgetty-m4vt`. The config tab has grown past its name and its axis.

```
  today                                    wanted
  [changes] [specs] [config]               [changes] [specs] [properties]
   repo  store  store details              ╭──────────────╮ ╭──────────────╮
  ╭──────────────────────────────╮         │ project      │ │              │
  │ schema: spec-driven          │         │              │ │              │
  │ store: nivis-tunnel          │         │ SCHEMAS      │ │              │
  ╰──────────────────────────────╯         │  spec-driven │ │              │
   horizontal, 1 to 3 sub-tabs,            │  tinychange  │ │              │
   config files only                       │              │ │              │
                                           │ store        │ │              │
                                           ╰──────────────╯ ╰──────────────╯
```

Three things are wrong with the tab as it stands.

**It is misnamed.** It already shows a store's git state, which is not
configuration, and it is about to show workflow schemas, which are not either.
`config` also collides with specgetty's own `~/.config/specgetty/config.yml`,
which the UI never shows.

**It shows a file that does nothing.** `follow-the-store` gave a store-backed
project two config panes on the belief that the repo's own `context:`, `rules:`
and `operations:` still applied. They do not. openspec reads the pointing repo's
config for exactly one key, `store:`, and takes everything else from the
resolved root. Proven in a sandbox: a store and a repo with distinctly marked
context and rules, and `openspec instructions` injects the store's and never the
repo's. openspec's own term for the rest is `inertPointerDeclarations`.

**Schemas are invisible.** A change's workflow schema decides which artifacts it
needs, and specgetty has never read one. This project runs 24 changes on
`spec-driven` and 10 on `tinychange` and the UI says nothing about either.

## What Changes

- The tab is called `properties`, and its sub-navigation becomes a vertical list
  beside the content, the same split the specs tab already uses
- `project` shows the one configuration that applies, which is the resolved
  root's. For a store-backed repo that is the store's file, not the repo's
- `store` shows the store's identity, root, git state and the file that declared
  it, and reports what that file declares in vain. openspec warns about inert
  `references` and says nothing about inert `context`, `rules` or `operations`,
  so this is the only place a person would ever find out
- `store` reads `local` for a project that holds its own content, which is what
  tells someone who has never met a store that the concept exists
- One row per schema the project actually uses, with its source, artifact chain,
  apply rule and how many changes are on it. Used, not available: specgetty has
  already read every change, so it is the only surface that can say which
  schemas are in play
- A schema's definition is located by `openspec schema which <name> --json` and
  then read by specgetty's own YAML parser. Shell out to locate, parse to read,
  which is the rule the store work already follows with the registry standing in
  for the subprocess
- An open change names the schema it uses
- The sentence `openspec-scanning` picked up from `follow-the-store`, that a
  pointing repo is "the only place that repo's own context and rules can be read
  from", is corrected. It is true and it endorses something inert

Not in scope: the change list's `schema` column, which belongs with the three
other columns waiting in `specgetty-vru8` and shares its field-registry
mechanics.

Not in scope: surfacing specgetty's own configuration. The tab is about the
project, and mixing the tool's settings into it is what made `config` ambiguous.

Not in scope: rendering a schema's `instruction` prose. `spec-driven` carries
about twenty lines per artifact, which is guidance for whoever writes the
artifact, not a property of the project.

## Capabilities

### Added Capabilities
- `schema-inspection`: which schemas a project's changes use, how a schema
  definition is located and read, and when that is done

### Modified Capabilities
- `config-tab-display`: the tab is renamed and restructured. The capability
  keeps its directory name because OpenSpec's RENAMED operates on requirements
  and has no capability-level form; the Purpose says what the tab is called
- `detail-tabs`: the third tab is `properties`
- `change-list-view`: an open change names its schema
- `openspec-scanning`: the inert-config sentence is corrected

## Impact

- `src/ui/configpanes.go`: panes become rows, and the pane set stops depending
  on store-backedness
- `src/ui/ui.go`: a second tab becomes a list/content split, so
  `focusSpecsList` and `focusSpecsContent` are renamed to a generic pair
- `src/scanner/`: `.openspec.yaml` read per change for its schema, and a new
  resolver that shells out once per used schema
- The properties tab is the first tab whose content arrives asynchronously

## Rollback

Its own commit. Reverting restores the horizontal sub-tabs and the two config
panes, and takes the schema rows with it. Nothing outside the tab and the
scanner's change parsing is touched.
