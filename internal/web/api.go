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
	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/engine"
	"gopkg.in/yaml.v3"
)

var validNameRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

func (s *Server) dotgenDir() (string, error) {
	return config.FindDotgenDir()
}

func (s *Server) projectDir() (string, error) {
	return config.GetProjectDir()
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
	projectDir, err := config.GetProjectDir()
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
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"synced":  len(results),
		"results": results,
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
	projectDir, err := config.GetProjectDir()
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
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"synced":  len(results),
		"target":  targetName,
		"results": results,
	})
}

func (s *Server) handleClean(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	projectDir, err := config.GetProjectDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.lock()
	defer s.unlock()

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
	if err := engine.RemoveGeneratedContents(dotgenDir); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("removed %d symlinks but failed to clean generated: %v", removed, err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"removed": removed})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	projectDir, err := config.GetProjectDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	links, err := engine.FindDotagenSymlinks(projectDir, dotgenDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type linkStatus struct {
		Path     string `json:"path"`
		Agent    string `json:"agent"`
		Platform string `json:"platform"`
		Broken   bool   `json:"broken"`
	}
	var statuses []linkStatus
	for _, l := range links {
		rel, _ := filepath.Rel(projectDir, l.Path)
		statuses = append(statuses, linkStatus{
			Path:     rel,
			Agent:    l.Agent,
			Platform: l.Platform,
			Broken:   l.Broken,
		})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"symlinks": statuses,
	})
}
