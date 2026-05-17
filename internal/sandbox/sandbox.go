package sandbox

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/everglowlabs/codedna/internal/parser"
	"github.com/everglowlabs/codedna/internal/schema"
)

// Sandbox provides a way to validate code against DNA standards.
type Sandbox struct{}

func NewSandbox() *Sandbox {
	return &Sandbox{}
}

type Violation struct {
	StandardID string
	FilePath   string
	Message    string
}

// ValidationResult represents the outcome of a linting/validation check.
type ValidationResult struct {
	Passed     bool
	Issues     []string
	Summary    string
	Violations []Violation
}

// ValidateSnippet checks if a code snippet adheres to the provided standards.
func (s *Sandbox) ValidateSnippet(snippet string, standards []schema.Standard) ValidationResult {
	result := ValidationResult{Passed: true}
	
	for _, std := range standards {
		// Basic heuristic: if the pattern is a keyword that should be present/absent
		if strings.Contains(std.Category, "Error Handling") && !strings.Contains(snippet, "fmt.Errorf") && !strings.Contains(snippet, "unwrap") && !strings.Contains(snippet, "throw ") && !strings.Contains(snippet, "except ") {
			result.Passed = false
			msg := "Missing recommended error wrapping/handling"
			result.Issues = append(result.Issues, fmt.Sprintf("Standard '%s' violation: %s", std.Title, msg))
			result.Violations = append(result.Violations, Violation{
				StandardID: std.ID,
				Message:    fmt.Sprintf("Standard '%s' (ID: %s) violation: %s", std.Title, std.ID, msg),
			})
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
	ext := filepath.Ext(path)
	p := parser.NewParser(ext)
	if p == nil {
		return ValidationResult{Passed: true, Summary: "File type not supported for parsing"}, nil
	}

	blocks, err := p.ExtractBlocks(path)
	if err != nil {
		return ValidationResult{}, err
	}
	
	allIssues := []string{}
	allViolations := []Violation{}
	passed := true
	
	for _, block := range blocks {
		res := s.ValidateSnippet(block.Content, standards)
		if !res.Passed {
			passed = false
			for _, issue := range res.Issues {
				allIssues = append(allIssues, fmt.Sprintf("%s: %s", path, issue))
			}
			for _, v := range res.Violations {
				allViolations = append(allViolations, Violation{
					StandardID: v.StandardID,
					FilePath:   path,
					Message:    v.Message,
				})
			}
		}
	}
	
	return ValidationResult{
		Passed:     passed,
		Issues:     allIssues,
		Violations: allViolations,
		Summary:    fmt.Sprintf("Linted %s: %d issues", path, len(allIssues)),
	}, nil
}
