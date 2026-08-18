---
sidebar_position: 1
---

# Dashboard Overview

Start the embedded dashboard with:

```bash
dotagen serve
```

It provides an overview, agent and skill library views, rendered previews,
sync controls, and symlink status. It operates on the same `~/.dotagen/`
store as the CLI.

![Dashboard overview](/img/screenshots/overview.png)

![Agent library](/img/screenshots/agent-library.png)

![Skill library](/img/screenshots/skill-library.png)

![Symlink status](/img/screenshots/status.png)

## Security limitation

The dashboard REST service has no authentication and listens on `:<port>`, not
loopback-only. Although the CLI prints a localhost URL, hosts on a reachable
network can connect to it. Use only on a trusted local network until dotagen
adds stronger binding or authentication.

The [REST API reference](./rest-api.md) lists the available endpoints.
