package entity

import (
	"github.com/dohr-michael/storyline-api/pkg/domain/universe"
	"github.com/dohr-michael/storyline-api/pkg/domain/user"
	"time"
)

type Type string

const (
	Character    = Type("Character")
	Place        = Type("Place")
	Organization = Type("Organization")
)

type Entity struct {
	Id          string             `json:"id" mapstructure:"_key"`
	Name        string             `json:"name" mapstructure:"name"`
	Description string             `json:"description" mapstructure:"description"`
	BelongTo    *universe.Universe `json:"-" mapstructure:"belongTo"`
	Owner       *user.User         `json:"-" mapstructure:"owner"`
	CreatedAt   time.Time          `json:"createdAt" mapstructure:"createdAt"`
	CreatedBy   *user.User         `json:"-" mapstructure:"createdBy"`
}

type CharacterEntity struct {
	Entity
	Origin string `json:"origin" mapstructure:"origin"`
	Age    int    `json:"age" mapstructure:"age"`
}

type PlaceEntity struct {
	Entity
}

type OrganizationEntity struct {
	Entity
}


