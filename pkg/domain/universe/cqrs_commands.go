package universe

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

const (
	CreateCmd      = eh.CommandType("storyline:universe:create")
	ChangeOwnerCmd = eh.CommandType("storyline:universe:change-owner")
)

func init() {
	eh.RegisterCommand(func() eh.Command { return &Create{} })
	eh.RegisterCommand(func() eh.Command { return &ChangeOwner{} })
}

var _ = eh.Command(&Create{})
var _ = eh.Command(&ChangeOwner{})

type Create struct {
	Id          uuid.UUID
	Name        string   `mapstructure:"name"`
	Description string   `mapstructure:"description"`
	Picture     string   `mapstructure:"picture,omitempty"`
	Tags        []string `mapstructure:"tags"`
}

func (c *Create) AggregateID() uuid.UUID        { return c.Id }
func (*Create) AggregateType() eh.AggregateType { return AggregateType }
func (*Create) CommandType() eh.CommandType     { return CreateCmd }

type ChangeOwner struct {
	Id       uuid.UUID
	NewOwner string `mapstructure:"newOwner"`
}

func (c *ChangeOwner) AggregateID() uuid.UUID        { return c.Id }
func (*ChangeOwner) AggregateType() eh.AggregateType { return AggregateType }
func (*ChangeOwner) CommandType() eh.CommandType     { return ChangeOwnerCmd }
