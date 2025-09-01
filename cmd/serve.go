package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gupta-nu/devlink/internal/gitutils"
	"github.com/gupta-nu/devlink/internal/signal"
	"github.com/gupta-nu/devlink/internal/zrok"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run inside a git repo: create a bare mirror, run git daemon, expose via zrok",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) detect repo
		repoPath, err := gitutils.DetectRepo()
		if err != nil {
			return err
		}
		fmt.Println("Serving repo from:", repoPath)

		// 2) create bare mirror
		tmpDir := filepath.Join(os.TempDir(), "devlink-git-shares")
		barePath, err := gitutils.CreateBareMirror(repoPath, tmpDir)
		if err != nil {
			return err
		}
		fmt.Println("Created bare mirror at:", barePath)

		// 3) start git daemon bound to localhost:9418 serving the base-path parent of barePath
		gitCmd, err := gitutils.StartGitDaemon(filepath.Dir(barePath), 9418)
		if err != nil {
			_ = os.RemoveAll(barePath)
			return err
		}
		fmt.Println("git daemon started on localhost:9418")

		// 4) start zrok to share :9418
		zcmd, url, err := zrok.StartTunnel(":9418")
		if err != nil {
			// best-effort cleanup
			if gitCmd != nil && gitCmd.Process != nil {
				_ = gitCmd.Process.Kill()
			}
			_ = os.RemoveAll(barePath)
			return err
		}

		if url != "" {
			fmt.Println("Tunnel available at:", url)
		} else {
			fmt.Println("Tunnel started but URL not yet detected. Check `zrok` output for details.")
		}

		// Print friendly clone/push instructions (guessing repo name)
		repoName := filepath.Base(barePath) // e.g., myrepo.git
		fmt.Println("\n=== Share instructions for your teammate ===")
		// If url starts with tcp:// -> git:// clone; if https:// -> http(s) clone
		fmt.Printf("Clone (raw): git clone %s/%s\n", url, repoName)
		fmt.Printf("Or (explicit): git clone <protocol>://<host>:<port>/%s\n", repoName)
		fmt.Println("When done, press Ctrl+C here to stop sharing and clean up.\n")

		// 5) wait for interrupt, then cleanup
		signal.WaitForInterrupt()
		fmt.Println("\nReceived interrupt — shutting down...")

		// kill zrok first
		if zcmd != nil && zcmd.Process != nil {
			_ = zcmd.Process.Kill()
			time.Sleep(150 * time.Millisecond)
		}
		// kill git daemon
		if gitCmd != nil && gitCmd.Process != nil {
			_ = gitCmd.Process.Kill()
			time.Sleep(150 * time.Millisecond)
		}

		// remove bare mirror
		if err := os.RemoveAll(barePath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to remove bare mirror: %v\n", err)
		}

		fmt.Println("Stopped.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
