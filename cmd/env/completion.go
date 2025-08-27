package env

import (
	"fmt"

	"github.com/devlink/internal/util"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "⌨️  Generate shell completion scripts",
	Long: `⌨️  Shell completion for DevLink

Generate completion scripts for bash and zsh shells.

Examples:
  devlink env completion bash > ~/.bash_completion
  devlink env completion zsh > ~/.zsh_completion
  source <(devlink env completion bash)`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completion script",
	RunE:  runCompletionBash,
}

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completion script",
	RunE:  runCompletionZsh,
}

var completionSuggestCmd = &cobra.Command{
	Use:   "suggest [type] [partial]",
	Short: "💡 Get completion suggestions",
	Args:  cobra.ExactArgs(2),
	RunE:  runCompletionSuggest,
}

func runCompletionBash(cmd *cobra.Command, args []string) error {
	completionManager := util.NewCompletionManager()
	script := completionManager.GenerateBashCompletion()
	fmt.Print(script)
	return nil
}

func runCompletionZsh(cmd *cobra.Command, args []string) error {
	completionManager := util.NewCompletionManager()
	script := completionManager.GenerateZshCompletion()
	fmt.Print(script)
	return nil
}

func runCompletionSuggest(cmd *cobra.Command, args []string) error {
	completionType := args[0]
	partial := args[1]

	completionManager := util.NewCompletionManager()
	var suggestions []string

	switch completionType {
	case "sharecode":
		suggestions = completionManager.GetShareCodeSuggestions(partial)
	case "file":
		suggestions = completionManager.GetFileSuggestions(partial)
	case "template":
		suggestions = completionManager.GetTemplateSuggestions(partial)
	case "category":
		suggestions = completionManager.GetCategorySuggestions(partial)
	case "command":
		suggestions = completionManager.GetCommandSuggestions(partial)
	default:
		return fmt.Errorf("❌ unknown completion type: %s", completionType)
	}

	if len(suggestions) == 0 {
		fmt.Printf("💡 No suggestions for '%s' (%s)\n", partial, completionType)
		return nil
	}

	fmt.Printf("💡 Suggestions for '%s' (%s):\n", partial, completionType)
	for _, suggestion := range suggestions {
		fmt.Printf("   %s\n", suggestion)
	}

	return nil
}

func init() {
	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionZshCmd)
	completionCmd.AddCommand(completionSuggestCmd)
}
