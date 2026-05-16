# Pattern Detected in internal/compiler/compiler_test.go

Category: General
Rationale: Automatically extracted from codebase structure.

## Guidelines
- Follow the structural pattern demonstrated in the sample.

## Examples
```go
// File: internal/compiler/compiler_test.go
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
```

