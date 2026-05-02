package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
)

// ParseSkillsDir reads all skill directories under the given path.
// Each subdirectory must contain a SKILL.md file.
func ParseSkillsDir(skillsDir string) ([]Skill, error) {
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read skills directory: %w", err)
	}

	skills := make([]Skill, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirPath := filepath.Join(skillsDir, entry.Name())
		sk, err := ParseSkillDir(dirPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse skill %s: %w", entry.Name(), err)
		}
		if sk != nil {
			skills = append(skills, *sk)
		}
	}

	return skills, nil
}

// ParseSkillDir parses a single skill directory containing a SKILL.md file
// and optional reference files.
func ParseSkillDir(dirPath string) (*Skill, error) {
	skillFile := filepath.Join(dirPath, "SKILL.md")
	data, err := os.ReadFile(skillFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Skip directories without SKILL.md
		}
		return nil, fmt.Errorf("failed to read SKILL.md: %w", err)
	}

	name := filepath.Base(dirPath)

	var fm map[string]string
	body, err := frontmatter.Parse(strings.NewReader(string(data)), &fm)
	if err != nil {
		body = data
		fm = make(map[string]string)
	}

	content := strings.TrimSpace(string(body))

	refs, err := collectReferences(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to collect references: %w", err)
	}

	return &Skill{
		Name:        name,
		Content:     content,
		Frontmatter: fm,
		DirPath:     dirPath,
		References:  refs,
	}, nil
}

// collectReferences walks the skill directory and collects all files
// that are not SKILL.md as references.
func collectReferences(dirPath string) ([]Reference, error) {
	var refs []Reference

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Skip the main SKILL.md file
		if filepath.Base(path) == "SKILL.md" && filepath.Dir(path) == dirPath {
			return nil
		}

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read reference %s: %w", relPath, err)
		}

		refs = append(refs, Reference{
			Name:     relPath,
			FilePath: path,
			Content:  string(content),
		})
		return nil
	})

	return refs, err
}

// ExtractDescription returns the description from frontmatter or first non-empty line.
func ExtractDescription(sk Skill) string {
	if desc, ok := sk.Frontmatter["description"]; ok && desc != "" {
		return desc
	}
	lines := strings.Split(sk.Content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			return strings.TrimPrefix(trimmed, "# ")
		}
	}
	return ""
}

// ScaffoldSkillContent returns a template SKILL.md content for a new skill.
func ScaffoldSkillContent(name string) string {
	cleanName := strings.TrimPrefix(name, "ds-")
	title := strings.ReplaceAll(cleanName, "-", " ")
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	displayName := strings.Join(words, " ")

	return fmt.Sprintf(`---
name: %s
description: TODO — describe when this skill should trigger.
category: custom
---

# %s

## When to Use

TODO: Describe when to invoke this skill.

## Workflow

1. Step one
2. Step two
3. Step three

## Checklist

- [ ] Item 1
- [ ] Item 2
`, cleanName, displayName)
}
