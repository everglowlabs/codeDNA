package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

		failOn, _ := cmd.Flags().GetString("fail-on")
		failOn = strings.ToLower(failOn)
		if failOn != "error" && failOn != "warning" && failOn != "info" {
			fmt.Printf("Error: invalid value for --fail-on flag: %s. Allowed values: error, warning, info.\n", failOn)
			os.Exit(1)
		}

		severityWeights := map[string]int{
			"info":    1,
			"warning": 2,
			"error":   3,
		}

		thresholdWeight := severityWeights[failOn]

		// Build a map of standard ID to severity
		stdSeverities := make(map[string]string)
		for _, std := range dna.Standards {
			sev := strings.ToLower(std.Severity)
			if sev == "" {
				sev = "warning"
			}
			stdSeverities[std.ID] = sev
		}

		fmt.Printf("🔍 Linting project against DNA: %s (v%s) [fail-on: %s]\n", dna.ProjectName, dna.Version, failOn)

		sb := sandbox.NewSandbox()
		s := scanner.NewScanner(path)
		files, err := s.Scan(nil)
		if err != nil {
			fmt.Printf("Error scanning files: %v\n", err)
			os.Exit(1)
		}

		blockingViolations := 0
		totalViolations := 0

		for _, f := range files {
			res, err := sb.LintFile(f, dna.Standards)
			if err != nil {
				continue
			}
			for _, v := range res.Violations {
				totalViolations++
				sev := stdSeverities[v.StandardID]
				if sev == "" {
					sev = "warning"
				}
				vWeight := severityWeights[sev]
				if vWeight >= thresholdWeight {
					fmt.Printf("  [FAIL] %s:%s (Severity: %s)\n", v.FilePath, v.Message, sev)
					blockingViolations++
				} else {
					fmt.Printf("  [SUGGESTION] %s:%s (Severity: %s)\n", v.FilePath, v.Message, sev)
				}
			}
		}

		if blockingViolations > 0 {
			fmt.Printf("\n❌ Found %d exit-blocking violations (out of %d total violations).\n", blockingViolations, totalViolations)
			os.Exit(1)
		} else if totalViolations > 0 {
			fmt.Printf("\n✅ All critical standards passed (%d suggestions ignored based on threshold '%s').\n", totalViolations, failOn)
		} else {
			fmt.Println("\n✅ All files adhere to DNA standards.")
		}
	},
}

func init() {
	lintCmd.Flags().StringP("config", "c", "codedna.json", "Path to codedna.json")
	lintCmd.Flags().String("fail-on", "warning", "Severity threshold to fail the build (error, warning, info)")
	rootCmd.AddCommand(lintCmd)
}
