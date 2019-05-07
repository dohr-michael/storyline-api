package graphql

import (
	"github.com/dohr-michael/storyline-api/pkg/domain/user"
	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "",
	IsTypeOf: func(p graphql.IsTypeOfParams) bool {
		_, ok := p.Value.(*user.User)
		return ok
	},
	Fields: graphql.Fields{
		"email": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Email of the user",
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				if user, ok := p.Source.(*user.User); ok {
					return user.Email, nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Name of the user",
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				if user, ok := p.Source.(*user.User); ok {
					return user.Name, nil
				}
				return nil, nil
			},
		},
		"picture": &graphql.Field{
			Type:        graphql.String,
			Description: "Picture of the user, can be url of base64 image",
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				if user, ok := p.Source.(*user.User); ok {
					return user.Picture, nil
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.DateTime),
			Description: "Creation date",
			Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
				if user, ok := p.Source.(*user.User); ok {
					return user.CreatedAt, nil
				}
				return nil, nil
			},
		},
	},
})
