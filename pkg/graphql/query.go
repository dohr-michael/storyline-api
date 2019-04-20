package graphql

import (
	lgraphql "github.com/dohr-michael/go-libs/graphql"
	"github.com/dohr-michael/storyline-api/pkg/domain/universe"
	"github.com/dohr-michael/storyline-api/pkg/domain/user"
	"github.com/graphql-go/graphql"
	"log"
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
		log.Printf("%v", p.Info.Path.AsArray())
		return &universe.Universe{
			Owner: &user.User{},
		}, nil
	}
	return base
}

var query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"searchUsers":     searchUsersQuery(),
		"searchUniverses": searchUniversesQuery(),
		"userByEmail":     userByEmailQuery(),
		"universeById":    universeByIdQuery(),
	},
})
