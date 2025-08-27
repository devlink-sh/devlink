package env

import (
	"fmt"
	"time"

	"github.com/devlink/internal/env"
	"github.com/spf13/cobra"
)

var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "📦 Bulk environment file operations",
	Long: `📦 Bulk operations for DevLink

Share multiple environment files at once with advanced options.

Examples:
  devlink env bulk share file1.env file2.env     # Share multiple files
  devlink env bulk share *.env --prefix myproject # Share with prefix
  devlink env bulk share .env* --expiry 24h      # Share with custom expiry`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var bulkShareCmd = &cobra.Command{
	Use:   "share [files...]",
	Short: "📦 Share multiple environment files",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runBulkShare,
}

func runBulkShare(cmd *cobra.Command, args []string) error {
	files := args
	expiry, _ := cmd.Flags().GetString("expiry")
	readonly, _ := cmd.Flags().GetBool("readonly")
	prefix, _ := cmd.Flags().GetString("prefix")
	groupBy, _ := cmd.Flags().GetString("group-by")

	expiryDuration, err := time.ParseDuration(expiry)
	if err != nil {
		return fmt.Errorf("❌ invalid expiry format: %w", err)
	}

	request := env.BulkShareRequest{
		Files:    files,
		Expiry:   expiryDuration,
		ReadOnly: readonly,
		Prefix:   prefix,
		GroupBy:  groupBy,
	}

	bulkManager := env.NewBulkManager()

	if err := bulkManager.ValidateBulkRequest(request); err != nil {
		return fmt.Errorf("❌ %w", err)
	}

	fmt.Printf("📦 Bulk sharing %d files\n", len(files))
	fmt.Printf("⏰ Expiry: %s\n", expiry)
	if readonly {
		fmt.Println("🔒 Read-only: enabled")
	}
	if prefix != "" {
		fmt.Printf("🏷️  Prefix: %s\n", prefix)
	}
	fmt.Println()

	results := bulkManager.ShareFiles(request)

	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
			fmt.Printf("✅ %s → %s\n", result.File, result.ShareCode)
		} else {
			fmt.Printf("❌ %s → %s\n", result.File, result.Error)
		}
	}

	fmt.Println()
	stats := bulkManager.GetBulkStatistics(results)
	fmt.Printf("📊 Summary: %d/%d successful (%d%%)\n", 
		stats["successful_shares"].(int), 
		stats["total_files"].(int),
		int(stats["success_rate"].(float64)*100))

	return nil
}

func init() {
	bulkShareCmd.Flags().StringP("expiry", "e", "1h", "Share expiry time")
	bulkShareCmd.Flags().BoolP("readonly", "r", false, "Make shares read-only")
	bulkShareCmd.Flags().StringP("prefix", "p", "", "Prefix for share codes")
	bulkShareCmd.Flags().StringP("group-by", "g", "", "Group results by category")

	bulkCmd.AddCommand(bulkShareCmd)
}
