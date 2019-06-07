package universe

import (
	"github.com/dohr-michael/storyline-api/pkg/domain/user"
	"time"
)

// Read model
type Tag struct {
	Name string `json:"name" mapstructure:"name"`
}

type Universe struct {
	Id          string     `json:"id" mapstructure:"_key"`
	Name        string     `json:"name" mapstructure:"name"`
	Description string     `json:"description" mapstructure:"description"`
	Picture     string     `json:"picture" mapstructure:"picture"`
	Tags        []string   `json:"tags" mapstructure:"tags"`
	Owner       *user.User `json:"-" mapstructure:"owner"`
	CreatedAt   time.Time  `json:"createdAt" mapstructure:"createdAt"`
	CreatedBy   *user.User `json:"-" mapstructure:"createdBy"`
}

type Universes []*Universe

type PagedUniverse struct {
	Items  *Universes
	Total  int64
	Limit  int64
	Offset int64
}

func (t *PagedUniverse) GetItems() interface{} {
	return t.Items
}

func (t *PagedUniverse) GetTotal() int64 {
	return t.Total
}

func (t *PagedUniverse) GetLimit() int64 {
	return t.Limit
}

func (t *PagedUniverse) GetOffset() int64 {
	return t.Offset
}
