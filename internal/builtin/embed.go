package builtin

import (
	"embed"
	"io/fs"
	"sort"
	"strings"
)

//go:embed agents/*.md
var DefaultAgents embed.FS

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
