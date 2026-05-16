package optimizer

import (
	"fmt"
	"time"

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
// GenerateDNA creates a full DNA_Schema from samples.
func (o *Optimizer) GenerateDNA(projectName string, samples []schema.Sample) schema.DNA_Schema {
	dna := schema.DNA_Schema{
		ProjectName: projectName,
		Version:     "0.1.0",
		GeneratedAt: time.Now(),
		Standards:   []schema.Standard{},
	}

	// Group samples by some category (mocked for now)
	for i, s := range samples {
		std := schema.Standard{
			ID:          fmt.Sprintf("STD-%d", i),
			Title:       fmt.Sprintf("Pattern Detected in %s", s.FilePath),
			Category:    "General",
			Rationale:   "Automatically extracted from codebase structure.",
			Samples:     []schema.Sample{s},
			Rules:       []string{"Follow the structural pattern demonstrated in the sample."},
		}
		dna.Standards = append(dna.Standards, std)
	}

	return dna
}
