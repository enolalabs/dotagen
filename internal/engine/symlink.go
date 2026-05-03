package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/enolalabs/dotagen/v2/internal/config"
)

func CreateSymlink(src, dst string) error {
	if info, err := os.Lstat(dst); err == nil {
		if info.Mode()&os.ModeSymlink == 0 {
			return fmt.Errorf("refusing to overwrite non-symlink file: %s", dst)
		}
		existing, err := os.Readlink(dst)
		if err != nil {
			return fmt.Errorf("failed to read existing symlink: %w", err)
		}
		if existing == src {
			return nil
		}
		if err := os.Remove(dst); err != nil {
			return fmt.Errorf("failed to remove existing symlink: %w", err)
		}
	}

	return os.Symlink(src, dst)
}

func RemoveSymlink(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return os.Remove(path)
	}
	return fmt.Errorf("%s is not a symlink", path)
}

func IsSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}

type SymlinkInfo struct {
	Path     string
	Target   string
	Agent    string
	Platform string
	Broken   bool
}

func ListManagedSymlinks(projectDir string, managedPaths []string) ([]SymlinkInfo, error) {
	var links []SymlinkInfo
	for _, p := range managedPaths {
		fullPath := filepath.Join(projectDir, p)
		isLink, err := IsSymlink(fullPath)
		if err != nil {
			continue
		}
		if !isLink {
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
			Path:   fullPath,
			Target: target,
			Broken: broken,
		})
	}
	return links, nil
}

func RemoveGeneratedContents(dotgenDir string) error {
	generatedDir := filepath.Join(dotgenDir, ".generated")
	entries, err := os.ReadDir(generatedDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(generatedDir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func FindDotagenSymlinks(projectDir string, dotgenDir string) ([]SymlinkInfo, error) {
	platformDirs := map[string]string{
		config.ClaudeCodeRootPath: "claude-code",
		config.CodexRootPath:     "codex",
		config.GeminiCliRootPath: "gemini-cli",
		config.OpenCodeRootPath:  "opencode",
	}

	var links []SymlinkInfo
	for dir, platform := range platformDirs {
		fullDir := filepath.Join(projectDir, dir)
		entries, err := os.ReadDir(fullDir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			if !strings.HasPrefix(name, "da-") {
				continue
			}
			fullPath := filepath.Join(fullDir, entry.Name())
			isLink, err := IsSymlink(fullPath)
			if err != nil || !isLink {
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
			if !strings.HasPrefix(resolvedTarget, dotgenDir) {
				continue
			}
			broken := false
			if _, err := os.Stat(resolvedTarget); err != nil {
				broken = true
			}
			links = append(links, SymlinkInfo{
				Path:     fullPath,
				Target:   target,
				Agent:    name,
				Platform: platform,
				Broken:   broken,
			})
		}
	}
	return links, nil
}

func RemoveStaleSymlinks(projectDir string, dotgenDir string, activeAgentNames []string, syncTargets []string) ([]string, error) {
	links, err := FindDotagenSymlinks(projectDir, dotgenDir)
	if err != nil {
		return nil, err
	}

	activeSet := make(map[string]bool)
	for _, name := range activeAgentNames {
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
		if err := os.Remove(link.Path); err != nil {
			continue
		}
		rel, _ := filepath.Rel(projectDir, link.Path)
		removed = append(removed, rel)
	}
	return removed, nil
}
