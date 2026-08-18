# Hallmark built-in skill design

## Goal

Package the MIT-licensed Hallmark design skill from `Nutlope/hallmark` as a
Dotagent built-in skill. It must be available after `dotagen init` but remain
disabled until the user assigns targets in `.dotagen/config.yaml`.

## Scope

- Add a point-in-time snapshot of the Hallmark skill and its required
  `references/` files under `skillsrc/data/dotagent-nutlope-hallmark/`.
- Normalize the skill frontmatter for Dotagent discovery: `name` is
  `dotagent:nutlope:hallmark`, `vendor` is `nutlope`, `category` is
  `Frontend & UI`, and the description and license metadata identify the
  imported skill and its MIT source. Retain the upstream operational
  instructions and relative reference links.
- Keep the snapshot manual: do not alter `scripts/fetch-official-skills.py` or
  introduce a network dependency for builds, installs, or runtime use.
- Update public catalog documentation and built-in-skill counts to include the
  Nutlope vendor and Hallmark in the Frontend & UI category.
- Add focused tests covering embedded discovery, copied reference files, and
  default disabled configuration, metadata, snapshot integrity, and size.

## Non-goals

- Automatically checking, downloading, or updating the upstream Hallmark
  repository.
- Bundling Hallmark's website, demos, or other development-only repository
  content.
- Enabling Hallmark automatically for any platform.

## Architecture and data flow

`skillsrc.DefaultSkills` embeds all contents of `skillsrc/data`. The new
`dotagent-nutlope-hallmark` directory therefore becomes discoverable through
`skillsrc.ListSkills()` without a runtime code change. During `dotagen init`,
`runInit` copies every embedded file from the skill directory into
`~/.dotagen/skills/dotagent-nutlope-hallmark/` and writes a matching config
entry with `targets: []`.

The snapshot preserves the layout expected by Hallmark:

```text
skillsrc/data/dotagent-nutlope-hallmark/
├── SKILL.md
└── references/
    └── ... upstream reference files ...
```

This layout keeps relative links from `SKILL.md` valid after initialization and
allows all supported platform renderers to consume the same source material.
The files are copied only during `dotagen init`, using the existing one-pass
embedded-file copy; normal CLI commands do not copy Hallmark again.

## Snapshot rules

Import only `skills/hallmark/SKILL.md` and files it requires below
`skills/hallmark/references/` from an immutable commit on the upstream
`https://github.com/Nutlope/hallmark` repository. Add a committed `SOURCE.md`
alongside the snapshot that records the repository URL, source directory,
commit SHA, snapshot date, MIT license/copyright notice, exact intentional
frontmatter normalization, complete imported path list, and SHA-256 digest of
each imported file. The manifest is the source of truth for the snapshot's
file count and total byte size; tests must enforce those recorded values.

Before merging a new snapshot, a maintainer must fetch the pinned revision,
review the changed file set and manifest, confirm the MIT license and
attribution, and inspect imported text for executable content, remote prompt
loading, credential collection/exfiltration instructions, and unsafe command
execution. They then update the manifest and catalog counts, run the named Go
tests, and review the resulting diff. This is the complete manual update
procedure; no updater or network dependency is added to Dotagent.

If a required source file is binary, exceeds the snapshot's reviewed manifest
size, or contains an absolute path that prevents use after init, the import or
PR validation must fail before merge and report the path and reason. Do not
silently omit required references or rewrite Hallmark's operational guidance.

## Documentation

Update `README.md` to reflect the new total built-in skill and vendor counts,
the Frontend & UI category count, and a Nutlope vendor entry. Preserve
`docs/CATALOG.vi.md` as a legacy catalog: add a short scope note that directs
readers to the README for the current built-in library, rather than adding a
Hallmark entry or attempting to modernize its legacy counts and naming scheme.

## Validation

Add or extend Go tests to verify that:

1. `skillsrc.ListSkills()` returns `dotagent-nutlope-hallmark`.
2. The embedded skill and the `runInit` output each expose every path listed in
   `SOURCE.md`; their SHA-256 digests, file count, and total bytes match the
   manifest.
3. The frontmatter identifies the skill as `dotagent:nutlope:hallmark`, with
   the normalized Nutlope, Frontend & UI, description, and license metadata.
4. `runInit` copies the skill layout and emits the skill's config entry with
   `targets: []`.

Run `go test ./skillsrc ./internal/cli` and `go test ./...`; CI must enforce
these checks. Keep unrelated working-tree changes untouched.
