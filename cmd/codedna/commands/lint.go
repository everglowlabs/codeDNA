package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/everglowlabs/codedna/internal/sandbox"
	"github.com/everglowlabs/codedna/internal/scanner"
	"github.com/everglowlabs/codedna/internal/schema"
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint [path]",
	Short: "Check if the codebase adheres to the generated CodeDNA",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		dnaPath, _ := cmd.Flags().GetString("config")
		if _, err := os.Stat(dnaPath); os.IsNotExist(err) {
			fmt.Printf("Error: DNA config file not found at %s. Run 'scan' first.\n", dnaPath)
			os.Exit(1)
		}

		// Load DNA
		dnaData, err := os.ReadFile(dnaPath)
		if err != nil {
			fmt.Printf("Error reading DNA: %v\n", err)
			os.Exit(1)
		}

		var dna schema.DNA_Schema
		if err := json.Unmarshal(dnaData, &dna); err != nil {
			fmt.Printf("Error parsing DNA: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🔍 Linting project against DNA: %s (v%s)\n", dna.ProjectName, dna.Version)

		sb := sandbox.NewSandbox()
		s := scanner.NewScanner(path)
		files, err := s.Scan(nil)
		if err != nil {
			fmt.Printf("Error scanning files: %v\n", err)
			os.Exit(1)
		}

		violations := 0
		for _, f := range files {
			res, err := sb.LintFile(f, dna.Standards)
			if err != nil {
				continue
			}
			if !res.Passed {
				for _, issue := range res.Issues {
					fmt.Printf("  [FAIL] %s\n", issue)
					violations++
				}
			}
		}

		if violations > 0 {
			fmt.Printf("\n❌ Found %d violations.\n", violations)
			os.Exit(1)
		} else {
			fmt.Println("\n✅ All files adhere to DNA standards.")
		}
	},
}

func init() {
	lintCmd.Flags().StringP("config", "c", "codedna.json", "Path to codedna.json")
	rootCmd.AddCommand(lintCmd)
}
