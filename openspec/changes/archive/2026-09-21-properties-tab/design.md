## Context

Three decisions from the exploration behind `specgetty-m4vt`, and one fact that
changed what the tab should hold.

The fact first, because it removes work rather than adding it. `follow-the-store`
gave a store-backed project two configuration panes on the belief that the
pointing repo's `context:`, `rules:` and `operations:` still applied to it. They
do not. A sandbox with a store and a repo, each carrying a distinctly marked
context and rules, run through `openspec instructions proposal`:

```
  1 STORE_CONTEXT_MARKER        1 STORE_RULE_MARKER
  0 REPO_CONTEXT_MARKER         0 REPO_RULE_MARKER
```

The source says the same. Every consumer calls `readProjectConfig(root.path)`
with the resolved root; the pointing repo's file is read in one place,
`readStorePointer`, which looks at one key. OpenSpec's own name for the
remainder, in `commands/doctor.js`, is `inertPointerDeclarations`.

So there is only ever one configuration in force, and the second pane was
showing a file that does nothing.

## Decisions

### The list is a fixed shape, and `store: local` is doing work

```
  project        the one configuration that applies
  SCHEMAS
    spec-driven  24 changes      <- one row per schema the project uses
    tinychange   10 changes
  store          local, or the store this content comes from
```

The pane set used to be one, two or three rows depending on whether the project
had a store. Now it is always the same shape. A person who has never met a store
opens the tab, sees `store: local`, and has learned that content can come from
somewhere else. A blank or absent row teaches nothing.

The schema rows are the exception to fixed: there are as many as the project
uses. A sectioned list rather than one long document, because `spec-driven`'s
definition runs to about eighty lines of artifact prose, and stacking two
schemas in one pane would mean scrolling past the first to reach the second.

### Shell out to locate, parse to read

Built-in schemas live inside the OpenSpec installation:

```
  schema: tinychange                     schema: spec-driven
  source: project                        source: package
  path:   openspec/schemas/tinychange/   path:   /nix/store/phgx...-openspec-1.10.0/
          ^ specgetty can find this              lib/openspec/schemas/spec-driven/
                                                 ^ a hash that changes on upgrade
```

`openspec schema which <name> --json` returns that path, and specgetty reads
`schema.yaml` from it with the parser it already has. That is the same shape the
store work settled on: something else locates, specgetty parses. For stores the
locator was a file at a documented XDG path; here it is a subprocess, because
nothing about the package's location is documented or stable.

`openspec schemas --json` was the obvious alternative and is the wrong call. It
returns artifact ids and no path, so the definition stays out of reach, and it
does not report `shadows`, which is how a project that has forked a built-in
would be recognised.

### Once per project, concurrently, and only if the tab is opened

The subprocess costs about 0.7 seconds, measured at 0.79, 0.69 and 0.72. That is
node cold start, and it is not small next to a tool whose startup deliberately
never walks the disk. Three rules keep it affordable:

```
  when      first time the properties tab is opened for a project
  how many  one per used schema, concurrently, so the wait is the slowest
  until     the project changes; nothing else invalidates it
```

Most `spg` runs never open this tab and pay nothing. A project on two schemas
waits once, for about as long as a project on one.

Nothing re-reads a schema while a project is open, and that is deliberate rather
than unfinished. The watcher covers `<root>/openspec` recursively, so editing a
schema file already triggers a rescan and the changes and specs tabs will
refresh while these rows do not. Schemas are workflow definitions: they change
once in a project's life, where changes and specs move all day. Paying for live
invalidation here would be machinery for an event that does not happen. Opening
the project again is the way to see an edit.

### The capability keeps its old name

The tab is renamed and the capability is not. OpenSpec's RENAMED operation works
on requirements and has no capability-level form, and removing every requirement
from `config-tab-display` to recreate them elsewhere would leave a spec with an
empty requirements section, which `validate --strict` rejects. Renaming the
directory by hand would make the archived deltas point at a capability that no
longer exists.

So `openspec/specs/config-tab-display/` keeps its path and its Purpose says what
the tab is called. This is a wart, written down here so the next reader knows it
was a choice rather than an oversight.

## Risks

### The first asynchronous tab

Every tab today renders synchronously out of `ProjectInfo`. A schema row has
states that no pane has needed before:

```
   unread ──open tab──▶ reading ──┬──▶ read
                                  ├──▶ no openspec binary
                                  ├──▶ schema does not resolve (with alternatives)
                                  ├──▶ timed out
                                  └──▶ schema.yaml unparseable
```

All five of the failures are reachable and all five are specified. The one to
get wrong is the timeout: a hung subprocess must not take the interface with it.

### A second split tab, and the focus constants

`focusSpecsList` and `focusSpecsContent` were named when one tab had halves.
Two do now, so they become a generic pair. Three places currently compare
against them without also checking which tab is active; those are the lines
where this goes wrong quietly rather than loudly.

### The list steals width from paths

The specs tab gives its list thirty percent, which is right for capability names
and wrong for three short words:

```
  terminal   specs rule   label-sized   what the content holds
  60         35           40            store paths, wrapped either way
  92         57           72            store paths, unwrapped at 72
```

The store row shows absolute paths, which is the content that suffers most from
a narrow column, so this tab sizes its list to its labels.

### The inert warning could nag

Every store-backed repo in this author's tree carries inert declarations right
now, so the warning will fire on all of them at once. That is the point, but it
means the first impression of the feature is a row full of complaint. It belongs
under `store`, where someone goes to ask how the project is wired, rather than
anywhere it would be seen without asking.

### Rollback

Its own commit. Reverting restores the horizontal sub-tabs and the two
configuration panes, and takes the schema rows with it. The scanner's change
parsing gains one field that nothing else reads, so a revert of the UI alone
would leave it harmless.
