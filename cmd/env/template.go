package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/devlink/internal/env"
	"github.com/devlink/internal/util"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "📋 Manage environment templates",
	Long: `📋 Environment templates for DevLink

Create and manage pre-built environment templates for common development scenarios.

Examples:
  devlink env template list                    # List all templates
  devlink env template show nodejs             # Show template details
  devlink env template create nodejs           # Create from template
  devlink env template search backend          # Search templates`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "📋 List all available templates",
	RunE:  runTemplateList,
}

var templateShowCmd = &cobra.Command{
	Use:   "show [template]",
	Short: "📋 Show template details",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateShow,
}

var templateCreateCmd = &cobra.Command{
	Use:   "create [template]",
	Short: "📋 Create environment file from template",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateCreate,
}

var templateSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "🔍 Search templates",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateSearch,
}

func runTemplateList(cmd *cobra.Command, args []string) error {
	templateManager := util.NewTemplateManager()
	templates := templateManager.ListTemplates()

	fmt.Println("📋 Available Templates")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, template := range templates {
		fmt.Printf("📄 %s\n", template.Name)
		fmt.Printf("   Description: %s\n", template.Description)
		fmt.Printf("   Category: %s\n", template.Category)
		fmt.Printf("   Variables: %d\n", len(template.Variables))
		fmt.Printf("   Tags: %s\n", strings.Join(template.Tags, ", "))
		fmt.Println()
	}

	return nil
}

func runTemplateShow(cmd *cobra.Command, args []string) error {
	templateName := args[0]
	templateManager := util.NewTemplateManager()

	template, err := templateManager.GetTemplate(templateName)
	if err != nil {
		return fmt.Errorf("❌ %w", err)
	}

	fmt.Printf("📋 Template: %s\n", template.Name)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Description: %s\n", template.Description)
	fmt.Printf("Category: %s\n", template.Category)
	fmt.Printf("Tags: %s\n", strings.Join(template.Tags, ", "))
	fmt.Printf("Variables: %d\n", len(template.Variables))
	fmt.Println()

	fmt.Println("🔧 Variables:")
	for _, variable := range template.Variables {
		sensitive := ""
		if variable.IsSensitive {
			sensitive = " 🔒"
		}
		fmt.Printf("   %s=%s%s\n", variable.Key, variable.Value, sensitive)
	}

	return nil
}

func runTemplateCreate(cmd *cobra.Command, args []string) error {
	templateName := args[0]
	outputFile, _ := cmd.Flags().GetString("output")

	templateManager := util.NewTemplateManager()
	envFile, err := templateManager.CreateEnvFileFromTemplate(templateName)
	if err != nil {
		return fmt.Errorf("❌ %w", err)
	}

	if outputFile == "" {
		outputFile = fmt.Sprintf("%s.env", templateName)
	}

	formatter := env.NewFormatter()
	formatOptions := &util.FormatOptions{
		MaskSensitive:   false,
		ShowComments:    true,
		ShowLineNumbers: false,
		OutputFormat:    "text",
	}

	formattedContent, err := formatter.Format(envFile, formatOptions)
	if err != nil {
		return fmt.Errorf("❌ failed to format template: %w", err)
	}

	if err := os.WriteFile(outputFile, []byte(formattedContent), 0600); err != nil {
		return fmt.Errorf("❌ failed to write file: %w", err)
	}

	fmt.Printf("✅ Created environment file from template '%s'\n", templateName)
	fmt.Printf("📁 Output: %s\n", outputFile)
	fmt.Printf("📊 Variables: %d\n", len(envFile.Variables))

	return nil
}

func runTemplateSearch(cmd *cobra.Command, args []string) error {
	query := args[0]
	templateManager := util.NewTemplateManager()

	results := templateManager.SearchTemplates(query)

	if len(results) == 0 {
		fmt.Printf("🔍 No templates found for query: %s\n", query)
		return nil
	}

	fmt.Printf("🔍 Search results for: %s\n", query)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, template := range results {
		fmt.Printf("📄 %s\n", template.Name)
		fmt.Printf("   Description: %s\n", template.Description)
		fmt.Printf("   Category: %s\n", template.Category)
		fmt.Printf("   Variables: %d\n", len(template.Variables))
		fmt.Println()
	}

	return nil
}

func init() {
	templateCreateCmd.Flags().StringP("output", "o", "", "Output file path")

	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateCreateCmd)
	templateCmd.AddCommand(templateSearchCmd)
}
