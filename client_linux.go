//go:build !windows
// +build !windows

package discordrichpresence

import (
	"fmt"
	"net"
	"os"
	"time"
)

// connectToPipe connects to the appropriate IPC mechanism based on the platform.
func (c *Client) connectToPipe(pipeName string) (net.Conn, error) {
	// Unix-like systems use Unix domain sockets
	address := c.getSocketAddress(pipeName)
	return net.DialTimeout("unix", address, 2*time.Second)
}

// getSocketAddress returns the appropriate socket address for Unix-like systems.
func (c *Client) getSocketAddress(pipeName string) string {
	tmpDir := os.Getenv("XDG_RUNTIME_DIR")
	if tmpDir == "" {
		tmpDir = os.Getenv("TMPDIR")
	}
	if tmpDir == "" {
		tmpDir = "/tmp"
	}
	return fmt.Sprintf("%s/%s", tmpDir, pipeName)
}
