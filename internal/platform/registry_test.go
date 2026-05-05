package platform

import (
	"testing"

	"github.com/enolalabs/dotagen/v2/internal/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	r := NewRegistry()
	assert.Equal(t, 5, len(r.List()))

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

func TestCodexAdapter(t *testing.T) {
	a := NewCodexAdapter()
	assert.Equal(t, "codex", a.Name())

	ag := agent.Agent{Name: "test", Content: "# Test\nHello"}
	out, err := a.Render(ag)
	require.NoError(t, err)
	assert.Contains(t, out, "---")
	assert.Contains(t, out, "name: test")
	assert.Contains(t, out, "# Test\nHello")
	assert.Equal(t, "codex/test.md", a.OutputPath("test"))
	assert.Equal(t, ".codex/agents/test.md", a.SymlinkPath("test"))
}

func TestCodexAdapterSkillAdapter(t *testing.T) {
	a := NewCodexAdapter()
	// Verify it implements SkillAdapter
	var _ SkillAdapter = a

	assert.Equal(t, "codex/skills/ds-my-skill", a.SkillOutputDir("ds-my-skill"))
	assert.Equal(t, ".codex/skills/ds-my-skill", a.SkillSymlinkDir("ds-my-skill"))
}



func TestGeminiCLIAdapter(t *testing.T) {
	a := NewGeminiCLIAdapter()
	assert.Equal(t, "gemini-cli", a.Name())

	ag := agent.Agent{Name: "test", Content: "# Test\nHello"}
	out, err := a.Render(ag)
	require.NoError(t, err)
	assert.Contains(t, out, "---")
	assert.Contains(t, out, "name: test")
	assert.Contains(t, out, "# Test\nHello")
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
	assert.Equal(t, ".config/opencode/agents/test.md", a.SymlinkPath("test"))
}

func TestAntigravityAdapter(t *testing.T) {
	a := NewAntigravityAdapter()
	assert.Equal(t, "antigravity", a.Name())

	ag := agent.Agent{Name: "test", Content: "# Test\nHello"}
	out, err := a.Render(ag)
	require.NoError(t, err)
	// Antigravity uses plain markdown, no frontmatter
	assert.Equal(t, "# Test\nHello", out)
	assert.Equal(t, "antigravity/test.md", a.OutputPath("test"))
	assert.Equal(t, ".agents/test.md", a.SymlinkPath("test"))
}

func TestAntigravityAdapterSkillAdapter(t *testing.T) {
	a := NewAntigravityAdapter()
	// Verify it implements SkillAdapter
	var _ SkillAdapter = a

	assert.Equal(t, "antigravity/skills/ds-my-skill", a.SkillOutputDir("ds-my-skill"))
	assert.Equal(t, ".agents/skills/ds-my-skill", a.SkillSymlinkDir("ds-my-skill"))
}

