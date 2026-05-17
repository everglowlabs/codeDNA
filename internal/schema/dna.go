package schema

import "time"

// DNA_Schema represents the structured engineering standards of a project.
type DNA_Schema struct {
	ProjectName string     `json:"project_name"`
	Version     string     `json:"version"`
	GeneratedAt time.Time  `json:"generated_at"`
	Standards   []Standard `json:"standards"`
}

// Standard defines a specific engineering category or rule.
type Standard struct {
	ID          string    `json:"id"`
	Severity    string    `json:"severity"` // "error", "warning", "info"
	Category    string    `json:"category"` // e.g., "Error Handling", "Logging", "Testing"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Rationale   string    `json:"rationale"`
	Patterns    []string  `json:"patterns"` // Abstracted patterns detected
	Samples     []Sample  `json:"samples"`  // Golden samples of code
	Rules       []string  `json:"rules"`    // Concrete instructions for the AI agent
}

// Sample represents a "Golden Sample" of code that demonstrates a standard.
type Sample struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
	Language string `json:"language"`
}
