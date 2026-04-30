package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var ValidTargets = []string{"claude-code", "cursor", "gemini-cli", "opencode"}

type StringOrSlice []string

func (s *StringOrSlice) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		*s = []string{value.Value}
		return nil
	}
	if value.Kind == yaml.SequenceNode {
		var items []string
		if err := value.Decode(&items); err != nil {
			return err
		}
		*s = items
		return nil
	}
	return fmt.Errorf("expected string or sequence")
}

type AgentConfig struct {
	Targets  StringOrSlice `yaml:"targets" json:"targets"`
	Disabled bool          `yaml:"disabled" json:"disabled"`
}

type Config struct {
	Targets []string              `yaml:"targets" json:"targets"`
	Agents  map[string]AgentConfig `yaml:"agents" json:"agents"`
}

func LoadConfig(dotgenDir string) (*Config, error) {
	configPath := filepath.Join(dotgenDir, "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	validMap := make(map[string]bool)
	for _, t := range ValidTargets {
		validMap[t] = true
	}

	for _, t := range c.Targets {
		if !validMap[t] {
			return fmt.Errorf("invalid target %q in config, valid targets: %v", t, ValidTargets)
		}
	}

	for name, agent := range c.Agents {
		if len(agent.Targets) == 1 && agent.Targets[0] == "all" {
			continue
		}
		for _, t := range agent.Targets {
			if !validMap[t] {
				return fmt.Errorf("invalid target %q for agent %q, valid targets: %v", t, name, ValidTargets)
			}
		}
	}

	return nil
}

func (c *Config) ResolveTargets(agentName string) []string {
	agent, ok := c.Agents[agentName]
	if !ok {
		return nil
	}
	if agent.Disabled {
		return nil
	}
	if len(agent.Targets) == 1 && agent.Targets[0] == "all" {
		return c.Targets
	}
	return agent.Targets
}

func (c *Config) AddAgent(name string, targets []string) {
	if len(targets) == 1 && targets[0] == "all" {
		c.Agents[name] = AgentConfig{Targets: StringOrSlice{"all"}}
	} else {
		c.Agents[name] = AgentConfig{Targets: StringOrSlice(targets)}
	}
}

func SaveConfig(dotgenDir string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	configPath := filepath.Join(dotgenDir, "config.yaml")
	return os.WriteFile(configPath, data, 0o644)
}

func AddAgentToConfig(dotgenDir string, name string, targets []string) error {
	cfg, err := LoadConfig(dotgenDir)
	if err != nil {
		return err
	}
	cfg.AddAgent(name, targets)
	if err := cfg.Validate(); err != nil {
		return err
	}
	return SaveConfig(dotgenDir, cfg)
}

func FindDotgenDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := cwd
	for {
		candidate := filepath.Join(dir, ".dotagen")
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf(".dotagen directory not found; run 'dotagen init' first")
		}
		dir = parent
	}
}
