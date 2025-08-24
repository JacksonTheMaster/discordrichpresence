//go:build windows
// +build windows

package discordrichpresence

import (
	"fmt"
	"net"
	"syscall"
	"time"
)

// Windows-specific constants for named pipe access
const (
	GENERIC_READ          = 0x80000000
	GENERIC_WRITE         = 0x40000000
	OPEN_EXISTING         = 3
	FILE_ATTRIBUTE_NORMAL = 0x80
)

// WindowsPipeConn wraps a Windows named pipe handle to implement net.Conn
type WindowsPipeConn struct {
	handle syscall.Handle
}

// connectToPipe connects to the appropriate IPC mechanism based on the platform.
func (c *Client) connectToPipe(pipeName string) (net.Conn, error) {
	return c.connectWindowsNamedPipe(pipeName)
}

// connectWindowsNamedPipe connects to a Windows named pipe.
func (c *Client) connectWindowsNamedPipe(pipeName string) (net.Conn, error) {
	pipePath := `\\.\pipe\` + pipeName

	// Convert string to UTF-16 for Windows API
	pipePathPtr, err := syscall.UTF16PtrFromString(pipePath)
	if err != nil {
		return nil, fmt.Errorf("failed to convert pipe path: %w", err)
	}

	// Try to open the named pipe
	handle, err := syscall.CreateFile(
		pipePathPtr,
		GENERIC_READ|GENERIC_WRITE,
		0,   // dwShareMode
		nil, // lpSecurityAttributes
		OPEN_EXISTING,
		FILE_ATTRIBUTE_NORMAL,
		0, // hTemplateFile
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open named pipe %s: %w", pipePath, err)
	}

	return &WindowsPipeConn{handle: handle}, nil
}

// WindowsPipeConn implementation of net.Conn interface
func (w *WindowsPipeConn) Read(b []byte) (n int, err error) {
	var bytesRead uint32
	err = syscall.ReadFile(w.handle, b, &bytesRead, nil)
	return int(bytesRead), err
}

func (w *WindowsPipeConn) Write(b []byte) (n int, err error) {
	var bytesWritten uint32
	err = syscall.WriteFile(w.handle, b, &bytesWritten, nil)
	return int(bytesWritten), err
}

func (w *WindowsPipeConn) Close() error {
	return syscall.CloseHandle(w.handle)
}

func (w *WindowsPipeConn) LocalAddr() net.Addr {
	return &net.UnixAddr{Name: "discord-pipe", Net: "pipe"}
}

func (w *WindowsPipeConn) RemoteAddr() net.Addr {
	return &net.UnixAddr{Name: "discord-pipe", Net: "pipe"}
}

func (w *WindowsPipeConn) SetDeadline(t time.Time) error {
	// Named pipes don't support deadlines in the same way
	return nil
}

func (w *WindowsPipeConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (w *WindowsPipeConn) SetWriteDeadline(t time.Time) error {
	return nil
}
