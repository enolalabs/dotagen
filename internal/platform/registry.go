package platform

import "fmt"

type Registry struct {
	adapters map[string]Adapter
}

func NewRegistry() *Registry {
	r := &Registry{adapters: make(map[string]Adapter)}
	r.Register(NewAntigravityAdapter())
	r.Register(NewClaudeCodeAdapter())
	r.Register(NewCodexAdapter())
	r.Register(NewGeminiCLIAdapter())
	r.Register(NewOpenCodeAdapter())
	return r
}

func (r *Registry) Register(a Adapter) {
	r.adapters[a.Name()] = a
}

func (r *Registry) Get(name string) (Adapter, error) {
	a, ok := r.adapters[name]
	if !ok {
		return nil, fmt.Errorf("unknown platform: %s", name)
	}
	return a, nil
}

// GetSkillAdapter returns the SkillAdapter for a platform, if it implements one.
func (r *Registry) GetSkillAdapter(name string) (SkillAdapter, error) {
	a, ok := r.adapters[name]
	if !ok {
		return nil, fmt.Errorf("unknown platform: %s", name)
	}
	sa, ok := a.(SkillAdapter)
	if !ok {
		return nil, fmt.Errorf("platform %s does not support skills", name)
	}
	return sa, nil
}

func (r *Registry) List() []string {
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	return names
}
