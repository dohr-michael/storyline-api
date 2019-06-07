package graphql

import (
	"context"
	"github.com/dohr-michael/go-libs/filters"
	"github.com/graphql-go/graphql"
)

func universesQuery() *graphql.Field {
	return pagedQuery("Universes", universeType, func(query *filters.Query, ctx context.Context) (interface{}, error) {
		rep, err := universeRepo(ctx)
		if err != nil {
			return nil, err
		}
		return rep.FetchMany(query, ctx)
	})
}

func universeQuery() *graphql.Field {
	return byIdQuery("id", universeType, func(id string, ctx context.Context) (interface{}, error) {
		rep, err := universeRepo(ctx)
		if err != nil {
			return nil, err
		}
		return rep.FetchOne(id, ctx)
	})
}

func universeTagsQuery() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String))),
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			rep, err := universeRepo(p.Context)
			if err != nil {
				return nil, err
			}
			tags, err := rep.FetchTags(p.Context)
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
		"universe":     universeQuery(),
		"universes":    universesQuery(),
		"universeTags": universeTagsQuery(),
	},
})
