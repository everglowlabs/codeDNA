# Pattern Detected in internal/compiler/compiler.go

Category: General
Rationale: Automatically extracted from codebase structure.

## Guidelines
- Follow the structural pattern demonstrated in the sample.

## Examples
```go
// File: internal/compiler/compiler.go
func (c *CopilotCompiler) Compile(dna schema.DNA_Schema) (map[string]string, error) {
	outputs := make(map[string]string)

	for _, s := range dna.Standards {
		var sb strings.Builder

		// Determine globs from samples
		globMap := make(map[string]bool)
		for _, sample := range s.Samples {
			dir := filepath.Dir(sample.FilePath)
			ext := filepath.Ext(sample.FilePath)
			if dir == "." {
				globMap["*"+ext] = true
			} else {
				globMap[dir+"/**/*"+ext] = true
			}
		}

		var globs []string
		for g := range globMap {
			globs = append(globs, g)
		}

		// YAML Frontmatter
		sb.WriteString("---\n")
		sb.WriteString("globs:\n")
		if len(globs) == 0 {
			sb.WriteString("  - \"**/*\"\n")
		} else {
			for _, g := range globs {
				sb.WriteString(fmt.Sprintf("  - \"%s\"\n", g))
			}
		}
		sb.WriteString("---\n\n")

		sb.WriteString("# " + s.Title + "\n\n")
		sb.WriteString(s.Description + "\n\n")
		sb.WriteString("## Rationale\n")
		sb.WriteString(s.Rationale + "\n\n")
		sb.WriteString("## Guidelines\n")
		for _, rule := range s.Rules {
			sb.WriteString(fmt.Sprintf("- %s\n", rule))
		}

		// Sanitize title for filename
		filename := strings.ToLower(s.Title)
		filename = strings.ReplaceAll(filename, " ", "_")
		filename = strings.ReplaceAll(filename, "/", "_")
		fullPath := fmt.Sprintf(".github/instructions/%s.instructions.md", filename)

		outputs[fullPath] = sb.String()
	}

	// Also add the legacy root file if there are global instructions?
	// For now, let's just stick to modularity as it's cleaner.

	return outputs, nil
}
```

