package compiler

import (
	"strings"
	"testing"

	"github.com/everglowlabs/codedna/internal/schema"
)

func TestAntigravityCompiler(t *testing.T) {
	c := &AntigravityCompiler{}
	dna := schema.DNA_Schema{
		ProjectName: "TestProject",
		Standards: []schema.Standard{
			{
				Title:    "Logging Pattern",
				Category: "Logging",
				Rationale: "Consistent logging is good.",
				Rules:    []string{"Use JSON logging", "Include trace ID"},
				Samples: []schema.Sample{
					{
						FilePath: "main.go",
						Language: "go",
						Content:  "log.Info().Msg(\"hello\")",
					},
				},
			},
			{
				Title:    "Error Handling",
				Category: "Errors",
				Rationale: "Wrap errors.",
				Rules:    []string{"Wrap with %w"},
			},
		},
	}

	files, err := c.Compile(dna)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}

	expectedFiles := []string{
		".agents/rules/logging_pattern.md",
		".agents/rules/error_handling.md",
	}

	for _, ef := range expectedFiles {
		if _, ok := files[ef]; !ok {
			t.Errorf("Missing expected file: %s", ef)
		}
	}

	// Check content of one file
	content := files[".agents/rules/logging_pattern.md"]
	if len(content) == 0 {
		t.Error("Empty content for logging_pattern.md")
	}
	if len(content) > 12000 {
		t.Error("Content exceeds 12,000 characters")
	}
}

func TestCopilotCompiler(t *testing.T) {
	c := &CopilotCompiler{}
	dna := schema.DNA_Schema{
		ProjectName: "TestProject",
		Standards: []schema.Standard{
			{
				Title: "DB Standards",
				Samples: []schema.Sample{
					{FilePath: "internal/db/postgres.go"},
				},
				Rules: []string{"No raw SQL"},
			},
		},
	}

	files, err := c.Compile(dna)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	expectedFile := ".github/instructions/db_standards.instructions.md"
	content, ok := files[expectedFile]
	if !ok {
		t.Fatalf("Missing expected file: %s", expectedFile)
	}

	if !strings.Contains(content, "globs:") {
		t.Error("Missing YAML frontmatter")
	}
	if !strings.Contains(content, "internal/db/**/*.go") {
		t.Errorf("Expected glob internal/db/**/*.go, got: %s", content)
	}
}

func TestClaudeCompiler(t *testing.T) {
	c := &ClaudeCompiler{}
	dna := schema.DNA_Schema{
		ProjectName: "TestProject",
		Standards: []schema.Standard{
			{
				Title: "Logging",
				Samples: []schema.Sample{
					{FilePath: "cmd/main.go"},
				},
				Rules: []string{"Use zap"},
			},
		},
	}

	files, err := c.Compile(dna)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Check CLAUDE.md
	if _, ok := files["CLAUDE.md"]; !ok {
		t.Error("Missing CLAUDE.md")
	}

	// Check modular rule
	expectedFile := ".claude/rules/logging.md"
	content, ok := files[expectedFile]
	if !ok {
		t.Fatalf("Missing expected file: %s", expectedFile)
	}

	if !strings.Contains(content, "paths:") {
		t.Error("Missing YAML frontmatter with 'paths'")
	}
	if !strings.Contains(content, "cmd/**/*") {
		t.Errorf("Expected path cmd/**/*, got: %s", content)
	}
}

func TestCursorCompiler(t *testing.T) {
	c := &CursorCompiler{}
	dna := schema.DNA_Schema{
		ProjectName: "TestProject",
		Standards: []schema.Standard{
			{
				Title:       "API Standards",
				Description: "API design rules",
				Samples: []schema.Sample{
					{FilePath: "api/v1/user.go"},
				},
				Rules: []string{"Use REST"},
			},
		},
	}

	files, err := c.Compile(dna)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	expectedFile := ".cursor/rules/api_standards.mdc"
	content, ok := files[expectedFile]
	if !ok {
		t.Fatalf("Missing expected file: %s", expectedFile)
	}

	if !strings.Contains(content, "description: API design rules") {
		t.Error("Missing description in frontmatter")
	}
	if !strings.Contains(content, "globs: api/v1/**/*.go") {
		t.Errorf("Expected glob api/v1/**/*.go, got: %s", content)
	}
}

func TestContinueCompiler(t *testing.T) {
	c := &ContinueCompiler{}
	dna := schema.DNA_Schema{
		ProjectName: "TestProject",
		Standards: []schema.Standard{
			{
				Title:       "Error Handling",
				Description: "Error rules",
				Samples: []schema.Sample{
					{FilePath: "internal/errs/errors.go"},
				},
				Rules: []string{"Wrap errors"},
			},
		},
	}

	files, err := c.Compile(dna)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	expectedFile := ".continue/rules/error_handling.md"
	content, ok := files[expectedFile]
	if !ok {
		t.Fatalf("Missing expected file: %s", expectedFile)
	}

	if !strings.Contains(content, "name: Error Handling") {
		t.Error("Missing name in frontmatter")
	}
	if !strings.Contains(content, "internal/errs/**/*.go") {
		t.Errorf("Expected glob internal/errs/**/*.go, got: %s", content)
	}
}

func TestWindsurfCompiler(t *testing.T) {
	c := &WindsurfCompiler{}
	dna := schema.DNA_Schema{
		ProjectName: "TestProject",
		Standards: []schema.Standard{
			{
				Title:       "Testing Standards",
				Description: "Unit testing rules",
				Samples: []schema.Sample{
					{FilePath: "internal/service/user_test.go"},
				},
				Rules: []string{"Use testify"},
			},
		},
	}

	files, err := c.Compile(dna)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	expectedFile := ".windsurf/rules/testing_standards.md"
	content, ok := files[expectedFile]
	if !ok {
		t.Fatalf("Missing expected file: %s", expectedFile)
	}

	if !strings.Contains(content, "trigger: glob") {
		t.Error("Missing trigger in frontmatter")
	}
	if !strings.Contains(content, "internal/service/**/*.go") {
		t.Errorf("Expected glob internal/service/**/*.go, got: %s", content)
	}
}

func TestJetBrainsCompiler(t *testing.T) {
	c := &JetBrainsCompiler{}
	dna := schema.DNA_Schema{
		ProjectName: "TestProject",
		Standards: []schema.Standard{
			{
				Title:       "Format Standards",
				Description: "Code formatting rules",
				Samples: []schema.Sample{
					{FilePath: "pkg/utils/format.go"},
				},
				Rules: []string{"Use gofmt"},
			},
		},
	}

	files, err := c.Compile(dna)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	expectedFile := ".aiassistant/rules/format_standards.md"
	content, ok := files[expectedFile]
	if !ok {
		t.Fatalf("Missing expected file: %s", expectedFile)
	}

	if !strings.Contains(content, "type: file patterns") {
		t.Error("Missing type in frontmatter")
	}
	if !strings.Contains(content, "pattern: \"pkg/utils/**/*.go\"") {
		t.Errorf("Expected pattern pkg/utils/**/*.go, got: %s", content)
	}
}
