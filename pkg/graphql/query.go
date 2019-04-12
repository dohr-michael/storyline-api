package graphql

import (
	"github.com/graphql-go/graphql"
)

var query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"userByEmail": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Description: "Email of the user",
					Type:        graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				repo, err := userRepo(p.Context)
				if err != nil {
					return nil, err
				}

				return repo.FetchOne(p.Args["email"].(string), p.Context)
			},
		},
	},
})
