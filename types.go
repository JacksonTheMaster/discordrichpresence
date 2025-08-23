package discordrichpresence

import (
	"fmt"
	"time"
)

type RPCMessage struct {
	Cmd   string      `json:"cmd,omitempty"`
	Args  interface{} `json:"args,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Evt   string      `json:"evt,omitempty"`
	Nonce string      `json:"nonce,omitempty"`
}

type Activity struct {
	State      string     `json:"state,omitempty"`
	Details    string     `json:"details,omitempty"`
	Timestamps Timestamps `json:"timestamps,omitempty"`
	Assets     Assets     `json:"assets,omitempty"`
	Type       int        `json:"type"`
}

type Timestamps struct {
	Start int64 `json:"start,omitempty"`
	End   int64 `json:"end,omitempty"`
}

type Assets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

// ActivityBuilder provides a fluent interface for building activities
type ActivityBuilder struct {
	activity Activity
}

// NewActivity creates a new ActivityBuilder
func NewActivity() *ActivityBuilder {
	return &ActivityBuilder{
		activity: Activity{Type: 0}, // Default to "Playing"
	}
}

func (ab *ActivityBuilder) State(state string) *ActivityBuilder {
	ab.activity.State = state
	return ab
}

func (ab *ActivityBuilder) Details(details string) *ActivityBuilder {
	ab.activity.Details = details
	return ab
}

func (ab *ActivityBuilder) StartTime(start time.Time) *ActivityBuilder {
	ab.activity.Timestamps.Start = start.Unix()
	return ab
}

func (ab *ActivityBuilder) EndTime(end time.Time) *ActivityBuilder {
	ab.activity.Timestamps.End = end.Unix()
	return ab
}

func (ab *ActivityBuilder) LargeImage(image, text string) *ActivityBuilder {
	ab.activity.Assets.LargeImage = image
	ab.activity.Assets.LargeText = text
	return ab
}

func (ab *ActivityBuilder) SmallImage(image, text string) *ActivityBuilder {
	ab.activity.Assets.SmallImage = image
	ab.activity.Assets.SmallText = text
	return ab
}

func (ab *ActivityBuilder) Type(activityType int) *ActivityBuilder {
	ab.activity.Type = activityType
	return ab
}

func (ab *ActivityBuilder) Build() Activity {
	return ab.activity
}

// Utility function
func generateNonce() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
