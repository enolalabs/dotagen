package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func isValidAgentName(name string) bool {
	clean := strings.TrimPrefix(name, "da-")
	return validNameRe.MatchString(clean)
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

	var buf bytes.Buffer
	buf.WriteString("---\n")
	for k, v := range fm {
		buf.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	buf.WriteString("---\n\n")
	buf.WriteString(content)
	return buf.String()
}

func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
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
		Frontmatter map[string]string `json:"frontmatter,omitempty"`
	}
	var summaries []agentSummary
	for _, a := range agents {
		desc := a.Frontmatter["description"]
		cat := a.Frontmatter["category"]
		summaries = append(summaries, agentSummary{
			Name:        a.Name,
			Description: desc,
			Category:    cat,
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
		content = scaffoldAgentContent(body.Name)
	}
	content = buildContentWithFrontmatter(content, body.Description, body.Category)

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

func scaffoldAgentContent(name string) string {
	title := strings.TrimPrefix(name, "da-")
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.Title(title)
	return fmt.Sprintf("# %s\n\n## Role\n\nTODO: Describe the agent's role.\n\n## Guidelines\n\n- Guideline 1\n- Guideline 2\n", title)
}

func (s *Server) handleUpdateAgent(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	name := r.PathValue("name")
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

	links, err := engine.FindDotagenSymlinks(projectDir)
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
	engine.RemoveGeneratedContents(dotgenDir)
	writeJSON(w, http.StatusOK, map[string]int{"removed": removed})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	projectDir, err := config.GetProjectDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	links, err := engine.FindDotagenSymlinks(projectDir)
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
