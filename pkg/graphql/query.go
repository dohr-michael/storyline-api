package graphql

import (
	lgraphql "github.com/dohr-michael/go-libs/graphql"
	"github.com/dohr-michael/storyline-api/pkg/domain/universe"
	"github.com/dohr-michael/storyline-api/pkg/domain/user"
	"github.com/graphql-go/graphql"
)

func searchUsersQuery() *graphql.Field {
	return lgraphql.PagedQuery("Users", userType, UserRepoKey)
}

func userByEmailQuery() *graphql.Field {
	return lgraphql.ById("email", userType, UserRepoKey)
}

func searchUniversesQuery() *graphql.Field {
	return lgraphql.PagedQuery("Universes", universeType, UniverseRepoKey)
}

func universeByIdQuery() *graphql.Field {
	base := lgraphql.ById("id", universeType, UniverseRepoKey)
	// baseResolver := base.Resolve
	base.Resolve = func(p graphql.ResolveParams) (i interface{}, e error) {
		return &universe.Universe{
			Owner: &user.User{},
		}, nil
	}
	return base
}

func listUniverseTags() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String))),
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			repo, err := universeRepo(p.Context)
			if err != nil {
				return nil, err
			}
			tags, err := repo.FetchTags(p.Context)
			if err != nil {
				return nil, err
			}
			var res []string
			for _, r := range tags {
				res = append(res, r.Name)
			}
			return res, nil
		},
	}
}

var query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"users":        searchUsersQuery(),
		"universes":    searchUniversesQuery(),
		"universeTags": listUniverseTags(),
		"user":         userByEmailQuery(),
		"universe":     universeByIdQuery(),
	},
})
