package agent

type Agent struct {
	Name        string            `json:"name"`
	Content     string            `json:"content"`
	Frontmatter map[string]string `json:"frontmatter"`
	FilePath    string            `json:"filePath"`
}
