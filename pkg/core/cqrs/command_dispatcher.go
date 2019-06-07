package cqrs

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

type CommandDispatcher interface {
	Apply(id uuid.UUID, state interface{}, command eh.Command, ctx context.Context) (eh.EventType, interface{}, error)
}

type commandDispatcher struct {
	handlers []CommandHandler
}

func (c *commandDispatcher) Apply(id uuid.UUID, state interface{}, command eh.Command, ctx context.Context) (eh.EventType, interface{}, error) {
	for _, handler := range c.handlers {
		if handler.IsDefinedAt(state, command) {
			if err := handler.IsAvailableAt(state, ctx); err != nil {
				return eh.EventType(""), nil, err
			}
			return handler.Apply(id, state, command, ctx)
		}
	}
	return eh.EventType(""), nil, fmt.Errorf("no command handler found for %s", command.CommandType())
}

func NewCommandDIspatcher(handlers ...CommandHandler) CommandDispatcher {
	return &commandDispatcher{
		handlers: handlers,
	}
}
