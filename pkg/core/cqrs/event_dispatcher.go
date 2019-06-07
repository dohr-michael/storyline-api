package cqrs

import (
	eh "github.com/looplab/eventhorizon"
)

type EventDispatcher interface {
	Apply(state StateIn, event eh.Event) (StateOut, error)
}
