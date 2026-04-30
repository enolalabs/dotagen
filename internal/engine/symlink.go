package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateSymlink(src, dst string) error {
	if info, err := os.Lstat(dst); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			existing, _ := os.Readlink(dst)
			if existing == src {
				return nil
			}
		}
		os.Remove(dst)
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

func FindAllSymlinkPaths(projectDir string, dotgenDir string, agentNames []string, cfgTargets []string, registry interface{ Get(string) (interface{ SymlinkPath(string) string }, error) }) []string {
	return nil
}

func CollectSymlinkPaths(projectDir string, agentNames []string, targets []string, adapters map[string]struct {
	SymlinkPath func(string) string
}) []string {
	var paths []string
	for _, name := range agentNames {
		for _, target := range targets {
			a, ok := adapters[target]
			if !ok {
				continue
			}
			paths = append(paths, a.SymlinkPath(name))
		}
	}
	return paths
}

func FindDotagenSymlinks(projectDir string) ([]SymlinkInfo, error) {
	var links []SymlinkInfo
	dotgenDir := filepath.Join(projectDir, ".dotagen")

	platformDirs := map[string]string{
		".claude/agents":  "claude-code",
		".cursor/rules":   "cursor",
		".gemini/agents":  "gemini-cli",
		".opencode/agents": "opencode",
	}

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
			isLink, _ := IsSymlink(fullPath)
			if !isLink {
				continue
			}
			target, _ := os.Readlink(fullPath)
			broken := false
			if !strings.HasPrefix(target, dotgenDir) && !strings.Contains(target, ".dotagen") {
				continue
			}
			if _, err := os.Stat(target); err != nil {
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

func RemoveStaleSymlinks(projectDir string, activeAgentNames []string, syncTargets []string) ([]string, error) {
	links, err := FindDotagenSymlinks(projectDir)
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
		os.Remove(link.Path)
		rel, _ := filepath.Rel(projectDir, link.Path)
		removed = append(removed, rel)
	}
	return removed, nil
}
