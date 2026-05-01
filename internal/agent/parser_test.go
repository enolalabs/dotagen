package agent

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAgentFile(t *testing.T) {
	dir := t.TempDir()
	content := `---
description: Test agent
---

# Test Agent

This is a test.`
	err := os.WriteFile(filepath.Join(dir, "test.md"), []byte(content), 0o644)
	require.NoError(t, err)

	a, err := ParseAgentFile(filepath.Join(dir, "test.md"))
	require.NoError(t, err)
	assert.Equal(t, "test", a.Name)
	assert.Equal(t, "Test agent", a.Frontmatter["description"])
	assert.Contains(t, a.Content, "# Test Agent")
}

func TestParseAgentFileNoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	content := "# Simple Agent\n\nNo frontmatter here."
	err := os.WriteFile(filepath.Join(dir, "simple.md"), []byte(content), 0o644)
	require.NoError(t, err)

	a, err := ParseAgentFile(filepath.Join(dir, "simple.md"))
	require.NoError(t, err)
	assert.Equal(t, "simple", a.Name)
	assert.Contains(t, a.Content, "# Simple Agent")
}

func TestParseAgentsDir(t *testing.T) {
	dir := t.TempDir()
	agentsDir := filepath.Join(dir, "agents")
	require.NoError(t, os.MkdirAll(agentsDir, 0o755))

	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "a.md"), []byte("# A"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "b.md"), []byte("# B"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "readme.txt"), []byte("not md"), 0o644))

	agents, err := ParseAgentsDir(agentsDir)
	require.NoError(t, err)
	assert.Equal(t, 2, len(agents))
}
