package cqrs

import (
	"context"
	"github.com/dohr-michael/storyline-api/pkg/core"
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
)

type Aggregate interface {
	eh.Aggregate
	State() interface{}
	ChangeState(interface{})
}

type aggregate struct {
	*events.AggregateBase
	commandDispatcher CommandDispatcher
	eventDispatcher   EventDispatcher
	state             interface{}
}

var _ = eh.CommandHandler(&aggregate{})
var _ = events.Aggregate(&aggregate{})

func NewAggregate(
	t eh.AggregateType,
	id uuid.UUID,
	initialState interface{},
	dispatcher CommandDispatcher,
) Aggregate {
	return &aggregate{
		AggregateBase:     events.NewAggregateBase(t, id),
		commandDispatcher: dispatcher,
		state:             initialState,
	}
}

func (a *aggregate) HandleCommand(ctx context.Context, cmd eh.Command) error {
	evt, evtData, err := a.commandDispatcher.Apply(a.EntityID(), a.State(), cmd, ctx)
	if err != nil {
		return err
	}
	now, _ := core.Now(ctx)
	a.StoreEvent(evt, evtData, now)
	return nil
}

func (a *aggregate) ApplyEvent(ctx context.Context, evt eh.Event) error {
	out, err := a.eventDispatcher.Apply(a.State(), evt)
	if err != nil {
		return err
	}
	a.ChangeState(out)
	return nil
}

func (a *aggregate) State() interface{}            { return a.state }
func (a *aggregate) ChangeState(value interface{}) { a.state = value }
