---
sidebar_position: 2
---

# Initialize

Run:

```bash
dotagen init
```

This creates `~/.dotagen/`, imports the built-in agents and skills, creates an
empty `.generated/` directory, and writes `config.yaml`. All imported agent and
skill entries start with empty `targets`, so they are disabled until you opt in.

The command checks for existing platform directories in your home directory and
uses the detected platform IDs as the initial top-level `targets` list. If none
are detected, it writes all eight supported target IDs.

## Existing stores

If `~/.dotagen/` already exists, `init` asks `Overwrite? (y/N)`. Answering
`y` or `yes` removes dotagen-managed symlinks, generated contents, agent and
skill source directories, and `config.yaml` before reinitializing. Answering
anything else leaves the store unchanged.

Most commands auto-initialize a missing `~/.dotagen/` store. Run `init`
explicitly when you need to inspect the result or intentionally overwrite an
existing store.

Next, [configure targets](./configuration.md) and run `dotagen sync`.
