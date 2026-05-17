package optimizer

import (
	"fmt"
	"path/filepath"
	"strings"
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
		best := blocks[0]
		for _, b := range blocks {
			if len(b.Content) > len(best.Content) {
				best = b
			}
		}

		ext := filepath.Ext(best.FilePath)
		lang := "go"
		switch ext {
		case ".go":
			lang = "go"
		case ".py":
			lang = "python"
		case ".js":
			lang = "javascript"
		case ".ts", ".tsx":
			lang = "typescript"
		case ".rs":
			lang = "rust"
		case ".java":
			lang = "java"
		}

		samples = append(samples, schema.Sample{
			FilePath: best.FilePath,
			Content:  best.Content,
			Language: lang,
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

	// Generate realistic multi-lingual standards based on the sample content
	for i, s := range samples {
		var std schema.Standard
		std.ID = fmt.Sprintf("std-%03d", i+1)
		std.Samples = []schema.Sample{s}

		contentLower := strings.ToLower(s.Content)

		switch s.Language {
		case "go":
			if strings.Contains(contentLower, "fmt.errorf") || strings.Contains(contentLower, "errors.new") || strings.Contains(contentLower, "err !=") {
				std.Title = "Wrapped Errors with Context"
				std.Category = "Error Handling"
				std.Severity = "error"
				std.Rationale = "Error wrapping enables callers to use errors.Is / errors.As and preserves the full stack of context through the call chain."
				std.Patterns = []string{"fmt.Errorf(\"...: %w\", err)", "errors.Is(err, target)"}
				std.Rules = []string{
					"Always wrap errors with fmt.Errorf(\"context: %w\", err) — never return bare errors.",
					"Use errors.Is for comparison, never == on error values.",
					"Log errors at the boundary where they are handled, not where they are created.",
				}
			} else if strings.Contains(contentLower, "log.") || strings.Contains(contentLower, "zap.") || strings.Contains(contentLower, "logger") {
				std.Title = "Structured Logging Pattern"
				std.Category = "Logging"
				std.Severity = "warning"
				std.Rationale = "Standardized structured logging enables reliable log parsing and aggregation in production environments."
				std.Patterns = []string{"logger.Info(\"...\", zap.String(...))"}
				std.Rules = []string{
					"Always use structured log fields instead of string formatting inside log messages.",
					"Include contextual fields like request_id, user_id, or transaction_id where available.",
				}
			} else {
				std.Title = "Go idiomatic Function and Method Structure"
				std.Category = "General Code Structure"
				std.Severity = "info"
				std.Rationale = "Adhering to Go standard formatting rules improves codebase readability and consistency."
				std.Patterns = []string{"func (...) ..."}
				std.Rules = []string{
					"Keep functions and methods below 50 lines where possible.",
					"Accept interfaces and return concrete types.",
				}
			}

		case "typescript", "javascript":
			if strings.Contains(contentLower, "try {") || strings.Contains(contentLower, "catch (") {
				std.Title = "TS/JS Exception Handling"
				std.Category = "Error Handling"
				std.Severity = "error"
				std.Rationale = "Robust try-catch blocks prevent unhandled promise rejections and process crashes."
				std.Patterns = []string{"try { ... } catch (error) { ... }"}
				std.Rules = []string{
					"Always type cast caught errors as Error (e.g. error instanceof Error) before accessing message.",
					"Provide clear, user-friendly fallback behaviors inside catch blocks.",
				}
			} else if strings.Contains(contentLower, "interface ") || strings.Contains(contentLower, "type ") {
				std.Title = "TypeScript Type Safety"
				std.Category = "Interface & Typing"
				std.Severity = "info"
				std.Rationale = "Strong type contracts ensure compile-time safety and reduce runtime exceptions."
				std.Patterns = []string{"interface User { ... }"}
				std.Rules = []string{
					"Avoid using 'any' - prefer strict type annotations or 'unknown' for dynamic data.",
					"Mark optional properties explicitly using '?' instead of allowing undefined values.",
				}
			} else {
				std.Title = "Modern ES6+ Code Standard"
				std.Category = "General Code Structure"
				std.Severity = "warning"
				std.Rationale = "Consistent ES6 syntax improves runtime performance and maintains modern formatting styles."
				std.Patterns = []string{"const x = () => { ... }"}
				std.Rules = []string{
					"Always use 'const' or 'let' - never use 'var'.",
					"Prefer async/await over raw Promise chain callbacks.",
				}
			}

		case "python":
			if strings.Contains(contentLower, "except ") || strings.Contains(contentLower, "raise ") {
				std.Title = "Python Exception Scope Standard"
				std.Category = "Error Handling"
				std.Severity = "error"
				std.Rationale = "Catching specific exceptions prevents silent failures and aids debugging."
				std.Patterns = []string{"except ValueError as e:"}
				std.Rules = []string{
					"Never catch a bare Exception - always catch specific exceptions (e.g., ValueError, KeyError).",
					"Always use logger.exception() when handling exceptions to preserve the stack trace.",
				}
			} else if strings.Contains(contentLower, "class ") {
				std.Title = "Python OOP and Class Standard"
				std.Category = "Structure"
				std.Severity = "info"
				std.Rationale = "Consistent object-oriented layouts align with PEP-8 guidelines."
				std.Patterns = []string{"class MyService:"}
				std.Rules = []string{
					"Always define class methods with explicit 'self' or 'cls' parameters.",
					"Use type hints for class properties and method signatures.",
				}
			} else {
				std.Title = "Python Idiomatic Function Design"
				std.Category = "General Code Structure"
				std.Severity = "warning"
				std.Rationale = "Clean, PEP-8 compliant functions reduce technical debt and simplify unit testing."
				std.Patterns = []string{"def my_func(a: int) -> str:"}
				std.Rules = []string{
					"Use explicit type hints for all function parameters and return values.",
					"Include a descriptive docstring explaining arguments, returns, and raises.",
				}
			}

		case "rust":
			if strings.Contains(contentLower, "result<") || strings.Contains(contentLower, "option<") || strings.Contains(contentLower, "unwrap") {
				std.Title = "Rust Safe Error Propagation"
				std.Category = "Error Handling"
				std.Severity = "error"
				std.Rationale = "Explicit Result handling and propagation guarantees memory and runtime safety without panics."
				std.Patterns = []string{"Result<T, E>", "unwrap_or_else"}
				std.Rules = []string{
					"Avoid calling '.unwrap()' or '.expect()' in production code - prefer explicit error propagation with the '?' operator.",
					"Define context-aware custom error enums using 'thiserror' or 'anyhow'.",
				}
			} else if strings.Contains(contentLower, "struct ") || strings.Contains(contentLower, "impl ") {
				std.Title = "Rust Data Modeling and Impls"
				std.Category = "Structure"
				std.Severity = "info"
				std.Rationale = "Separating structural representation from implementation blocks is Rust's primary layout."
				std.Patterns = []string{"struct UserData", "impl UserData"}
				std.Rules = []string{
					"Derive standard traits (Debug, Clone, Serialize) on structs where helpful.",
					"Encapsulate fields with proper crate-level visibility (pub, pub(crate)).",
				}
			} else {
				std.Title = "Rust Function Integrity"
				std.Category = "General Code Structure"
				std.Severity = "warning"
				std.Rationale = "Strict ownership rules require functions to carefully design borrowing and references."
				std.Patterns = []string{"fn process(&self)"}
				std.Rules = []string{
					"Prefer taking references (&str, &[]) over owned values (String, Vec) in function signatures.",
					"Keep functions strictly focused with single responsibility.",
				}
			}

		case "java":
			if strings.Contains(contentLower, "try {") || strings.Contains(contentLower, "catch (") || strings.Contains(contentLower, "throw new") {
				std.Title = "Java Checked and Unchecked Exceptions"
				std.Category = "Error Handling"
				std.Severity = "error"
				std.Rationale = "Explicit checked exceptions make interface error handling contracts clear and robust."
				std.Patterns = []string{"try { ... } catch (CustomException e)"}
				std.Rules = []string{
					"Never catch raw Throwable or Exception unless performing top-level framework recovery.",
					"Close system resources in a try-with-resources statement or finally block.",
				}
			} else if strings.Contains(contentLower, "class ") || strings.Contains(contentLower, "interface ") {
				std.Title = "Java Class and Interface Layout"
				std.Category = "Structure"
				std.Severity = "info"
				std.Rationale = "Strict class layouts align with standard enterprise style guidelines."
				std.Patterns = []string{"public class UserService"}
				std.Rules = []string{
					"Group members logically: fields first, then constructors, then methods.",
					"Use camelCase for method names and PascalCase for class names.",
				}
			} else {
				std.Title = "Java Clean Method Guidelines"
				std.Category = "General Code Structure"
				std.Severity = "warning"
				std.Rationale = "Well-structured Java methods improve clarity, maintainability, and code coverage."
				std.Patterns = []string{"public void process()"}
				std.Rules = []string{
					"Annotate methods with @Override, @Nullable or @NonNull to express constraints clearly.",
					"Limit methods to a single level of abstraction where possible.",
				}
			}

		default:
			std.Title = fmt.Sprintf("Pattern Detected in %s", s.FilePath)
			std.Category = "General"
			std.Severity = "warning"
			std.Rationale = "Automatically extracted from codebase structure."
			std.Patterns = []string{"general pattern"}
			std.Rules = []string{"Follow the structural pattern demonstrated in the sample."}
		}

		dna.Standards = append(dna.Standards, std)
	}

	return dna
}
