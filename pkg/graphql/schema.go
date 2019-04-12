package graphql

import (
	"context"
	"errors"
	"github.com/dohr-michael/storyline-api/pkg/repo"
	"github.com/graphql-go/graphql"
)

const (
	UserRepoKey = "UserRepoKey"
)

func userRepo(ctx context.Context) (repo.UserRepo, error) {
	res, ok := ctx.Value(UserRepoKey).(repo.UserRepo)
	if !ok {
		return nil, errors.New("user repo not set")
	}
	return res, nil
}

func NewSchema() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
}
