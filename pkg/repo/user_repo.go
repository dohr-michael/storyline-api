package repo

import (
	"github.com/arangodb/go-driver"
	"github.com/dohr-michael/storyline-api/pkg/core/data"
	"github.com/dohr-michael/storyline-api/pkg/core/data/arango"
	"github.com/dohr-michael/storyline-api/pkg/model"
)

const UserCollection = "Users"

type UserRepo interface {
	data.Repository
}

func NewUserRepo() (UserRepo, error) {
	base, err := arango.NewRepository(
		UserCollection,
		driver.CollectionTypeDocument,
		func() interface{} {
			return &model.User{}
		}, func() interface{} {
			return &model.Users{}
		},
	)
	if err != nil {
		return nil, err
	}
	return &userRepo{
		Repository: base,
	}, nil
}

var _ = UserRepo(&userRepo{})

type userRepo struct {
	*arango.Repository
}
