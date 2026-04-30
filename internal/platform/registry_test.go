package platform

import (
	"testing"

	"github.com/k0walski/dotagen/internal/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	r := NewRegistry()
	assert.Equal(t, 4, len(r.List()))

	_, err := r.Get("claude-code")
	assert.NoError(t, err)

	_, err = r.Get("unknown")
	assert.Error(t, err)
}

func TestClaudeCodeAdapter(t *testing.T) {
	a := NewClaudeCodeAdapter()
	assert.Equal(t, "claude-code", a.Name())

	ag := agent.Agent{Name: "test", Content: "# Test\nHello"}
	out, err := a.Render(ag)
	require.NoError(t, err)
	assert.Contains(t, out, "---")
	assert.Contains(t, out, "name: test")
	assert.Contains(t, out, "# Test\nHello")
	assert.Equal(t, "claude-code/test.md", a.OutputPath("test"))
	assert.Equal(t, ".claude/agents/test.md", a.SymlinkPath("test"))
}

func TestCursorAdapter(t *testing.T) {
	a := NewCursorAdapter()
	assert.Equal(t, "cursor", a.Name())

	ag := agent.Agent{Name: "test", Content: "# My Agent\nHello"}
	out, err := a.Render(ag)
	require.NoError(t, err)
	assert.Contains(t, out, "---")
	assert.Contains(t, out, "description: My Agent")
	assert.Contains(t, out, "alwaysApply: true")
	assert.Contains(t, out, "# My Agent")
	assert.Equal(t, "cursor/test.mdc", a.OutputPath("test"))
	assert.Equal(t, ".cursor/rules/test.mdc", a.SymlinkPath("test"))
}

func TestCursorAdapterWithDescription(t *testing.T) {
	a := NewCursorAdapter()
	ag := agent.Agent{
		Name:    "review",
		Content: "# Review\nContent",
		Frontmatter: map[string]string{"description": "Code reviewer"},
	}
	out, err := a.Render(ag)
	require.NoError(t, err)
	assert.Contains(t, out, "description: Code reviewer")
}

func TestGeminiCLIAdapter(t *testing.T) {
	a := NewGeminiCLIAdapter()
	assert.Equal(t, "gemini-cli", a.Name())

	ag := agent.Agent{Name: "test", Content: "# Test\nHello"}
	out, err := a.Render(ag)
	require.NoError(t, err)
	assert.Equal(t, "# Test\nHello", out)
	assert.Equal(t, "gemini-cli/test.md", a.OutputPath("test"))
	assert.Equal(t, ".gemini/agents/test.md", a.SymlinkPath("test"))
}

func TestOpenCodeAdapter(t *testing.T) {
	a := NewOpenCodeAdapter()
	assert.Equal(t, "opencode", a.Name())

	ag := agent.Agent{Name: "test", Content: "# My Agent\nHello"}
	out, err := a.Render(ag)
	require.NoError(t, err)
	assert.Contains(t, out, "---")
	assert.Contains(t, out, "description: My Agent")
	assert.Contains(t, out, "mode: subagent")
	assert.Contains(t, out, "# My Agent")
	assert.Equal(t, "opencode/test.md", a.OutputPath("test"))
	assert.Equal(t, ".opencode/agents/test.md", a.SymlinkPath("test"))
}
