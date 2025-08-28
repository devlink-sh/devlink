package git

import (
	"fmt"
	"net"
	"os/exec"
	"time"
)

func execLookPath(name string) (string, error) {
	return exec.LookPath(name)
}

func findFreePort() (int, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer ln.Close()
	addr := ln.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

func generateName() string {
	return fmt.Sprintf("devlink-%d", time.Now().Unix())
}
