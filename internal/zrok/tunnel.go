package zrok

import (
	"fmt"
	"os/exec"
)

// StartTunnel starts a zrok tunnel for the given repo path
func StartTunnel(repoPath string) (string, error) {
	cmd := exec.Command("zrok", "share", "public", repoPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to start zrok: %w\n%s", err, out)
	}
	// TODO: parse actual URL from zrok output
	return string(out), nil
}
