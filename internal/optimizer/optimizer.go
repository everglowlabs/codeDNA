package optimizer

import (
	"github.com/everglowlabs/codedna/internal/parser"
	"github.com/everglowlabs/codedna/internal/schema"
)

type Optimizer struct {
	MaxTokens int
}

func NewOptimizer(maxTokens int) *Optimizer {
	return &Optimizer{
		MaxTokens: maxTokens,
	}
}

// SelectGoldenSamples picks the best representatives from each cluster.
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

// TrimPayload ensures the total content doesn't exceed token limits.
// This is a placeholder for more advanced token-aware trimming.
func (o *Optimizer) TrimPayload(samples []schema.Sample) []schema.Sample {
	totalChars := 0
	maxChars := o.MaxTokens * 4 // Rough estimate: 4 chars per token

	var trimmed []schema.Sample
	for _, s := range samples {
		if totalChars+len(s.Content) > maxChars {
			break
		}
		trimmed = append(trimmed, s)
		totalChars += len(s.Content)
	}

	return trimmed
}
