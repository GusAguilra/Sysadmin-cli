package models

type Database struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Commands    []Command `json:"commands"`
}

type Command struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Command     string   `json:"command"`
	Example     string   `json:"example"`
	Tags        []string `json:"tags"`
	Notes       string   `json:"notes"`
}

func (c Command) Matches(query string) bool {
	if query == "" {
		return true
	}
	q := query
	for _, tag := range c.Tags {
		if contains(tag, q) {
			return true
		}
	}
	if contains(c.Title, q) {
		return true
	}
	if contains(c.Description, q) {
		return true
	}
	if contains(c.Command, q) {
		return true
	}
	if contains(c.Notes, q) {
		return true
	}
	return false
}

func (c Category) Matches(query string) bool {
	if query == "" {
		return true
	}
	q := query
	if contains(c.Name, q) {
		return true
	}
	if contains(c.Description, q) {
		return true
	}
	for _, cmd := range c.Commands {
		if cmd.Matches(q) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(substr) > 0 && containsLower(s, substr)
}

func containsLower(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	s = toLower(s)
	substr = toLower(substr)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}
