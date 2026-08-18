#!/usr/bin/env python3
"""
Import agents and skills from donchitos/claude-code-game-studios into dotagen.

Source: https://github.com/donchitos/claude-code-game-studios (MIT)

  .claude/agents/*.md        -> internal/builtin/agents/dotagent-game-studios-{name}.md
  .claude/skills/*/SKILL.md  -> skillsrc/data/dotagent-game-studios-{name}/SKILL.md

Frontmatter is normalised to dotagen conventions:
  - agents:  name=dotagent-game-studios-{name}, category=game-development,
             list-valued keys (e.g. `skills: [...]`) are dropped because dotagen
             parses frontmatter as flat string->string.
  - skills:  name="dotagent:game-studios:{name}", category="Game Development",
             vendor="game-studios"; `agent:` references are rewritten to the
             prefixed agent name.

Usage:
    python3 scripts/import-game-studios.py [path/to/claude-code-game-studios]

If no path is given the repo is shallow-cloned into a temp directory.
"""

import re
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path

REPO_URL = "https://github.com/donchitos/claude-code-game-studios.git"
VENDOR = "game-studios"
AGENT_CATEGORY = "game-development"
SKILL_CATEGORY = "Game Development"

SCRIPT_DIR = Path(__file__).resolve().parent
ROOT = SCRIPT_DIR.parent
AGENTS_OUT = ROOT / "internal" / "builtin" / "agents"
SKILLS_OUT = ROOT / "skillsrc" / "data"

FM_RE = re.compile(r"^---\s*\n(.*?)\n---\s*\n?(.*)$", re.DOTALL)


def yaml_escape(s: str) -> str:
    return '"' + s.replace("\\", "\\\\").replace('"', '\\"') + '"'


def split_frontmatter(content: str) -> tuple[list[tuple[str, str]], str]:
    """Return ordered (key, raw_value) pairs from the leading frontmatter block + body."""
    m = FM_RE.match(content)
    if not m:
        return [], content
    pairs: list[tuple[str, str]] = []
    for line in m.group(1).split("\n"):
        if not line or line.startswith(" ") or line.startswith("#"):
            continue
        if ":" not in line:
            continue
        k, v = line.split(":", 1)
        pairs.append((k.strip(), v.strip()))
    return pairs, m.group(2)


def is_scalar(v: str) -> bool:
    return not (v.startswith("[") or v.startswith("{") or v == "")


def unquote(v: str) -> str:
    if len(v) >= 2 and v[0] == v[-1] and v[0] in "\"'":
        return v[1:-1]
    return v


def prefixed_agent(name: str) -> str:
    return f"dotagent-{VENDOR}-{name}"


def convert_agent(src: Path) -> str:
    base = src.stem
    pairs, body = split_frontmatter(src.read_text(encoding="utf-8"))
    fm = dict(pairs)
    desc = unquote(fm.get("description", ""))
    lines = [
        f"name: {prefixed_agent(base)}",
        f"description: {yaml_escape(desc)}",
        f"category: {AGENT_CATEGORY}",
    ]
    for k, v in pairs:
        if k in ("name", "description", "category"):
            continue
        if not is_scalar(v):
            continue  # e.g. skills: [a, b] — dotagen frontmatter is flat strings only
        lines.append(f"{k}: {v}")
    return "---\n" + "\n".join(lines) + "\n---\n\n" + body.lstrip("\n")


def convert_skill(src: Path, agent_names: set[str]) -> str:
    base = src.parent.name
    pairs, body = split_frontmatter(src.read_text(encoding="utf-8"))
    fm = dict(pairs)
    desc = unquote(fm.get("description", ""))
    lines = [
        f"name: {yaml_escape(f'dotagent:{VENDOR}:{base}')}",
        f"description: {yaml_escape(desc)}",
        f"category: {yaml_escape(SKILL_CATEGORY)}",
        f"vendor: {yaml_escape(VENDOR)}",
    ]
    for k, v in pairs:
        if k in ("name", "description", "category", "vendor"):
            continue
        if not is_scalar(v):
            continue
        if k == "agent" and unquote(v) in agent_names:
            v = prefixed_agent(unquote(v))
        lines.append(f"{k}: {v}")
    return "---\n" + "\n".join(lines) + "\n---\n\n" + body.lstrip("\n")


def main() -> int:
    tmp: str | None = None
    if len(sys.argv) > 1:
        repo = Path(sys.argv[1]).resolve()
    else:
        tmp = tempfile.mkdtemp(prefix="ccgs-")
        repo = Path(tmp) / "claude-code-game-studios"
        subprocess.run(["git", "clone", "--depth", "1", REPO_URL, str(repo)], check=True)

    agents_dir = repo / ".claude" / "agents"
    skills_dir = repo / ".claude" / "skills"
    if not agents_dir.is_dir() or not skills_dir.is_dir():
        print(f"error: {repo} does not look like claude-code-game-studios", file=sys.stderr)
        return 1

    agent_files = sorted(agents_dir.glob("*.md"))
    agent_names = {p.stem for p in agent_files}

    # Remove previously imported entries so deletions upstream propagate.
    for p in AGENTS_OUT.glob(f"dotagent-{VENDOR}-*.md"):
        p.unlink()
    for p in SKILLS_OUT.glob(f"dotagent-{VENDOR}-*"):
        if p.is_dir():
            shutil.rmtree(p)

    for src in agent_files:
        out = AGENTS_OUT / f"{prefixed_agent(src.stem)}.md"
        out.write_text(convert_agent(src), encoding="utf-8", newline="\n")

    n_skills = 0
    for skill_md in sorted(skills_dir.glob("*/SKILL.md")):
        out_dir = SKILLS_OUT / f"dotagent-{VENDOR}-{skill_md.parent.name}"
        out_dir.mkdir(parents=True, exist_ok=True)
        (out_dir / "SKILL.md").write_text(
            convert_skill(skill_md, agent_names), encoding="utf-8", newline="\n"
        )
        # copy any supporting files alongside SKILL.md
        for extra in skill_md.parent.rglob("*"):
            if extra.is_file() and extra.name != "SKILL.md":
                dst = out_dir / extra.relative_to(skill_md.parent)
                dst.parent.mkdir(parents=True, exist_ok=True)
                shutil.copy2(extra, dst)
        n_skills += 1

    print(f"imported {len(agent_files)} agents -> {AGENTS_OUT}")
    print(f"imported {n_skills} skills -> {SKILLS_OUT}")

    if tmp:
        shutil.rmtree(tmp, ignore_errors=True)
    return 0


if __name__ == "__main__":
    sys.exit(main())
