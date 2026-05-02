package skill

// Skill represents a parsed skill definition from a SKILL.md file.
type Skill struct {
	Name        string            `json:"name"`
	Content     string            `json:"content"`
	Frontmatter map[string]string `json:"frontmatter"`
	DirPath     string            `json:"dirPath"`
	References  []Reference       `json:"references,omitempty"`
}

// Reference represents a supporting file bundled with a skill.
type Reference struct {
	Name     string `json:"name"`     // relative path, e.g. "tests.md" or "scripts/hitl-loop.template.sh"
	FilePath string `json:"filePath"` // absolute path on disk
	Content  string `json:"content"`  // file content (loaded on demand)
}
