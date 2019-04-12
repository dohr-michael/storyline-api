package graphql

import (
	"github.com/dohr-michael/storyline-api/pkg/model"
	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "",
	Fields: graphql.Fields{
		"email": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Email of the user",
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				if user, ok := p.Source.(*model.User); ok {
					return user.Id, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Name of the user",
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				if user, ok := p.Source.(*model.User); ok {
					return user.Name, nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.DateTime),
			Description: "Creation date",
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				if user, ok := p.Source.(*model.User); ok {
					return user.CreatedAt, nil
				}
				return nil, nil
			},
		},
	},
})
