package graphql

import (
	"context"
	"errors"
	"github.com/dohr-michael/storyline-api/pkg/domain/universe"
	"github.com/dohr-michael/storyline-api/pkg/domain/user"
	"github.com/graphql-go/graphql"
)

const (
	UserRepoKey     = "UserRepoKey"
	UniverseRepoKey = "UniverseRepoKey"
)

func userRepo(ctx context.Context) (user.Handlers, error) {
	res, ok := ctx.Value(UserRepoKey).(user.Handlers)
	if !ok {
		return nil, errors.New(UserRepoKey + " is not set")
	}
	return res, nil
}

func universeRepo(ctx context.Context) (universe.Handlers, error) {
	res, ok := ctx.Value(UniverseRepoKey).(universe.Handlers)
	if !ok {
		return nil, errors.New(UniverseRepoKey + " is not set")
	}
	return res, nil
}

func NewSchema() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
}
