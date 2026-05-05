package skill

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Skill represents a parsed skill definition from a SKILL.md file.
//
// A skill is a directory containing a SKILL.md and optional reference files.
// Reference files (scripts, templates, docs) are deployed alongside SKILL.md
// in the same directory so that AI agents can access them as real files.
type Skill struct {
	Name        string            `json:"name"`
	Content     string            `json:"content"`
	Frontmatter map[string]string `json:"frontmatter"`
	DirPath     string            `json:"dirPath"`
	References  []Reference       `json:"references,omitempty"`
}

// Reference represents a supporting file bundled with a skill.
// Reference files are deployed as real files alongside SKILL.md in the
// skill directory, preserving their relative paths.
type Reference struct {
	Name     string `json:"name"`     // relative path, e.g. "tests.md" or "scripts/hitl-loop.template.sh"
	FilePath string `json:"filePath"` // absolute path on disk
	Content  string `json:"content"`  // file content (loaded on demand)
}

// ReferenceNames returns the relative paths of all reference files.
func (sk Skill) ReferenceNames() []string {
	names := make([]string, len(sk.References))
	for i, ref := range sk.References {
		names[i] = ref.Name
	}
	return names
}

// langFromExt returns a code fence language hint for common file extensions.
func langFromExt(ext string) string {
	switch ext {
	case ".sh", ".bash":
		return "bash"
	case ".py":
		return "python"
	case ".go":
		return "go"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	case ".rb":
		return "ruby"
	case ".rs":
		return "rust"
	case ".java":
		return "java"
	case ".sql":
		return "sql"
	default:
		return ""
	}
}

// ContentWithReferences returns the skill content with all reference file
// contents appended. This is a fallback for platforms that cannot read
// skill directories and can only consume a single file.
//
// For the standard directory-based flow, use sk.Content directly — reference
// files are deployed as real files alongside SKILL.md.
func (sk Skill) ContentWithReferences(symlinkDir string) string {
	if len(sk.References) == 0 {
		return sk.Content
	}

	var sb strings.Builder
	sb.WriteString(sk.Content)
	sb.WriteString("\n\n---\n\n## Bundled Reference Files\n")
	sb.WriteString("\nThe following files are bundled with this skill and available ")
	sb.WriteString("in the skill directory alongside this SKILL.md.\n")

	for _, ref := range sk.References {
		sb.WriteString(fmt.Sprintf("\n### %s\n\n", ref.Name))

		ext := strings.ToLower(filepath.Ext(ref.Name))
		if ext == ".md" || ext == ".markdown" {
			sb.WriteString(ref.Content)
		} else {
			if symlinkDir != "" {
				filePath := filepath.Join(symlinkDir, ref.Name)
				sb.WriteString(fmt.Sprintf("> **File path**: `%s`\n", filePath))
				sb.WriteString("> This is an actual file on disk — use the path above when copying or executing.\n\n")
			}
			lang := langFromExt(ext)
			sb.WriteString(fmt.Sprintf("```%s\n%s\n```\n", lang, ref.Content))
		}
	}

	return sb.String()
}
