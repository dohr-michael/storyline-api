package graphql

import (
	"github.com/dohr-michael/storyline-api/pkg/model"
	"github.com/graphql-go/graphql"
)

var mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Description: "Name of the user",
					Type:        graphql.NewNonNull(graphql.String),
				},
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
				_, res, err := repo.Create(
					&model.CreateUser{
						Email: p.Args["email"].(string),
						Name:  p.Args["name"].(string),
					},
					p.Context)
				return res, err
			},
		},
	},
})
