package graphql

import (
	"github.com/dohr-michael/storyline-api/pkg/domain/universe"
	"github.com/graphql-go/graphql"
)

var mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createUniverse": &graphql.Field{
			Type: graphql.NewNonNull(universeType),
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Description: "Name of the universe",
					Type:        graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				repo, err := universeRepo(p.Context)
				if err != nil {
					return nil, err
				}
				return repo.Create(
					&universe.Create{
						Name: p.Args["name"].(string),
					},
					p.Context,
				)
			},
		},
	},
})
