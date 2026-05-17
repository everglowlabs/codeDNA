package cluster

import (
	"context"
	"fmt"

	"github.com/everglowlabs/codedna/internal/parser"
	"github.com/philippgille/chromem-go"
)

type ClusterManager struct {
	db         *chromem.DB
	collection *chromem.Collection
}

func NewClusterManager(collectionName string, embedFunc chromem.EmbeddingFunc) (*ClusterManager, error) {
	db := chromem.NewDB()
	col, err := db.CreateCollection(collectionName, nil, embedFunc)
	if err != nil {
		return nil, err
	}

	return &ClusterManager{
		db:         db,
		collection: col,
	}, nil
}

func (cm *ClusterManager) AddBlocks(ctx context.Context, blocks []parser.CodeBlock) error {
	ids := make([]string, len(blocks))
	metadatas := make([]map[string]string, len(blocks))
	contents := make([]string, len(blocks))

	for i, block := range blocks {
		ids[i] = fmt.Sprintf("%s-%d", block.FilePath, i)
		metadatas[i] = map[string]string{
			"file_path": block.FilePath,
			"kind":      block.Kind,
		}
		contents[i] = block.Content
	}

	return cm.collection.Add(ctx, ids, nil, metadatas, contents)
}

// FindSimilar finds code blocks similar to the given content.
func (cm *ClusterManager) FindSimilar(ctx context.Context, content string, limit int) ([]parser.CodeBlock, error) {
	results, err := cm.collection.Query(ctx, content, limit, nil, nil)
	if err != nil {
		return nil, err
	}

	var blocks []parser.CodeBlock
	for _, res := range results {
		blocks = append(blocks, parser.CodeBlock{
			Content:  res.Content,
			Kind:     res.Metadata["kind"],
			FilePath: res.Metadata["file_path"],
		})
	}
	return blocks, nil
}

// GroupByPattern groups blocks by their semantic similarity.
// This is a simplified version of clustering for Phase 2.
func (cm *ClusterManager) GroupByPattern(
	ctx context.Context, blocks []parser.CodeBlock) (map[string][]parser.CodeBlock, error) {
	clusters := make(map[string][]parser.CodeBlock)
	seen := make(map[string]bool)

	for _, block := range blocks {
		if seen[block.Content] {
			continue
		}

		limit := 5
		if len(blocks) < limit {
			limit = len(blocks)
		}
		if limit == 0 {
			continue
		}
		similar, err := cm.FindSimilar(ctx, block.Content, limit)
		if err != nil {
			continue
		}

		clusterID := fmt.Sprintf("pattern-%d", len(clusters))
		for _, s := range similar {
			if !seen[s.Content] {
				clusters[clusterID] = append(clusters[clusterID], s)
				seen[s.Content] = true
			}
		}
	}

	return clusters, nil
}
