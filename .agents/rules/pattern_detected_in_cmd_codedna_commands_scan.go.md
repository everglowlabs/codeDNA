# Pattern Detected in cmd/codedna/commands/scan.go

Category: General
Rationale: Automatically extracted from codebase structure.

## Guidelines
- Follow the structural pattern demonstrated in the sample.

## Examples
```go
// File: cmd/codedna/commands/scan.go
func init() {
	scanCmd.Flags().StringSliceP("format", "f", []string{"json"}, "Output formats (json, markdown, cursor, copilot, antigravity, claude, continue, windsurf, jetbrains)")
	rootCmd.AddCommand(scanCmd)
}
```

