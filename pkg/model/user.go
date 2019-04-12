package model

import (
	"github.com/dohr-michael/storyline-api/pkg/core/data/arango"
	"time"
)

type User struct {
	// Correspond to email.
	Id        string            `json:"email" mapstructure:"_key"`
	Name      string            `json:"name" mapstructure:"name"`
	CreatedAt time.Time         `json:"createdAt" mapstructure:"createdAt"`
	Meta      arango.Identifier `json:"-"`
}

type Users []*User

type CreateUser struct {
	Email string `json:"_key"`
	Name  string `json:"name"`
}
