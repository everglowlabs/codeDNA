package parser

import (
	"context"
	"os"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

type Parser struct {
	Language *sitter.Language
	Query    string
}

func NewParser(ext string) *Parser {
	switch ext {
	case ".go":
		return &Parser{
			Language: golang.GetLanguage(),
			Query: `
				(function_declaration) @func
				(method_declaration) @method
			`,
		}
	case ".py":
		return &Parser{
			Language: python.GetLanguage(),
			Query: `
				(function_definition) @func
				(class_definition) @class
			`,
		}
	case ".ts", ".tsx":
		return &Parser{
			Language: typescript.GetLanguage(),
			Query: `
				(function_declaration) @func
				(method_definition) @method
				(interface_declaration) @interface
				(class_declaration) @class
			`,
		}
	case ".js":
		return &Parser{
			Language: javascript.GetLanguage(),
			Query: `
				(function_declaration) @func
				(method_definition) @method
				(class_declaration) @class
			`,
		}
	default:
		return nil
	}
}

func NewGoParser() *Parser {
	return NewParser(".go")
}

type CodeBlock struct {
	Content  string
	Kind     string // e.g., "function", "interface", "method", "class"
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
