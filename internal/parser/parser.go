package parser

import (
	"context"
	"os"

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

func (p *Parser) ExtractFunctions(path string) ([]string, error) {
	content, err := os.ReadFile(path)
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
	
	// Very basic query to find function declarations
	queryStr := `(function_declaration name: (identifier) @func.name)`
	q, err := sitter.NewQuery([]byte(queryStr), p.Language)
	if err != nil {
		return nil, err
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	var functions []string
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, c := range m.Captures {
			functions = append(functions, c.Node.Content(content))
		}
	}

	return functions, nil
}
