package zrok

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"time"
)

// StartTunnel starts a zrok share for the given target (e.g. ":9418" or "http://localhost:8080")
// Returns running *exec.Cmd and the parsed public endpoint (tcp://...:port or https://...), or empty string if not found yet.
func StartTunnel(target string) (*exec.Cmd, string, error) {
	// zrok CLI expects: zrok share public <target> --headless
	cmd := exec.Command("zrok", "share", "public", target, "--headless")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, "", fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, "", fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, "", fmt.Errorf("failed to start zrok: %w", err)
	}

	// capture either tcp:// or https:// endpoints
	urlCh := make(chan string, 1)
	re := regexp.MustCompile(`([a-z]+://[^\s]+)`)

	go func() {
		sc := bufio.NewScanner(stdout)
		for sc.Scan() {
			line := sc.Text()
			if m := re.FindString(line); m != "" {
				select {
				case urlCh <- m:
				default:
				}
			}
		}
	}()

	go func() {
		sc := bufio.NewScanner(stderr)
		for sc.Scan() {
			line := sc.Text()
			if m := re.FindString(line); m != "" {
				select {
				case urlCh <- m:
				default:
				}
			}
		}
	}()

	// Wait up to 10s for a URL to appear
	select {
	case url := <-urlCh:
		return cmd, url, nil
	case <-time.After(10 * time.Second):
		return cmd, "", nil
	}
}
