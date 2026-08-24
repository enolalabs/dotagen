package engine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/enolalabs/dotagen/v2/internal/config"
)

func CreateSymlink(src, dst string) error {
	if info, err := os.Lstat(dst); err == nil {
		if info.Mode()&os.ModeSymlink == 0 {
			if runtime.GOOS == "windows" {
				if err := os.RemoveAll(dst); err != nil {
					return fmt.Errorf("failed to remove existing file: %w", err)
				}
			} else {
				return fmt.Errorf("refusing to overwrite non-symlink file: %s", dst)
			}
		}
		if info.Mode()&os.ModeSymlink != 0 {
			existing, err := os.Readlink(dst)
			if err != nil {
				return fmt.Errorf("failed to read existing symlink: %w", err)
			}
			if existing == src {
				return nil
			}
		}
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("failed to remove existing symlink: %w", err)
		}
	}

	err := os.Symlink(src, dst)
	if err == nil {
		return nil
	}

	if runtime.GOOS == "windows" {
		srcInfo, statErr := os.Stat(src)
		if statErr != nil {
			return err
		}
		if srcInfo.IsDir() {
			return copyDir(src, dst)
		}
		if linkErr := os.Link(src, dst); linkErr == nil {
			return nil
		}
		return copyFile(src, dst)
	}

	return err
}

func copyFile(src, dst string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer srcF.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	dstF, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	if err := dstF.Close(); err != nil {
		return err
	}

	srcInfo, err := srcF.Stat()
	if err == nil {
		os.Chmod(dst, srcInfo.Mode())
	}

	return nil
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source dir: %w", err)
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination dir: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source dir: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
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
		return os.RemoveAll(path)
	}
	if runtime.GOOS == "windows" {
		return os.RemoveAll(path)
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
		fullPath := config.ResolvePath(projectDir, p)
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
		config.ClaudeCodeRootPath:  "claude-code",
		config.CodexRootPath:       "codex",
		config.GeminiCliRootPath:   "gemini-cli",
		config.OpenCodeRootPath:    "opencode",
		config.AntigravityRootPath: "antigravity",
		config.CursorRootPath:      "cursor",
		config.CopilotRootPath:     "github-copilot",
		config.WindsurfRootPath:    "windsurf",
	}

	var links []SymlinkInfo
	for dir, platform := range platformDirs {
		fullDir := config.ResolvePath(projectDir, dir)
		entries, err := os.ReadDir(fullDir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			if !strings.HasPrefix(name, "dotagent-") && !strings.HasPrefix(name, "da-") {
				continue
			}
			fullPath := filepath.Join(fullDir, entry.Name())
			isLink, err := IsSymlink(fullPath)
			if err != nil {
				continue
			}
			if !isLink && runtime.GOOS != "windows" {
				continue
			}
			var target string
			broken := false
			if isLink {
				target, err = os.Readlink(fullPath)
				if err != nil {
					continue
				}
				resolvedTarget := target
				if !filepath.IsAbs(resolvedTarget) {
					resolvedTarget = filepath.Join(filepath.Dir(fullPath), resolvedTarget)
				}
				if _, err := os.Stat(resolvedTarget); err != nil {
					broken = true
				}
			}
			if !strings.HasPrefix(filepath.ToSlash(fullPath), filepath.ToSlash(dotgenDir)) {
				if target != "" {
					resolvedTarget := target
					if !filepath.IsAbs(resolvedTarget) {
						resolvedTarget = filepath.Join(filepath.Dir(fullPath), resolvedTarget)
					}
					if !strings.HasPrefix(resolvedTarget, dotgenDir) {
						continue
					}
				}
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
