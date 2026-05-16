# Pattern Detected in internal/optimizer/optimizer.go

Category: General
Rationale: Automatically extracted from codebase structure.

## Guidelines
- Follow the structural pattern demonstrated in the sample.

## Examples
```go
// File: internal/optimizer/optimizer.go
func (o *Optimizer) SelectGoldenSamples(clusters map[string][]parser.CodeBlock) []schema.Sample {
	var samples []schema.Sample

	for _, blocks := range clusters {
		if len(blocks) == 0 {
			continue
		}

		// Heuristic: Pick the longest block as it likely has the most context/patterns.
		// Or pick the one closest to the centroid (not implemented here).
		best := blocks[0]
		for _, b := range blocks {
			if len(b.Content) > len(best.Content) {
				best = b
			}
		}

		samples = append(samples, schema.Sample{
			FilePath: best.FilePath,
			Content:  best.Content,
			Language: "go", // TODO: Detect from extension
		})
	}

	return samples
}
```

