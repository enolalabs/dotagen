package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/enolalabs/dotagen/v2/internal/config"
	"github.com/enolalabs/dotagen/v2/internal/skill"
)

type SkillRenderResult struct {
	SkillName     string
	Target        string
	GeneratedPath string
	SymlinkPath   string
}

// RenderAllSkills renders all skills to their platform-specific output directories and creates symlinks.
func (r *Renderer) RenderAllSkills(skills []skill.Skill, cfg *config.Config, dotgenDir string, projectDir string) ([]SkillRenderResult, error) {
	var results []SkillRenderResult
	generatedDir := filepath.Join(dotgenDir, ".generated")

	for _, sk := range skills {
		targets := cfg.ResolveSkillTargets(sk.Name)
		for _, target := range targets {
			sa, err := r.registry.GetSkillAdapter(target)
			if err != nil {
				continue // Platform doesn't support skills — skip silently
			}

			rendered, err := sa.RenderSkill(sk)
			if err != nil {
				return nil, fmt.Errorf("failed to render skill %q for %s: %w", sk.Name, target, err)
			}

			// Create output directory: .generated/<platform>/skills/<name>/
			outDir := filepath.Join(generatedDir, sa.SkillOutputDir(sk.Name))
			if err := os.MkdirAll(outDir, 0o755); err != nil {
				return nil, fmt.Errorf("failed to create skill output directory: %w", err)
			}

			// Write SKILL.md
			skillMDPath := filepath.Join(outDir, "SKILL.md")
			if err := os.WriteFile(skillMDPath, []byte(rendered), 0o644); err != nil {
				return nil, fmt.Errorf("failed to write skill SKILL.md: %w", err)
			}

			// Copy references
			for _, ref := range sk.References {
				refPath := filepath.Join(outDir, ref.Name)
				if err := os.MkdirAll(filepath.Dir(refPath), 0o755); err != nil {
					return nil, fmt.Errorf("failed to create reference directory: %w", err)
				}
				if err := os.WriteFile(refPath, []byte(ref.Content), 0o644); err != nil {
					return nil, fmt.Errorf("failed to write reference %s: %w", ref.Name, err)
				}
			}

			// Create directory symlink
			absGenerated, err := filepath.Abs(outDir)
			if err != nil {
				return nil, err
			}

			symlinkDir := filepath.Join(projectDir, sa.SkillSymlinkDir(sk.Name))
			if err := sa.EnsureSkillDirectories(projectDir); err != nil {
				return nil, fmt.Errorf("failed to ensure skill directories for %s: %w", target, err)
			}

			if err := CreateSymlink(absGenerated, symlinkDir); err != nil {
				return nil, fmt.Errorf("failed to create skill symlink %s: %w", symlinkDir, err)
			}

			results = append(results, SkillRenderResult{
				SkillName:     sk.Name,
				Target:        target,
				GeneratedPath: outDir,
				SymlinkPath:   symlinkDir,
			})
		}
	}

	return results, nil
}

// FindDotagenSkillSymlinks finds all dotagen-managed skill symlinks in platform directories.
func FindDotagenSkillSymlinks(projectDir string, dotgenDir string) ([]SymlinkInfo, error) {
	type skillDirEntry struct {
		dir      string
		platform string
	}
	entries := []skillDirEntry{
		{config.AntigravitySkillPath, "antigravity"},
		{config.ClaudeCodeSkillPath, "claude-code"},
		{config.CodexSkillPath, "codex"},
		{config.GeminiCliSkillPath, "gemini-cli"},
		{config.OpenCodeSkillPath, "opencode"},
	}

	// Deduplicate directories — if multiple platforms share a path, use the first one
	seen := make(map[string]bool)
	var skillDirs []skillDirEntry
	for _, e := range entries {
		if !seen[e.dir] {
			seen[e.dir] = true
			skillDirs = append(skillDirs, e)
		}
	}

	var links []SymlinkInfo
	for _, entry := range skillDirs {
		dir, plat := entry.dir, entry.platform
		fullDir := filepath.Join(projectDir, dir)
		entries, err := os.ReadDir(fullDir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			name := entry.Name()
			if !entry.IsDir() {
				// Could be a symlink to a directory
				fullPath := filepath.Join(fullDir, name)
				isLink, err := IsSymlink(fullPath)
				if err != nil || !isLink {
					continue
				}
				if !hasDotagenSkillPrefix(name) {
					continue
				}
				target, err := os.Readlink(fullPath)
				if err != nil {
					continue
				}
				resolvedTarget := target
				if !filepath.IsAbs(resolvedTarget) {
					resolvedTarget = filepath.Join(filepath.Dir(fullPath), resolvedTarget)
				}
				broken := false
				if _, err := os.Stat(resolvedTarget); err != nil {
					broken = true
				}
				links = append(links, SymlinkInfo{
					Path:     fullPath,
					Target:   target,
					Agent:    name,
					Platform: plat,
					Broken:   broken,
				})
			} else {
				// Check if it's a symlinked directory
				fullPath := filepath.Join(fullDir, name)
				isLink, err := IsSymlink(fullPath)
				if err != nil || !isLink {
					continue
				}
				if !hasDotagenSkillPrefix(name) {
					continue
				}
				target, err := os.Readlink(fullPath)
				if err != nil {
					continue
				}
				broken := false
				if _, err := os.Stat(target); err != nil {
					broken = true
				}
				links = append(links, SymlinkInfo{
					Path:     fullPath,
					Target:   target,
					Agent:    name,
					Platform: plat,
					Broken:   broken,
				})
			}
		}
	}
	return links, nil
}

func hasDotagenSkillPrefix(name string) bool {
	return len(name) > 3 && name[:3] == "ds-"
}

// RemoveStaleSkillSymlinks removes skill symlinks that are no longer active.
func RemoveStaleSkillSymlinks(projectDir string, dotgenDir string, activeSkillNames []string, syncTargets []string) ([]string, error) {
	links, err := FindDotagenSkillSymlinks(projectDir, dotgenDir)
	if err != nil {
		return nil, err
	}

	activeSet := make(map[string]bool)
	for _, name := range activeSkillNames {
		for _, target := range syncTargets {
			activeSet[name+"|"+target] = true
		}
	}

	var removed []string
	for _, link := range links {
		key := link.Agent + "|" + link.Platform
		if activeSet[key] {
			continue
		}
		if err := os.RemoveAll(link.Path); err != nil {
			continue
		}
		rel, _ := filepath.Rel(projectDir, link.Path)
		removed = append(removed, rel)
	}
	return removed, nil
}
