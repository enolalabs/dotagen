package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/engine"
	"github.com/enolalabs/dotagen/v2/internal/skill"
)

func isValidSkillName(name string) bool {
	clean := strings.TrimPrefix(name, "ds-")
	return validNameRe.MatchString(clean)
}

func (s *Server) handleListSkills(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	skills, err := skill.ParseSkillsDir(filepath.Join(dotgenDir, "skills"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	type skillSummary struct {
		Name        string            `json:"name"`
		Description string            `json:"description,omitempty"`
		Category    string            `json:"category,omitempty"`
		Categories  []string          `json:"categories,omitempty"`
		Frontmatter map[string]string `json:"frontmatter,omitempty"`
		RefCount    int               `json:"refCount"`
	}
	var summaries []skillSummary
	for _, sk := range skills {
		desc := sk.Frontmatter["description"]
		cat := sk.Frontmatter["category"]
		var cats []string
		if cat != "" {
			for _, c := range strings.Split(cat, ",") {
				c = strings.TrimSpace(c)
				if c != "" {
					cats = append(cats, c)
				}
			}
		}
		summaries = append(summaries, skillSummary{
			Name:        sk.Name,
			Description: desc,
			Category:    cat,
			Categories:  cats,
			Frontmatter: sk.Frontmatter,
			RefCount:    len(sk.References),
		})
	}
	writeJSON(w, http.StatusOK, summaries)
}

func (s *Server) handleGetSkill(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	name := r.PathValue("name")
	if !isValidSkillName(name) {
		writeError(w, http.StatusBadRequest, "invalid skill name")
		return
	}
	dirPath := filepath.Join(dotgenDir, "skills", name)
	sk, err := skill.ParseSkillDir(dirPath)
	if err != nil || sk == nil {
		writeError(w, http.StatusNotFound, "skill not found")
		return
	}
	writeJSON(w, http.StatusOK, sk)
}

func (s *Server) handleCreateSkill(w http.ResponseWriter, r *http.Request) {
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
	if !isValidSkillName(body.Name) {
		writeError(w, http.StatusBadRequest, "name must contain only alphanumeric characters, hyphens, and underscores")
		return
	}

	if !strings.HasPrefix(body.Name, "ds-") {
		body.Name = "ds-" + body.Name
	}

	skillDir := filepath.Join(dotgenDir, "skills", body.Name)
	if _, err := os.Stat(skillDir); err == nil {
		writeError(w, http.StatusConflict, "skill already exists")
		return
	}

	content := body.Content
	if content == "" {
		content = skill.ScaffoldSkillContent(body.Name)
	}

	// Build frontmatter
	if body.Description != "" || body.Category != "" {
		content = buildSkillContentWithFrontmatter(content, body.Description, body.Category)
	}

	s.lock()
	defer s.unlock()

	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	skillFile := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillFile, []byte(content), 0o644); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	targets := body.Targets
	if len(targets) == 0 {
		targets = config.ValidTargets
	}
	if err := config.AddSkillToConfig(dotgenDir, body.Name, targets); err != nil {
		os.RemoveAll(skillDir)
		writeError(w, http.StatusInternalServerError, "failed to update config: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"name":    body.Name,
		"targets": targets,
	})
}

func (s *Server) handleUpdateSkill(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	name := r.PathValue("name")
	if !isValidSkillName(name) {
		writeError(w, http.StatusBadRequest, "invalid skill name")
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

	skillFile := filepath.Join(dotgenDir, "skills", name, "SKILL.md")
	if _, err := os.Stat(skillFile); os.IsNotExist(err) {
		writeError(w, http.StatusNotFound, "skill not found")
		return
	}

	content := body.Content
	if body.Description != "" || body.Category != "" {
		content = buildSkillContentWithFrontmatter(content, body.Description, body.Category)
	}

	if strings.TrimSpace(content) == "" {
		writeError(w, http.StatusBadRequest, "skill content cannot be empty")
		return
	}

	s.lock()
	defer s.unlock()

	if err := os.WriteFile(skillFile, []byte(content), 0o644); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(body.Targets) > 0 {
		if err := config.AddSkillToConfig(dotgenDir, name, body.Targets); err != nil {
			log.Printf("failed to update config targets for skill %s: %v", name, err)
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"name": name})
}

func (s *Server) handleDeleteSkill(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	name := r.PathValue("name")
	if !isValidSkillName(name) {
		writeError(w, http.StatusBadRequest, "invalid skill name")
		return
	}
	s.lock()
	defer s.unlock()
	skillDir := filepath.Join(dotgenDir, "skills", name)
	if err := os.RemoveAll(skillDir); err != nil {
		writeError(w, http.StatusNotFound, "skill not found")
		return
	}
	if err := config.RemoveSkillFromConfig(dotgenDir, name); err != nil {
		log.Printf("failed to remove %s from config: %v", name, err)
	}
	writeJSON(w, http.StatusOK, map[string]string{"deleted": name})
}

func (s *Server) handlePreviewSkill(w http.ResponseWriter, r *http.Request) {
	dotgenDir, err := s.dotgenDir()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	skillName := r.PathValue("skill")
	targetName := r.PathValue("target")

	if !isValidSkillName(skillName) {
		writeError(w, http.StatusBadRequest, "invalid skill name")
		return
	}
	if !isValidTarget(targetName) {
		writeError(w, http.StatusBadRequest, "invalid target")
		return
	}

	dirPath := filepath.Join(dotgenDir, "skills", skillName)
	sk, err := skill.ParseSkillDir(dirPath)
	if err != nil || sk == nil {
		writeError(w, http.StatusNotFound, "skill not found")
		return
	}
	adapter, err := s.registry.GetSkillAdapter(targetName)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	rendered, err := adapter.RenderSkill(*sk)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"skill":   skillName,
		"target":  targetName,
		"content": rendered,
	})
}

