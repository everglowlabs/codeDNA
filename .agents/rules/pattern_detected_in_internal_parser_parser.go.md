# Pattern Detected in internal/parser/parser.go

Category: General
Rationale: Automatically extracted from codebase structure.

## Guidelines
- Follow the structural pattern demonstrated in the sample.

## Examples
```go
// File: internal/parser/parser.go
func (p *Parser) ExtractBlocks(path string) ([]CodeBlock, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	cpp := sitter.NewParser()
	cpp.SetLanguage(p.Language)

	tree, err := cpp.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, err
	}

	n := tree.RootNode()

	q, err := sitter.NewQuery([]byte(p.Query), p.Language)
	if err != nil {
		return nil, err
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	var blocks []CodeBlock
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, c := range m.Captures {
			kind := q.CaptureNameForId(c.Index)
			// Map internal capture names to user-friendly kinds
			switch kind {
			case "func":
				kind = "function"
			case "method":
				kind = "method"
			case "class":
				kind = "class"
			case "interface":
				kind = "interface"
			}

			blocks = append(blocks, CodeBlock{
				Content:  c.Node.Content(content),
				Kind:     kind,
				FilePath: path,
			})
		}
	}

	return blocks, nil
}
```

