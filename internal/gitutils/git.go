package gitutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DetectRepo finds the current git repo root
func DetectRepo() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not inside a git repository")
		}
		dir = parent
	}
}

// RunStreaming runs a command and streams output
func RunStreaming(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// RunSilently runs a command without streaming
func RunSilently(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// CloneOrFetch clones a repo if missing, else fetches updates
func CloneOrFetch(url string) error {
	dir := filepath.Base(url)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return RunStreaming("git", "clone", url)
	}
	return RunStreaming("git", "-C", dir, "fetch", "--all")
}

// ServeRepo exposes the current repo over a zrok tunnel
func ServeRepo() error {
	repo, err := DetectRepo()
	if err != nil {
		return err
	}

	fmt.Println("Detected repo at:", repo)

	// Run zrok share and capture output
	output, err := RunSilently("zrok", "share", "public", repo)
	if err != nil {
		return fmt.Errorf("failed to start zrok tunnel: %v\n%s", err, string(output))
	}

	// Extract and print the tunnel URL
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "https://") {
			fmt.Println("Repo available at:", line)
			return nil
		}
	}

	// If no URL found, just print raw output
	fmt.Println(string(output))
	return nil
}
