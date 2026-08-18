---
sidebar_position: 2
---

# Troubleshooting

## A configured target was not synced

Run `dotagen sync <target>` explicitly. Without an argument, sync uses
detected existing home-directory platform folders before it uses `config.yaml`.

## Sync refuses to overwrite a file

dotagen only replaces an existing symlink. Move or remove a manually managed
platform file, or choose a different agent name, then sync again.

## Status says missing or broken symlink

Run `dotagen sync` to recreate output and links. Use `dotagen clean` first if
you want to remove all managed links and generated output before recreating
them.

## A skill does not appear

Confirm the directory contains `SKILL.md`, its name appears under `skills` in
`config.yaml`, and its targets are not empty or disabled. `dotagen skill list`
shows the parsed installed skills and resolved targets.

## `init` would overwrite my setup

Answer anything other than `y` or `yes` at the prompt. Back up
`~/.dotagen/` before intentionally reinitializing; the overwrite path removes
the local source directories and config file.

## Dashboard is reachable from another machine

Stop `dotagen serve`. The server binds to all interfaces and has no
authentication. Do not run it on an untrusted or shared network.
