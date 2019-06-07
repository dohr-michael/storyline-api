package universe

import (
	"github.com/dohr-michael/storyline-api/pkg/core/cqrs"
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

const AggregateType = eh.AggregateType("storyline:universe")

func init() {
	eh.RegisterAggregate(func(id uuid.UUID) eh.Aggregate {
		return &Aggregate{Aggregate: cqrs.NewAggregate(AggregateType, id, nil, nil)}
	})
}

var _ = eh.Aggregate(&Aggregate{})

type Aggregate struct {
	cqrs.Aggregate
	owner   string
	created bool
}

var createHandler = cqrs.NewCommandHandler()
