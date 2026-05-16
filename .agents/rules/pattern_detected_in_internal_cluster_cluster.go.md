# Pattern Detected in internal/cluster/cluster.go

Category: General
Rationale: Automatically extracted from codebase structure.

## Guidelines
- Follow the structural pattern demonstrated in the sample.

## Examples
```go
// File: internal/cluster/cluster.go
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
```

