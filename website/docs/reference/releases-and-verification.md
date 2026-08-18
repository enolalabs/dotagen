---
sidebar_position: 3
---

# Releases And Verification

Use the GitHub Releases page to obtain release binaries. After installation or
an update, verify the executable and supported target list:

```bash
dotagen version
dotagen help
```

`dotagen update` queries the latest GitHub release, downloads the binary that
matches the running OS and architecture, and replaces the current executable.
If it cannot find a matching release asset or cannot replace the binary, it
returns an error.

For a source checkout, run the same checks used for documentation and runtime
validation:

```bash
npm run build --prefix website
go test ./...
```

The Docusaurus build is configured to fail on broken links. The Go suite covers
the CLI, configuration, rendering, and platform behavior.
