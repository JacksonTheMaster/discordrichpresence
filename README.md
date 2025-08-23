# discordrichpresence

A Go module for integrating Discord Rich Presence into your applications. This module provides a simple, thread-safe, and fluent API to set and manage Discord Rich Presence activities, allowing you to display custom status messages, images, and timestamps in Discord.

## Features

- **Fluent API**: Build Discord Rich Presence activities with a chainable, easy-to-use `ActivityBuilder`.
- **Automatic Updates**: Set activities with periodic updates to keep the presence active.
- **Cross-Platform Support**: Works on Windows, Linux, and macOS by connecting to Discord’s IPC socket.
- **Thread-Safe**: Uses mutexes to ensure safe concurrent access to the client.
- **Lightweight**: No external dependencies beyond the Go standard library.
- **Well-Documented**: Includes GoDoc comments for all public functions and types.

## Installation

To use this module in your Go project, run:

```bash
go get github.com/jacksonthemaster/discordrichpresence
```

Ensure you have Go 1.21 or higher installed. You can check your Go version with:

```bash
go version
```

## Requirements

- **Go**: Version 1.21 or higher.
- **Discord**: The Discord desktop client must be running on the same machine.
- **Discord Application ID**: You need a valid Discord application ID from the [Discord Developer Portal](https://discord.com/developers/applications).

## Usage

Below is a complete example of how to use the `discordrichpresence` module to set a Discord Rich Presence status.

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	discordrichpresence "github.com/jacksonthemaster/discordrichpresence"
)

func main() {
	// Initialize a new Discord RPC client with your application ID
	client := discordrichpresence.NewClient("YOUR_APPLICATION_ID")
	defer client.Close()

	// Build an activity using the fluent ActivityBuilder
	activity := discordrichpresence.NewActivity().
		State("Playing a game").
		Details("In the main menu").
		StartTime(time.Now()).
		LargeImage("logo", "Game Logo").
		SmallImage("status", "Online").
		Type(0). // Playing
		Build()

	// Start the client with the activity, updating every 30 seconds
	if err := client.StartWithActivity(activity, 30*time.Second); err != nil {
		log.Fatalf("Failed to start Discord RPC: %v", err)
	}

	log.Println("Discord Rich Presence started successfully!")
	log.Println("Press Ctrl+C to stop...")

	// Wait for an interrupt signal to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down Discord RPC...")
}
```

### Steps to Use

1. **Create a Discord Application**:
   - Go to the [Discord Developer Portal](https://discord.com/developers/applications).
   - Create a new application and note its **Application ID**.
   - Optionally, upload images for `LargeImage` and `SmallImage` in the "Rich Presence" -> "Art Assets" section.

2. **Replace `YOUR_APPLICATION_ID`**:
   - In the example above, replace `"YOUR_APPLICATION_ID"` with your actual Discord application ID.

3. **Run Your Program**:
   - Ensure the Discord desktop client is running.
   - Execute your program with `go run main.go`.
   - Your Discord profile should display the custom Rich Presence.

### Advanced Usage

#### Setting a Custom Activity Type
The `Type` field in the `Activity` struct supports different activity types:
- `0`: Playing
- `1`: Streaming
- `2`: Listening
- `3`: Watching
- `4`: Custom
- `5`: Competing

Example:
```go
activity := discordrichpresence.NewActivity().
	State("Watching a movie").
	Details("In a theater").
	Type(3). // Watching
	Build()
```

#### Adding an End Time
You can set an end time for the activity to show a countdown in Discord:
```go
endTime := time.Now().Add(2 * time.Hour)
activity := discordrichpresence.NewActivity().
	State("Playing a match").
	Details("Competitive Mode").
	StartTime(time.Now()).
	EndTime(endTime).
	Build()
```

#### Clearing the Activity
To remove the current Rich Presence:
```go
if err := client.ClearActivity(); err != nil {
	log.Printf("Failed to clear activity: %v", err)
}
```

#### Checking Client Readiness
To verify if the client is connected and ready:
```go
if client.IsReady() {
	log.Println("Discord RPC client is ready")
} else {
	log.Println("Discord RPC client is not ready")
}
```

## API Reference

The module’s API is documented on [pkg.go.dev](https://pkg.go.dev/github.com/jacksonthemaster/discordrichpresence). Key components include:

- **`NewClient(appID string) *Client`**: Creates a new Discord RPC client.
- **`Client.StartWithActivity(activity Activity, updateInterval time.Duration) error`**: Starts the client with a periodic activity update.
- **`Client.SetActivity(activity Activity) error`**: Sets a new activity.
- **`Client.ClearActivity() error`**: Clears the current activity.
- **`Client.Close() error`**: Closes the connection and stops updates.
- **`NewActivity() *ActivityBuilder`**: Creates a new activity builder for fluent configuration.
- **`ActivityBuilder` Methods**: Chainable methods like `State`, `Details`, `StartTime`, `LargeImage`, etc.

## Development

### Building the Module
Since this is a library, there’s no binary to build. To use it in your project, simply include it as a dependency and import it as shown in the usage example.

### Directory Structure
```
.
├── client.go        # Core Discord RPC client logic
├── client_test.go   # Unit tests (optional, if added)
├── go.mod           # Module definition
├── LICENSE          # MIT License
├── README.md        # This file
└── types.go         # Types and ActivityBuilder
```

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Make your changes and add tests if applicable.
4. Commit your changes (`git commit -m "Add your feature"`).
5. Push to your branch (`git push origin feature/your-feature`).
6. Open a pull request on GitHub.

Please ensure your code follows Go conventions!
## Troubleshooting

- **"Failed to connect to Discord" Error**:
  - Ensure the Discord desktop client is running.
  - Verify the application ID is correct.
  - Check that your system supports Unix sockets (Linux/macOS) or named pipes (Windows).

- **Activity Not Updating**:
  - Ensure the update interval is reasonable (e.g., `30*time.Second`).
  - Check that the Discord client is not in a restricted mode (e.g., "Do Not Disturb").

- **Images Not Showing**:
  - Upload image assets to your Discord application in the Developer Portal.
  - Use the exact image keys (e.g., `logo`, `status`) as defined in the portal.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- None, **yet...**

## Contact

For questions or issues, open a n issue on the [GitHub Issues page](https://github.com/jacksonthemaster/discordrichpresence/issues).