package cqrs

import (
	"context"
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"log"
)

type CommandHandler interface {
	IsDefinedAt(state interface{}, command eh.Command) bool
	IsAvailableAt(state interface{}, ctx context.Context) error
	Apply(id uuid.UUID, state interface{}, command eh.Command, ctx context.Context) (eh.EventType, interface{}, error)
}

type CommandHandlerBuilder func(*commandHandler) *commandHandler

func NewCommandHandler(
	commandType eh.CommandType,
	apply func(uuid.UUID, interface{}, eh.Command, context.Context) (eh.EventType, interface{}, error),
	configs ...CommandHandlerBuilder,
) CommandHandler {
	res := &commandHandler{
		commandType: commandType,
		apply:       apply,
	}
	for _, config := range configs {
		res = config(res)
	}
	if res.apply == nil {
		log.Panic("apply function in not initialized")
	}
	if res.isAvailableAt == nil {
		res.isAvailableAt = func(interface{}, context.Context) error { return nil }
	}
	if res.isDefinedAt == nil {
		res.isDefinedAt = func(interface{}, eh.Command) bool { return true }
	}
	return res
}

type commandHandler struct {
	commandType   eh.CommandType
	isDefinedAt   func(state interface{}, command eh.Command) bool
	isAvailableAt func(state interface{}, ctx context.Context) error
	apply         func(id uuid.UUID, state interface{}, command eh.Command, ctx context.Context) (eh.EventType, interface{}, error)
}

func (c *commandHandler) IsDefinedAt(state interface{}, command eh.Command) bool {
	return c.commandType == command.CommandType() && c.isDefinedAt(state, command)
}
func (c *commandHandler) IsAvailableAt(state interface{}, ctx context.Context) error {
	return c.isAvailableAt(state, ctx)
}
func (c *commandHandler) Apply(id uuid.UUID, state interface{}, command eh.Command, ctx context.Context) (eh.EventType, interface{}, error) {
	return c.apply(id, state, command, ctx)
}
