package universe

import (
	eh "github.com/looplab/eventhorizon"
	"time"
)

const (
	Created      = eh.EventType("storyline:universe:created")
	OwnerChanged = eh.EventType("storyline:universe:owner-changed")
)

func init() {
	eh.RegisterEventData(Created, func() eh.EventData { return &CreatedData{} })
	eh.RegisterEventData(OwnerChanged, func() eh.EventData { return &OwnerChangedData{} })
}

var _ = eh.EventData(&Create{})
var _ = eh.EventData(&OwnerChangedData{})

type EventLog struct {
	By string    `mapstructure:"by"`
	At time.Time `mapstructure:"at"`
}

type CreatedData struct {
	EventLog
	Name        string   `mapstructure:"name"`
	Description string   `mapstructure:"description"`
	Picture     string   `mapstructure:"picture,omitempty"`
	Tags        []string `mapstructure:"tags"`
}

type OwnerChangedData struct {
	EventLog
	NewOwner string `mapstructure:"newOwner"`
}
