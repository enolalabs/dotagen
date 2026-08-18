---
sidebar_position: 2
---

# REST API

The dashboard serves JSON from the same port as the UI. There is no
authentication. Requests can read and modify `~/.dotagen/`, create or remove
links, and initialize the store, so do not expose this service beyond a trusted
local network.

| Method | Endpoint | Action |
| --- | --- | --- |
| `GET` | `/api/config` | Read configuration |
| `PUT` | `/api/config` | Replace configuration |
| `GET` | `/api/agents` | List agents |
| `GET` | `/api/agents/{name}` | Read an agent |
| `POST` | `/api/agents` | Create an agent |
| `PUT` | `/api/agents/{name}` | Update an agent |
| `DELETE` | `/api/agents/{name}` | Delete an agent |
| `GET` | `/api/skills` | List skills |
| `GET` | `/api/skills/{name}` | Read a skill |
| `POST` | `/api/skills` | Create a skill |
| `PUT` | `/api/skills/{name}` | Update a skill |
| `DELETE` | `/api/skills/{name}` | Delete a skill |
| `GET` | `/api/targets` | List supported target IDs |
| `GET` | `/api/preview/{agent}/{target}` | Preview an agent render |
| `GET` | `/api/preview/skill/{skill}/{target}` | Preview a skill render |
| `POST` | `/api/sync` | Sync agents and skills |
| `POST` | `/api/sync/{target}` | Sync one target |
| `POST` | `/api/sync-skills` | Sync skills only |
| `POST` | `/api/clean` | Remove managed links and generated output, then clear all entry targets |
| `POST` | `/api/clean-broken` | Repair or remove broken links |
| `GET` | `/api/status` | List managed link status |
| `POST` | `/api/toggle` | Enable or disable agent or skill targets |
| `POST` | `/api/init` | Reinitialize the store |

The API uses JSON request and response bodies. Invalid JSON and invalid target
or name values return HTTP 400; missing definitions return HTTP 404. Several
mutating operations serialize access with the server's internal lock. Unlike
the CLI `dotagen clean`, `POST /api/clean` also sets every configured agent and
skill target list to empty.
