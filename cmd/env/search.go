package env

import (
	"fmt"
	"os"

	"github.com/devlink/internal/env"
	"github.com/devlink/internal/util"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "🔍 Search environment variables",
	Long: `🔍 Search and filter environment variables

Search through environment files for specific variables or patterns.

Examples:
  devlink env search "DATABASE"                    # Search for variables containing "DATABASE"
  devlink env search --sensitive                   # Show only sensitive variables
  devlink env search --category database           # Search by category
  devlink env search --regex "API_.*"              # Use regex pattern`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	files, _ := cmd.Flags().GetStringSlice("files")
	sensitive, _ := cmd.Flags().GetBool("sensitive")
	categories, _ := cmd.Flags().GetStringSlice("categories")
	exactMatch, _ := cmd.Flags().GetBool("exact")
	caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
	useRegex, _ := cmd.Flags().GetBool("regex")
	outputFormat, _ := cmd.Flags().GetString("output")

	if len(files) == 0 {
		files = []string{".env"}
	}

	searchManager := util.NewSearchManager()

	for _, filePath := range files {
		if _, err := os.Stat(filePath); err == nil {
			parser := env.NewParser()
			envFile, err := parser.ParseFile(filePath)
			if err != nil {
				fmt.Printf("⚠️  Warning: Could not parse %s: %v\n", filePath, err)
				continue
			}
			searchManager.AddEnvFile(envFile)
		}
	}

	var sensitivePtr *bool
	if cmd.Flags().Changed("sensitive") {
		sensitivePtr = &sensitive
	}

	filter := util.SearchFilter{
		Query:         query,
		Categories:    categories,
		Sensitive:     sensitivePtr,
		ExactMatch:    exactMatch,
		CaseSensitive: caseSensitive,
		Regex:         useRegex,
	}

	results := searchManager.Search(filter)

	if len(results) == 0 {
		fmt.Println("🔍 No results found")
		return nil
	}

	fmt.Printf("🔍 Found %d results\n", len(results))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, result := range results {
		if i >= 10 {
			fmt.Printf("... and %d more results\n", len(results)-10)
			break
		}

		sensitive := ""
		if result.Variable.IsSensitive {
			sensitive = " 🔒"
		}

		fmt.Printf("%d. %s (Line %d) - %s%s\n",
			i+1, result.Variable.Key, result.LineNumber, result.File, sensitive)

		if outputFormat == "detailed" {
			fmt.Printf("   Value: %s\n", result.Variable.Value)
			fmt.Printf("   Match: %s (Score: %.1f)\n", result.MatchType, result.Score)
		}
	}

	if outputFormat == "stats" {
		stats := searchManager.GetStatistics()
		fmt.Println("\n📊 Search Statistics:")
		fmt.Printf("   Total variables: %d\n", stats["total_variables"])
		fmt.Printf("   Sensitive variables: %d\n", stats["sensitive_variables"])
		fmt.Printf("   Files searched: %d\n", stats["files_count"])
	}

	return nil
}

func runSuggestions(cmd *cobra.Command, args []string) error {
	partial := ""
	if len(args) > 0 {
		partial = args[0]
	}

	files, _ := cmd.Flags().GetStringSlice("files")

	if len(files) == 0 {
		files = []string{".env"}
	}

	searchManager := util.NewSearchManager()

	for _, filePath := range files {
		if _, err := os.Stat(filePath); err == nil {
			parser := env.NewParser()
			envFile, err := parser.ParseFile(filePath)
			if err != nil {
				continue
			}
			searchManager.AddEnvFile(envFile)
		}
	}

	suggestions := searchManager.GetVariableSuggestions(partial)

	if len(suggestions) == 0 {
		fmt.Println("💡 No suggestions found")
		return nil
	}

	fmt.Printf("💡 Suggestions for '%s':\n", partial)
	for _, suggestion := range suggestions {
		fmt.Printf("   %s\n", suggestion)
	}

	return nil
}

var suggestionsCmd = &cobra.Command{
	Use:   "suggest [partial]",
	Short: "💡 Get variable name suggestions",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSuggestions,
}

func init() {
	searchCmd.Flags().StringSliceP("files", "f", nil, "Files to search in")
	searchCmd.Flags().BoolP("sensitive", "s", false, "Show only sensitive variables")
	searchCmd.Flags().StringSliceP("categories", "c", nil, "Filter by categories")
	searchCmd.Flags().BoolP("exact", "e", false, "Exact match only")
	searchCmd.Flags().BoolP("case-sensitive", "C", false, "Case sensitive search")
	searchCmd.Flags().BoolP("regex", "r", false, "Use regex pattern")
	searchCmd.Flags().StringP("output", "o", "simple", "Output format (simple, detailed, stats)")

	suggestionsCmd.Flags().StringSliceP("files", "f", nil, "Files to search in")

	searchCmd.RunE = runSearch
	searchCmd.AddCommand(suggestionsCmd)
}
