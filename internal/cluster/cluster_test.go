package cluster

import (
	"context"
	"testing"

	"github.com/everglowlabs/codedna/internal/parser"
)

func TestClusterManager(t *testing.T) {
	cm, err := NewClusterManager("test_col", dummyEmbedder)
	if err != nil {
		t.Fatal(err)
	}

	blocks := []parser.CodeBlock{
		{
			Content:  "fn hello() {}",
			Kind:     "function",
			FilePath: "test.rs",
		},
		{
			Content:  "struct UserProfile {}",
			Kind:     "class",
			FilePath: "test.rs",
		},
	}

	ctx := context.Background()
	err = cm.AddBlocks(ctx, blocks)
	if err != nil {
		t.Fatalf("AddBlocks failed: %v", err)
	}

	clusters, err := cm.GroupByPattern(ctx, blocks)
	if err != nil {
		t.Fatalf("GroupByPattern failed: %v", err)
	}

	t.Logf("Found %d clusters", len(clusters))
	for id, b := range clusters {
		t.Logf("  Cluster %s has %d blocks", id, len(b))
	}

	if len(clusters) == 0 {
		t.Error("Expected to get at least one cluster, got 0")
	}
}

func dummyEmbedder(ctx context.Context, text string) ([]float32, error) {
	res := make([]float32, 384)
	res[0] = 1.0
	return res, nil
}
