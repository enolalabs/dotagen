---
sidebar_position: 4
---

# Built-In Catalog

`dotagen init` imports the built-in agents and the built-in skills embedded in
the binary into `~/.dotagen/agents/` and `~/.dotagen/skills/`. Imported entries
are listed in `config.yaml` with empty target lists, so initialization alone
does not deploy them.

Use the installed catalog rather than a published count as the source of truth:

```bash
dotagen skill list
```

The command shows each installed skill, its frontmatter category, resolved
targets, and reference-file count. You can edit imported definitions locally or
create new ones.

`docs/CATALOG.vi.md` is a legacy, historical catalog. Its counts must not be
treated as current facts about the built-in catalog.
