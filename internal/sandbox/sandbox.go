package sandbox

import (
	"fmt"
	"strings"

	"github.com/everglowlabs/codedna/internal/parser"
	"github.com/everglowlabs/codedna/internal/schema"
)

// Sandbox provides a way to validate code against DNA standards.
type Sandbox struct {
	parser *parser.Parser
}

func NewSandbox() *Sandbox {
	return &Sandbox{
		parser: parser.NewGoParser(),
	}
}

// ValidationResult represents the outcome of a linting/validation check.
type ValidationResult struct {
	Passed  bool
	Issues  []string
	Summary string
}

// ValidateSnippet checks if a code snippet adheres to the provided standards.
func (s *Sandbox) ValidateSnippet(snippet string, standards []schema.Standard) ValidationResult {
	// For now, this is a structural heuristic.
	// In the future, this would use the LLM + AST comparison.
	
	result := ValidationResult{Passed: true}
	
	for _, std := range standards {
		// Basic heuristic: if the pattern is a keyword that should be present/absent
		// This is very rudimentary for now.
		if strings.Contains(std.Category, "Error Handling") && !strings.Contains(snippet, "fmt.Errorf") {
			result.Passed = false
			result.Issues = append(result.Issues, fmt.Sprintf("Standard '%s' violation: Missing recommended error wrapping (fmt.Errorf)", std.Title))
		}
	}
	
	if result.Passed {
		result.Summary = "Code adheres to project DNA standards."
	} else {
		result.Summary = fmt.Sprintf("Detected %d violations of project DNA.", len(result.Issues))
	}
	
	return result
}

// LintFile checks a file on disk against DNA standards.
func (s *Sandbox) LintFile(path string, standards []schema.Standard) (ValidationResult, error) {
	blocks, err := s.parser.ExtractBlocks(path)
	if err != nil {
		return ValidationResult{}, err
	}
	
	allIssues := []string{}
	passed := true
	
	for _, block := range blocks {
		res := s.ValidateSnippet(block.Content, standards)
		if !res.Passed {
			passed = false
			for _, issue := range res.Issues {
				allIssues = append(allIssues, fmt.Sprintf("%s: %s", path, issue))
			}
		}
	}
	
	return ValidationResult{
		Passed:  passed,
		Issues:  allIssues,
		Summary: fmt.Sprintf("Linted %s: %d issues", path, len(allIssues)),
	}, nil
}
