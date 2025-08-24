package discordrichpresence

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// Constants for Discord IPC opcodes
const (
	opHandshake = 0
	opFrame     = 1
	opClose     = 2
	opPing      = 3
	opPong      = 4
)

// Client represents a Discord Rich Presence client.
type Client struct {
	conn         net.Conn
	appID        string
	ready        bool
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	updateTicker *time.Ticker
}

// NewClient creates a new Discord RPC client with the specified application ID.
func NewClient(appID string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		appID:  appID,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Connect establishes a connection to Discord's IPC socket.
func (c *Client) Connect() error {
	conn, err := c.findDiscordSocket()
	if err != nil {
		return fmt.Errorf("failed to connect to Discord: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	if err := c.handshake(); err != nil {
		c.Close()
		return fmt.Errorf("handshake failed: %w", err)
	}

	//log.Println("Discord rich presence RPC connected successfully")
	return nil
}

// StartWithActivity starts the RPC client with the given activity, updating at the specified interval.
// It runs in a goroutine and does not block the caller.
func (c *Client) StartWithActivity(activity Activity, updateInterval time.Duration) error {
	if err := c.Connect(); err != nil {
		return err
	}

	// Set initial activity
	if err := c.SetActivity(activity); err != nil {
		return fmt.Errorf("failed to set initial activity: %w", err)
	}

	// Start update routine in a goroutine
	go c.updateRoutine(activity, updateInterval)

	return nil
}

// updateRoutine runs in a goroutine to keep the activity updated.
func (c *Client) updateRoutine(activity Activity, interval time.Duration) {
	c.updateTicker = time.NewTicker(interval)
	defer c.updateTicker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.updateTicker.C:
			if err := c.SetActivity(activity); err != nil {
				//log.Printf("Failed to update activity: %v", err)
			}
		}
	}
}

// findDiscordSocket tries to connect to Discord's IPC socket.
func (c *Client) findDiscordSocket() (net.Conn, error) {
	pipePaths := []string{
		"discord-ipc-0", "discord-ipc-1", "discord-ipc-2", "discord-ipc-3", "discord-ipc-4",
		"discord-ipc-5", "discord-ipc-6", "discord-ipc-7", "discord-ipc-8", "discord-ipc-9",
	}

	for _, pipeName := range pipePaths {
		conn, err := c.connectToPipe(pipeName)
		if err == nil {
			//log.Printf("Connected to Discord via %s", pipeName)
			return conn, nil
		}
	}

	return nil, fmt.Errorf("could not connect to any Discord RPC socket")
}

// handshake performs the initial handshake with Discord.
func (c *Client) handshake() error {
	handshakeData := map[string]interface{}{
		"v":         1,
		"client_id": c.appID,
	}

	if err := c.sendMessage(opHandshake, handshakeData); err != nil {
		return err
	}

	// Wait for READY event
	for i := 0; i < 10; i++ {
		msg, err := c.readMessage()
		if err != nil {
			return err
		}

		if evt, ok := msg["evt"].(string); ok && evt == "READY" {
			c.mu.Lock()
			c.ready = true
			c.mu.Unlock()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for READY event")
}

// sendMessage sends a message to Discord RPC.
func (c *Client) sendMessage(op int, data interface{}) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("not connected")
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	header := make([]byte, 8)
	binary.LittleEndian.PutUint32(header[0:4], uint32(op))
	binary.LittleEndian.PutUint32(header[4:8], uint32(len(payload)))

	if _, err := conn.Write(header); err != nil {
		return err
	}
	if _, err := conn.Write(payload); err != nil {
		return err
	}

	return nil
}

// readMessage reads a message from Discord RPC.
func (c *Client) readMessage() (map[string]interface{}, error) {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	header := make([]byte, 8)
	if _, err := conn.Read(header); err != nil {
		return nil, err
	}

	length := binary.LittleEndian.Uint32(header[4:8])
	payload := make([]byte, length)
	if _, err := conn.Read(payload); err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(payload, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// SetActivity sets the Discord Rich Presence activity.
func (c *Client) SetActivity(activity Activity) error {
	c.mu.RLock()
	ready := c.ready
	c.mu.RUnlock()

	if !ready {
		return fmt.Errorf("client not ready")
	}

	args := map[string]interface{}{
		"pid":      os.Getpid(),
		"activity": activity,
	}

	msg := RPCMessage{
		Cmd:   "SET_ACTIVITY",
		Args:  args,
		Nonce: generateNonce(),
	}

	return c.sendMessage(opFrame, msg)
}

// ClearActivity clears the current Discord Rich Presence activity.
func (c *Client) ClearActivity() error {
	c.mu.RLock()
	ready := c.ready
	c.mu.RUnlock()

	if !ready {
		return fmt.Errorf("client not ready")
	}

	args := map[string]interface{}{
		"pid": os.Getpid(),
	}

	msg := RPCMessage{
		Cmd:   "SET_ACTIVITY",
		Args:  args,
		Nonce: generateNonce(),
	}

	return c.sendMessage(opFrame, msg)
}

// Close closes the Discord RPC connection and stops all routines.
func (c *Client) Close() error {
	c.cancel() // Stop all goroutines

	if c.updateTicker != nil {
		c.updateTicker.Stop()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		c.ready = false
		return err
	}
	return nil
}

// IsReady returns whether the client is ready to send commands.
func (c *Client) IsReady() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ready
}
