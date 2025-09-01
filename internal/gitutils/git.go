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
	base := filepath.Base(url)
	if idx := strings.IndexAny(base, "?#"); idx != -1 {
		base = base[:idx]
	}
	base = strings.TrimSuffix(base, ".git")
	dir := base
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return RunStreaming("git", "clone", url)
	}
	return RunStreaming("git", "-C", dir, "fetch", "--all")
}

// CreateBareMirror creates a mirror/bare repo from src path into outDir and returns dest path
func CreateBareMirror(srcRepoPath string, outDir string) (string, error) {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", err
	}
	name := filepath.Base(srcRepoPath) + ".git"
	dest := filepath.Join(outDir, name)
	if _, err := os.Stat(dest); err == nil {
		if err := os.RemoveAll(dest); err != nil {
			return "", fmt.Errorf("failed to remove existing mirror %s: %w", dest, err)
		}
	}
	cmd := exec.Command("git", "clone", "--mirror", srcRepoPath, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create bare mirror: %w", err)
	}
	return dest, nil
}

// StartGitDaemon starts git daemon serving basePath on given port (9418) and returns the cmd
func StartGitDaemon(basePath string, port int) (*exec.Cmd, error) {
	// git daemon --reuseaddr --base-path=<basePath> --export-all --enable=receive-pack --port=<port>
	portStr := fmt.Sprintf("%d", port)
	cmd := exec.Command("git", "daemon",
		"--reuseaddr",
		"--base-path="+basePath,
		"--export-all",
		"--enable=receive-pack",
		"--verbose",
		"--port="+portStr,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Start without waiting; caller will kill it later
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start git daemon: %w", err)
	}
	return cmd, nil
}