func (s *Server) handleSyncWithSkills(w http.ResponseWriter, r *http.Request) {
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

	// Render skills
	skills, err := skill.ParseSkillsDir(filepath.Join(dotgenDir, "skills"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	renderer := engine.NewRenderer(s.registry)
	skillResults, err := renderer.RenderAllSkills(skills, cfg, dotgenDir, projectDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"synced":  len(skillResults),
		"results": skillResults,
	})
}

func buildSkillContentWithFrontmatter(content, description, category string) string {
	if description == "" && category == "" {
		return content
	}

	var fm strings.Builder
	fm.WriteString("---\n")

	// If content already has frontmatter, parse and merge
	if strings.HasPrefix(strings.TrimSpace(content), "---") {
		parts := strings.SplitN(strings.TrimSpace(content), "---", 3)
		if len(parts) >= 3 {
			// Preserve existing frontmatter fields
			existing := strings.TrimSpace(parts[1])
			lines := strings.Split(existing, "\n")
			hasDesc := false
			hasCat := false
			for _, line := range lines {
				if strings.HasPrefix(line, "description:") && description != "" {
					fm.WriteString(fmt.Sprintf("description: %s\n", description))
					hasDesc = true
				} else if strings.HasPrefix(line, "category:") && category != "" {
					fm.WriteString(fmt.Sprintf("category: %s\n", category))
					hasCat = true
				} else {
					fm.WriteString(line + "\n")
				}
			}
			if !hasDesc && description != "" {
				fm.WriteString(fmt.Sprintf("description: %s\n", description))
			}
			if !hasCat && category != "" {
				fm.WriteString(fmt.Sprintf("category: %s\n", category))
			}
			fm.WriteString("---\n\n")
			return fm.String() + strings.TrimSpace(parts[2])
		}
	}

	// No existing frontmatter
	if description != "" {
		fm.WriteString(fmt.Sprintf("description: %s\n", description))
	}
	if category != "" {
		fm.WriteString(fmt.Sprintf("category: %s\n", category))
	}
	fm.WriteString("---\n\n")
	return fm.String() + content
}
