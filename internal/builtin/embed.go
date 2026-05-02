package builtin

import (
	"embed"
	"io/fs"
	"sort"
	"strings"
)

//go:embed agents/*.md
var DefaultAgents embed.FS

//go:embed all:skills
var DefaultSkills embed.FS

func ListAgents() []string {
	entries, err := fs.ReadDir(DefaultAgents, "agents")
	if err != nil {
		return nil
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		names = append(names, strings.TrimSuffix(entry.Name(), ".md"))
	}
	sort.Strings(names)
	return names
}

func ReadAgent(name string) ([]byte, error) {
	return fs.ReadFile(DefaultAgents, "agents/"+name+".md")
}

// ListSkills returns all skill directory names (e.g. "ds-tdd", "ds-diagnose").
func ListSkills() []string {
	entries, err := fs.ReadDir(DefaultSkills, "skills")
	if err != nil {
		return nil
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Verify it has a SKILL.md
		if _, err := fs.Stat(DefaultSkills, "skills/"+entry.Name()+"/SKILL.md"); err == nil {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names
}

// ReadSkillFile reads a single file from a builtin skill directory.
// path is relative to skills/, e.g. "ds-tdd/SKILL.md" or "ds-tdd/tests.md"
func ReadSkillFile(path string) ([]byte, error) {
	return fs.ReadFile(DefaultSkills, "skills/"+path)
}

// ListSkillFiles returns all files in a skill directory (relative paths).
func ListSkillFiles(skillName string) []string {
	var files []string
	dir := "skills/" + skillName
	fs.WalkDir(DefaultSkills, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := relPath(dir, path)
		if err != nil {
			return nil
		}
		files = append(files, rel)
		return nil
	})
	return files
}

func relPath(base, target string) (string, error) {
	if !strings.HasPrefix(target, base+"/") {
		return target, nil
	}
	return target[len(base)+1:], nil
}
