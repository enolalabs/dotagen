---
sidebar_position: 2
---

# Skills

A skill is a directory below `~/.dotagen/skills/` with a required `SKILL.md`
entry point. The directory name is the skill identifier.

```text
~/.dotagen/skills/ds-release-check/
├── SKILL.md
└── references/
    └── checklist.md
```

`SKILL.md` may have YAML frontmatter such as `name`, `description`, and
`category`. Every other file below the skill directory is collected as a
reference and copied into the generated skill directory during sync.

```markdown
---
name: ds-release-check
description: Verify a release before publishing it.
category: engineering
---

# Release Check

Run the project's required build and test commands before release.
```

Skills are enabled independently from agents in the `skills` section of
`config.yaml`. A platform without a skill adapter is skipped during sync; all
eight current platform adapters implement skill support.
