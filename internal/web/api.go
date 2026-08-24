package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/enolalabs/dotagen/v2/internal/agent"
	"github.com/enolalabs/dotagen/v2/internal/builtin"
	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/engine"
	"github.com/enolalabs/dotagen/v2/internal/skill"
	"github.com/enolalabs/dotagen/v2/skillsrc"
	"gopkg.in/yaml.v3"
)

var validNameRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

func (s *Server) dotgenDir() (string, error) {
	return config.FindDotgenDir()
}

func (s *Server) projectDir() (string, error) {
	return os.UserHomeDir()
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func isValidAgentName(name string) bool {
	clean := strings.TrimPrefix(name, "da-")
	return validNameRe.MatchString(clean)
}

func isValidTarget(name string) bool {
	for _, t := range config.ValidTargets {
		if t == name {
			return true
		}
	}
	return false
}

func buildContentWithFrontmatter(content, description, category string) string {
	fm := make(map[string]string)

	if strings.HasPrefix(strings.TrimSpace(content), "---") {
		var existing map[string]string
		body, err := frontmatter.Parse(strings.NewReader(content), &existing)
		if err == nil {
			fm = existing
			content = strings.TrimSpace(string(body))
		}
	}

	if description != "" {
		fm["description"] = description
	}
	if category != "" {
		fm["category"] = category
	}

	if len(fm) == 0 {
		return content
	}

	sortedKeys := make([]string, 0, len(fm))
	for k := range fm {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	fmData := make(map[string]string)
	for _, k := range sortedKeys {
		fmData[k] = fm[k]
	}

	fmBytes, err := yaml.Marshal(fmData)
	if err != nil {
		return content
	}

	return fmt.Sprintf("---\n%s---\n\n%s", string(fmBytes), content)
}

func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	s.lock()
	defer s.unlock()
	cfg, err := config.LoadConfig(dotgenDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

func (s *Server) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	var cfg config.Config
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := cfg.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.lock()
	defer s.unlock()
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to marshal config")
		return
	}
	configPath := filepath.Join(dotgenDir, "config.yaml")
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

func (s *Server) handleListAgents(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	agents, err := agent.ParseAgentsDir(filepath.Join(dotgenDir, "agents"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	type agentSummary struct {
		Name        string            `json:"name"`
		Description string            `json:"description,omitempty"`
		Category    string            `json:"category,omitempty"`
		Categories  []string          `json:"categories,omitempty"`
		Frontmatter map[string]string `json:"frontmatter,omitempty"`
	}
	var summaries []agentSummary
	for _, a := range agents {
		desc := a.Frontmatter["description"]
		cat := a.Frontmatter["category"]
		var cats []string
		if cat != "" {
			for _, c := range strings.Split(cat, ",") {
				c = strings.TrimSpace(c)
				if c != "" {
					cats = append(cats, c)
				}
			}
		}
		summaries = append(summaries, agentSummary{
			Name:        a.Name,
			Description: desc,
			Category:    cat,
			Categories:  cats,
			Frontmatter: a.Frontmatter,
		})
	}
	writeJSON(w, http.StatusOK, summaries)
}

func (s *Server) handleGetAgent(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	name := r.PathValue("name")
	if !isValidAgentName(name) {
		writeError(w, http.StatusBadRequest, "invalid agent name")
		return
	}
	filePath := filepath.Join(dotgenDir, "agents", name+".md")
	a, err := agent.ParseAgentFile(filePath)
	if err != nil {
		writeError(w, http.StatusNotFound, "agent not found")
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (s *Server) handleCreateAgent(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	var body struct {
		Name        string   `json:"name"`
		Content     string   `json:"content"`
		Description string   `json:"description"`
		Category    string   `json:"category"`
		Targets     []string `json:"targets"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if !isValidAgentName(body.Name) {
		writeError(w, http.StatusBadRequest, "name must contain only alphanumeric characters, hyphens, and underscores")
		return
	}

	if !strings.HasPrefix(body.Name, "da-") {
		body.Name = "da-" + body.Name
	}

	agentsDir := filepath.Join(dotgenDir, "agents")
	filePath := filepath.Join(agentsDir, body.Name+".md")
	if _, err := os.Stat(filePath); err == nil {
		writeError(w, http.StatusConflict, "agent already exists")
		return
	}

	content := body.Content
	if content == "" {
		content = agent.ScaffoldContent(body.Name)
	}
	content = buildContentWithFrontmatter(content, body.Description, body.Category)

	s.lock()
	defer s.unlock()

	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	targets := body.Targets
	if len(targets) == 0 {
		targets = config.ValidTargets
	}
	if err := config.AddAgentToConfig(dotgenDir, body.Name, targets); err != nil {
		os.Remove(filePath)
		writeError(w, http.StatusInternalServerError, "failed to update config: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"name":    body.Name,
		"targets": targets,
	})
}

func (s *Server) handleUpdateAgent(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	name := r.PathValue("name")
	if !isValidAgentName(name) {
		writeError(w, http.StatusBadRequest, "invalid agent name")
		return
	}
	var body struct {
		Content     string   `json:"content"`
		Description string   `json:"description"`
		Category    string   `json:"category"`
		Targets     []string `json:"targets"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	filePath := filepath.Join(dotgenDir, "agents", name+".md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		writeError(w, http.StatusNotFound, "agent not found")
		return
	}

	content := body.Content
	if body.Description != "" || body.Category != "" {
		content = buildContentWithFrontmatter(content, body.Description, body.Category)
	}

	if strings.TrimSpace(content) == "" {
		writeError(w, http.StatusBadRequest, "agent content cannot be empty")
		return
	}

	s.lock()
	defer s.unlock()

	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(body.Targets) > 0 {
		if err := config.AddAgentToConfig(dotgenDir, name, body.Targets); err != nil {
			log.Printf("failed to update config targets for %s: %v", name, err)
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"name": name})
}

func (s *Server) handleDeleteAgent(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	name := r.PathValue("name")
	if !isValidAgentName(name) {
		writeError(w, http.StatusBadRequest, "invalid agent name")
		return
	}
	s.lock()
	defer s.unlock()
	filePath := filepath.Join(dotgenDir, "agents", name+".md")
	if err := os.Remove(filePath); err != nil {
		writeError(w, http.StatusNotFound, "agent not found")
		return
	}
	if err := config.RemoveAgentFromConfig(dotgenDir, name); err != nil {
		log.Printf("failed to remove %s from config: %v", name, err)
	}
	writeJSON(w, http.StatusOK, map[string]string{"deleted": name})
}

func (s *Server) handleListTargets(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string][]string{
		"targets": config.ValidTargets,
	})
}

func (s *Server) handlePreview(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	agentName := r.PathValue("agent")
	targetName := r.PathValue("target")

	if !isValidAgentName(agentName) {
		writeError(w, http.StatusBadRequest, "invalid agent name")
		return
	}
	if !isValidTarget(targetName) {
		writeError(w, http.StatusBadRequest, "invalid target")
		return
	}

	a, err := agent.ParseAgentFile(filepath.Join(dotgenDir, "agents", agentName+".md"))
	if err != nil {
		writeError(w, http.StatusNotFound, "agent not found")
		return
	}
	adapter, err := s.registry.Get(targetName)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	rendered, err := adapter.Render(*a)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"agent":   agentName,
		"target":  targetName,
		"content": rendered,
	})
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	projectDir, err := s.projectDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.lock()
	defer s.unlock()

	cfg, err := config.LoadConfig(dotgenDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	detected := config.DetectPlatforms(projectDir)
	if len(detected) > 0 {
		cfg.Targets = detected
	}

	agents, err := agent.ParseAgentsDir(filepath.Join(dotgenDir, "agents"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	renderer := engine.NewRenderer(s.registry)
	results, err := renderer.RenderAll(agents, cfg, dotgenDir, projectDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Also sync skills
	skills, err := skill.ParseSkillsDir(filepath.Join(dotgenDir, "skills"))
	if err != nil {
		log.Printf("failed to parse skills: %v", err)
		skills = nil
	}
	var skillResults []engine.SkillRenderResult
	if len(skills) > 0 {
		skillResults, err = renderer.RenderAllSkills(skills, cfg, dotgenDir, projectDir)
		if err != nil {
			log.Printf("failed to render skills: %v", err)
		}
	}

	// Sync global workflows for Antigravity
	var workflowsSynced int
	for _, t := range cfg.Targets {
		if t == "antigravity" {
			workflowsSynced, err = engine.SyncGlobalWorkflows(agents, projectDir)
			if err != nil {
				log.Printf("failed to sync global workflows: %v", err)
			}
			break
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"synced":          len(results) + len(skillResults),
		"agentsSynced":    len(results),
		"skillsSynced":    len(skillResults),
		"workflowsSynced": workflowsSynced,
		"results":         results,
		"skillResults":    skillResults,
	})
}

func (s *Server) handleSyncTarget(w http.ResponseWriter, r *http.Request) {
	targetName := r.PathValue("target")
	if !isValidTarget(targetName) {
		writeError(w, http.StatusBadRequest, "invalid target")
		return
	}
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	projectDir, err := s.projectDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.lock()
	defer s.unlock()

	cfg, err := config.LoadConfig(dotgenDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	filteredCfg := &config.Config{
		Targets: []string{targetName},
		Agents:  make(map[string]config.AgentConfig),
		Skills:  make(map[string]config.SkillConfig),
	}
	for name, ac := range cfg.Agents {
		resolved := cfg.ResolveTargets(name)
		for _, t := range resolved {
			if t == targetName {
				filteredCfg.Agents[name] = config.AgentConfig{Targets: config.StringOrSlice{targetName}, Disabled: ac.Disabled}
				break
			}
		}
	}
	for name, sc := range cfg.Skills {
		resolved := cfg.ResolveSkillTargets(name)
		for _, t := range resolved {
			if t == targetName {
				filteredCfg.Skills[name] = config.SkillConfig{Targets: config.StringOrSlice{targetName}, Disabled: sc.Disabled}
				break
			}
		}
	}

	agents, err := agent.ParseAgentsDir(filepath.Join(dotgenDir, "agents"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	renderer := engine.NewRenderer(s.registry)
	results, err := renderer.RenderAll(agents, filteredCfg, dotgenDir, projectDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Also sync skills for this target
	skills, err := skill.ParseSkillsDir(filepath.Join(dotgenDir, "skills"))
	if err != nil {
		log.Printf("failed to parse skills: %v", err)
		skills = nil
	}
	var skillResults []engine.SkillRenderResult
	if len(skills) > 0 {
		skillResults, err = renderer.RenderAllSkills(skills, filteredCfg, dotgenDir, projectDir)
		if err != nil {
			log.Printf("failed to render skills: %v", err)
		}
	}

	// Sync global workflows for Antigravity
	var workflowsSynced int
	if targetName == "antigravity" {
		workflowsSynced, err = engine.SyncGlobalWorkflows(agents, projectDir)
		if err != nil {
			log.Printf("failed to sync global workflows: %v", err)
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"synced":          len(results) + len(skillResults),
		"agentsSynced":    len(results),
		"skillsSynced":    len(skillResults),
		"workflowsSynced": workflowsSynced,
		"target":          targetName,
		"results":         results,
		"skillResults":    skillResults,
	})
}

func (s *Server) handleClean(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	projectDir, err := s.projectDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.lock()
	defer s.unlock()

	cfg, err := config.LoadConfig(dotgenDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	links, err := engine.FindDotagenSymlinks(projectDir, dotgenDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	removed := 0
	for _, link := range links {
		if err := os.Remove(link.Path); err == nil {
			removed++
		}
	}

	skillLinks, err := engine.FindDotagenSkillSymlinks(projectDir, dotgenDir)
	if err != nil {
		log.Printf("failed to find skill symlinks: %v", err)
	} else {
		for _, link := range skillLinks {
			if err := os.RemoveAll(link.Path); err == nil {
				removed++
			}
		}
	}

	if err := engine.RemoveGeneratedContents(dotgenDir); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("removed %d symlinks but failed to clean generated: %v", removed, err))
		return
	}

	// Clear all targets in config
	for name := range cfg.Agents {
		cfg.Agents[name] = config.AgentConfig{Targets: config.StringOrSlice{}, Disabled: false}
	}
	for name := range cfg.Skills {
		cfg.Skills[name] = config.SkillConfig{Targets: config.StringOrSlice{}, Disabled: false}
	}
	if err := config.SaveConfig(dotgenDir, cfg); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"removed": removed,
		"config":  cfg,
	})
}

func (s *Server) handleCleanBroken(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	projectDir, err := s.projectDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.lock()
	defer s.unlock()

	cfg, err := config.LoadConfig(dotgenDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fixed := 0
	removed := 0

	// Fix/remove broken agent symlinks
	links, err := engine.FindDotagenSymlinks(projectDir, dotgenDir)
	if err == nil {
		for _, link := range links {
			if !link.Broken {
				continue
			}
			adapter, err := s.registry.Get(link.Platform)
			if err != nil {
				os.Remove(link.Path)
				removeFromAgentConfig(cfg, link.Agent, link.Platform)
				removed++
				continue
			}
			ag, err := agent.ParseAgentFile(filepath.Join(dotgenDir, "agents", link.Agent+".md"))
			if err != nil {
				os.Remove(link.Path)
				removeFromAgentConfig(cfg, link.Agent, link.Platform)
				removed++
				continue
			}
			rendered, err := adapter.Render(*ag)
			if err != nil {
				os.Remove(link.Path)
				removeFromAgentConfig(cfg, link.Agent, link.Platform)
				removed++
				continue
			}
			outPath := filepath.Join(dotgenDir, ".generated", adapter.OutputPath(link.Agent))
			os.MkdirAll(filepath.Dir(outPath), 0o755)
			if err := os.WriteFile(outPath, []byte(rendered), 0o644); err != nil {
				os.Remove(link.Path)
				removeFromAgentConfig(cfg, link.Agent, link.Platform)
				removed++
				continue
			}
			fixed++
		}
	}

	// Fix/remove broken skill symlinks
	skillLinks, err := engine.FindDotagenSkillSymlinks(projectDir, dotgenDir)
	if err == nil {
		for _, link := range skillLinks {
			if !link.Broken {
				continue
			}
			sa, err := s.registry.GetSkillAdapter(link.Platform)
			if err != nil {
				os.RemoveAll(link.Path)
				removeFromSkillConfig(cfg, link.Agent, link.Platform)
				removed++
				continue
			}
			sk, err := skill.ParseSkillDir(filepath.Join(dotgenDir, "skills", link.Agent))
			if err != nil || sk == nil {
				os.RemoveAll(link.Path)
				removeFromSkillConfig(cfg, link.Agent, link.Platform)
				removed++
				continue
			}
			rendered, err := sa.RenderSkill(*sk)
			if err != nil {
				os.RemoveAll(link.Path)
				removeFromSkillConfig(cfg, link.Agent, link.Platform)
				removed++
				continue
			}
			outDir := filepath.Join(dotgenDir, ".generated", sa.SkillOutputDir(link.Agent))
			os.MkdirAll(outDir, 0o755)
			os.WriteFile(filepath.Join(outDir, "SKILL.md"), []byte(rendered), 0o644)
			for _, ref := range sk.References {
				refPath := filepath.Join(outDir, ref.Name)
				os.MkdirAll(filepath.Dir(refPath), 0o755)
				os.WriteFile(refPath, []byte(ref.Content), 0o644)
			}
			fixed++
		}
	}

	config.SaveConfig(dotgenDir, cfg)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"fixed":   fixed,
		"removed": removed,
		"config":  cfg,
	})
}

func removeFromAgentConfig(cfg *config.Config, name, platform string) {
	entry, ok := cfg.Agents[name]
	if !ok {
		return
	}
	entry.Targets = removeFromStringSlice(entry.Targets, platform)
	cfg.Agents[name] = entry
}

func removeFromSkillConfig(cfg *config.Config, name, platform string) {
	entry, ok := cfg.Skills[name]
	if !ok {
		return
	}
	entry.Targets = removeFromStringSlice(entry.Targets, platform)
	cfg.Skills[name] = entry
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	projectDir, err := s.projectDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	links, err := engine.FindDotagenSymlinks(projectDir, dotgenDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	skillLinks, err := engine.FindDotagenSkillSymlinks(projectDir, dotgenDir)
	if err != nil {
		log.Printf("failed to find skill symlinks: %v", err)
	}

	type linkStatus struct {
		Path     string `json:"path"`
		Agent    string `json:"agent"`
		Platform string `json:"platform"`
		Type     string `json:"type"`
		Broken   bool   `json:"broken"`
	}
	var statuses []linkStatus
	for _, l := range links {
		rel, _ := filepath.Rel(projectDir, l.Path)
		statuses = append(statuses, linkStatus{
			Path:     rel,
			Agent:    l.Agent,
			Platform: l.Platform,
			Type:     "agent",
			Broken:   l.Broken,
		})
	}
	for _, l := range skillLinks {
		rel, _ := filepath.Rel(projectDir, l.Path)
		statuses = append(statuses, linkStatus{
			Path:     rel,
			Agent:    l.Agent,
			Platform: l.Platform,
			Type:     "skill",
			Broken:   l.Broken,
		})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"symlinks": statuses,
	})
}

func (s *Server) handleToggle(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	projectDir, err := s.projectDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req struct {
		Type  string `json:"type"`
		Items []struct {
			Name     string `json:"name"`
			Platform string `json:"platform"`
			Enable   bool   `json:"enable"`
		} `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Type != "agent" && req.Type != "skill" {
		writeError(w, http.StatusBadRequest, "type must be 'agent' or 'skill'")
		return
	}

	s.lock()
	defer s.unlock()

	cfg, err := config.LoadConfig(dotgenDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 1. Validate items and update config
	for _, item := range req.Items {
		if !isValidTarget(item.Platform) {
			continue
		}
		if req.Type == "agent" {
			// Validate source file exists before touching config
			if item.Enable {
				if _, err := os.Stat(filepath.Join(dotgenDir, "agents", item.Name+".md")); err != nil {
					continue
				}
			}
			entry := cfg.Agents[item.Name]
			entry.Disabled = false
			if item.Enable {
				if !containsString(entry.Targets, item.Platform) {
					entry.Targets = append(entry.Targets, item.Platform)
				}
			} else {
				entry.Targets = removeFromStringSlice(entry.Targets, item.Platform)
			}
			if cfg.Agents == nil {
				cfg.Agents = make(map[string]config.AgentConfig)
			}
			cfg.Agents[item.Name] = entry
		} else {
			// Validate source dir exists before touching config
			if item.Enable {
				if _, err := os.Stat(filepath.Join(dotgenDir, "skills", item.Name)); err != nil {
					continue
				}
			}
			entry := cfg.Skills[item.Name]
			entry.Disabled = false
			if item.Enable {
				if !containsString(entry.Targets, item.Platform) {
					entry.Targets = append(entry.Targets, item.Platform)
				}
			} else {
				entry.Targets = removeFromStringSlice(entry.Targets, item.Platform)
			}
			if cfg.Skills == nil {
				cfg.Skills = make(map[string]config.SkillConfig)
			}
			cfg.Skills[item.Name] = entry
		}
	}

	if err := config.SaveConfig(dotgenDir, cfg); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 2. Create/remove symlinks
	synced := 0
	removed := 0
	for _, item := range req.Items {
		if !isValidTarget(item.Platform) {
			continue
		}
		if req.Type == "agent" {
			adapter, err := s.registry.Get(item.Platform)
			if err != nil {
				continue
			}
			if item.Enable {
				ag, err := agent.ParseAgentFile(filepath.Join(dotgenDir, "agents", item.Name+".md"))
				if err != nil {
					continue
				}
				rendered, err := adapter.Render(*ag)
				if err != nil {
					continue
				}
				outPath := filepath.Join(dotgenDir, ".generated", adapter.OutputPath(item.Name))
				os.MkdirAll(filepath.Dir(outPath), 0o755)
				os.WriteFile(outPath, []byte(rendered), 0o644)
				absGenerated, _ := filepath.Abs(outPath)
				adapter.EnsureDirectories(projectDir)
				symlinkPath := config.ResolvePath(projectDir, adapter.SymlinkPath(item.Name))
				engine.CreateSymlink(absGenerated, symlinkPath)
				synced++
			} else {
				symlinkPath := config.ResolvePath(projectDir, adapter.SymlinkPath(item.Name))
				if err := engine.RemoveSymlink(symlinkPath); err == nil {
					removed++
				}
			}
		} else {
			sa, err := s.registry.GetSkillAdapter(item.Platform)
			if err != nil {
				continue
			}
			if item.Enable {
				sk, err := skill.ParseSkillDir(filepath.Join(dotgenDir, "skills", item.Name))
				if err != nil || sk == nil {
					continue
				}
				rendered, err := sa.RenderSkill(*sk)
				if err != nil {
					continue
				}
				outDir := filepath.Join(dotgenDir, ".generated", sa.SkillOutputDir(item.Name))
				os.MkdirAll(outDir, 0o755)
				os.WriteFile(filepath.Join(outDir, "SKILL.md"), []byte(rendered), 0o644)
				for _, ref := range sk.References {
					refPath := filepath.Join(outDir, ref.Name)
					os.MkdirAll(filepath.Dir(refPath), 0o755)
					os.WriteFile(refPath, []byte(ref.Content), 0o644)
				}
				absGenerated, _ := filepath.Abs(outDir)
				sa.EnsureSkillDirectories(projectDir)
				symlinkDir := config.ResolvePath(projectDir, sa.SkillSymlinkDir(item.Name))
				engine.CreateSymlink(absGenerated, symlinkDir)
				synced++
			} else {
				symlinkDir := config.ResolvePath(projectDir, sa.SkillSymlinkDir(item.Name))
				os.RemoveAll(symlinkDir)
				removed++
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"synced":  synced,
		"removed": removed,
		"config":  cfg,
	})
}

func containsString(slice config.StringOrSlice, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func removeFromStringSlice(slice config.StringOrSlice, s string) config.StringOrSlice {
	var result config.StringOrSlice
	for _, v := range slice {
		if v != s {
			result = append(result, v)
		}
	}
	return result
}

func (s *Server) handleInit(w http.ResponseWriter, r *http.Request) {
	home, err := os.UserHomeDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get home directory: "+err.Error())
		return
	}
	dotgenDir, err := config.FindDotgenDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.lock()
	defer s.unlock()

	if _, err := os.Stat(dotgenDir); err == nil {
		links, _ := engine.FindDotagenSymlinks(home, dotgenDir)
		for _, link := range links {
			os.Remove(link.Path)
		}
		engine.RemoveGeneratedContents(dotgenDir)
		os.RemoveAll(filepath.Join(dotgenDir, "agents"))
		os.RemoveAll(filepath.Join(dotgenDir, "skills"))
		os.Remove(filepath.Join(dotgenDir, "config.yaml"))
	}

	for _, dir := range []string{
		filepath.Join(dotgenDir, "agents"),
		filepath.Join(dotgenDir, "skills"),
		filepath.Join(dotgenDir, ".generated"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create directory: "+err.Error())
			return
		}
	}

	agentNames := builtin.ListAgents()
	for _, name := range agentNames {
		data, err := builtin.ReadAgent(name)
		if err != nil {
			continue
		}
		os.WriteFile(filepath.Join(dotgenDir, "agents", name+".md"), data, 0o644)
	}

	skillNames := skillsrc.ListSkills()
	for _, name := range skillNames {
		files := skillsrc.ListSkillFiles(name)
		for _, file := range files {
			data, err := skillsrc.ReadSkillFile(name + "/" + file)
			if err != nil {
				continue
			}
			outPath := filepath.Join(dotgenDir, "skills", name, file)
			os.MkdirAll(filepath.Dir(outPath), 0o755)
			os.WriteFile(outPath, data, 0o644)
		}
	}

	detected := config.DetectPlatforms(home)
	allTargets := config.ValidTargets
	var cfgTargets []string
	if len(detected) > 0 {
		cfgTargets = detected
	} else {
		cfgTargets = allTargets
	}

	var sb strings.Builder
	sb.WriteString("# dotagen configuration\n")
	sb.WriteString("# Docs: https://github.com/enolalabs/dotagen\n")
	sb.WriteString("#\n")
	sb.WriteString("# All agents and skills are listed with empty targets (disabled).\n")
	sb.WriteString("# Set targets to enable them. Examples:\n")
	sb.WriteString("#   targets: all                    — all platforms\n")
	sb.WriteString("#   targets: [claude-code, gemini-cli]  — specific platforms\n")
	sb.WriteString("#\n")
	sb.WriteString("# Platforms are auto-detected from $HOME.\n\n")
	sb.WriteString("targets:\n")
	for _, t := range cfgTargets {
		sb.WriteString(fmt.Sprintf("  - %s\n", t))
	}
	sb.WriteString("\nagents:\n")
	for _, name := range agentNames {
		sb.WriteString(fmt.Sprintf("  %s:\n    targets: []\n", name))
	}
	sb.WriteString("\nskills:\n")
	for _, name := range skillNames {
		sb.WriteString(fmt.Sprintf("  %s:\n    targets: []\n", name))
	}
	os.WriteFile(filepath.Join(dotgenDir, "config.yaml"), []byte(sb.String()), 0o644)
	os.WriteFile(filepath.Join(dotgenDir, ".gitignore"), []byte(".generated/\n"), 0o644)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"agents":    len(agentNames),
		"skills":    len(skillNames),
		"platforms": detected,
		"reinit":    true,
	})
}
