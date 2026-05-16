package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codedna",
	Short: "CodeDNA extracts engineering patterns from your codebase",
	Long: `CodeDNA is a forensic tool that analyzes your codebase using Tree-sitter
to identify implicit engineering standards and generate AI agent rules.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Root flags if any
}
