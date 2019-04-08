package graphql

import (
	"github.com/dohr-michael/storyline-api/pkg/model"
	"github.com/graph-gophers/graphql-go"
	"time"
)

type userResolver struct {
	user *model.User
}

func (u *userResolver) ID() graphql.ID {
	return graphql.ID(u.user.ID)
}

func (u *userResolver) Name() string {
	return u.user.Name
}

func (u *userResolver) Email() string {
	return u.user.Email
}

func (u *userResolver) CreatedAt() (graphql.Time, error) {
	t, err := time.Parse(time.RFC3339, u.user.CreatedAt)
	return graphql.Time{Time: t}, err
}
