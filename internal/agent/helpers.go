package agent

import (
	"fmt"
	"strings"
)

func ExtractDescription(ag Agent) string {
	if desc, ok := ag.Frontmatter["description"]; ok && desc != "" {
		return desc
	}
	lines := strings.Split(ag.Content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			return strings.TrimPrefix(trimmed, "# ")
		}
	}
	return ""
}

func ScaffoldContent(name string) string {
	title := DisplayName(name)
	return fmt.Sprintf(`# %s

## Role

TODO: Describe the agent's role and expertise.

## Guidelines

- Guideline 1
- Guideline 2

## Examples

TODO: Add examples of expected behavior.
`, title)
}

func DisplayName(name string) string {
	title := strings.ReplaceAll(name, "-", " ")
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}
