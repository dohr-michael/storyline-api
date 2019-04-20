package graphql

import (
	"github.com/dohr-michael/storyline-api/pkg/domain/universe"
	"github.com/graphql-go/graphql"
)

var universeType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Universe",
	Description: "",
	IsTypeOf: func(p graphql.IsTypeOfParams) bool {
		_, ok := p.Value.(*universe.Universe)
		return ok
	},
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Id of the universe",
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				if universe, ok := p.Source.(*universe.Universe); ok {
					return universe.Id, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Name of the user",
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				if universe, ok := p.Source.(*universe.Universe); ok {
					return universe.Name, nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.DateTime),
			Description: "Creation date",
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				if universe, ok := p.Source.(*universe.Universe); ok {
					return universe.CreatedAt, nil
				}
				return nil, nil
			},
		},
		"owner": &graphql.Field{
			Type:        graphql.NewNonNull(userType),
			Description: "Owner",
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				if universe, ok := p.Source.(*universe.Universe); ok {
					return universe.Owner, nil
				}
				return nil, nil
			},
		},
	},
})
