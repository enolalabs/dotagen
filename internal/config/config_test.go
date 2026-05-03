package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	content := `
targets:
  - claude-code
  - gemini-cli
agents:
  review:
    targets: all
  testing:
    targets:
      - claude-code
`
	err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0o644)
	require.NoError(t, err)

	cfg, err := LoadConfig(dir)
	require.NoError(t, err)
	assert.Equal(t, []string{"claude-code", "gemini-cli"}, cfg.Targets)
	assert.Equal(t, 2, len(cfg.Agents))
}

func TestConfigValidate(t *testing.T) {
	cfg := &Config{
		Targets: []string{"claude-code", "gemini-cli"},
		Agents: map[string]AgentConfig{
			"test": {Targets: []string{"claude-code"}},
		},
	}
	assert.NoError(t, cfg.Validate())

	invalid := &Config{
		Targets: []string{"invalid-platform"},
	}
	assert.Error(t, invalid.Validate())
}

func TestResolveTargets(t *testing.T) {
	cfg := &Config{
		Targets: []string{"claude-code", "gemini-cli", "opencode"},
		Agents: map[string]AgentConfig{
			"all-agent":  {Targets: []string{"all"}},
			"some-agent": {Targets: []string{"claude-code", "gemini-cli"}},
			"disabled":   {Targets: []string{"all"}, Disabled: true},
		},
	}

	assert.Equal(t, []string{"claude-code", "gemini-cli", "opencode"}, cfg.ResolveTargets("all-agent"))
	assert.Equal(t, []string{"claude-code", "gemini-cli"}, cfg.ResolveTargets("some-agent"))
	assert.Nil(t, cfg.ResolveTargets("disabled"))
	assert.Nil(t, cfg.ResolveTargets("nonexistent"))
}
