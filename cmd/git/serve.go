package git

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/devlink-sh/devlink/internal"
	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

func getFreePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func canonicalRepoName(repoPath string) (string, string, error) {
	abs, err := filepath.Abs(repoPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to resolve repo path: %w", err)
	}

	// If user gave ".", normalize to repo root using `git rev-parse --show-toplevel`
	cmdCheck := exec.Command("git", "rev-parse", "--show-toplevel")
	cmdCheck.Dir = abs
	if out, err := cmdCheck.Output(); err == nil {
		abs = strings.TrimSpace(string(out))
	}

	// If it's a working repo, redirect to .git folder
	gitDir := filepath.Join(abs, ".git")
	if stat, err := os.Stat(gitDir); err == nil && stat.IsDir() {
		abs = gitDir
	}

	// Validate repo (must contain HEAD)
	if _, err := os.Stat(filepath.Join(abs, "HEAD")); err != nil {
		return "", "", fmt.Errorf("error: %s is not a valid git repository", repoPath)
	}

	// Derive canonical repo name <basename>.git and parent dir used as base-path for git-daemon
	repoRoot := filepath.Dir(abs) // abs points to .../.git
	repoName := filepath.Base(repoRoot) + ".git"
	parentDir := filepath.Dir(repoRoot)
	return repoName, parentDir, nil
}

var gitServeCmd = &cobra.Command{
	Use:   "serve <repo-path>",
	Short: "Share a local Git repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := args[0]

		repoName, parentDir, err := canonicalRepoName(repoPath)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Ensure repo exportability (create temporarily; delete on shutdown if we created it)
		repoGitDir := filepath.Join(parentDir, strings.TrimSuffix(repoName, ".git"), ".git")
		exportOk := filepath.Join(repoGitDir, "git-daemon-export-ok")
		exportCreated := false
		if _, err := os.Stat(exportOk); os.IsNotExist(err) {
			if f, err := os.Create(exportOk); err == nil {
				_ = f.Close()
				exportCreated = true
			} else {
				log.Fatalf("failed to create git-daemon-export-ok: %v", err)
			}
		}

		// Pick a free local port for git-daemon
		gitPort, err := getFreePort()
		if err != nil {
			log.Fatalf("failed to allocate port for git daemon: %v", err)
		}

		// Start git daemon bound to loopback and chosen port
		gitDaemon := exec.Command("git", "daemon",
			"--reuseaddr",
			"--listen=127.0.0.1",
			fmt.Sprintf("--port=%d", gitPort),
			"--base-path="+parentDir,
			"--export-all",
			"--verbose",
			"--enable=upload-pack",  // allow clone/pull
			"--enable=receive-pack", // allow push
		)
		gitDaemon.Stdout = os.Stdout
		gitDaemon.Stderr = os.Stderr

		if err := gitDaemon.Start(); err != nil {
			log.Fatalf("failed to start git daemon: %v", err)
		}
		log.Printf("git daemon started for %s (127.0.0.1:%d)", repoName, gitPort)

		// Setup zrok share (directional: this side is the server)
		root, err := environment.LoadRoot()
		if err != nil {
			_ = gitDaemon.Process.Kill()
			log.Fatal(err)
		}

		share, err := sdk.CreateShare(root, &sdk.ShareRequest{
			BackendMode: sdk.TcpTunnelBackendMode,
			ShareMode:   sdk.PrivateShareMode,
			Target:      "git",
		})
		if err != nil {
			_ = gitDaemon.Process.Kill()
			log.Fatal(err)
		}
		log.Printf("Git share ready!")
		log.Printf("Share this command with your teammate:\n\n  devlink git connect %s %s\n", share.Token, repoName)

		// Accept zrok connections and forward to local git-daemon
		listener, err := sdk.NewListener(share.Token, root)
		if err != nil {
			_ = sdk.DeleteShare(root, share)
			_ = gitDaemon.Process.Kill()
			log.Fatal(err)
		}

		// Graceful shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Shutting down git serve...")
			if exportCreated {
				_ = os.Remove(exportOk)
			}
			_ = listener.Close()
			_ = sdk.DeleteShare(root, share)
			_ = gitDaemon.Process.Kill()
			os.Exit(0)
		}()

		// Accept loop with backoff on temporary errors
		var tempDelay time.Duration
		for {
			remote, err := listener.Accept()
			if err != nil {
				if ne, ok := err.(interface{ Temporary() bool }); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 50 * time.Millisecond
					} else {
						tempDelay *= 2
						if tempDelay > time.Second {
							tempDelay = time.Second
						}
					}
					log.Printf("temporary accept error: %v; retrying in %v", err, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				// Permanent error
				log.Printf("fatal accept error: %v", err)
				break
			}
			tempDelay = 0

			go func(remoteConn net.Conn) {
				local, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", gitPort))
				if err != nil {
					log.Printf("error dialing local git daemon: %v", err)
					_ = remoteConn.Close()
					return
				}
				log.Println("Forwarding git traffic (remote<->local)...")
				internal.Pipe(remoteConn, local)
			}(remote)
		}

		// Cleanup on loop exit
		if exportCreated {
			_ = os.Remove(exportOk)
		}
		_ = listener.Close()
		_ = sdk.DeleteShare(root, share)
		_ = gitDaemon.Process.Kill()
	},
}
