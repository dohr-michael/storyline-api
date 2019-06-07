package graphql

import "github.com/graphql-go/graphql"

var entityType = graphql.NewInterface(graphql.InterfaceConfig{
	Name: "Entity",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"description": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"belongTo": &graphql.Field{
			Type: graphql.NewNonNull(universeType),
		},
		"owner": &graphql.Field{
			Type: graphql.NewNonNull(userType),
		},
	},
})
