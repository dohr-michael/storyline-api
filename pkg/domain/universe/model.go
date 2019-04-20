package universe

import (
	"github.com/dohr-michael/storyline-api/pkg/domain/user"
	"time"
)

// Read model
type Universe struct {
	Id        string     `json:"id" mapstructure:"_key"`
	Name      string     `json:"name" mapstructure:"name"`
	CreatedAt time.Time  `json:"createdAt" mapstructure:"createdAt"`
	Owner     *user.User `json:"-" mapstructure:"owner"`
	CreatedBy *user.User `json:"-" mapstructure:"createdBy"`
}
type Universes []*Universe
