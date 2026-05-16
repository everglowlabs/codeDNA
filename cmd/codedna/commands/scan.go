package commands

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/everglowlabs/codedna/internal/compiler"
	"github.com/everglowlabs/codedna/internal/schema"
	"github.com/everglowlabs/codedna/internal/tui"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan a directory for patterns",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		p := tea.NewProgram(tui.InitialModel(path))
		m, err := p.Run()
		if err != nil {
			fmt.Printf("Error running TUI: %v", err)
			os.Exit(1)
		}

		// After TUI finishes, if we have a result, export it
		if finalModel, ok := m.(interface{ GetResult() schema.DNA_Schema }); ok {
			dna := finalModel.GetResult()
			if dna.ProjectName != "" {
				formats, _ := cmd.Flags().GetStringSlice("format")
				exportDNA(dna, formats)
			}
		}
	},
}

func exportDNA(dna schema.DNA_Schema, formats []string) {
	var compilers []compiler.Compiler
	for _, f := range formats {
		switch f {
		case "json":
			compilers = append(compilers, &compiler.JSONCompiler{})
		case "markdown":
			compilers = append(compilers, &compiler.MarkdownCompiler{})
		case "cursor":
			compilers = append(compilers, &compiler.CursorCompiler{})
		case "copilot":
			compilers = append(compilers, &compiler.CopilotCompiler{})
		case "antigravity":
			compilers = append(compilers, &compiler.AntigravityCompiler{})
		case "claude":
			compilers = append(compilers, &compiler.ClaudeCompiler{})
		case "continue":
			compilers = append(compilers, &compiler.ContinueCompiler{})
		case "windsurf":
			compilers = append(compilers, &compiler.WindsurfCompiler{})
		case "jetbrains":
			compilers = append(compilers, &compiler.JetBrainsCompiler{})
		}
	}

	// Always export JSON if not specified? Or just follow flags.
	// Default to JSON if nothing else.
	if len(compilers) == 0 {
		compilers = append(compilers, &compiler.JSONCompiler{})
	}

	for _, c := range compilers {
		files, err := c.Compile(dna)
		if err != nil {
			fmt.Printf("Error compiling: %v\n", err)
			continue
		}

		for filename, content := range files {
			// Ensure directory exists
			dir := filepath.Dir(filename)
			if dir != "." {
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					fmt.Printf("Error creating directory %s: %v\n", dir, err)
					continue
				}
			}

			err = os.WriteFile(filename, []byte(content), 0644)
			if err != nil {
				fmt.Printf("Error writing %s: %v\n", filename, err)
			} else {
				fmt.Printf("✅ Exported standards to %s\n", filename)
			}
		}
	}
}

func init() {
	scanCmd.Flags().StringSliceP("format", "f", []string{"json"}, "Output formats (json, markdown, cursor, copilot, antigravity, claude, continue, windsurf, jetbrains)")
	rootCmd.AddCommand(scanCmd)
}
