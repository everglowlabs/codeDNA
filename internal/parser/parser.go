package parser

import (
	"context"
	"os"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type Parser struct {
	Language *sitter.Language
}

func NewGoParser() *Parser {
	return &Parser{
		Language: golang.GetLanguage(),
	}
}

type CodeBlock struct {
	Content  string
	Kind     string // e.g., "function", "interface", "method"
	FilePath string
}

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

	// Query to find function and method declarations
	queryStr := `
		(function_declaration) @func
		(method_declaration) @method
	`
	q, err := sitter.NewQuery([]byte(queryStr), p.Language)
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
			kind := "function"
			if q.CaptureNameForId(c.Index) == "method" {
				kind = "method"
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
